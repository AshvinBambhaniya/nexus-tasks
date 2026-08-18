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

// TeamController handles team related requests.
type TeamController struct {
	teamService services.TeamService
	logger      *zap.Logger
}

// NewTeamController creates a new instance of TeamController.
func NewTeamController(teamService services.TeamService, logger *zap.Logger) (*TeamController, error) {
	return &TeamController{
		teamService: teamService,
		logger:      logger,
	}, nil
}

// Create handles the creation of a new team in a workspace.
//
/*
 swagger:operation POST /workspaces/{workspaceId}/teams teams createTeam

 # Create a team

 Creates a new team within the specified workspace. Requires workspace admin privileges.

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
        "$ref": "#/definitions/ReqCreateTeam"

 responses:

	201:
	  description: Team created.
	  schema:
      "$ref": "#/definitions/ResTeam"
	400:
	  description: Validation error.
	401:
	  description: Not authenticated.
	403:
	  description: Insufficient permissions.
	500:
	  description: Internal server error.
*/
func (ctrl *TeamController) Create(c *fiber.Ctx) error {
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

	var req structs.ReqCreateTeam
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	team, err := ctrl.teamService.CreateTeam(uid, wsID, req)
	if err != nil {
		if err.Error() == "unauthorized: only workspace admins can create teams" {
			return utils.JSONFail(c, http.StatusForbidden, err.Error())
		}
		ctrl.logger.Error("failed to create team", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "failed to create team")
	}

	return utils.JSONSuccess(c, http.StatusCreated, structs.ResTeam{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		WorkspaceID: team.WorkspaceID,
	})
}

// List handles listing all teams in a workspace.
//
/*
 swagger:operation GET /workspaces/{workspaceId}/teams teams listTeams

 # List teams in a workspace

 Returns all teams that belong to the specified workspace.

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
	  description: List of teams.
	  schema:
	    type: array
	    items:
	      "$ref": "#/definitions/ResTeam"
	400:
	  description: Invalid workspace ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
func (ctrl *TeamController) List(c *fiber.Ctx) error {
	wsIDStr := c.Params(constants.ParamWorkspaceID)
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidWorkspaceID)
	}

	teams, err := ctrl.teamService.ListTeamsByWorkspaceID(wsID)
	if err != nil {
		ctrl.logger.Error("failed to list teams", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "failed to list teams")
	}

	var res []structs.ResTeam
	for _, t := range teams {
		res = append(res, structs.ResTeam{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			WorkspaceID: t.WorkspaceID,
		})
	}

	return utils.JSONSuccess(c, http.StatusOK, res)
}

// Get handles fetching a single team by its ID, including its projects.
//
/*
 swagger:operation GET /workspaces/{workspaceId}/teams/{teamId} teams getTeam

 # Get a team by ID

 Returns the details of a single team, including its associated projects.

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
   - name: teamId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the team.

 responses:

	200:
	  description: Team details with associated projects.
	  schema:
      "$ref": "#/definitions/ResTeamWithProjects"
	400:
	  description: Invalid team ID.
	401:
	  description: Not authenticated.
	404:
	  description: Team not found.
	500:
	  description: Internal server error.
*/
func (ctrl *TeamController) Get(c *fiber.Ctx) error {
	teamIDStr := c.Params(constants.ParamTeamID)
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTeamID)
	}

	teamWithProjects, err := ctrl.teamService.GetTeam(teamID)
	if err != nil {
		return utils.JSONFail(c, http.StatusNotFound, "Team not found")
	}

	return utils.JSONSuccess(c, http.StatusOK, teamWithProjects)
}

