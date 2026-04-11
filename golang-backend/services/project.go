package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProjectService struct {
	projectModel   *models.ProjectModel
	workspaceModel *models.WorkspaceModel
	teamModel      *models.TeamModel
	userModel      *models.UserModel
	db             *goqu.Database
	logger         *zap.Logger
}

func NewProjectService(db *goqu.Database, logger *zap.Logger, projectModel *models.ProjectModel, workspaceModel *models.WorkspaceModel, teamModel *models.TeamModel, userModel *models.UserModel) *ProjectService {
	return &ProjectService{
		projectModel:   projectModel,
		workspaceModel: workspaceModel,
		teamModel:      teamModel,
		userModel:      userModel,
		db:             db,
		logger:         logger,
	}
}

func (s *ProjectService) CreateProject(userID, workspaceID uuid.UUID, req structs.ReqCreateProject) (models.Project, error) {
	// 1. Verify Workspace Admin
	wsMember, err := s.workspaceModel.GetMember(workspaceID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Project{}, errors.New("unauthorized: not a member of this workspace")
		}
		return models.Project{}, err
	}
	if wsMember.Role != models.WorkspaceRoleAdmin {
		return models.Project{}, errors.New("unauthorized: only workspace admins can create projects")
	}

	project := models.Project{
		Name:        req.Name,
		Description: req.Description,
		WorkspaceID: workspaceID,
		IsArchived:  false,
		CreatedAt:   time.Now(),
	}

	isOk := false
	transaction, err := s.db.Begin()
	if err != nil {
		return models.Project{}, err
	}
	defer func() {
		if isOk {
			transaction.Commit()
		} else {
			transaction.Rollback()
		}
	}()

	createdProject, err := s.projectModel.Create(transaction, project)
	if err != nil {
		return models.Project{}, err
	}

	// Add Creator as Admin
	err = s.projectModel.AddMember(transaction, models.ProjectMember{
		ProjectID: createdProject.ID,
		UserID:    userID,
		Role:      models.ProjectRoleAdmin,
	})
	if err != nil {
		return models.Project{}, err
	}

	isOk = true
	return createdProject, nil
}

func (s *ProjectService) GetProject(userID, projectID uuid.UUID) (models.Project, error) {
	err := s.validateProjectAccess(projectID, userID, false)
	if err != nil {
		return models.Project{}, err
	}
	return s.projectModel.GetByID(projectID)
}

func (s *ProjectService) UpdateProject(userID, projectID uuid.UUID, req structs.ReqUpdateProject) (models.Project, error) {
	err := s.validateProjectAccess(projectID, userID, true)
	if err != nil {
		return models.Project{}, err
	}

	project, err := s.projectModel.GetByID(projectID)
	if err != nil {
		return models.Project{}, err
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

	isOk := false
	transaction, err := s.db.Begin()
	if err != nil {
		return models.Project{}, err
	}
	defer func() {
		if isOk {
			transaction.Commit()
		} else {
			transaction.Rollback()
		}
	}()

	updatedProject, err := s.projectModel.Update(transaction, project)
	if err != nil {
		return models.Project{}, err
	}

	isOk = true
	return updatedProject, nil
}

func (s *ProjectService) AddMember(userID, projectID uuid.UUID, req structs.ReqAddProjectMember) (models.ProjectMember, error) {
	err := s.validateProjectAccess(projectID, userID, true)
	if err != nil {
		return models.ProjectMember{}, err
	}

	// Find User
	userToAdd, err := s.userModel.GetByEmail(req.Email)
	if err != nil {
		return models.ProjectMember{}, errors.New("user not found")
	}

	// Get Project to check workspace
	project, err := s.projectModel.GetByID(projectID)
	if err != nil {
		return models.ProjectMember{}, err
	}

	// Check if user in workspace
	_, err = s.workspaceModel.GetMember(project.WorkspaceID, userToAdd.ID)
	if err != nil {
		return models.ProjectMember{}, errors.New("user must be a member of the workspace first")
	}

	// Check if already in project
	_, err = s.projectModel.GetMember(projectID, userToAdd.ID)
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

	isOk := false
	transaction, err := s.db.Begin()
	if err != nil {
		return models.ProjectMember{}, err
	}
	defer func() {
		if isOk {
			transaction.Commit()
		} else {
			transaction.Rollback()
		}
	}()

	err = s.projectModel.AddMember(transaction, member)
	if err != nil {
		return models.ProjectMember{}, err
	}

	isOk = true
	return member, nil
}

func (s *ProjectService) RemoveMember(userID, projectID, targetUserID uuid.UUID) error {
	err := s.validateProjectAccess(projectID, userID, true)
	if err != nil {
		return err
	}

	isOk := false
	transaction, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if isOk {
			transaction.Commit()
		} else {
			transaction.Rollback()
		}
	}()

	err = s.projectModel.RemoveMember(transaction, projectID, targetUserID)
	if err != nil {
		return err
	}

	isOk = true
	return nil
}

