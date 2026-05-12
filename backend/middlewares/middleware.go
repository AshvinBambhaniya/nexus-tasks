// Package middlewares provides HTTP middleware for the application.
package middlewares

import (
	"github.com/AshvinBambhaniya/nexus-tasks/v2/config"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

// Middleware provides various middleware handlers for the application.
type Middleware struct {
	workspaceModel models.WorkspaceRepository
	config         *config.AppConfig
	logger         *zap.Logger
}

// NewMiddleware creates a new Middleware instance.
func NewMiddleware(goqu *goqu.Database, cfg *config.AppConfig, logger *zap.Logger) Middleware {
	workspaceModel := models.InitWorkspaceModel(goqu)

	return Middleware{
		workspaceModel: workspaceModel,
		config:         cfg,
		logger:         logger,
	}
}
