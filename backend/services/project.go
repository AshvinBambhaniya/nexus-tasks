package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/structs"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ProjectService defines the interface for project-related business logic
type ProjectService interface {
	// Project Management
	CreateProject(userID, workspaceID uuid.UUID, req structs.ReqCreateProject) (models.Project, error)
	GetProject(userID, projectID uuid.UUID) (models.Project, error)
	UpdateProject(userID, projectID uuid.UUID, req structs.ReqUpdateProject) (models.Project, error)
	AddMember(userID, projectID uuid.UUID, req structs.ReqAddProjectMember) (models.ProjectMember, error)
	RemoveMember(userID, projectID, targetUserID uuid.UUID) error
	ListMembers(userID, projectID uuid.UUID) ([]structs.ResProjectMember, error)
	ListByWorkspaceID(workspaceID uuid.UUID) ([]models.Project, error)

	// Team Management
	AddTeam(userID, projectID, teamID uuid.UUID) (structs.ResProjectTeam, error)
	RemoveTeam(userID, projectID, teamID uuid.UUID) error
	ListTeams(userID, projectID uuid.UUID) ([]structs.ResProjectTeam, error)

	// Helpers
	ValidateProjectAccess(projectID, userID uuid.UUID, requireAdmin bool) error
}

type projectService struct {
	storage models.Storage
	logger  *zap.Logger
}

// NewProjectService creates a new project service instance
func NewProjectService(storage models.Storage, logger *zap.Logger) ProjectService {
	return &projectService{
		storage: storage,
		logger:  logger,
	}
}

