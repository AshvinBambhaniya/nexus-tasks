package services

import (
	"context"

	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"go.uber.org/zap"
)

type HealthService interface {
	CheckDatabaseHealth(ctx context.Context) error
}

type healthService struct {
	storage models.Storage
	logger  *zap.Logger
}

func NewHealthService(storage models.Storage, logger *zap.Logger) HealthService {
	return &healthService{
		storage: storage,
		logger:  logger,
	}
}

func (s *healthService) CheckDatabaseHealth(ctx context.Context) error {
	return s.storage.CheckHealth(ctx)
}
