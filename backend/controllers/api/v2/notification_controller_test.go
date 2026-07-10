package v2

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupNotificationControllerTest() (*fiber.App, *mockNotificationService, *NotificationController) {
	app := fiber.New()
	mockSvc := new(mockNotificationService)
	logger := zap.NewNop()
	ctrl := NewNotificationController(mockSvc, logger)

	app.Get("/api/v2/notifications", func(c *fiber.Ctx) error {
		c.Locals("userId", c.Get("X-User-Id"))
		return ctrl.GetInbox(c)
	})
	app.Put("/api/v2/notifications/:notificationId/read", func(c *fiber.Ctx) error {
		c.Locals("userId", c.Get("X-User-Id"))
		return ctrl.MarkAsRead(c)
	})
	app.Delete("/api/v2/notifications/:notificationId/clear", func(c *fiber.Ctx) error {
		c.Locals("userId", c.Get("X-User-Id"))
		return ctrl.MarkAsCleared(c)
	})
	app.Delete("/api/v2/notifications", func(c *fiber.Ctx) error {
		c.Locals("userId", c.Get("X-User-Id"))
		return ctrl.ClearAll(c)
	})

	return app, mockSvc, ctrl
}

func TestNotificationController_GetInbox(t *testing.T) {
	app, ms, _ := setupNotificationControllerTest()

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		ms.On("GetInbox", userID).Return([]models.Notification{{ID: uuid.New(), Title: "Test"}}, nil).Once()

		req := httptest.NewRequest("GET", "/api/v2/notifications", nil)
		req.Header.Set("X-User-Id", userID.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("invalid user context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v2/notifications", nil)
		req.Header.Set("X-User-Id", "invalid-id")

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 401, resp.StatusCode)
	})

	t.Run("service failure", func(t *testing.T) {
		userID := uuid.New()
		ms.On("GetInbox", userID).Return([]models.Notification{}, errors.New("db error")).Once()

		req := httptest.NewRequest("GET", "/api/v2/notifications", nil)
		req.Header.Set("X-User-Id", userID.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 500, resp.StatusCode)
	})
}

func TestNotificationController_MarkAsRead(t *testing.T) {
	app, ms, _ := setupNotificationControllerTest()

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		notifID := uuid.New()
		ms.On("MarkAsRead", notifID, userID).Return(nil).Once()

		req := httptest.NewRequest("PUT", "/api/v2/notifications/"+notifID.String()+"/read", nil)
		req.Header.Set("X-User-Id", userID.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("invalid user context", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/v2/notifications/"+uuid.New().String()+"/read", nil)
		req.Header.Set("X-User-Id", "bad")

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 401, resp.StatusCode)
	})

	t.Run("invalid notification id", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/v2/notifications/invalid-id/read", nil)
		req.Header.Set("X-User-Id", uuid.New().String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("service failure", func(t *testing.T) {
		userID := uuid.New()
		notifID := uuid.New()
		ms.On("MarkAsRead", notifID, userID).Return(errors.New("db error")).Once()

		req := httptest.NewRequest("PUT", "/api/v2/notifications/"+notifID.String()+"/read", nil)
		req.Header.Set("X-User-Id", userID.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 500, resp.StatusCode)
	})
}

func TestNotificationController_MarkAsCleared(t *testing.T) {
	app, ms, _ := setupNotificationControllerTest()

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		notifID := uuid.New()
		ms.On("MarkAsCleared", notifID, userID).Return(nil).Once()

		req := httptest.NewRequest("DELETE", "/api/v2/notifications/"+notifID.String()+"/clear", nil)
		req.Header.Set("X-User-Id", userID.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("invalid user context", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v2/notifications/"+uuid.New().String()+"/clear", nil)
		req.Header.Set("X-User-Id", "bad")

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 401, resp.StatusCode)
	})

	t.Run("invalid notification id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v2/notifications/invalid-id/clear", nil)
		req.Header.Set("X-User-Id", uuid.New().String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("service failure", func(t *testing.T) {
		userID := uuid.New()
		notifID := uuid.New()
		ms.On("MarkAsCleared", notifID, userID).Return(errors.New("db err")).Once()

		req := httptest.NewRequest("DELETE", "/api/v2/notifications/"+notifID.String()+"/clear", nil)
		req.Header.Set("X-User-Id", userID.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 500, resp.StatusCode)
	})
}

func TestNotificationController_ClearAll(t *testing.T) {
	app, ms, _ := setupNotificationControllerTest()

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		ms.On("ClearAll", userID).Return(nil).Once()

		req := httptest.NewRequest("DELETE", "/api/v2/notifications", nil)
		req.Header.Set("X-User-Id", userID.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("invalid user context", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v2/notifications", nil)
		req.Header.Set("X-User-Id", "bad")

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 401, resp.StatusCode)
	})

	t.Run("service failure", func(t *testing.T) {
		userID := uuid.New()
		ms.On("ClearAll", userID).Return(errors.New("db error")).Once()

		req := httptest.NewRequest("DELETE", "/api/v2/notifications", nil)
		req.Header.Set("X-User-Id", userID.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 500, resp.StatusCode)
	})
}
