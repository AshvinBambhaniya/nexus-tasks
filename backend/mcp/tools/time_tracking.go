package tools

import (
	"context"
	"fmt"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/structs"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/services"
	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// RegisterTimeTrackingTools registers all time tracking related tools.
//
//nolint:gocyclo
func RegisterTimeTrackingTools(s *server.MCPServer, service services.TimeTrackingService, logger *zap.Logger) {
	// start_task_timer
	startTimerTool := mcp.NewTool("start_task_timer",
		mcp.WithDescription("Starts a timer for a specific task."),
		mcp.WithString("task_id",
			mcp.Required(),
			mcp.Description("The UUID of the task to start the timer for."),
		),
	)
	s.AddTool(startTimerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
		if !ok {
			return mcp.NewToolResultError("Unauthorized"), nil
		}

		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments"), nil
		}

		taskIDStr, ok := args["task_id"].(string)
		if !ok {
			return mcp.NewToolResultError("task_id is required and must be a string"), nil
		}
		taskID, err := uuid.Parse(taskIDStr)
		if err != nil {
			return mcp.NewToolResultError("Invalid task_id format"), nil
		}

		timer, err := service.StartTimer(userID, taskID)
		if err != nil {
			logger.Error("MCP: failed to start timer", zap.Error(err))
			return mcp.NewToolResultError(fmt.Sprintf("Failed to start timer: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Timer started successfully for task %s at %s", timer.TaskTitle, timer.StartTime.Format("2006-01-02 15:04:05"))), nil
	})

	// stop_task_timer
	stopTimerTool := mcp.NewTool("stop_task_timer",
		mcp.WithDescription("Stops the active timer for a specific task."),
		mcp.WithString("task_id",
			mcp.Required(),
			mcp.Description("The UUID of the task."),
		),
		mcp.WithString("description",
			mcp.Description("Optional description of the work done."),
		),
	)
	s.AddTool(stopTimerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
		if !ok {
			return mcp.NewToolResultError("Unauthorized"), nil
		}

		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments"), nil
		}

		taskIDStr, ok := args["task_id"].(string)
		if !ok {
			return mcp.NewToolResultError("task_id is required and must be a string"), nil
		}
		taskID, err := uuid.Parse(taskIDStr)
		if err != nil {
			return mcp.NewToolResultError("Invalid task_id format"), nil
		}

		req := structs.ReqStopTimer{}
		if desc, ok := args["description"].(string); ok {
			req.Description = desc
		}

		entry, err := service.StopTimer(userID, taskID, req)
		if err != nil {
			logger.Error("MCP: failed to stop timer", zap.Error(err))
			return mcp.NewToolResultError(fmt.Sprintf("Failed to stop timer: %v", err)), nil
		}

		duration := 0
		if entry.DurationMinutes != nil {
			duration = *entry.DurationMinutes
		}
		return mcp.NewToolResultText(fmt.Sprintf("Timer stopped successfully. Logged %d minutes.", duration)), nil
	})

	// log_task_time
	logTimeTool := mcp.NewTool("log_task_time",
		mcp.WithDescription("Manually logs time for a specific task."),
		mcp.WithString("task_id",
			mcp.Required(),
			mcp.Description("The UUID of the task."),
		),
		mcp.WithNumber("duration_minutes",
			mcp.Required(),
			mcp.Description("The duration in minutes to log."),
		),
		mcp.WithString("description",
			mcp.Description("Optional description of the work done."),
		),
	)
	s.AddTool(logTimeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
		if !ok {
			return mcp.NewToolResultError("Unauthorized"), nil
		}

		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments"), nil
		}

		taskIDStr, ok := args["task_id"].(string)
		if !ok {
			return mcp.NewToolResultError("task_id is required and must be a string"), nil
		}
		taskID, err := uuid.Parse(taskIDStr)
		if err != nil {
			return mcp.NewToolResultError("Invalid task_id format"), nil
		}

		durationFloat, ok := args["duration_minutes"].(float64)
		if !ok {
			return mcp.NewToolResultError("duration_minutes is required and must be a number"), nil
		}

		req := structs.ReqLogManualTime{
			DurationMinutes: int(durationFloat),
		}
		if desc, ok := args["description"].(string); ok {
			req.Description = desc
		}

		entry, err := service.LogManualTime(userID, taskID, req)
		if err != nil {
			logger.Error("MCP: failed to log time manually", zap.Error(err))
			return mcp.NewToolResultError(fmt.Sprintf("Failed to log time: %v", err)), nil
		}

		duration := 0
		if entry.DurationMinutes != nil {
			duration = *entry.DurationMinutes
		}
		return mcp.NewToolResultText(fmt.Sprintf("Time logged successfully. Logged %d minutes.", duration)), nil
	})

	// get_project_time_analytics
	getProjectAnalyticsTool := mcp.NewTool("get_project_time_analytics",
		mcp.WithDescription("Gets time analytics for a specific project."),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("The UUID of the project."),
		),
	)
	s.AddTool(getProjectAnalyticsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
		if !ok {
			return mcp.NewToolResultError("Unauthorized"), nil
		}

		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments"), nil
		}

		projectIDStr, ok := args["project_id"].(string)
		if !ok {
			return mcp.NewToolResultError("project_id is required and must be a string"), nil
		}
		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			return mcp.NewToolResultError("Invalid project_id format"), nil
		}

		analytics, err := service.GetProjectAnalytics(userID, projectID)
		if err != nil {
			logger.Error("MCP: failed to get project analytics", zap.Error(err))
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get project analytics: %v", err)), nil
		}

		summary := fmt.Sprintf("Project Analytics:\n- Total Time Logged: %d minutes\n- Estimated Time: %d minutes\n\nUsers Time Distribution:\n",
			analytics.TotalLoggedMinutes, analytics.TotalEstimatedMinutes)

		for _, userStat := range analytics.ByMember {
			summary += fmt.Sprintf("  - %s: %d minutes\n", userStat.FullName, userStat.LoggedMinutes)
		}

		return mcp.NewToolResultText(summary), nil
	})
}
