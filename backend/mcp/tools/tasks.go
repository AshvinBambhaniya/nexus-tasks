package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

type contextKey string

// UserIDKey is the context key for the user ID
const UserIDKey contextKey = "userID"

// RegisterTaskTools registers all task-related tools.
func RegisterTaskTools(s *server.MCPServer, storage models.Storage, logger *zap.Logger) {
	// list_project_tasks
	listProjectTasksTool := mcp.NewTool("list_project_tasks",
		mcp.WithDescription("Lists tasks in a specific project."),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("The UUID of the project."),
		),
		mcp.WithString("status",
			mcp.Description("Optional. Filter by task status: 'TODO', 'IN_PROGRESS', 'DONE', 'BACKLOG'"),
		),
	)
	s.AddTool(listProjectTasksTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleListProjectTasks(ctx, req, storage, logger)
	})

	// list_my_tasks
	listMyTasksTool := mcp.NewTool("list_my_tasks",
		mcp.WithDescription("Lists tasks assigned to the currently authenticated user."),
		mcp.WithString("status",
			mcp.Description("Optional. Filter by task status: 'TODO', 'IN_PROGRESS', 'DONE', 'BACKLOG'"),
		),
	)
	s.AddTool(listMyTasksTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleListMyTasks(ctx, req, storage, logger)
	})

	// get_task_details
	getTaskDetailsTool := mcp.NewTool("get_task_details",
		mcp.WithDescription("Gets the full details of a specific task by its UUID."),
		mcp.WithString("task_id",
			mcp.Required(),
			mcp.Description("The UUID of the task to retrieve."),
		),
	)
	s.AddTool(getTaskDetailsTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetTaskDetails(ctx, req, storage, logger)
	})

	// update_task
	updateTaskTool := mcp.NewTool("update_task",
		mcp.WithDescription("Updates an existing task. All fields are optional except task_id."),
		mcp.WithString("task_id",
			mcp.Required(),
			mcp.Description("The UUID of the task to update."),
		),
		mcp.WithString("status",
			mcp.Description("The new status: 'TODO', 'IN_PROGRESS', 'DONE', 'BACKLOG'"),
		),
		mcp.WithString("priority",
			mcp.Description("The new priority: 'P0', 'P1', 'P2', 'P3'"),
		),
		mcp.WithString("assignee_id",
			mcp.Description("The UUID of the user to assign the task to."),
		),
		mcp.WithString("due_date",
			mcp.Description("The due date in RFC3339 format (e.g., '2026-12-31T23:59:59Z')."),
		),
		mcp.WithString("title",
			mcp.Description("The new title for the task."),
		),
		mcp.WithString("description",
			mcp.Description("The new description for the task."),
		),
	)
	s.AddTool(updateTaskTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleUpdateTask(ctx, req, storage, logger)
	})

	// create_task
	createTaskTool := mcp.NewTool("create_task",
		mcp.WithDescription("Creates a new task in a specific project."),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("The UUID of the project to create the task in."),
		),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("The title of the task."),
		),
		mcp.WithString("description",
			mcp.Description("Optional. The detailed description of the task."),
		),
		mcp.WithString("priority",
			mcp.Description("Optional. The priority: 'P0', 'P1', 'P2', 'P3'"),
		),
		mcp.WithString("assignee_id",
			mcp.Description("Optional. The UUID of the user to assign the task to."),
		),
		mcp.WithString("due_date",
			mcp.Description("Optional. The due date in RFC3339 format (e.g., '2026-12-31T23:59:59Z')."),
		),
	)
	s.AddTool(createTaskTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleCreateTask(ctx, req, storage, logger)
	})
}

func getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return uuid.Nil, fmt.Errorf("unauthorized: no user ID in context")
	}
	userID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("internal error: invalid user ID format in context")
	}
	return userID, nil
}

