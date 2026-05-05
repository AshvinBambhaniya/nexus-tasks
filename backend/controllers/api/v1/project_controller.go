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

type ProjectController struct {
	projectService services.ProjectService
	logger         *zap.Logger
}

func NewProjectController(projectService services.ProjectService, logger *zap.Logger) (*ProjectController, error) {

	return &ProjectController{
		projectService: projectService,
		logger:         logger,
	}, nil
}

func (ctrl *ProjectController) Create(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUid))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, "invalid user id")
	}

	wsIDStr := c.Params(constants.ParamWorkspaceID)
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid workspace id")
	}

	var req structs.ReqCreateProject
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	project, err := ctrl.projectService.CreateProject(uid, wsID, req)
	if err != nil {
		if err.Error() == "unauthorized: only workspace admins can create projects" {
			return utils.JSONFail(c, http.StatusForbidden, err.Error())
		}
		ctrl.logger.Error("failed to create project", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "failed to create project")
	}

	return utils.JSONSuccess(c, http.StatusCreated, structs.ResProject{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		IsArchived:  project.IsArchived,
		WorkspaceID: project.WorkspaceID,
		CreatedAt:   project.CreatedAt,
	})
}

func (ctrl *ProjectController) List(c *fiber.Ctx) error {
	wsIDStr := c.Params(constants.ParamWorkspaceID)
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid workspace id")
	}

	projects, err := ctrl.projectService.ListByWorkspaceID(wsID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	var res []structs.ResProject
	for _, p := range projects {
		res = append(res, structs.ResProject{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			IsArchived:  p.IsArchived,
			WorkspaceID: p.WorkspaceID,
			CreatedAt:   p.CreatedAt,
		})
	}

	return utils.JSONSuccess(c, http.StatusOK, res)
}

func (ctrl *ProjectController) Get(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUid))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, "invalid user id")
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid project id")
	}

	project, err := ctrl.projectService.GetProject(uid, projectID)
	if err != nil {
		return utils.JSONFail(c, http.StatusNotFound, "Project not found or unauthorized")
	}

	return utils.JSONSuccess(c, http.StatusOK, structs.ResProject{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		IsArchived:  project.IsArchived,
		WorkspaceID: project.WorkspaceID,
		CreatedAt:   project.CreatedAt,
	})
}

func (ctrl *ProjectController) Update(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUid))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, "invalid user id")
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid project id")
	}

	var req structs.ReqUpdateProject
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	project, err := ctrl.projectService.UpdateProject(uid, projectID, req)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, structs.ResProject{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		IsArchived:  project.IsArchived,
		WorkspaceID: project.WorkspaceID,
		CreatedAt:   project.CreatedAt,
	})
}

func (ctrl *ProjectController) AddMember(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUid))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, "invalid user id")
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid project id")
	}

	var req structs.ReqAddProjectMember
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	member, err := ctrl.projectService.AddMember(uid, projectID, req)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	// Fetch user email for response
	return utils.JSONSuccess(c, http.StatusOK, structs.ResProjectMember{
		UserID:   member.UserID,
		Email:    req.Email, // Echo back
		Role:     member.Role,
		IsDirect: true,
	})
}

func (ctrl *ProjectController) RemoveMember(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUid))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, "invalid user id")
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid project id")
	}

	targetUserIDStr := c.Params(constants.ParamUid)
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid target user id")
	}

	err = ctrl.projectService.RemoveMember(uid, projectID, targetUserID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{"message": "Member removed"})
}

func (ctrl *ProjectController) ListMembers(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUid))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, "invalid user id")
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid project id")
	}

	members, err := ctrl.projectService.ListMembers(uid, projectID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, members)
}

func (ctrl *ProjectController) AddTeam(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUid))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, "invalid user id")
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid project id")
	}

	var req structs.ReqAddProjectTeam
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	projectTeam, err := ctrl.projectService.AddTeam(uid, projectID, req.TeamID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, projectTeam)
}

func (ctrl *ProjectController) RemoveTeam(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUid))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, "invalid user id")
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid project id")
	}

	teamIDStr := c.Params(constants.ParamTeamID)
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid team id")
	}

	err = ctrl.projectService.RemoveTeam(uid, projectID, teamID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{"message": "Team removed"})
}

func (ctrl *ProjectController) ListTeams(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUid))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, "invalid user id")
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "invalid project id")
	}

	teams, err := ctrl.projectService.ListTeams(uid, projectID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, teams)
}
