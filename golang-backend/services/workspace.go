package services

import (
	"database/sql"
	"errors"

	"github.com/AshvinBambhaniya/nexus-tasks/cli/workers"
	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/watermill"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WorkspaceService struct {
	workspaceModel *models.WorkspaceModel
	userModel      *models.UserModel
	publisher      *watermill.WatermillPublisher
	db             *goqu.Database
	logger         *zap.Logger
}

func NewWorkspaceService(db *goqu.Database, logger *zap.Logger, workspaceModel *models.WorkspaceModel, userModel *models.UserModel, publisher *watermill.WatermillPublisher) *WorkspaceService {
	return &WorkspaceService{
		workspaceModel: workspaceModel,
		userModel:      userModel,
		publisher:      publisher,
		db:             db,
		logger:         logger,
	}
}

// CreateWorkspace creates a TEAM workspace and adds the creator as ADMIN
func (s *WorkspaceService) CreateWorkspace(ownerID uuid.UUID, req structs.ReqCreateWorkspace) (models.Workspace, error) {
	ws := models.Workspace{
		Name:    req.Name,
		Type:    models.WorkspaceTypeTeam,
		OwnerID: ownerID,
	}

	isOk := false
	transaction, err := s.db.Begin()
	if err != nil {
		return models.Workspace{}, err
	}

	defer func() {
		if isOk {
			err := transaction.Commit()
			if err != nil {
				s.logger.Error("error during commit in create workspace", zap.Error(err))
			}
		} else {
			err := transaction.Rollback()
			if err != nil {
				s.logger.Error("error during rollback in create workspace", zap.Error(err))
			}
		}
	}()

	// 1. Create Workspace
	createdWs, err := s.workspaceModel.CreateWorkspace(transaction, ws)
	if err != nil {
		return createdWs, err
	}

	// 2. Add Owner as Admin
	err = s.workspaceModel.AddMemberTx(transaction, models.WorkspaceMember{
		WorkspaceID: createdWs.ID,
		UserID:      ownerID,
		Role:        models.WorkspaceRoleAdmin,
	})
	if err != nil {
		return createdWs, err
	}

	isOk = true
	return createdWs, nil
}

func (s *WorkspaceService) InviteMember(requestorID, workspaceID uuid.UUID, email string) error {
	// 1. Validate Admin Access (Internal check required here as role check is specific)
	member, err := s.ValidateAccess(requestorID, workspaceID)
	if err != nil {
		return err
	}
	if member.Role != models.WorkspaceRoleAdmin {
		s.logger.Error("unauthorized invite attempt", zap.Any("requestorID", requestorID), zap.Any("workspaceID", workspaceID))
		return errors.New("unauthorized: only admins can invite members")
	}

	// 2. Find User by Email
	user, err := s.userModel.GetByEmail(email)
	if err != nil {
		s.logger.Error("user not found for email", zap.String("email", email), zap.Error(err))
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	// 3. Check if already member
	_, err = s.workspaceModel.GetMember(workspaceID, user.ID)
	if err == nil {
		s.logger.Error("user is already a member of the workspace", zap.Any("workspaceID", workspaceID), zap.Any("userID", user.ID))
		return errors.New("user is already a member")
	}

	// 4. Add Member
	err = s.workspaceModel.AddMember(models.WorkspaceMember{
		WorkspaceID: workspaceID,
		UserID:      user.ID,
		Role:        models.WorkspaceRoleMember,
	})

	if err != nil {
		s.logger.Error("failed to add member to workspace", zap.Any("workspaceID", workspaceID), zap.Any("userID", user.ID), zap.Error(err))
		return err
	}

	// 5. Send Notification (Async)
	workspace, err := s.workspaceModel.GetByID(workspaceID)
	if err == nil && s.publisher != nil {
		err = s.publisher.Publish(constants.TopicWorkspaceInvites, workers.WorkspaceInvitationMail{
			Email:         email,
			WorkspaceName: workspace.Name,
		})
		if err != nil {
			s.logger.Error("failed to publish workspace invitation", zap.Error(err))
			// We don't return error here because the DB operation was successful
		}
	}

	return nil
}

func (s *WorkspaceService) RemoveMember(requestorID, workspaceID, userID uuid.UUID) error {
	// 1. Validate Admin Access
	member, err := s.ValidateAccess(requestorID, workspaceID)
	if err != nil {
		return err
	}
	if member.Role != models.WorkspaceRoleAdmin {
		return errors.New("unauthorized: only admins can remove members")
	}

	// 2. Remove
	return s.workspaceModel.RemoveMember(nil, workspaceID, userID)
}

// ValidateAccess checks if a user is a member of the workspace (Internal Helper)
func (s *WorkspaceService) ValidateAccess(userID, workspaceID uuid.UUID) (models.WorkspaceMember, error) {
	member, err := s.workspaceModel.GetMember(workspaceID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.WorkspaceMember{}, errors.New("unauthorized: not a member of this workspace")
		}
		return models.WorkspaceMember{}, err
	}
	return member, nil
}
