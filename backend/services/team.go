package services

import (
	"context"
	"errors"

	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TeamService interface {
	CreateTeam(userID uuid.UUID, workspaceID uuid.UUID, req structs.ReqCreateTeam) (models.Team, error)
	GetTeam(teamID uuid.UUID) (structs.ResTeamWithProjects, error)
	ListTeamsByWorkspaceID(workspaceID uuid.UUID) ([]models.Team, error)
	UpdateTeam(requestorID, workspaceID, teamID uuid.UUID, req structs.ReqUpdateTeam) (models.Team, error)
	DeleteTeam(requestorID, workspaceID, teamID uuid.UUID) error
	AddMember(requestorID, workspaceID, teamID uuid.UUID, email, role string) error
	RemoveMember(requestorID, workspaceID, teamID, userID uuid.UUID) error
	ListMembersByTeamId(teamID uuid.UUID) ([]models.TeamMemberWithUser, error)
}

type teamService struct {
	storage models.Storage
	logger  *zap.Logger
}

func NewTeamService(storage models.Storage, logger *zap.Logger) TeamService {
	return &teamService{
		storage: storage,
		logger:  logger,
	}
}

// CreateTeam creates a new team in a workspace and adds the creator as admin
func (s *teamService) CreateTeam(userID uuid.UUID, workspaceID uuid.UUID, req structs.ReqCreateTeam) (models.Team, error) {
	// 1. Verify User is Workspace Admin (Requirement from Python logic)
	member, err := s.storage.Workspaces().GetMember(workspaceID, userID)
	if err != nil {
		return models.Team{}, err
	}
	if member.Role != models.WorkspaceRoleAdmin {
		return models.Team{}, errors.New("unauthorized: only workspace admins can create teams")
	}

	var createdTeam models.Team
	err = s.storage.Atomic(context.Background(), func(txStorage models.Storage) error {
		team := models.Team{
			Name:        req.Name,
			Description: req.Description,
			WorkspaceID: workspaceID,
		}

		var err error
		createdTeam, err = txStorage.Teams().CreateTeam(team)
		if err != nil {
			s.logger.Error("failed to create team", zap.Error(err))
			return err
		}

		// Add creator as Team Admin
		err = txStorage.Teams().AddMember(models.TeamMember{
			TeamID: createdTeam.ID,
			UserID: userID,
			Role:   models.TeamRoleAdmin,
		})
		if err != nil {
			s.logger.Error("failed to add team member", zap.Error(err))
			return err
		}
		return nil
	})

	if err != nil {
		return models.Team{}, err
	}

	return createdTeam, nil
}

func (s *teamService) GetTeam(teamID uuid.UUID) (structs.ResTeamWithProjects, error) {

	team, err := s.storage.Teams().GetByID(teamID)
	if err != nil {
		s.logger.Error("error while get team by teamID", zap.Error(err))
		return structs.ResTeamWithProjects{}, err
	}

	projects, err := s.storage.Projects().ListByTeamID(teamID)
	if err != nil {
		s.logger.Error("error while list project by teamID", zap.Error(err))
		return structs.ResTeamWithProjects{}, err
	}

	return structs.ResTeamWithProjects{
		ResTeam: structs.ResTeam{
			ID:          team.ID,
			Name:        team.Name,
			Description: team.Description,
			WorkspaceID: team.WorkspaceID,
		},
		Projects: projects,
	}, err
}

func (s *teamService) ListTeamsByWorkspaceID(workspaceID uuid.UUID) ([]models.Team, error) {
	return s.storage.Teams().ListTeamsByWorkspaceID(workspaceID)
}

func (s *teamService) UpdateTeam(requestorID, workspaceID, teamID uuid.UUID, req structs.ReqUpdateTeam) (models.Team, error) {
	// Verify Access (Team Admin or Workspace Admin)
	// For simplicity, let's allow Workspace Admin or Team Admin.
	// Check Team Admin
	isTeamAdmin := false
	teamMember, err := s.storage.Teams().GetMember(teamID, requestorID)
	if err == nil && teamMember.Role == models.TeamRoleAdmin {
		isTeamAdmin = true
	}

	if !isTeamAdmin {
		// Check Workspace Admin
		wsMember, err := s.storage.Workspaces().GetMember(workspaceID, requestorID)
		if err != nil || wsMember.Role != models.WorkspaceRoleAdmin {
			return models.Team{}, errors.New("unauthorized: team admin or workspace admin required")
		}
	}

	team, err := s.storage.Teams().GetByID(teamID)
	if err != nil {
		return models.Team{}, err
	}

	if req.Name != "" {
		team.Name = req.Name
	}
	if req.Description != "" {
		team.Description = req.Description
	}

	updatedTeam, err := s.storage.Teams().UpdateTeam(team)
	if err != nil {
		s.logger.Error("error while update the team", zap.Error(err))
		return models.Team{}, err
	}

	return updatedTeam, nil
}

func (s *teamService) DeleteTeam(requestorID, workspaceID, teamID uuid.UUID) error {
	isAuthorized := false

	// 1. Check Team Admin
	teamMember, err := s.storage.Teams().GetMember(teamID, requestorID)
	if err == nil && teamMember.Role == models.TeamRoleAdmin {
		isAuthorized = true
	}

	// 2. Check Workspace Admin (Override)
	if !isAuthorized {
		wsMember, err := s.storage.Workspaces().GetMember(workspaceID, requestorID)
		if err == nil && wsMember.Role == models.WorkspaceRoleAdmin {
			isAuthorized = true
		}
	}

	if !isAuthorized {
		return errors.New("unauthorized")
	}

	err = s.storage.Teams().DeleteTeam(teamID)
	if err != nil {
		s.logger.Error("error while deleting the team", zap.Error(err))
		return err
	}

	return nil
}

func (s *teamService) AddMember(requestorID, workspaceID, teamID uuid.UUID, email, role string) error {
	// Verify Admin
	isAuthorized := false
	teamMember, err := s.storage.Teams().GetMember(teamID, requestorID)
	if err == nil && teamMember.Role == models.TeamRoleAdmin {
		isAuthorized = true
	}
	if !isAuthorized {
		wsMember, err := s.storage.Workspaces().GetMember(workspaceID, requestorID)
		if err == nil && wsMember.Role == models.WorkspaceRoleAdmin {
			isAuthorized = true
		}
	}
	if !isAuthorized {
		return errors.New("unauthorized")
	}

	// Find User
	user, err := s.storage.Users().GetByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if already in team
	_, err = s.storage.Teams().GetMember(teamID, user.ID)
	if err == nil {
		return errors.New("user already in team")
	}

	// Validate Role
	if role != string(models.TeamRoleAdmin) && role != string(models.TeamRoleMember) {
		return errors.New("invalid role")
	}

	return s.storage.Atomic(context.Background(), func(txStorage models.Storage) error {
		return txStorage.Teams().AddMember(models.TeamMember{
			TeamID: teamID,
			UserID: user.ID,
			Role:   models.TeamRole(role),
		})
	})
}

func (s *teamService) RemoveMember(requestorID, workspaceID, teamID, userID uuid.UUID) error {
	// Verify Admin
	isAuthorized := false
	teamMember, err := s.storage.Teams().GetMember(teamID, requestorID)
	if err == nil && teamMember.Role == models.TeamRoleAdmin {
		isAuthorized = true
	}
	if !isAuthorized {
		wsMember, err := s.storage.Workspaces().GetMember(workspaceID, requestorID)
		if err == nil && wsMember.Role == models.WorkspaceRoleAdmin {
			isAuthorized = true
		}
	}
	if !isAuthorized {
		return errors.New("unauthorized")
	}

	err = s.storage.Teams().RemoveMember(teamID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *teamService) ListMembersByTeamId(teamID uuid.UUID) ([]models.TeamMemberWithUser, error) {
	return s.storage.Teams().ListMembersByTeamId(teamID)
}
