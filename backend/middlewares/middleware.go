// Package middlewares provides HTTP middleware for the application.
package middlewares

import (
	"github.com/AshvinBambhaniya/nexus-tasks/v2/config"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/services"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

// Middleware provides various middleware handlers for the application.
type Middleware struct {
	workspaceModel models.WorkspaceRepository
	config         *config.AppConfig
	logger         *zap.Logger
	apiKeyService  services.APIKeyService
}

// NewMiddleware creates a new Middleware instance.
func NewMiddleware(goqu *goqu.Database, cfg *config.AppConfig, logger *zap.Logger, apiKeyService services.APIKeyService) Middleware {
	workspaceModel := models.InitWorkspaceModel(goqu)

	return Middleware{
		workspaceModel: workspaceModel,
		config:         cfg,
		logger:         logger,
		apiKeyService:  apiKeyService,
	}
}