func handleListMyTasks(ctx context.Context, req mcp.CallToolRequest, storage models.Storage, logger *zap.Logger) (*mcp.CallToolResult, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := req.Params.Arguments.(map[string]any)
	var statusFilter *models.TaskStatus
	if st, ok := args["status"].(string); ok && st != "" {
		s := models.TaskStatus(st)
		statusFilter = &s
	}

	tasks, err := storage.Tasks().ListByAssigneeID(userID)
	if err != nil {
		logger.Error("Failed to list tasks", zap.Error(err))
		return mcp.NewToolResultError("Failed to fetch tasks from database"), nil
	}

	// Filter by status manually since ListByAssigneeID doesn't take status filter in this version
	var filtered []models.TaskWithAssignee
	for _, t := range tasks {
		if statusFilter == nil || t.Status == *statusFilter {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) == 0 {
		return mcp.NewToolResultText("No tasks found."), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d tasks:\n\n", len(filtered))
	for _, t := range filtered {
		fmt.Fprintf(&sb, "- [%s] %s (ID: %s)\n", t.Status, t.Title, t.ID)
		if t.Description != "" {
			desc := t.Description
			if len(desc) > 100 {
				desc = desc[:97] + "..."
			}
			fmt.Fprintf(&sb, "  Description: %s\n", desc)
		}
		fmt.Fprintf(&sb, "  Priority: %s\n", t.Priority)
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func handleGetTaskDetails(_ context.Context, req mcp.CallToolRequest, storage models.Storage, _ *zap.Logger) (*mcp.CallToolResult, error) {
	// Not strictly verifying access for demo, but normally check if user has access to task's workspace
	args := req.Params.Arguments.(map[string]any)
	taskIDStr, ok := args["task_id"].(string)
	if !ok || taskIDStr == "" {
		return mcp.NewToolResultError("task_id is required"), nil
	}

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return mcp.NewToolResultError("invalid task_id format"), nil
	}

	task, err := storage.Tasks().GetByID(taskID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Task not found or error: %v", err)), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# %s\n\n", task.Title)
	fmt.Fprintf(&sb, "**ID:** %s\n", task.ID)
	fmt.Fprintf(&sb, "**Status:** %s\n", task.Status)
	fmt.Fprintf(&sb, "**Priority:** %s\n", task.Priority)
	if task.AssigneeFullName != nil {
		fmt.Fprintf(&sb, "**Assignee:** %s\n", *task.AssigneeFullName)
	}
	if task.DueDate != nil {
		fmt.Fprintf(&sb, "**Due Date:** %s\n", task.DueDate.Format(time.RFC3339))
	}
	fmt.Fprintf(&sb, "\n## Description\n%s\n", task.Description)

	return mcp.NewToolResultText(sb.String()), nil
}

//nolint:gocyclo
func handleUpdateTask(_ context.Context, req mcp.CallToolRequest, storage models.Storage, _ *zap.Logger) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments.(map[string]any)

	taskIDStr, ok := args["task_id"].(string)
	if !ok || taskIDStr == "" {
		return mcp.NewToolResultError("task_id is required"), nil
	}

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return mcp.NewToolResultError("invalid task_id format"), nil
	}

	task, err := storage.Tasks().GetByID(taskID)
	if err != nil {
		return mcp.NewToolResultError("task not found"), nil
	}

	if statusStr, ok := args["status"].(string); ok && statusStr != "" {
		task.Status = models.TaskStatus(statusStr)
	}

	if priorityStr, ok := args["priority"].(string); ok && priorityStr != "" {
		task.Priority = models.TaskPriority(priorityStr)
	}

	if titleStr, ok := args["title"].(string); ok && titleStr != "" {
		task.Title = titleStr
	}

	if descStr, ok := args["description"].(string); ok && descStr != "" {
		task.Description = descStr
	}

	if assigneeIDStr, ok := args["assignee_id"].(string); ok && assigneeIDStr != "" {
		parsed, err := uuid.Parse(assigneeIDStr)
		if err == nil {
			task.AssigneeID = &parsed
		} else {
			return mcp.NewToolResultError("invalid assignee_id format"), nil
		}
	}

	if dueDateStr, ok := args["due_date"].(string); ok && dueDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, dueDateStr)
		if err == nil {
			task.DueDate = &parsed
		} else {
			return mcp.NewToolResultError("invalid due_date format, expected RFC3339"), nil
		}
	}

	_, err = storage.Tasks().Update(task.Task)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to update task: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully updated task %s", taskID)), nil
}

