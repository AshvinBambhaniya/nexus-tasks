// Package middlewares provides HTTP middleware for the application.
package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/jwt"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/utils"
	"github.com/gofiber/fiber/v2"
	j "github.com/lestrrat-go/jwx/v2/jwt"
	"go.uber.org/zap"
)

// Authenticated is a middleware that validates the JWT token in the cookie
// or a Bearer personal access token in the Authorization header.
func (m *Middleware) Authenticated(c *fiber.Ctx) error {
	// 1. Try JWT cookie (existing browser session flow)
	token := c.Cookies(constants.CookieUser, "")
	if token != "" {
		claims, err := jwt.ParseToken(m.config.Secret, token)
		if err == nil {
			c.Locals(constants.ContextUID, claims.Subject())
			return c.Next()
		}

		// Only log non-expiry/non-invalid JWT errors
		if !errors.Is(err, j.ErrInvalidJWT()) && !errors.Is(err, j.ErrTokenExpired()) {
			m.logger.Error("error while checking user identity", zap.Error(err))
		}
	}

	// 2. Try Bearer PAT (for MCP servers, CI/CD, external integrations)
	authHeader := c.Get("Authorization")
	var rawToken string
	if strings.HasPrefix(authHeader, "Bearer "+constants.APIKeyPrefix) {
		rawToken = strings.TrimPrefix(authHeader, "Bearer ")
	} else if apiKey := c.Query("apiKey"); strings.HasPrefix(apiKey, constants.APIKeyPrefix) {
		// Fallback to query parameter (required for MCP SSE clients that drop headers on POST)
		rawToken = apiKey
	}

	if rawToken != "" {
		if m.apiKeyService != nil {
			user, _, err := m.apiKeyService.ValidateToken(rawToken)
			if err == nil {
				c.Locals(constants.ContextUID, user.ID.String())
				return c.Next()
			}
			m.logger.Debug("api key validation failed", zap.Error(err))
		}
	}

	// 3. Neither authentication method succeeded
	return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
}
