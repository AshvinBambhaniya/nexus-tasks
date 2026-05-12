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