func handleCreateTask(ctx context.Context, req mcp.CallToolRequest, storage models.Storage, _ *zap.Logger) (*mcp.CallToolResult, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := req.Params.Arguments.(map[string]any)

	projectIDStr, ok := args["project_id"].(string)
	if !ok || projectIDStr == "" {
		return mcp.NewToolResultError("project_id is required"), nil
	}
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return mcp.NewToolResultError("invalid project_id format"), nil
	}

	title, ok := args["title"].(string)
	if !ok || title == "" {
		return mcp.NewToolResultError("title is required"), nil
	}

	description, _ := args["description"].(string)

	priorityStr, ok := args["priority"].(string)
	if !ok || priorityStr == "" {
		priorityStr = "P2" // Default priority
	}

	// Get next task number
	nextNumber, err := storage.Tasks().GetNextTaskNumber(projectID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to generate task number: %v", err)), nil
	}

	newTask := models.Task{
		Title:       title,
		Description: description,
		Status:      models.TaskStatusTodo,
		Priority:    models.TaskPriority(priorityStr),
		ProjectID:   projectID,
		AuthorID:    &userID,
		Number:      nextNumber,
	}

	assigneeIDStr, _ := args["assignee_id"].(string)
	if assigneeIDStr != "" {
		parsed, err := uuid.Parse(assigneeIDStr)
		if err == nil {
			newTask.AssigneeID = &parsed
		} else {
			return mcp.NewToolResultError("invalid assignee_id format"), nil
		}
	}

	dueDateStr, _ := args["due_date"].(string)
	if dueDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, dueDateStr)
		if err == nil {
			newTask.DueDate = &parsed
		} else {
			return mcp.NewToolResultError("invalid due_date format, expected RFC3339"), nil
		}
	}

	created, err := storage.Tasks().Create(newTask)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create task: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully created task '%s' (ID: %s)", created.Title, created.ID)), nil
}

func handleListProjectTasks(ctx context.Context, req mcp.CallToolRequest, storage models.Storage, logger *zap.Logger) (*mcp.CallToolResult, error) {
	_, err := getUserIDFromContext(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := req.Params.Arguments.(map[string]any)
	projectIDStr, ok := args["project_id"].(string)
	if !ok || projectIDStr == "" {
		return mcp.NewToolResultError("project_id is required"), nil
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return mcp.NewToolResultError("invalid project_id format"), nil
	}

	var statusFilter *models.TaskStatus
	if st, ok := args["status"].(string); ok && st != "" {
		s := models.TaskStatus(st)
		statusFilter = &s
	}

	tasks, err := storage.Tasks().ListByProjectID(projectID, statusFilter, nil)
	if err != nil {
		logger.Error("Failed to list project tasks", zap.Error(err))
		return mcp.NewToolResultError("Failed to fetch tasks from database"), nil
	}

	if len(tasks) == 0 {
		return mcp.NewToolResultText("No tasks found."), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d tasks in project %s:\n\n", len(tasks), projectID)
	for _, t := range tasks {
		assignee := "Unassigned"
		if t.AssigneeFullName != nil {
			assignee = *t.AssigneeFullName
		}
		fmt.Fprintf(&sb, "- [%s] %s (ID: %s)\n", t.Status, t.Title, t.ID)
		fmt.Fprintf(&sb, "  Assignee: %s\n", assignee)
		fmt.Fprintf(&sb, "  Priority: %s\n", t.Priority)
	}

	return mcp.NewToolResultText(sb.String()), nil
}
