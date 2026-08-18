package v2

import (
	"net/http"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/structs"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/services"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

// WorkspaceController handles workspace-related requests.
type WorkspaceController struct {
	workspaceService services.WorkspaceService
	logger           *zap.Logger
}

// NewWorkspaceController creates a new instance of WorkspaceController.
func NewWorkspaceController(workspaceService services.WorkspaceService, logger *zap.Logger) (*WorkspaceController, error) {
	return &WorkspaceController{
		workspaceService: workspaceService,
		logger:           logger,
	}, nil
}

// Create handles the creation of a new workspace.
//
/*
 swagger:operation POST /workspaces workspaces createWorkspace

 # Create a workspace

 Creates a new workspace owned by the authenticated user.

 ---
 consumes:
 - application/json
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: body
     in: body
     required: true
     schema:
        "$ref": "#/definitions/ReqCreateWorkspace"

 responses:

	201:
	  description: Workspace created.
	  schema:
      "$ref": "#/definitions/ResWorkspace"
	400:
	  description: Validation error.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
func (ctrl *WorkspaceController) Create(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		ctrl.logger.Error(constants.ErrInvalidUserID, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
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

// List returns a list of workspaces for the authenticated user.
//
/*
 swagger:operation GET /workspaces workspaces listWorkspaces

 # List workspaces

 Returns all workspaces the authenticated user is a member of.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 responses:

	200:
	  description: List of workspaces.
	  schema:
	    type: array
	    items:
	      "$ref": "#/definitions/ResWorkspace"
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
func (ctrl *WorkspaceController) List(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		ctrl.logger.Error(constants.ErrInvalidUserID, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
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

// ListMembers returns a list of members for a specific workspace.
//
/*
 swagger:operation GET /workspaces/{workspaceId}/members workspaces listWorkspaceMembers

 # List workspace members

 Returns all users who are members of the given workspace.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: workspaceId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the workspace.

 responses:

	200:
	  description: List of workspace members.
	  schema:
	    type: array
	    items:
	      "$ref": "#/definitions/ResWorkspaceMember"
	400:
	  description: Invalid workspace ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
func (ctrl *WorkspaceController) ListMembers(c *fiber.Ctx) error {

	wsIDStr := c.Params(constants.ParamWorkspaceID)
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidWorkspaceID)
	}

	members, err := ctrl.workspaceService.ListMembersByWorkspaceID(wsID)
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

// InviteMember invites a new member to the workspace.
//
/*
 swagger:operation POST /workspaces/{workspaceId}/members workspaces inviteWorkspaceMember

 # Invite a user to a workspace

 Sends an invitation (or directly adds) a user to the workspace by email address.

 ---
 consumes:
 - application/json
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: workspaceId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the workspace.
   - name: body
     in: body
     required: true
     schema:
        "$ref": "#/definitions/ReqInviteWorkspaceMember"

 responses:

	200:
	  description: Member invited successfully.
	400:
	  description: Validation error or invalid workspace ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
func (ctrl *WorkspaceController) InviteMember(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	wsIDStr := c.Params(constants.ParamWorkspaceID)
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidWorkspaceID)
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

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{constants.PropMessage: constants.MsgMemberInvited})
}

// RemoveMember removes a member from the workspace.
//
/*
 swagger:operation DELETE /workspaces/{workspaceId}/members workspaces removeWorkspaceMember

 # Remove a member from a workspace

 Removes the specified user from the workspace. Only workspace admins may perform this action.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 parameters:
   - name: workspaceId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the workspace.
   - name: uid
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the user to remove.

 responses:

	200:
	  description: Member removed successfully.
	400:
	  description: Invalid workspace or user ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
func (ctrl *WorkspaceController) RemoveMember(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	wsIDStr := c.Params(constants.ParamWorkspaceID)
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidWorkspaceID)
	}

	targetUIDStr := c.Params(constants.ParamUID)
	targetUID, err := uuid.Parse(targetUIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTargetUserID)
	}

	err = ctrl.workspaceService.RemoveMember(uid, wsID, targetUID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{constants.PropMessage: constants.MsgMemberRemoved})
}
