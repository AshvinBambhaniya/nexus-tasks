package middlewares

import (
	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type Middleware struct {
	workspaceModel *models.WorkspaceModel
	config         config.AppConfig
	logger         *zap.Logger
}

func NewMiddleware(goqu *goqu.Database, cfg config.AppConfig, logger *zap.Logger) (Middleware, error) {
	workspaceModel, err := models.InitWorkspaceModel(goqu)
	if err != nil {
		return Middleware{}, err
	}

	return Middleware{
		workspaceModel: &workspaceModel,
		config:         cfg,
		logger:         logger,
	}, nil
}
