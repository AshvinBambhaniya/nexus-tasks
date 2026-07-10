package v2

import (
	"net/http"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/services"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// NotificationController handles HTTP requests for user inbox and notifications.
type NotificationController struct {
	notificationService services.NotificationService
	logger              *zap.Logger
}

// NewNotificationController initializes and returns a new NotificationController.
func NewNotificationController(ns services.NotificationService, logger *zap.Logger) *NotificationController {
	return &NotificationController{
		notificationService: ns,
		logger:              logger,
	}
}

// GetInbox handles GET /api/v2/inbox
func (nc *NotificationController) GetInbox(c *fiber.Ctx) error {
	userIDStr := utils.GetString(c.Locals(constants.ContextUID))
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		nc.logger.Error("invalid user id", zap.Error(err))
		return utils.JSONError(c, http.StatusUnauthorized, "invalid user context")
	}

	notifications, err := nc.notificationService.GetInbox(userID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "failed to get inbox")
	}

	return utils.JSONSuccess(c, http.StatusOK, notifications)
}

// MarkAsRead handles PATCH /api/v2/inbox/:notificationId/read
func (nc *NotificationController) MarkAsRead(c *fiber.Ctx) error {
	userIDStr := utils.GetString(c.Locals(constants.ContextUID))
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JSONError(c, http.StatusUnauthorized, "invalid user context")
	}

	notificationIDStr := c.Params("notificationId")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "invalid notification id")
	}

	err = nc.notificationService.MarkAsRead(notificationID, userID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "failed to mark as read")
	}

	return utils.JSONSuccess(c, http.StatusOK, nil)
}

// MarkAsCleared handles PATCH /api/v2/inbox/:notificationId/clear
func (nc *NotificationController) MarkAsCleared(c *fiber.Ctx) error {
	userIDStr := utils.GetString(c.Locals(constants.ContextUID))
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JSONError(c, http.StatusUnauthorized, "invalid user context")
	}

	notificationIDStr := c.Params("notificationId")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "invalid notification id")
	}

	err = nc.notificationService.MarkAsCleared(notificationID, userID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "failed to mark as cleared")
	}

	return utils.JSONSuccess(c, http.StatusOK, nil)
}

// ClearAll handles PATCH /api/v2/inbox/clear-all
func (nc *NotificationController) ClearAll(c *fiber.Ctx) error {
	userIDStr := utils.GetString(c.Locals(constants.ContextUID))
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JSONError(c, http.StatusUnauthorized, "invalid user context")
	}

	err = nc.notificationService.ClearAll(userID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "failed to clear all read notifications")
	}

	return utils.JSONSuccess(c, http.StatusOK, nil)
}
