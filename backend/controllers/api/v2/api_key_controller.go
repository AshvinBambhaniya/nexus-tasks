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
type createAPIKeyRequest struct {
	Name string `json:"name"`
}

// createAPIKeyResponse represents the response after creating an API key.
// The raw_token is only included in this response and never shown again.
type createAPIKeyResponse struct {
	RawToken string      `json:"raw_token"`
	Key      interface{} `json:"key"`
}

// Create handles POST /api/v2/auth/api-keys
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

// List handles GET /api/v2/auth/api-keys
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

// Revoke handles DELETE /api/v2/auth/api-keys/:keyId
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
