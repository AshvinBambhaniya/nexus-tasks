package services

import (
	"errors"

	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TeamService struct {
	teamModel      *models.TeamModel
	projectModel   *models.ProjectModel
	workspaceModel *models.WorkspaceModel
	userModel      *models.UserModel
	db             *goqu.Database
	logger         *zap.Logger
}

func NewTeamService(db *goqu.Database, logger *zap.Logger, teamModel *models.TeamModel, projectModel *models.ProjectModel, workspaceModel *models.WorkspaceModel, userModel *models.UserModel) *TeamService {
	return &TeamService{
		teamModel:      teamModel,
		projectModel:   projectModel,
		workspaceModel: workspaceModel,
		userModel:      userModel,
		db:             db,
		logger:         logger,
	}
}

// CreateTeam creates a new team in a workspace and adds the creator as admin
func (s *TeamService) CreateTeam(userID uuid.UUID, workspaceID uuid.UUID, req structs.ReqCreateTeam) (models.Team, error) {
	// 1. Verify User is Workspace Admin (Requirement from Python logic)
	member, err := s.workspaceModel.GetMember(workspaceID, userID)
	if err != nil {
		return models.Team{}, err
	}
	if member.Role != models.WorkspaceRoleAdmin {
		return models.Team{}, errors.New("unauthorized: only workspace admins can create teams")
	}

	team := models.Team{
		Name:        req.Name,
		Description: req.Description,
		WorkspaceID: workspaceID,
	}

	isOk := false
	transaction, err := s.db.Begin()
	if err != nil {
		return models.Team{}, err
	}

	defer func() {
		if isOk {
			err := transaction.Commit()
			if err != nil {
				s.logger.Error("error during commit in create team", zap.Error(err))
			}
		} else {
			err := transaction.Rollback()
			if err != nil {
				s.logger.Error("error during rollback in create team", zap.Error(err))
			}
		}
	}()

	createdTeam, err := s.teamModel.CreateTeam(transaction, team)
	if err != nil {
		s.logger.Error("failed to create team", zap.Error(err))
		return createdTeam, err
	}

	// Add creator as Team Admin
	err = s.teamModel.AddMember(transaction, models.TeamMember{
		TeamID: createdTeam.ID,
		UserID: userID,
		Role:   models.TeamRoleAdmin,
	})
	if err != nil {
		s.logger.Error("failed to add team member", zap.Error(err))
		return createdTeam, err
	}

	isOk = true
	return createdTeam, nil
}

func (s *TeamService) GetTeam(teamID uuid.UUID) (structs.ResTeamWithProjects, error) {

	team, err := s.teamModel.GetByID(teamID)
	if err != nil {
		s.logger.Error("error while get team by teamID", zap.Error(err))
		return structs.ResTeamWithProjects{}, err
	}

	projects, err := s.projectModel.ListByTeamID(teamID)
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

func (s *TeamService) UpdateTeam(requestorID, workspaceID, teamID uuid.UUID, req structs.ReqUpdateTeam) (models.Team, error) {
	// Verify Access (Team Admin or Workspace Admin)
	// For simplicity, let's allow Workspace Admin or Team Admin.
	// Check Team Admin
	isTeamAdmin := false
	teamMember, err := s.teamModel.GetMember(teamID, requestorID)
	if err == nil && teamMember.Role == models.TeamRoleAdmin {
		isTeamAdmin = true
	}

	if !isTeamAdmin {
		// Check Workspace Admin
		wsMember, err := s.workspaceModel.GetMember(workspaceID, requestorID)
		if err != nil || wsMember.Role != models.WorkspaceRoleAdmin {
			return models.Team{}, errors.New("unauthorized: team admin or workspace admin required")
		}
	}

	team, err := s.teamModel.GetByID(teamID)
	if err != nil {
		return models.Team{}, err
	}

	if req.Name != "" {
		team.Name = req.Name
	}
	if req.Description != "" {
		team.Description = req.Description
	}

	updatedTeam, err := s.teamModel.UpdateTeam(team)
	if err != nil {
		s.logger.Error("error while update the team", zap.Error(err))
		return models.Team{}, err
	}

	return updatedTeam, nil
}

func (s *TeamService) DeleteTeam(requestorID, workspaceID, teamID uuid.UUID) error {
	isAuthorized := false

	// 1. Check Team Admin
	teamMember, err := s.teamModel.GetMember(teamID, requestorID)
	if err == nil && teamMember.Role == models.TeamRoleAdmin {
		isAuthorized = true
	}

	// 2. Check Workspace Admin (Override)
	if !isAuthorized {
		wsMember, err := s.workspaceModel.GetMember(workspaceID, requestorID)
		if err == nil && wsMember.Role == models.WorkspaceRoleAdmin {
			isAuthorized = true
		}
	}

	if !isAuthorized {
		return errors.New("unauthorized")
	}

	err = s.teamModel.DeleteTeam(teamID)
	if err != nil {
		s.logger.Error("error while deleting the team", zap.Error(err))
		return err
	}

	return nil
}

func (s *TeamService) AddMember(requestorID, workspaceID, teamID uuid.UUID, email, role string) error {
	// Verify Admin
	isAuthorized := false
	teamMember, err := s.teamModel.GetMember(teamID, requestorID)
	if err == nil && teamMember.Role == models.TeamRoleAdmin {
		isAuthorized = true
	}
	if !isAuthorized {
		wsMember, err := s.workspaceModel.GetMember(workspaceID, requestorID)
		if err == nil && wsMember.Role == models.WorkspaceRoleAdmin {
			isAuthorized = true
		}
	}
	if !isAuthorized {
		return errors.New("unauthorized")
	}

	// Find User
	user, err := s.userModel.GetByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if already in team
	_, err = s.teamModel.GetMember(teamID, user.ID)
	if err == nil {
		return errors.New("user already in team")
	}

	// Validate Role
	if role != string(models.TeamRoleAdmin) && role != string(models.TeamRoleMember) {
		return errors.New("invalid role")
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

	err = s.teamModel.AddMember(transaction, models.TeamMember{
		TeamID: teamID,
		UserID: user.ID,
		Role:   models.TeamRole(role),
	})
	if err != nil {
		return err
	}

	isOk = true
	return nil
}

func (s *TeamService) RemoveMember(requestorID, workspaceID, teamID, userID uuid.UUID) error {
	// Verify Admin
	isAuthorized := false
	teamMember, err := s.teamModel.GetMember(teamID, requestorID)
	if err == nil && teamMember.Role == models.TeamRoleAdmin {
		isAuthorized = true
	}
	if !isAuthorized {
		wsMember, err := s.workspaceModel.GetMember(workspaceID, requestorID)
		if err == nil && wsMember.Role == models.WorkspaceRoleAdmin {
			isAuthorized = true
		}
	}
	if !isAuthorized {
		return errors.New("unauthorized")
	}

	err = s.teamModel.RemoveMember(teamID, userID)
	if err != nil {
		return err
	}

	return nil
}