// Update handles updating an existing team.
//
/*
 swagger:operation PATCH /workspaces/{workspaceId}/teams/{teamId} teams updateTeam

 # Update a team

 Updates the name and/or description of a team. Requires workspace admin privileges.

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
   - name: teamId
     in: path
     required: true
     type: string
     format: uuid
   - name: body
     in: body
     required: true
     schema:
        "$ref": "#/definitions/ReqUpdateTeam"

 responses:

	200:
	  description: Team updated.
	  schema:
      "$ref": "#/definitions/ResTeam"
	400:
	  description: Invalid request.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
func (ctrl *TeamController) Update(c *fiber.Ctx) error {
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

	teamIDStr := c.Params(constants.ParamTeamID)
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTeamID)
	}

	var req structs.ReqUpdateTeam
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	team, err := ctrl.teamService.UpdateTeam(uid, wsID, teamID, req)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, structs.ResTeam{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		WorkspaceID: team.WorkspaceID,
	})
}

// Delete handles deleting a team.
//
/*
 swagger:operation DELETE /workspaces/{workspaceId}/teams/{teamId} teams deleteTeam

 # Delete a team

 Permanently removes a team from the workspace. Requires workspace admin privileges.

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
   - name: teamId
     in: path
     required: true
     type: string
     format: uuid

 responses:

	200:
	  description: Team deleted.
	400:
	  description: Invalid team or workspace ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
func (ctrl *TeamController) Delete(c *fiber.Ctx) error {
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

	teamIDStr := c.Params(constants.ParamTeamID)
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTeamID)
	}

	err = ctrl.teamService.DeleteTeam(uid, wsID, teamID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{constants.PropMessage: constants.MsgTeamDeleted})
}

// AddMember handles adding a new member to a team.
//
/*
 swagger:operation POST /workspaces/{workspaceId}/teams/{teamId}/members teams addTeamMember

 # Add a member to a team

 Adds an existing workspace member to the specified team with an assigned role.

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
   - name: teamId
     in: path
     required: true
     type: string
     format: uuid
   - name: body
     in: body
     required: true
     schema:
        "$ref": "#/definitions/ReqAddTeamMember"

 responses:

	200:
	  description: Member added successfully.
	400:
	  description: Validation error or invalid path parameter.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
func (ctrl *TeamController) AddMember(c *fiber.Ctx) error {
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

	teamIDStr := c.Params(constants.ParamTeamID)
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTeamID)
	}

	var req structs.ReqAddTeamMember
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	err = ctrl.teamService.AddMember(uid, wsID, teamID, req.Email, req.Role)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{constants.PropMessage: constants.MsgMemberAdded})
}

// RemoveMember handles removing a member from a team.
//
/*
 swagger:operation DELETE /workspaces/{workspaceId}/teams/{teamId}/members/{uid} teams removeTeamMember

 # Remove a member from a team

 Removes the specified user from the team. Their workspace membership is not affected.

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
   - name: teamId
     in: path
     required: true
     type: string
     format: uuid
   - name: uid
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the user to remove.

 responses:

	200:
	  description: Member removed.
	400:
	  description: Invalid path parameter.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
func (ctrl *TeamController) RemoveMember(c *fiber.Ctx) error {
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

	teamIDStr := c.Params(constants.ParamTeamID)
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTeamID)
	}

	targetUserIDStr := c.Params(constants.ParamUID)
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTargetUserID)
	}

	err = ctrl.teamService.RemoveMember(uid, wsID, teamID, targetUserID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{constants.PropMessage: constants.MsgMemberRemoved})
}

// ListMembers handles listing all members of a team.
//
/*
 swagger:operation GET /workspaces/{workspaceId}/teams/{teamId}/members teams listTeamMembers

 # List team members

 Returns all members of the specified team along with their roles.

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
   - name: teamId
     in: path
     required: true
     type: string
     format: uuid

 responses:

	200:
	  description: List of team members.
	  schema:
	    type: array
	    items:
	      "$ref": "#/definitions/ResTeamMember"
	400:
	  description: Invalid team ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
func (ctrl *TeamController) ListMembers(c *fiber.Ctx) error {
	teamIDStr := c.Params(constants.ParamTeamID)
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTeamID)
	}

	members, err := ctrl.teamService.ListMembersByTeamID(teamID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, members)
}
