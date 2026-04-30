package v1

import (
	"fmt"
	"strconv"

	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/jwt"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/realtime"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WebsocketController struct {
	hub            *realtime.Hub
	projectModel   models.ProjectModel
	workspaceModel models.WorkspaceModel
	teamModel      models.TeamModel
	config         config.AppConfig
	logger         *zap.Logger
}

func NewWebsocketController(hub *realtime.Hub, goqu *goqu.Database, cfg config.AppConfig, logger *zap.Logger) (*WebsocketController, error) {
	projectModel, err := models.InitProjectModel(goqu)
	if err != nil {
		return nil, err
	}
	workspaceModel, err := models.InitWorkspaceModel(goqu)
	if err != nil {
		return nil, err
	}
	teamModel, err := models.InitTeamModel(goqu)
	if err != nil {
		return nil, err
	}

	return &WebsocketController{
		hub:            hub,
		projectModel:   projectModel,
		workspaceModel: workspaceModel,
		teamModel:      teamModel,
		config:         cfg,
		logger:         logger,
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

	claims, err := jwt.ParseToken(ctrl.config, token)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	uidStr := claims.Subject()
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	// Set Locals for the websocket handler
	c.Locals(constants.ContextUid, uid)

	return c.Next()
}

// HandleWorkspaceConnection handles /ws/:workspaceId
// It subscribes the user to the workspace topic AND any relevant project topics if we want "Project" updates
// But strictly following Python: Client connects to Workspace ID.
// However, Tasks are broadcasting to Project ID.
// So we must subscribe to the project topics OR map incoming broadcast to workspace topics.
// BETTER: Subscribe to `workspace:{id}` AND `project:{id}` for all projects the user is in?
// OR: The client connects to specific channels?
// Python client connects to `/ws/{workspaceId}`.
// Python `create_task` broadcasts to `project_id`.
// This implies the Python client connects to a workspace channel but receives messages for projects?
// That only works if `project_id` broadcast also goes to `workspace_id` channel, OR if `manager` maps workspace -> projects.
//
// Let's implement robust logic:
// When connecting to /ws/:workspaceId, we verify access to workspace.
// Then we subscribe to `workspace:{id}`.
// AND we look up all projects this user is part of in that workspace and subscribe to `project:{id}`.
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
		c.WriteJSON(fiber.Map{"error": "Invalid workspace ID"})
		c.Close()
		return
	}

	// 1. Verify Workspace Access
	_, err = ctrl.workspaceModel.GetMember(workspaceID, uid)
	if err != nil {
		c.WriteJSON(fiber.Map{"error": "Unauthorized access to workspace"})
		c.Close()
		return
	}

	// 2. Subscribe to Workspace Topic
	wsTopic := fmt.Sprintf("workspace:%d", workspaceID)
	ctrl.hub.Subscribe(wsTopic, c)

	// 3. Auto-subscribe to all accessible projects in this workspace
	// (This bridges the gap between connecting to workspace and receiving project updates)
	projects, err := ctrl.projectModel.ListByWorkspaceID(workspaceID)
	if err == nil {
		for _, p := range projects {
			// Check project access (Direct, Team, or Workspace Admin)
			if ctrl.checkProjectAccess(p.ID, uid, workspaceID) {
				pTopic := fmt.Sprintf("project:%d", p.ID)
				ctrl.hub.Subscribe(pTopic, c)
			}
		}
	}

	// Loop to keep connection open and handle unsubscription
	defer func() {
		ctrl.hub.Unsubscribe(wsTopic, c)
		if err == nil {
			for _, p := range projects {
				pTopic := fmt.Sprintf("project:%d", p.ID)
				ctrl.hub.Unsubscribe(pTopic, c)
			}
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

func (ctrl *WebsocketController) checkProjectAccess(projectID, userID, workspaceID uuid.UUID) bool {
	// 1. Direct Member
	_, err := ctrl.projectModel.GetMember(projectID, userID)
	if err == nil {
		return true
	}

	// 2. Workspace Admin
	wsMember, err := ctrl.workspaceModel.GetMember(workspaceID, userID)
	if err == nil && wsMember.Role == models.WorkspaceRoleAdmin {
		return true
	}

	// 3. Team Member
	teams, err := ctrl.projectModel.GetTeams(projectID)
	if err == nil {
		for _, t := range teams {
			_, err := ctrl.teamModel.GetMember(t.TeamID, userID)
			if err == nil {
				return true
			}
		}
	}

	return false
}
