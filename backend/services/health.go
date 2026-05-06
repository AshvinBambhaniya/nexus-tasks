package services

import (
	"context"

	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"go.uber.org/zap"
)

// HealthService defines the interface for health-related operations
type HealthService interface {
	CheckDatabaseHealth(ctx context.Context) error
}

type healthService struct {
	storage models.Storage
	logger  *zap.Logger
}

// NewHealthService creates a new health service instance
func NewHealthService(storage models.Storage, logger *zap.Logger) HealthService {
	return &healthService{
		storage: storage,
		logger:  logger,
	}
}

func (s *healthService) CheckDatabaseHealth(ctx context.Context) error {
	return s.storage.CheckHealth(ctx)
}