func (s *projectService) CreateProject(userID, workspaceID uuid.UUID, req structs.ReqCreateProject) (models.Project, error) {
	// 1. Verify Workspace Admin
	wsMember, err := s.storage.Workspaces().GetMember(workspaceID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Project{}, errors.New("unauthorized: not a member of this workspace")
		}
		return models.Project{}, errors.New("failed to get workspace member")
	}
	if wsMember.Role != models.WorkspaceRoleAdmin {
		return models.Project{}, errors.New("unauthorized: only workspace admins can create projects")
	}

	var createdProject models.Project
	err = s.storage.Atomic(context.Background(), func(txStorage models.Storage) error {
		project := models.Project{
			Name:        req.Name,
			Description: req.Description,
			WorkspaceID: workspaceID,
			IsArchived:  false,
			CreatedAt:   time.Now(),
		}

		var err error
		createdProject, err = txStorage.Projects().Create(project)
		if err != nil {
			return err
		}

		// Add Creator as Admin
		err = txStorage.Projects().AddMember(models.ProjectMember{
			ProjectID: createdProject.ID,
			UserID:    userID,
			Role:      models.ProjectRoleAdmin,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return models.Project{}, err
	}

	return createdProject, nil
}

func (s *projectService) GetProject(userID, projectID uuid.UUID) (models.Project, error) {
	err := s.ValidateProjectAccess(projectID, userID, false)
	if err != nil {
		return models.Project{}, err
	}
	return s.storage.Projects().GetByID(projectID)
}

func (s *projectService) UpdateProject(userID, projectID uuid.UUID, req structs.ReqUpdateProject) (models.Project, error) {
	err := s.ValidateProjectAccess(projectID, userID, true)
	if err != nil {
		return models.Project{}, err
	}

	var updatedProject models.Project
	err = s.storage.Atomic(context.Background(), func(txStorage models.Storage) error {
		project, err := txStorage.Projects().GetByID(projectID)
		if err != nil {
			return err
		}

		if req.Name != "" {
			project.Name = req.Name
		}
		if req.Description != "" {
			project.Description = req.Description
		}
		if req.IsArchived != nil {
			project.IsArchived = *req.IsArchived
		}

		updatedProject, err = txStorage.Projects().Update(project)
		return err
	})

	if err != nil {
		return models.Project{}, err
	}

	return updatedProject, nil
}

func (s *projectService) AddMember(userID, projectID uuid.UUID, req structs.ReqAddProjectMember) (models.ProjectMember, error) {
	err := s.ValidateProjectAccess(projectID, userID, true)
	if err != nil {
		return models.ProjectMember{}, err
	}

	// Find User
	userToAdd, err := s.storage.Users().GetByEmail(req.Email)
	if err != nil {
		return models.ProjectMember{}, errors.New("user not found")
	}

	// Get Project to check workspace
	project, err := s.storage.Projects().GetByID(projectID)
	if err != nil {
		return models.ProjectMember{}, err
	}

	// Check if user in workspace
	_, err = s.storage.Workspaces().GetMember(project.WorkspaceID, userToAdd.ID)
	if err != nil {
		return models.ProjectMember{}, errors.New("user must be a member of the workspace first")
	}

	// Check if already in project
	_, err = s.storage.Projects().GetMember(projectID, userToAdd.ID)
	if err == nil {
		return models.ProjectMember{}, errors.New("user already in project")
	}

	role := req.Role
	if role == "" {
		role = models.ProjectRoleMember
	}

	member := models.ProjectMember{
		ProjectID: projectID,
		UserID:    userToAdd.ID,
		Role:      role,
	}

	err = s.storage.Atomic(context.Background(), func(txStorage models.Storage) error {
		return txStorage.Projects().AddMember(member)
	})

	if err != nil {
		return models.ProjectMember{}, err
	}

	return member, nil
}

func (s *projectService) RemoveMember(userID, projectID, targetUserID uuid.UUID) error {
	err := s.ValidateProjectAccess(projectID, userID, true)
	if err != nil {
		return err
	}

	return s.storage.Atomic(context.Background(), func(txStorage models.Storage) error {
		return txStorage.Projects().RemoveMember(projectID, targetUserID)
	})
}

func (s *projectService) ListMembers(userID, projectID uuid.UUID) ([]structs.ResProjectMember, error) {
	err := s.ValidateProjectAccess(projectID, userID, false)
	if err != nil {
		return nil, err
	}

	// 1. Direct Members
	directMembers, err := s.storage.Projects().GetMembers(projectID)
	if err != nil {
		return nil, err
	}

	// 2. Team Members
	projectTeams, err := s.storage.Projects().GetTeams(projectID)
	if err != nil {
		return nil, err
	}

	allMembers := make(map[uuid.UUID]structs.ResProjectMember)

	for _, m := range directMembers {
		allMembers[m.UserID] = structs.ResProjectMember{
			UserID:   m.UserID,
			Email:    m.Email,
			FullName: m.FullName,
			Role:     m.Role,
			IsDirect: true,
		}
	}

	for _, pt := range projectTeams {
		teamMembers, err := s.storage.Teams().ListMembersByTeamID(pt.TeamID)
		if err != nil {
			continue
		}
		for _, tm := range teamMembers {
			if _, exists := allMembers[tm.UserID]; !exists {
				allMembers[tm.UserID] = structs.ResProjectMember{
					UserID:   tm.UserID,
					Email:    tm.Email,
					FullName: tm.FullName,
					Role:     models.ProjectRoleMember, // Implicit role
					IsDirect: false,
				}
			}
		}
	}

	// Convert map to slice
	var res []structs.ResProjectMember
	for _, m := range allMembers {
		res = append(res, m)
	}

	return res, nil
}

func (s *projectService) ListByWorkspaceID(workspaceID uuid.UUID) ([]models.Project, error) {
	return s.storage.Projects().ListByWorkspaceID(workspaceID)
}

// Teams

func (s *projectService) AddTeam(userID, projectID, teamID uuid.UUID) (structs.ResProjectTeam, error) {
	err := s.ValidateProjectAccess(projectID, userID, true)
	if err != nil {
		return structs.ResProjectTeam{}, err
	}

	// Verify Team exists
	team, err := s.storage.Teams().GetByID(teamID)
	if err != nil {
		return structs.ResProjectTeam{}, errors.New("team not found")
	}

	// Verify Project
	project, err := s.storage.Projects().GetByID(projectID)
	if err != nil {
		return structs.ResProjectTeam{}, errors.New("project not found")
	}

	if team.WorkspaceID != project.WorkspaceID {
		return structs.ResProjectTeam{}, errors.New("team must belong to the same workspace")
	}

	// Check if already assigned
	_, err = s.storage.Projects().GetTeam(projectID, teamID)
	if err == nil {
		return structs.ResProjectTeam{}, errors.New("team already assigned to project")
	}

	err = s.storage.Atomic(context.Background(), func(txStorage models.Storage) error {
		return txStorage.Projects().AddTeam(models.ProjectTeam{
			ProjectID: projectID,
			TeamID:    teamID,
		})
	})

	if err != nil {
		return structs.ResProjectTeam{}, err
	}

	return structs.ResProjectTeam{
		ProjectID: projectID,
		TeamID:    teamID,
		TeamName:  team.Name,
	}, nil
}

func (s *projectService) RemoveTeam(userID, projectID, teamID uuid.UUID) error {
	err := s.ValidateProjectAccess(projectID, userID, true)
	if err != nil {
		return err
	}

	return s.storage.Atomic(context.Background(), func(txStorage models.Storage) error {
		return txStorage.Projects().RemoveTeam(projectID, teamID)
	})
}

func (s *projectService) ListTeams(userID, projectID uuid.UUID) ([]structs.ResProjectTeam, error) {
	err := s.ValidateProjectAccess(projectID, userID, false)
	if err != nil {
		return nil, err
	}

	teams, err := s.storage.Projects().GetTeams(projectID)
	if err != nil {
		return nil, err
	}

	var res []structs.ResProjectTeam
	for _, t := range teams {
		res = append(res, structs.ResProjectTeam{
			ProjectID: t.ProjectID,
			TeamID:    t.TeamID,
			TeamName:  t.TeamName,
		})
	}

	return res, nil
}

// Helpers

func (s *projectService) ValidateProjectAccess(projectID, userID uuid.UUID, requireAdmin bool) error {
	// 1. Direct Member
	member, err := s.storage.Projects().GetMember(projectID, userID)
	if err == nil {
		if !requireAdmin || member.Role == models.ProjectRoleAdmin {
			return nil // Authorized
		}
	}

	// 2. Workspace Admin
	project, err := s.storage.Projects().GetByID(projectID)
	if err != nil {
		return errors.New("project not found")
	}

	wsMember, err := s.storage.Workspaces().GetMember(project.WorkspaceID, userID)
	if err == nil && wsMember.Role == models.WorkspaceRoleAdmin {
		return nil // Authorized
	}

	// 3. Team Member (Implicit) - Only if not requiring admin
	if !requireAdmin {
		projectTeams, err := s.storage.Projects().GetTeams(projectID)
		if err == nil {
			for _, pt := range projectTeams {
				_, err := s.storage.Teams().GetMember(pt.TeamID, userID)
				if err == nil {
					return nil // Authorized as Member
				}
			}
		}
	}

	return errors.New("unauthorized")
}
