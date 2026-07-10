package cli

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/config"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/database"
	mcpserver "github.com/AshvinBambhaniya/nexus-tasks/v2/mcp"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/mcp/tools"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/services"
	mcp_server "github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// GetMCPCommandDef returns the command definition for starting the MCP server
func GetMCPCommandDef(cfg *config.AppConfig, logger *zap.Logger) cobra.Command {
	return cobra.Command{
		Use:   "mcp",
		Short: "Start the standalone MCP server",
		RunE: func(_ *cobra.Command, _ []string) error {
			// 1. Initialize Database
			db, err := database.Connect(cfg.DB)
			if err != nil {
				return err
			}
			storage := models.NewStorage(db)
			apiKeyService := services.NewAPIKeyService(storage, logger)

			// 2. Initialize MCP Server with domain services
			mcpServer := mcpserver.NewServer(storage, logger, cfg)

			// 3. Setup Streamable HTTP Transport
			// The client's configuration implies it is using the Stateless HTTP transport
			// (sending POST/DELETE requests) instead of the Server-Sent Events (SSE) transport.
			httpServer := mcp_server.NewStreamableHTTPServer(mcpServer,
				mcp_server.WithEndpointPath("/sse"), // Listen on /sse because of the config serverUrl
				mcp_server.WithStreamableHTTPCORS(), // Handle CORS OPTIONS requests natively
			)

			mux := http.NewServeMux()

			// The Antigravity CLI will connect to /sse to open the event stream.
			mux.Handle("/sse", httpServer)

			// 5. Authentication & Logging Middleware
			authMiddleware := func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					logger.Info("Incoming HTTP Request", zap.String("method", r.Method), zap.String("url", r.URL.String()))

					// We MUST allow OPTIONS requests to pass through without authentication
					// so the WithSSECORS middleware can reply with 200 OK headers.
					if r.Method == http.MethodOptions {
						next.ServeHTTP(w, r)
						return
					}

					authHeader := r.Header.Get("Authorization")
					if !strings.HasPrefix(authHeader, "Bearer ntx_") {
						logger.Warn("Request rejected: invalid or missing authorization header")
						http.Error(w, "Unauthorized", http.StatusUnauthorized)
						return
					}

					rawToken := strings.TrimPrefix(authHeader, "Bearer ")
					user, _, err := apiKeyService.ValidateToken(rawToken)
					if err != nil {
						logger.Warn("Request rejected: invalid token", zap.Error(err))
						http.Error(w, "Unauthorized", http.StatusUnauthorized)
						return
					}

					logger.Info("Request received with valid token", zap.String("userID", user.ID.String()))

					// Inject UserID into context
					ctx := context.WithValue(r.Context(), tools.UserIDKey, user.ID)
					next.ServeHTTP(w, r.WithContext(ctx))
				})
			}

			// 3. Graceful Shutdown Context
			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

			server := &http.Server{
				Addr:              ":8001",
				Handler:           authMiddleware(mux),
				ReadHeaderTimeout: 10 * time.Second,
			}

			go func() {
				logger.Info("Starting MCP Ser	ver (net/http)", zap.String("port", ":8001"))
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("MCP server error", zap.Error(err))
				}
			}()

			<-interrupt
			logger.Info("Gracefully shutting down MCP server...")
			return server.Close()
		},
	}
}
