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

// GetInbox handles GET /api/v2/inbox.
//
/*
 swagger:operation GET /inbox notifications getInbox

 # Get notification inbox

 Returns all unread and uncleared notifications for the authenticated user, ordered by most recent first.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 responses:

	200:
	  description: List of notifications.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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

// MarkAsRead handles PATCH /api/v2/inbox/:notificationId/read.
//
/*
 swagger:operation PATCH /inbox/{notificationId}/read notifications markNotificationRead

 # Mark a notification as read

 Sets the read flag on the specified notification. The notification remains visible until cleared.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: notificationId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the notification.

 responses:

	200:
	  description: Notification marked as read.
	400:
	  description: Invalid notification ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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

// MarkAsCleared handles PATCH /api/v2/inbox/:notificationId/clear.
//
/*
 swagger:operation PATCH /inbox/{notificationId}/clear notifications clearNotification

 # Clear a notification

 Marks a notification as cleared, removing it from the user's active inbox view.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: notificationId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the notification.

 responses:

	200:
	  description: Notification cleared.
	400:
	  description: Invalid notification ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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

// ClearAll handles PATCH /api/v2/inbox/clear-all.
//
/*
 swagger:operation PATCH /inbox/clear-all notifications clearAllNotifications

 # Clear all notifications

 Marks all of the authenticated user's notifications as cleared in a single operation.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 responses:

	200:
	  description: All notifications cleared.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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
