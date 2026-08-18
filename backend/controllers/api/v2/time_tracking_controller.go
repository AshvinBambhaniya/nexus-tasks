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

// GetActiveTimer returns the currently running timer for the authenticated user.
//
/*
 swagger:operation GET /timer/active timeTracking getActiveTimer

 # Get the active timer

 Returns the currently running timer for the authenticated user, or null if no timer is active.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 responses:

	200:
	  description: The active timer, or null.
	  schema:
      "$ref": "#/definitions/ResActiveTimer"
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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

// StartTimer starts a timer for the specified task.
//
/*
 swagger:operation POST /tasks/{taskId}/timer/start timeTracking startTimer

 # Start a task timer

 Begins a new timer session for the specified task. Only one timer may be active per user at a time.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: taskId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the task to time.

 responses:

	201:
	  description: Timer started.
	  schema:
      "$ref": "#/definitions/ResActiveTimer"
	400:
	  description: Invalid task ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error (e.g. timer already running).
*/
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

// StopTimer stops the active timer for the specified task.
//
/*
 swagger:operation POST /tasks/{taskId}/timer/stop timeTracking stopTimer

 # Stop a task timer

 Stops the currently running timer for the task and converts it to a time entry.
 An optional description and/or duration override may be provided.

 ---
 consumes:
 - application/json
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: taskId
     in: path
     required: true
     type: string
     format: uuid
   - name: body
     in: body
     description: Optional stop parameters. Body may be omitted entirely.
     schema:
        "$ref": "#/definitions/ReqStopTimer"

 responses:

	200:
	  description: Timer stopped and time entry created.
	  schema:
      "$ref": "#/definitions/ResTimeEntry"
	400:
	  description: Invalid task ID or malformed body.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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

// DiscardTimer discards the active timer without saving a time entry.
//
/*
 swagger:operation POST /tasks/{taskId}/timer/discard timeTracking discardTimer

 # Discard the active timer

 Cancels and discards the running timer without creating a time entry.
 Any elapsed time is permanently lost.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: taskId
     in: path
     required: true
     type: string
     format: uuid

 responses:

	200:
	  description: Timer discarded.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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

// LogManualTime manually logs time for a task.
//
/*
 swagger:operation POST /tasks/{taskId}/time-entries timeTracking logManualTime

 # Log manual time entry

 Creates a manual time entry for the specified task without using the start/stop timer flow.

 ---
 consumes:
 - application/json
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: taskId
     in: path
     required: true
     type: string
     format: uuid
   - name: body
     in: body
     required: true
     schema:
        "$ref": "#/definitions/ReqLogManualTime"

 responses:

	201:
	  description: Manual time entry created.
	  schema:
      "$ref": "#/definitions/ResTimeEntry"
	400:
	  description: Validation error or invalid task ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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

// ListTaskTimeEntries lists all time entries for a task.
//
/*
 swagger:operation GET /tasks/{taskId}/time-entries timeTracking listTaskTimeEntries

 # List time entries for a task

 Returns all time entries logged against the specified task, together with
 an aggregated total and the task's estimated duration.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: taskId
     in: path
     required: true
     type: string
     format: uuid

 responses:

	200:
	  description: Time entries with aggregate summary.
	  schema:
      "$ref": "#/definitions/ResTaskTimeEntries"
	400:
	  description: Invalid task ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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

// DeleteTimeEntry deletes a specific time entry.
//
/*
 swagger:operation DELETE /time-entries/{entryId} timeTracking deleteTimeEntry

 # Delete a time entry

 Permanently removes the specified time entry. Only the user who created the entry may delete it.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: entryId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the time entry.

 responses:

	200:
	  description: Time entry deleted.
	400:
	  description: Invalid entry ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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

// GetProjectAnalytics returns time tracking analytics for a project.
//
/*
 swagger:operation GET /workspaces/{workspaceId}/projects/{projectId}/time-analytics timeTracking getProjectAnalytics

 # Get project time analytics

 Returns aggregated time tracking statistics for all tasks within the specified project.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: workspaceId
     in: path
     required: true
     type: string
     format: uuid
   - name: projectId
     in: path
     required: true
     type: string
     format: uuid

 responses:

	200:
	  description: Project time analytics data.
	400:
	  description: Invalid project ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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

// ListProjectTimeEntries returns all time entries for a project with optional filters.
//
/*
 swagger:operation GET /workspaces/{workspaceId}/projects/{projectId}/time-entries timeTracking listProjectTimeEntries

 # List project time entries

 Returns all time entries logged against tasks in the project. Supports filtering by
 user, start date, and end date.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: workspaceId
     in: path
     required: true
     type: string
     format: uuid
   - name: projectId
     in: path
     required: true
     type: string
     format: uuid
   - name: user_id
     in: query
     type: string
     format: uuid
     description: Filter entries by a specific user UUID.
   - name: start_date
     in: query
     type: string
     format: date-time
     description: "Filter entries starting from this RFC3339 timestamp (inclusive)."
   - name: end_date
     in: query
     type: string
     format: date-time
     description: "Filter entries up to this RFC3339 timestamp (inclusive)."

 responses:

	200:
	  description: List of time entries with task context.
	  schema:
	    type: array
	    items:
	      "$ref": "#/definitions/ResTimeEntryWithTask"
	400:
	  description: Invalid project ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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
