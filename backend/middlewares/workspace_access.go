package middlewares

import (
	"database/sql"
	"net/http"

	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CheckAccess verifies if the authenticated user is a member of the workspace
func (m *Middleware) CheckAccess(c *fiber.Ctx) error {
	uidStr := c.Locals(constants.ContextUid).(string)
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		m.logger.Error("invalid user id in context", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "invalid user id in context")
	}

	wsIDStr := c.Params(constants.ParamWorkspaceID)
	if wsIDStr == "" {
		return utils.JSONFail(c, http.StatusBadRequest, "workspace id is required")
	}

	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid workspace id")
	}

	member, err := m.workspaceModel.GetMember(wsID, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusForbidden, "unauthorized: not a member of this workspace")
		}
		m.logger.Error("error while checking workspace access", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "error while checking workspace access")
	}

	// Store member info in locals for later use (e.g. role check)
	c.Locals(constants.ContextWorkspaceMember, member)

	return c.Next()
}
