package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/services"
	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// RegisterProjectTools registers all project and workspace related tools.
func RegisterProjectTools(s *server.MCPServer, storage models.Storage, logger *zap.Logger) {
	// list_workspaces
	listWorkspacesTool := mcp.NewTool("list_workspaces",
		mcp.WithDescription("Lists all workspaces the currently authenticated user belongs to."),
	)
	s.AddTool(listWorkspacesTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleListWorkspaces(ctx, req, storage, logger)
	})

	// list_projects
	listProjectsTool := mcp.NewTool("list_projects",
		mcp.WithDescription("Lists all active projects within a specific workspace."),
		mcp.WithString("workspace_id",
			mcp.Required(),
			mcp.Description("The UUID of the workspace to list projects for."),
		),
	)
	s.AddTool(listProjectsTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleListProjects(ctx, req, storage, logger)
	})

	// list_project_members
	listProjectMembersTool := mcp.NewTool("list_project_members",
		mcp.WithDescription("Lists all members of a specific project, including their user IDs and roles."),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("The UUID of the project to list members for."),
		),
	)
	s.AddTool(listProjectMembersTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleListProjectMembers(ctx, req, storage, logger)
	})
}

func handleListWorkspaces(ctx context.Context, _ mcp.CallToolRequest, storage models.Storage, logger *zap.Logger) (*mcp.CallToolResult, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	workspaces, err := storage.Workspaces().ListWorkspacesByUserID(userID)
	if err != nil {
		logger.Error("Failed to list workspaces", zap.Error(err))
		return mcp.NewToolResultError("Failed to fetch workspaces from database"), nil
	}

	if len(workspaces) == 0 {
		return mcp.NewToolResultText("No workspaces found for the user."), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d workspaces:\n\n", len(workspaces))
	for _, ws := range workspaces {
		fmt.Fprintf(&sb, "- %s (ID: %s)\n", ws.Name, ws.ID)
		fmt.Fprintf(&sb, "  Type: %s\n", ws.Type)
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func handleListProjects(ctx context.Context, req mcp.CallToolRequest, storage models.Storage, logger *zap.Logger) (*mcp.CallToolResult, error) {
	// Security check: verify the user is a member of the workspace
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := req.Params.Arguments.(map[string]any)
	workspaceIDStr, ok := args["workspace_id"].(string)
	if !ok || workspaceIDStr == "" {
		return mcp.NewToolResultError("workspace_id is required"), nil
	}

	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		return mcp.NewToolResultError("invalid workspace_id format"), nil
	}

	// Verify access
	_, err = storage.Workspaces().GetMember(workspaceID, userID)
	if err != nil {
		logger.Warn("Unauthorized access attempt to workspace", zap.String("userID", userID.String()), zap.String("workspaceID", workspaceID.String()))
		return mcp.NewToolResultError("Unauthorized: You do not have access to this workspace"), nil
	}

	projects, err := storage.Projects().ListByWorkspaceID(workspaceID)
	if err != nil {
		logger.Error("Failed to list projects", zap.Error(err))
		return mcp.NewToolResultError("Failed to fetch projects from database"), nil
	}

	if len(projects) == 0 {
		return mcp.NewToolResultText("No active projects found in this workspace."), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d active projects in workspace %s:\n\n", len(projects), workspaceID)
	for _, p := range projects {
		fmt.Fprintf(&sb, "- %s (ID: %s)\n", p.Name, p.ID)
		if p.Description != "" {
			fmt.Fprintf(&sb, "  Description: %s\n", p.Description)
		}
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func handleListProjectMembers(ctx context.Context, req mcp.CallToolRequest, storage models.Storage, logger *zap.Logger) (*mcp.CallToolResult, error) {
	// Security check: verify the user is authenticated
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

	// Use the ProjectService to fetch members (which already handles deduplication and access validation)
	projectSvc := services.NewProjectService(storage, logger)
	members, err := projectSvc.ListMembers(userID, projectID)
	if err != nil {
		logger.Error("Failed to list project members via service", zap.Error(err))
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch project members: %v", err)), nil
	}

	if len(members) == 0 {
		return mcp.NewToolResultText("No members found in this project."), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d members in project %s:\n\n", len(members), projectID)
	for _, m := range members {
		name := m.FullName
		if name == "" {
			name = m.Email
		}

		via := "Direct Member"
		if !m.IsDirect {
			via = "Via Team"
		}

		fmt.Fprintf(&sb, "- %s (ID: %s)\n", name, m.UserID)
		fmt.Fprintf(&sb, "  Role: %s\n", m.Role)
		fmt.Fprintf(&sb, "  Email: %s\n", m.Email)
		fmt.Fprintf(&sb, "  Access Via: %s\n", via)
	}

	return mcp.NewToolResultText(sb.String()), nil
}
