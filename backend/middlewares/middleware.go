package middlewares

import (
	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type Middleware struct {
	workspaceModel models.WorkspaceRepository
	config         *config.AppConfig
	logger         *zap.Logger
}

func NewMiddleware(goqu *goqu.Database, cfg *config.AppConfig, logger *zap.Logger) Middleware {
	workspaceModel := models.InitWorkspaceModel(goqu)

	return Middleware{
		workspaceModel: workspaceModel,
		config:         cfg,
		logger:         logger,
	}
}
