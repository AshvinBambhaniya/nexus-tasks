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

// APIKeyController handles API key management endpoints.
type APIKeyController struct {
	apiKeyService services.APIKeyService
	logger        *zap.Logger
}

// NewAPIKeyController creates a new instance of APIKeyController.
func NewAPIKeyController(apiKeyService services.APIKeyService, logger *zap.Logger) *APIKeyController {
	return &APIKeyController{
		apiKeyService: apiKeyService,
		logger:        logger,
	}
}

// createAPIKeyRequest represents the request body for creating an API key.
//
// swagger:model createAPIKeyRequest
type createAPIKeyRequest struct {
	// Human-readable label for this API key.
	// Required: true
	// example: CI/CD Pipeline Key
	Name string `json:"name"`
}

// createAPIKeyResponse represents the response after creating an API key.
// The raw_token is only included in this response and never shown again.
//
// swagger:model createAPIKeyResponse
type createAPIKeyResponse struct {
	// The plaintext API key token. Store this securely — it is only shown once.
	// example: nxt_abc123secrettoken
	RawToken string `json:"raw_token"`

	// The stored API key metadata (without the secret).
	Key interface{} `json:"key"`
}

// Create handles POST /api/v2/auth/api-keys.
//
/*
 swagger:operation POST /auth/api-keys apiKeys createAPIKey

 # Create an API key

 Generates a new long-lived API key for machine-to-machine access. The raw token
 is returned only once and must be stored securely by the caller.

 ---
 consumes:
 - application/json
 produces:
 - application/json
 security:
 - cookieAuth: []
 parameters:
   - name: body
     in: body
     required: true
     schema:
        "$ref": "#/definitions/createAPIKeyRequest"

 responses:

	201:
	  description: API key created. The raw_token is shown only once.
	  schema:
      "$ref": "#/definitions/createAPIKeyResponse"
	400:
	  description: Name is missing or request body is malformed.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
func (ac *APIKeyController) Create(c *fiber.Ctx) error {
	userIDStr := utils.GetString(c.Locals(constants.ContextUID))
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ac.logger.Error("invalid user id", zap.Error(err))
		return utils.JSONError(c, http.StatusUnauthorized, "invalid user context")
	}

	var req createAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid request body")
	}

	if req.Name == "" {
		return utils.JSONFail(c, http.StatusBadRequest, "name is required")
	}

	rawToken, key, err := ac.apiKeyService.GenerateKey(userID, req.Name)
	if err != nil {
		ac.logger.Error("failed to generate api key", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "failed to generate api key")
	}

	return utils.JSONSuccess(c, http.StatusCreated, createAPIKeyResponse{
		RawToken: rawToken,
		Key:      key,
	})
}

// List handles GET /api/v2/auth/api-keys.
//
/*
 swagger:operation GET /auth/api-keys apiKeys listAPIKeys

 # List API keys

 Returns all API keys belonging to the authenticated user. The raw tokens are never included.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 responses:

	200:
	  description: List of API key metadata objects.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
func (ac *APIKeyController) List(c *fiber.Ctx) error {
	userIDStr := utils.GetString(c.Locals(constants.ContextUID))
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ac.logger.Error("invalid user id", zap.Error(err))
		return utils.JSONError(c, http.StatusUnauthorized, "invalid user context")
	}

	keys, err := ac.apiKeyService.ListKeys(userID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "failed to list api keys")
	}

	return utils.JSONSuccess(c, http.StatusOK, keys)
}

// Revoke handles DELETE /api/v2/auth/api-keys/:keyId.
//
/*
 swagger:operation DELETE /auth/api-keys/{keyId} apiKeys revokeAPIKey

 # Revoke an API key

 Permanently invalidates the specified API key. Any requests authenticated with this
 key will be rejected immediately after revocation.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 parameters:
   - name: keyId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the API key to revoke.

 responses:

	200:
	  description: API key revoked.
	400:
	  description: Invalid key ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
func (ac *APIKeyController) Revoke(c *fiber.Ctx) error {
	userIDStr := utils.GetString(c.Locals(constants.ContextUID))
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ac.logger.Error("invalid user id", zap.Error(err))
		return utils.JSONError(c, http.StatusUnauthorized, "invalid user context")
	}

	keyIDStr := c.Params(constants.ParamKeyID)
	keyID, err := uuid.Parse(keyIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidKeyID)
	}

	err = ac.apiKeyService.RevokeKey(keyID, userID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "failed to revoke api key")
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.MsgKeyRevoked)
}
