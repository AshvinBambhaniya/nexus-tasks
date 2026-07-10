// Package mcp provides the Model Context Protocol (MCP) server implementation
package mcp

import (
	"github.com/AshvinBambhaniya/nexus-tasks/v2/config"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/mcp/tools"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// NewServer creates a new MCP server instance
func NewServer(storage models.Storage, logger *zap.Logger, _ *config.AppConfig) *server.MCPServer {
	// Create the server
	s := server.NewMCPServer(
		"nexus-mcp-server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Register all tools
	tools.RegisterTaskTools(s, storage, logger)
	tools.RegisterProjectTools(s, storage, logger)
	tools.RegisterCommentTools(s, storage, logger)

	return s
}
