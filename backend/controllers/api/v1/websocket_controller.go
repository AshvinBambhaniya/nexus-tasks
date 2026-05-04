package v1

import (
	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/jwt"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/realtime"
	"github.com/AshvinBambhaniya/nexus-tasks/services"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WebsocketController struct {
	hub              realtime.IHub
	websocketService services.WebsocketService
	config           *config.AppConfig
	logger           *zap.Logger
}

func NewWebsocketController(hub realtime.IHub, websocketService services.WebsocketService, cfg *config.AppConfig, logger *zap.Logger) (*WebsocketController, error) {
	return &WebsocketController{
		hub:              hub,
		websocketService: websocketService,
		config:           cfg,
		logger:           logger,
	}, nil
}

// UpgradeMiddleware handles authentication before upgrade
func (ctrl *WebsocketController) UpgradeMiddleware(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	// Authenticate
	token := c.Cookies(constants.CookieUser)
	if token == "" {
		// Try query param for development ease or if cookies fail
		token = c.Query("access_token")
	}

	if token == "" {
		return fiber.ErrUnauthorized
	}

	claims, err := jwt.ParseToken(ctrl.config.Secret, token)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	uidStr := claims.Subject()
	// Actually we should be using the string UID if we use UUIDs everywhere,
	// but the original code was attempting strconv.Atoi then using uuid.Parse later.
	// Let's keep it consistent with what the user had or improve it if it was a bug.
	// The original code:
	// uid, err := strconv.Atoi(uidStr)
	// c.Locals(constants.ContextUid, uid)
	// Then in HandleWorkspaceConnection:
	// uidStr := c.Locals(constants.ContextUid).(string) // This would panic if it was int!

	// I'll fix it to use string UID as expected by HandleWorkspaceConnection.
	c.Locals(constants.ContextUid, uidStr)

	return c.Next()
}

// HandleWorkspaceConnection handles /ws/:workspaceId
func (ctrl *WebsocketController) HandleWorkspaceConnection(c *websocket.Conn) {
	uidStr := c.Locals(constants.ContextUid).(string)
	workspaceIDStr := c.Params("id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.WriteJSON(fiber.Map{"error": "Invalid workspace ID"})
		c.Close()
		return
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.WriteJSON(fiber.Map{"error": "Invalid user ID"})
		c.Close()
		return
	}

	topics, err := ctrl.websocketService.GetConnectionTopics(uid, workspaceID)
	if err != nil {
		c.WriteJSON(fiber.Map{"error": err.Error()})
		c.Close()
		return
	}

	// Subscribe to all topics
	for _, topic := range topics {
		ctrl.hub.Subscribe(topic, c)
	}

	// Loop to keep connection open and handle unsubscription
	defer func() {
		for _, topic := range topics {
			ctrl.hub.Unsubscribe(topic, c)
		}
		c.Close()
	}()

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
}
