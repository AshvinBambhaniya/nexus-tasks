package services

import (
	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WorkspaceService struct {
	workspaceModel *models.WorkspaceModel
	db             *goqu.Database
	logger         *zap.Logger
}

func NewWorkspaceService(db *goqu.Database, logger *zap.Logger, workspaceModel *models.WorkspaceModel) *WorkspaceService {
	return &WorkspaceService{
		workspaceModel: workspaceModel,
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
		return models.Workspace{}, err
	}

	// 2. Add Owner as Admin
	err = s.workspaceModel.AddMember(transaction, models.WorkspaceMember{
		WorkspaceID: createdWs.ID,
		UserID:      ownerID,
		Role:        models.WorkspaceRoleAdmin,
	})
	if err != nil {
		return models.Workspace{}, err
	}

	isOk = true
	return createdWs, nil
}
