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

// ProjectController handles project related requests.
type ProjectController struct {
	projectService services.ProjectService
	logger         *zap.Logger
}

// NewProjectController creates a new instance of ProjectController.
func NewProjectController(projectService services.ProjectService, logger *zap.Logger) (*ProjectController, error) {

	return &ProjectController{
		projectService: projectService,
		logger:         logger,
	}, nil
}

// Create handles the creation of a new project.
func (ctrl *ProjectController) Create(c *fiber.Ctx) error {
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

// List handles listing all projects in a workspace.
func (ctrl *ProjectController) List(c *fiber.Ctx) error {
	wsIDStr := c.Params(constants.ParamWorkspaceID)
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidWorkspaceID)
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

// Get handles fetching a single project by its ID.
func (ctrl *ProjectController) Get(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidProjectID)
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

// Update handles updating an existing project.
func (ctrl *ProjectController) Update(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidProjectID)
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

// AddMember handles adding a new member to a project.
func (ctrl *ProjectController) AddMember(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidProjectID)
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

// RemoveMember handles removing a member from a project.
func (ctrl *ProjectController) RemoveMember(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidProjectID)
	}

	targetUserIDStr := c.Params(constants.ParamUID)
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTargetUserID)
	}

	err = ctrl.projectService.RemoveMember(uid, projectID, targetUserID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{constants.PropMessage: constants.MsgMemberRemoved})
}

// ListMembers handles listing all members of a project.
func (ctrl *ProjectController) ListMembers(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidProjectID)
	}

	members, err := ctrl.projectService.ListMembers(uid, projectID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}
	ctrl.logger.Info("ListMembers response", zap.Any("members", members))

	return utils.JSONSuccess(c, http.StatusOK, members)
}

// AddTeam handles adding a team to a project.
func (ctrl *ProjectController) AddTeam(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidProjectID)
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

// RemoveTeam handles removing a team from a project.
func (ctrl *ProjectController) RemoveTeam(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidProjectID)
	}

	teamIDStr := c.Params(constants.ParamTeamID)
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTeamID)
	}

	err = ctrl.projectService.RemoveTeam(uid, projectID, teamID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{constants.PropMessage: constants.MsgTeamRemoved})
}

// ListTeams handles listing all teams associated with a project.
func (ctrl *ProjectController) ListTeams(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	projectIDStr := c.Params(constants.ParamProjectID)
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidProjectID)
	}

	teams, err := ctrl.projectService.ListTeams(uid, projectID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, teams)
}
