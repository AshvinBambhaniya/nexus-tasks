package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AshvinBambhaniya/nexus-tasks/cli/workers"
	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/watermill"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// WorkspaceService defines the interface for workspace-related business logic
type WorkspaceService interface {
	CreateWorkspace(ownerID uuid.UUID, req structs.ReqCreateWorkspace) (models.Workspace, error)
	ListWorkspacesByUserID(userID uuid.UUID) ([]models.Workspace, error)
	ListMembersByWorkspaceID(workspaceID uuid.UUID) ([]models.WorkspaceMemberWithUser, error)
	InviteMember(requestorID, workspaceID uuid.UUID, email string) error
	RemoveMember(requestorID, workspaceID, userID uuid.UUID) error
	ValidateAccess(userID, workspaceID uuid.UUID) (models.WorkspaceMember, error)
}

type workspaceService struct {
	storage   models.Storage
	publisher watermill.IPublisher
	logger    *zap.Logger
}

// NewWorkspaceService creates a new workspace service instance
func NewWorkspaceService(storage models.Storage, logger *zap.Logger, publisher watermill.IPublisher) WorkspaceService {
	return &workspaceService{
		storage:   storage,
		publisher: publisher,
		logger:    logger,
	}
}

// CreateWorkspace creates a TEAM workspace and adds the creator as ADMIN
func (s *workspaceService) CreateWorkspace(ownerID uuid.UUID, req structs.ReqCreateWorkspace) (models.Workspace, error) {
	var createdWs models.Workspace

	err := s.storage.Atomic(context.Background(), func(txStorage models.Storage) error {
		ws := models.Workspace{
			Name:    req.Name,
			Type:    models.WorkspaceTypeTeam,
			OwnerID: ownerID,
		}

		// 1. Create Workspace
		var err error
		createdWs, err = txStorage.Workspaces().CreateWorkspace(ws)
		if err != nil {
			return err
		}

		// 2. Add Owner as Admin
		err = txStorage.Workspaces().AddMember(models.WorkspaceMember{
			WorkspaceID: createdWs.ID,
			UserID:      ownerID,
			Role:        models.WorkspaceRoleAdmin,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("failed to create workspace", zap.Error(err))
		return models.Workspace{}, err
	}

	return createdWs, nil
}

func (s *workspaceService) ListWorkspacesByUserID(userID uuid.UUID) ([]models.Workspace, error) {
	return s.storage.Workspaces().ListWorkspacesByUserID(userID)
}

func (s *workspaceService) ListMembersByWorkspaceID(workspaceID uuid.UUID) ([]models.WorkspaceMemberWithUser, error) {
	return s.storage.Workspaces().ListMembersByWorkspaceID(workspaceID)
}

func (s *workspaceService) InviteMember(requestorID, workspaceID uuid.UUID, email string) error {
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
	user, err := s.storage.Users().GetByEmail(email)
	if err != nil {
		s.logger.Error("user not found for email", zap.String("email", email), zap.Error(err))
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	// 3. Check if already member
	_, err = s.storage.Workspaces().GetMember(workspaceID, user.ID)
	if err == nil {
		s.logger.Error("user is already a member of the workspace", zap.Any("workspaceID", workspaceID), zap.Any("userID", user.ID))
		return errors.New("user is already a member")
	}

	// 4. Add Member
	err = s.storage.Workspaces().AddMember(models.WorkspaceMember{
		WorkspaceID: workspaceID,
		UserID:      user.ID,
		Role:        models.WorkspaceRoleMember,
	})

	if err != nil {
		s.logger.Error("failed to add member to workspace", zap.Any("workspaceID", workspaceID), zap.Any("userID", user.ID), zap.Error(err))
		return err
	}

	// 5. Send Notification (Async)
	workspace, err := s.storage.Workspaces().GetByID(workspaceID)
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

func (s *workspaceService) RemoveMember(requestorID, workspaceID, userID uuid.UUID) error {
	// 1. Validate Admin Access
	member, err := s.ValidateAccess(requestorID, workspaceID)
	if err != nil {
		return err
	}
	if member.Role != models.WorkspaceRoleAdmin {
		return errors.New("unauthorized: only admins can remove members")
	}

	// 2. Remove
	return s.storage.Workspaces().RemoveMember(workspaceID, userID)
}

// ValidateAccess checks if a user is a member of the workspace (Internal Helper)
func (s *workspaceService) ValidateAccess(userID, workspaceID uuid.UUID) (models.WorkspaceMember, error) {
	member, err := s.storage.Workspaces().GetMember(workspaceID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.WorkspaceMember{}, errors.New("unauthorized: not a member of this workspace")
		}
		return models.WorkspaceMember{}, err
	}
	return member, nil
}
