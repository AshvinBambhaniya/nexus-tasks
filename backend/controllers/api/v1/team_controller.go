package v1

import (
	"net/http"

	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/AshvinBambhaniya/nexus-tasks/services"
	"github.com/AshvinBambhaniya/nexus-tasks/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

type TeamController struct {
	teamService services.TeamService
	logger      *zap.Logger
}

func NewTeamController(teamService services.TeamService, logger *zap.Logger) (*TeamController, error) {
	return &TeamController{
		teamService: teamService,
		logger:      logger,
	}, nil
}

func (ctrl *TeamController) Create(c *fiber.Ctx) error {
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

func (ctrl *TeamController) List(c *fiber.Ctx) error {
	wsIDStr := c.Params(constants.ParamWorkspaceID)
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid workspace id")
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

func (ctrl *TeamController) Get(c *fiber.Ctx) error {
	teamIDStr := c.Params(constants.ParamTeamID)
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid team id")
	}

	teamWithProjects, err := ctrl.teamService.GetTeam(teamID)
	if err != nil {
		return utils.JSONFail(c, http.StatusNotFound, "Team not found")
	}

	return utils.JSONSuccess(c, http.StatusOK, teamWithProjects)
}

func (ctrl *TeamController) Update(c *fiber.Ctx) error {
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

	teamIDStr := c.Params(constants.ParamTeamID)
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid team id")
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

func (ctrl *TeamController) Delete(c *fiber.Ctx) error {
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

	teamIDStr := c.Params(constants.ParamTeamID)
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid team id")
	}

	err = ctrl.teamService.DeleteTeam(uid, wsID, teamID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{"message": "Team deleted"})
}

func (ctrl *TeamController) AddMember(c *fiber.Ctx) error {
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

	teamIDStr := c.Params(constants.ParamTeamID)
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid team id")
	}

	var req structs.ReqAddTeamMember
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	err = ctrl.teamService.AddMember(uid, wsID, teamID, req.Email, req.Role)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{"message": "Member added"})
}

func (ctrl *TeamController) RemoveMember(c *fiber.Ctx) error {
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

	teamIDStr := c.Params(constants.ParamTeamID)
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid team id")
	}

	targetUserIDStr := c.Params(constants.ParamUid)
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid target user id")
	}

	err = ctrl.teamService.RemoveMember(uid, wsID, teamID, targetUserID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{"message": "Member removed"})
}

func (ctrl *TeamController) ListMembers(c *fiber.Ctx) error {
	teamIDStr := c.Params(constants.ParamTeamID)
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid team id")
	}

	members, err := ctrl.teamService.ListMembersByTeamId(teamID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, members)
}
