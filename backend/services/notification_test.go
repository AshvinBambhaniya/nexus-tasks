package services

import (
	"errors"
	"testing"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupNotificationTest(_ *testing.T) (NotificationService, *mockNotificationRepository, *mockStorage) {
	mockNotifRepo := new(mockNotificationRepository)
	mockStor := new(mockStorage)
	logger := zap.NewNop()

	mockStor.On("Notifications").Return(mockNotifRepo)

	svc := NewNotificationService(mockStor, logger)

	return svc, mockNotifRepo, mockStor
}

func TestNotificationService_CreateNotification(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mn, _ := setupNotificationTest(t)
		notif := &models.Notification{ID: uuid.New()}

		mn.On("Create", notif).Return(nil)

		err := svc.CreateNotification(notif)
		assert.NoError(t, err)
		mn.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		svc, mn, _ := setupNotificationTest(t)
		notif := &models.Notification{ID: uuid.New()}

		mn.On("Create", notif).Return(errors.New("db error"))

		err := svc.CreateNotification(notif)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		mn.AssertExpectations(t)
	})
}

func TestNotificationService_GetInbox(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mn, _ := setupNotificationTest(t)
		userID := uuid.New()
		expected := []models.Notification{{ID: uuid.New()}}

		mn.On("GetInbox", userID).Return(expected, nil)

		res, err := svc.GetInbox(userID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, expected[0].ID, res[0].ID)
		mn.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		svc, mn, _ := setupNotificationTest(t)
		userID := uuid.New()

		mn.On("GetInbox", userID).Return([]models.Notification{}, errors.New("db error"))

		res, err := svc.GetInbox(userID)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "db error")
		mn.AssertExpectations(t)
	})
}

func TestNotificationService_MarkAsRead(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mn, _ := setupNotificationTest(t)
		id := uuid.New()
		userID := uuid.New()

		mn.On("MarkAsRead", id, userID).Return(nil)

		err := svc.MarkAsRead(id, userID)
		assert.NoError(t, err)
		mn.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		svc, mn, _ := setupNotificationTest(t)
		id := uuid.New()
		userID := uuid.New()

		mn.On("MarkAsRead", id, userID).Return(errors.New("db error"))

		err := svc.MarkAsRead(id, userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		mn.AssertExpectations(t)
	})
}

func TestNotificationService_MarkAsCleared(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mn, _ := setupNotificationTest(t)
		id := uuid.New()
		userID := uuid.New()

		mn.On("MarkAsCleared", id, userID).Return(nil)

		err := svc.MarkAsCleared(id, userID)
		assert.NoError(t, err)
		mn.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		svc, mn, _ := setupNotificationTest(t)
		id := uuid.New()
		userID := uuid.New()

		mn.On("MarkAsCleared", id, userID).Return(errors.New("db error"))

		err := svc.MarkAsCleared(id, userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		mn.AssertExpectations(t)
	})
}

func TestNotificationService_ClearAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mn, _ := setupNotificationTest(t)
		userID := uuid.New()

		mn.On("ClearAll", userID).Return(nil)

		err := svc.ClearAll(userID)
		assert.NoError(t, err)
		mn.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		svc, mn, _ := setupNotificationTest(t)
		userID := uuid.New()

		mn.On("ClearAll", userID).Return(errors.New("db error"))

		err := svc.ClearAll(userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		mn.AssertExpectations(t)
	})
}
