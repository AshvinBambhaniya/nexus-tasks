package v2

import (
	"net/http"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/structs"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/services"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

// TimeTrackingController handles HTTP requests for time tracking features.
type TimeTrackingController struct {
	service services.TimeTrackingService
	logger  *zap.Logger
}

// NewTimeTrackingController creates a new instance of TimeTrackingController.
func NewTimeTrackingController(service services.TimeTrackingService, logger *zap.Logger) (*TimeTrackingController, error) {
	return &TimeTrackingController{
		service: service,
		logger:  logger,
	}, nil
}

// GetActiveTimer GET /timer/active
func (c *TimeTrackingController) GetActiveTimer(ctx *fiber.Ctx) error {
	uidStr := ctx.Locals(constants.ContextUID).(string)
	userID, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(ctx, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	timer, err := c.service.GetActiveTimer(userID)
	if err != nil {
		c.logger.Error("failed to get active timer", zap.Error(err), zap.String("userID", userID.String()))
		return utils.JSONError(ctx, http.StatusInternalServerError, "Failed to get active timer")
	}

	return utils.JSONSuccess(ctx, http.StatusOK, timer)
}

// StartTimer POST /tasks/:taskId/timer/start
func (c *TimeTrackingController) StartTimer(ctx *fiber.Ctx) error {
	uidStr := ctx.Locals(constants.ContextUID).(string)
	userID, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(ctx, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	taskID, err := uuid.Parse(ctx.Params(constants.ParamTaskID))
	if err != nil {
		return utils.JSONFail(ctx, http.StatusBadRequest, constants.ErrInvalidTaskID)
	}

	timer, err := c.service.StartTimer(userID, taskID)
	if err != nil {
		c.logger.Error("failed to start timer", zap.Error(err), zap.String("userID", userID.String()), zap.String("taskID", taskID.String()))
		return utils.JSONError(ctx, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(ctx, http.StatusCreated, timer)
}

// StopTimer POST /tasks/:taskId/timer/stop
func (c *TimeTrackingController) StopTimer(ctx *fiber.Ctx) error {
	uidStr := ctx.Locals(constants.ContextUID).(string)
	userID, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(ctx, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	taskID, err := uuid.Parse(ctx.Params(constants.ParamTaskID))
	if err != nil {
		return utils.JSONFail(ctx, http.StatusBadRequest, constants.ErrInvalidTaskID)
	}

	var req structs.ReqStopTimer
	if err := ctx.BodyParser(&req); err != nil && len(ctx.Body()) > 0 { // body is optional but if provided it must parse
		return utils.JSONFail(ctx, http.StatusBadRequest, "Invalid request payload")
	}

	entry, err := c.service.StopTimer(userID, taskID, req)
	if err != nil {
		c.logger.Error("failed to stop timer", zap.Error(err), zap.String("userID", userID.String()), zap.String("taskID", taskID.String()))
		return utils.JSONError(ctx, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(ctx, http.StatusOK, entry)
}

// DiscardTimer POST /tasks/:taskId/timer/discard (Note: prompt says /tasks/:taskId but DiscardTimer only needs userID)
func (c *TimeTrackingController) DiscardTimer(ctx *fiber.Ctx) error {
	uidStr := ctx.Locals(constants.ContextUID).(string)
	userID, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(ctx, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	err = c.service.DiscardTimer(userID)
	if err != nil {
		c.logger.Error("failed to discard timer", zap.Error(err), zap.String("userID", userID.String()))
		return utils.JSONError(ctx, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(ctx, http.StatusOK, fiber.Map{"message": "Timer discarded successfully"})
}

// LogManualTime POST /tasks/:taskId/time-entries
func (c *TimeTrackingController) LogManualTime(ctx *fiber.Ctx) error {
	uidStr := ctx.Locals(constants.ContextUID).(string)
	userID, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(ctx, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	taskID, err := uuid.Parse(ctx.Params(constants.ParamTaskID))
	if err != nil {
		return utils.JSONFail(ctx, http.StatusBadRequest, constants.ErrInvalidTaskID)
	}

	var req structs.ReqLogManualTime
	if err := ctx.BodyParser(&req); err != nil {
		return utils.JSONFail(ctx, http.StatusBadRequest, "Invalid request payload")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.JSONFail(ctx, http.StatusBadRequest, err.Error())
	}

	entry, err := c.service.LogManualTime(userID, taskID, req)
	if err != nil {
		c.logger.Error("failed to log manual time", zap.Error(err), zap.String("userID", userID.String()), zap.String("taskID", taskID.String()))
		return utils.JSONError(ctx, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(ctx, http.StatusCreated, entry)
}

// ListTaskTimeEntries GET /tasks/:taskId/time-entries
func (c *TimeTrackingController) ListTaskTimeEntries(ctx *fiber.Ctx) error {
	uidStr := ctx.Locals(constants.ContextUID).(string)
	userID, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(ctx, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	taskID, err := uuid.Parse(ctx.Params(constants.ParamTaskID))
	if err != nil {
		return utils.JSONFail(ctx, http.StatusBadRequest, constants.ErrInvalidTaskID)
	}

	entries, totalLogged, estimated, err := c.service.ListTaskTimeEntries(userID, taskID)
	if err != nil {
		c.logger.Error("failed to list time entries", zap.Error(err), zap.String("userID", userID.String()), zap.String("taskID", taskID.String()))
		return utils.JSONError(ctx, http.StatusInternalServerError, err.Error())
	}

	resEntries := make([]structs.ResTimeEntry, len(entries))
	for i, e := range entries {
		resEntries[i] = structs.ResTimeEntry{
			ID:              e.ID,
			TaskID:          e.TaskID,
			UserID:          e.UserID,
			UserFullName:    e.UserFullName,
			Description:     e.Description,
			StartTime:       e.StartTime,
			EndTime:         e.EndTime,
			DurationMinutes: e.DurationMinutes,
			IsManual:        e.IsManual,
			CreatedAt:       e.CreatedAt,
		}
	}

	res := structs.ResTaskTimeEntries{
		Entries:            resEntries,
		TotalLoggedMinutes: totalLogged,
		EstimatedMinutes:   estimated,
	}

	return utils.JSONSuccess(ctx, http.StatusOK, res)
}

// DeleteTimeEntry DELETE /time-entries/:entryId
func (c *TimeTrackingController) DeleteTimeEntry(ctx *fiber.Ctx) error {
	uidStr := ctx.Locals(constants.ContextUID).(string)
	userID, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(ctx, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	entryID, err := uuid.Parse(ctx.Params(constants.ParamEntryID))
	if err != nil {
		return utils.JSONFail(ctx, http.StatusBadRequest, constants.ErrInvalidTaskID)
	}

	err = c.service.DeleteTimeEntry(userID, entryID)
	if err != nil {
		c.logger.Error("failed to delete time entry", zap.Error(err), zap.String("userID", userID.String()), zap.String("entryID", entryID.String()))
		return utils.JSONError(ctx, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(ctx, http.StatusOK, fiber.Map{"message": "Time entry deleted successfully"})
}

// GetProjectAnalytics GET /projects/:projectId/time-analytics (mounted under workspaces/:workspaceId/projects/:projectId/time-analytics)
func (c *TimeTrackingController) GetProjectAnalytics(ctx *fiber.Ctx) error {
	uidStr := ctx.Locals(constants.ContextUID).(string)
	userID, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(ctx, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	projectID, err := uuid.Parse(ctx.Params(constants.ParamProjectID))
	if err != nil {
		return utils.JSONFail(ctx, http.StatusBadRequest, constants.ErrInvalidTaskID)
	}

	analytics, err := c.service.GetProjectAnalytics(userID, projectID)
	if err != nil {
		c.logger.Error("failed to get project analytics", zap.Error(err), zap.String("userID", userID.String()), zap.String("projectID", projectID.String()))
		return utils.JSONError(ctx, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(ctx, http.StatusOK, analytics)
}

// ListProjectTimeEntries GET /workspaces/:workspaceId/projects/:projectId/time-entries
func (c *TimeTrackingController) ListProjectTimeEntries(ctx *fiber.Ctx) error {
	uidStr := ctx.Locals(constants.ContextUID).(string)
	userID, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(ctx, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	projectID, err := uuid.Parse(ctx.Params(constants.ParamProjectID))
	if err != nil {
		return utils.JSONFail(ctx, http.StatusBadRequest, "Invalid project ID")
	}

	var targetUserID *uuid.UUID
	if targetIDStr := ctx.Query("user_id"); targetIDStr != "" {
		if parsed, err := uuid.Parse(targetIDStr); err == nil {
			targetUserID = &parsed
		}
	}

	var startDate, endDate *time.Time
	if startStr := ctx.Query("start_date"); startStr != "" {
		// Try parsing RFC3339, or fallback to simple date if needed. The prompt mentions RFC3339.
		if parsed, err := time.Parse(time.RFC3339, startStr); err == nil {
			startDate = &parsed
		}
	}
	if endStr := ctx.Query("end_date"); endStr != "" {
		if parsed, err := time.Parse(time.RFC3339, endStr); err == nil {
			endDate = &parsed
		}
	}

	entries, err := c.service.ListProjectTimeEntries(userID, projectID, targetUserID, startDate, endDate)
	if err != nil {
		c.logger.Error("failed to list project time entries", zap.Error(err), zap.String("userID", userID.String()), zap.String("projectID", projectID.String()))
		return utils.JSONError(ctx, http.StatusInternalServerError, err.Error())
	}

	resEntries := make([]structs.ResTimeEntryWithTask, len(entries))
	for i, e := range entries {
		resEntries[i] = structs.ResTimeEntryWithTask{
			ResTimeEntry: structs.ResTimeEntry{
				ID:              e.ID,
				TaskID:          e.TaskID,
				UserID:          e.UserID,
				UserFullName:    e.UserFullName,
				Description:     e.Description,
				StartTime:       e.StartTime,
				EndTime:         e.EndTime,
				DurationMinutes: e.DurationMinutes,
				IsManual:        e.IsManual,
				CreatedAt:       e.CreatedAt,
			},
			TaskTitle:  e.TaskTitle,
			TaskNumber: e.TaskNumber,
		}
	}

	return utils.JSONSuccess(ctx, http.StatusOK, resEntries)
}
