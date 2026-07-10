package services

import (
	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// NotificationService provides methods for managing notifications.
type NotificationService interface {
	CreateNotification(n *models.Notification) error
	GetInbox(userID uuid.UUID) ([]models.Notification, error)
	MarkAsRead(id uuid.UUID, userID uuid.UUID) error
	MarkAsCleared(id uuid.UUID, userID uuid.UUID) error
	ClearAll(userID uuid.UUID) error
}

type notificationService struct {
	storage models.Storage
	logger  *zap.Logger
}

// NewNotificationService returns a new instance of NotificationService.
func NewNotificationService(storage models.Storage, logger *zap.Logger) NotificationService {
	return &notificationService{
		storage: storage,
		logger:  logger,
	}
}

func (s *notificationService) CreateNotification(n *models.Notification) error {
	err := s.storage.Notifications().Create(n)
	if err != nil {
		s.logger.Error("failed to create notification", zap.Error(err))
		return err
	}
	return nil
}

func (s *notificationService) GetInbox(userID uuid.UUID) ([]models.Notification, error) {
	notifications, err := s.storage.Notifications().GetInbox(userID)
	if err != nil {
		s.logger.Error("failed to get inbox notifications", zap.Error(err))
		return nil, err
	}
	return notifications, nil
}

func (s *notificationService) MarkAsRead(id uuid.UUID, userID uuid.UUID) error {
	err := s.storage.Notifications().MarkAsRead(id, userID)
	if err != nil {
		s.logger.Error("failed to mark notification as read", zap.Error(err))
		return err
	}
	return nil
}

func (s *notificationService) MarkAsCleared(id uuid.UUID, userID uuid.UUID) error {
	err := s.storage.Notifications().MarkAsCleared(id, userID)
	if err != nil {
		s.logger.Error("failed to mark notification as cleared", zap.Error(err))
		return err
	}
	return nil
}

func (s *notificationService) ClearAll(userID uuid.UUID) error {
	err := s.storage.Notifications().ClearAll(userID)
	if err != nil {
		s.logger.Error("failed to clear all read notifications", zap.Error(err))
		return err
	}
	return nil
}