func (s *ProjectService) ListMembers(userID, projectID uuid.UUID) ([]structs.ResProjectMember, error) {
	err := s.validateProjectAccess(projectID, userID, false)
	if err != nil {
		return nil, err
	}

	// 1. Direct Members
	directMembers, err := s.projectModel.GetMembers(projectID)
	if err != nil {
		return nil, err
	}

	// 2. Team Members
	projectTeams, err := s.projectModel.GetTeams(projectID)
	if err != nil {
		return nil, err
	}

	allMembers := make(map[uuid.UUID]structs.ResProjectMember)

	for _, m := range directMembers {
		allMembers[m.UserID] = structs.ResProjectMember{
			UserID:   m.UserID,
			Email:    m.Email,
			Role:     m.Role,
			IsDirect: true,
		}
	}

	for _, pt := range projectTeams {
		teamMembers, err := s.teamModel.ListMembersByTeamId(pt.TeamID)
		if err != nil {
			continue
		}
		for _, tm := range teamMembers {
			if _, exists := allMembers[tm.UserID]; !exists {
				allMembers[tm.UserID] = structs.ResProjectMember{
					UserID:   tm.UserID,
					Email:    tm.Email,
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

// Teams

func (s *ProjectService) AddTeam(userID, projectID, teamID uuid.UUID) (structs.ResProjectTeam, error) {
	err := s.validateProjectAccess(projectID, userID, true)
	if err != nil {
		return structs.ResProjectTeam{}, err
	}

	// Verify Team exists
	team, err := s.teamModel.GetByID(teamID)
	if err != nil {
		return structs.ResProjectTeam{}, errors.New("team not found")
	}

	// Verify Project
	project, err := s.projectModel.GetByID(projectID)
	if err != nil {
		return structs.ResProjectTeam{}, errors.New("project not found")
	}

	if team.WorkspaceID != project.WorkspaceID {
		return structs.ResProjectTeam{}, errors.New("team must belong to the same workspace")
	}

	// Check if already assigned
	_, err = s.projectModel.GetTeam(projectID, teamID)
	if err == nil {
		return structs.ResProjectTeam{}, errors.New("team already assigned to project")
	}

	isOk := false
	transaction, err := s.db.Begin()
	if err != nil {
		return structs.ResProjectTeam{}, err
	}
	defer func() {
		if isOk {
			transaction.Commit()
		} else {
			transaction.Rollback()
		}
	}()

	err = s.projectModel.AddTeam(transaction, models.ProjectTeam{
		ProjectID: projectID,
		TeamID:    teamID,
	})
	if err != nil {
		return structs.ResProjectTeam{}, err
	}

	isOk = true
	return structs.ResProjectTeam{
		ProjectID: projectID,
		TeamID:    teamID,
		TeamName:  team.Name,
	}, nil
}

func (s *ProjectService) RemoveTeam(userID, projectID, teamID uuid.UUID) error {
	err := s.validateProjectAccess(projectID, userID, true)
	if err != nil {
		return err
	}

	isOk := false
	transaction, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if isOk {
			transaction.Commit()
		} else {
			transaction.Rollback()
		}
	}()

	err = s.projectModel.RemoveTeam(transaction, projectID, teamID)
	if err != nil {
		return err
	}

	isOk = true
	return nil
}

func (s *ProjectService) ListTeams(userID, projectID uuid.UUID) ([]structs.ResProjectTeam, error) {
	err := s.validateProjectAccess(projectID, userID, false)
	if err != nil {
		return nil, err
	}

	teams, err := s.projectModel.GetTeams(projectID)
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

func (s *ProjectService) validateProjectAccess(projectID, userID uuid.UUID, requireAdmin bool) error {
	// 1. Direct Member
	member, err := s.projectModel.GetMember(projectID, userID)
	if err == nil {
		if requireAdmin && member.Role != models.ProjectRoleAdmin {
			// Check Workspace Admin below
		} else {
			return nil // Authorized
		}
	}

	// 2. Workspace Admin
	project, err := s.projectModel.GetByID(projectID)
	if err != nil {
		return errors.New("project not found")
	}

	wsMember, err := s.workspaceModel.GetMember(project.WorkspaceID, userID)
	if err == nil && wsMember.Role == models.WorkspaceRoleAdmin {
		return nil // Authorized
	}

	// 3. Team Member (Implicit) - Only if not requiring admin
	if !requireAdmin {
		projectTeams, err := s.projectModel.GetTeams(projectID)
		if err == nil {
			for _, pt := range projectTeams {
				_, err := s.teamModel.GetMember(pt.TeamID, userID)
				if err == nil {
					return nil // Authorized as Member
				}
			}
		}
	}

	return errors.New("unauthorized")
}
