// Package tools contains MCP tool implementations
package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// RegisterCommentTools registers tools for managing task comments.
func RegisterCommentTools(s *server.MCPServer, storage models.Storage, logger *zap.Logger) {
	// list_task_comments
	listTaskCommentsTool := mcp.NewTool("list_task_comments",
		mcp.WithDescription("Lists all comments for a specific task."),
		mcp.WithString("task_id",
			mcp.Required(),
			mcp.Description("The UUID of the task to retrieve comments for."),
		),
	)
	s.AddTool(listTaskCommentsTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleListTaskComments(ctx, req, storage, logger)
	})

	// add_comment
	addCommentTool := mcp.NewTool("add_comment",
		mcp.WithDescription("Adds a new comment to a specific task."),
		mcp.WithString("task_id",
			mcp.Required(),
			mcp.Description("The UUID of the task to add a comment to."),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("The content of the comment in Markdown format."),
		),
	)
	s.AddTool(addCommentTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleAddComment(ctx, req, storage, logger)
	})
}

func handleListTaskComments(ctx context.Context, req mcp.CallToolRequest, storage models.Storage, logger *zap.Logger) (*mcp.CallToolResult, error) {
	// Security: The actual implementation should verify the user has access to the task's project
	_, err := getUserIDFromContext(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := req.Params.Arguments.(map[string]any)
	taskIDStr, ok := args["task_id"].(string)
	if !ok || taskIDStr == "" {
		return mcp.NewToolResultError("task_id is required"), nil
	}

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return mcp.NewToolResultError("invalid task_id format"), nil
	}

	comments, err := storage.Comments().ListByTaskID(taskID)
	if err != nil {
		logger.Error("Failed to list comments", zap.Error(err))
		return mcp.NewToolResultError("Failed to fetch comments from database"), nil
	}

	if len(comments) == 0 {
		return mcp.NewToolResultText("No comments found on this task."), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d comments:\n\n", len(comments))
	for _, c := range comments {
		author := c.AuthorFullName
		if author == "" {
			author = "Unknown User"
		}
		fmt.Fprintf(&sb, "--- %s at %s\n%s\n\n", author, c.CreatedAt.Format("2006-01-02 15:04"), c.Content)
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func handleAddComment(ctx context.Context, req mcp.CallToolRequest, storage models.Storage, logger *zap.Logger) (*mcp.CallToolResult, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := req.Params.Arguments.(map[string]any)

	taskIDStr, ok := args["task_id"].(string)
	if !ok || taskIDStr == "" {
		return mcp.NewToolResultError("task_id is required"), nil
	}
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return mcp.NewToolResultError("invalid task_id format"), nil
	}

	content, ok := args["content"].(string)
	if !ok || content == "" {
		return mcp.NewToolResultError("content is required"), nil
	}

	// Verify the task exists
	_, err = storage.Tasks().GetByID(taskID)
	if err != nil {
		return mcp.NewToolResultError("task not found"), nil
	}

	newComment := models.Comment{
		TaskID:   taskID,
		AuthorID: userID,
		Content:  content,
	}

	created, err := storage.Comments().Create(newComment)
	if err != nil {
		logger.Error("Failed to create comment", zap.Error(err))
		return mcp.NewToolResultError("Failed to create comment"), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully added comment (ID: %s) to task %s", created.ID, taskIDStr)), nil
}
