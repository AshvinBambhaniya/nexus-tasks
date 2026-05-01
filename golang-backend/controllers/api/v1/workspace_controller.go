package v1

import (
	"net/http"

	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/watermill"
	"github.com/AshvinBambhaniya/nexus-tasks/services"
	"github.com/AshvinBambhaniya/nexus-tasks/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

type WorkspaceController struct {
	workspaceService services.WorkspaceService
	logger           *zap.Logger
}

func NewWorkspaceController(workspaceService services.WorkspaceService, logger *zap.Logger, publisher *watermill.WatermillPublisher) (*WorkspaceController, error) {
	return &WorkspaceController{
		workspaceService: workspaceService,
		logger:           logger,
	}, nil
}

func (ctrl *WorkspaceController) Create(c *fiber.Ctx) error {
	uidStr := c.Locals(constants.ContextUid).(string)
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		ctrl.logger.Error("invalid user id in context", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "invalid user id in context")
	}

	var req structs.ReqCreateWorkspace
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	ws, err := ctrl.workspaceService.CreateWorkspace(uid, req)
	if err != nil {
		ctrl.logger.Error("failed to create workspace", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "failed to create workspace")
	}

	return utils.JSONSuccess(c, http.StatusCreated, structs.ResWorkspace{
		ID:      ws.ID,
		Name:    ws.Name,
		Type:    ws.Type,
		OwnerID: ws.OwnerID,
	})
}

func (ctrl *WorkspaceController) List(c *fiber.Ctx) error {
	uidStr := c.Locals(constants.ContextUid).(string)
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		ctrl.logger.Error("invalid user id in context", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "invalid user id in context")
	}

	workspaces, err := ctrl.workspaceService.ListWorkspacesByUserID(uid)
	if err != nil {
		ctrl.logger.Error("failed to list workspaces", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "failed to list workspaces")
	}

	// Map to response struct (though simple here, good practice)
	var res []structs.ResWorkspace
	for _, ws := range workspaces {
		res = append(res, structs.ResWorkspace{
			ID:      ws.ID,
			Name:    ws.Name,
			Type:    ws.Type,
			OwnerID: ws.OwnerID,
		})
	}

	return utils.JSONSuccess(c, http.StatusOK, res)
}

func (ctrl *WorkspaceController) ListMembers(c *fiber.Ctx) error {

	wsIDStr := c.Params(constants.ParamWorkspaceID)
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid workspace id")
	}

	members, err := ctrl.workspaceService.ListMembersByWorkspaceId(wsID)
	if err != nil {
		ctrl.logger.Error("failed to list members", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "failed to list members")
	}

	var res []structs.ResWorkspaceMember
	for _, m := range members {
		res = append(res, structs.ResWorkspaceMember{
			WorkspaceID: m.WorkspaceID,
			UserID:      m.UserID,
			Role:        m.Role,
			User: structs.ResUser{
				ID:       m.UserID,
				Email:    m.Email,
				FullName: m.FullName,
			},
		})
	}

	return utils.JSONSuccess(c, http.StatusOK, res)
}

func (ctrl *WorkspaceController) InviteMember(c *fiber.Ctx) error {
	uidStr := c.Locals(constants.ContextUid).(string)
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, "invalid user id")
	}

	wsIDStr := c.Params(constants.ParamWorkspaceID)
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid workspace id")
	}

	var req structs.ReqInviteWorkspaceMember
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	err = ctrl.workspaceService.InviteMember(uid, wsID, req.Email)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{"message": "Member invited"})
}

func (ctrl *WorkspaceController) RemoveMember(c *fiber.Ctx) error {
	uidStr := c.Locals(constants.ContextUid).(string)
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, "invalid user id")
	}

	wsIDStr := c.Params(constants.ParamWorkspaceID)
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid workspace id")
	}

	targetUserIDStr := c.Params(constants.ParamUid)
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid target user id")
	}

	err = ctrl.workspaceService.RemoveMember(uid, wsID, targetUserID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{"message": "Member removed"})
}
