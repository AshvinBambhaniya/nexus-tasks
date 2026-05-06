package v1

import (
	"net/http"

	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/AshvinBambhaniya/nexus-tasks/services"
	"github.com/AshvinBambhaniya/nexus-tasks/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

// TaskController handles task related requests.
type TaskController struct {
	taskService    services.TaskService
	commentService services.CommentService
	logger         *zap.Logger
}

// NewTaskController creates a new instance of TaskController.
func NewTaskController(taskService services.TaskService, commentService services.CommentService, logger *zap.Logger) (*TaskController, error) {

	return &TaskController{
		taskService:    taskService,
		commentService: commentService,
		logger:         logger,
	}, nil
}

// CreateTask handles the creation of a new task in a project.
func (ctrl *TaskController) CreateTask(c *fiber.Ctx) error {
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

	var req structs.ReqCreateTask
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	task, err := ctrl.taskService.CreateTask(uid, projectID, req)
	if err != nil {
		ctrl.logger.Error("failed to create task", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusCreated, ctrl.mapTaskToRes(task))
}

// ListProjectTasks handles listing all tasks in a project.
func (ctrl *TaskController) ListProjectTasks(c *fiber.Ctx) error {
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

	// Filter
	statusStr := c.Query(constants.QueryStatus)
	var status *models.TaskStatus
	if statusStr != "" {
		s := models.TaskStatus(statusStr)
		status = &s
	}

	assigneeIDStr := c.Query(constants.QueryAssigneeID)
	var assigneeID *uuid.UUID
	if assigneeIDStr != "" {
		aid, err := uuid.Parse(assigneeIDStr)
		if err == nil {
			assigneeID = &aid
		}
	}

	tasks, err := ctrl.taskService.ListProjectTasks(uid, projectID, status, assigneeID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	var res []structs.ResTask
	for _, t := range tasks {
		res = append(res, ctrl.mapTaskToRes(t))
	}

	return utils.JSONSuccess(c, http.StatusOK, res)
}

// GetTask handles fetching a single task by its ID.
func (ctrl *TaskController) GetTask(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	taskIDStr := c.Params(constants.ParamTaskID)
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTaskID)
	}

	task, err := ctrl.taskService.GetTask(uid, taskID)
	if err != nil {
		return utils.JSONFail(c, http.StatusNotFound, "Task not found")
	}

	return utils.JSONSuccess(c, http.StatusOK, ctrl.mapTaskToRes(task))
}

// UpdateTask handles updating an existing task.
func (ctrl *TaskController) UpdateTask(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	taskIDStr := c.Params(constants.ParamTaskID)
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTaskID)
	}

	var req structs.ReqUpdateTask
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	task, err := ctrl.taskService.UpdateTask(uid, taskID, req)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, ctrl.mapTaskToRes(task))
}

// DeleteTask handles deleting a task.
func (ctrl *TaskController) DeleteTask(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	taskIDStr := c.Params(constants.ParamTaskID)
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTaskID)
	}

	err = ctrl.taskService.DeleteTask(uid, taskID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{constants.PropMessage: constants.MsgTaskDeleted})
}

// ListMyTasks handles listing all tasks assigned to the current user.
func (ctrl *TaskController) ListMyTasks(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	tasks, err := ctrl.taskService.ListMyTasks(uid)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	// For "My Tasks", Python returns task + project info.
	// Since I didn't add join in service yet, I'll return standard task list for now.
	// Enhancement: Add project details to response.

	var res []structs.ResTask
	for _, t := range tasks {
		res = append(res, ctrl.mapTaskToRes(t))
	}

	return utils.JSONSuccess(c, http.StatusOK, res)
}

// Comments

// CreateComment handles creating a new comment on a task.
func (ctrl *TaskController) CreateComment(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	taskIDStr := c.Params(constants.ParamTaskID)
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTaskID)
	}

	var req structs.ReqCreateComment
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	comment, err := ctrl.commentService.CreateComment(uid, taskID, req)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	// Need author info for response, maybe fetch or just return basics
	// Python returns CommentResponse which includes author: UserResponse
	// For now returning basic
	return utils.JSONSuccess(c, http.StatusCreated, structs.ResComment{
		ID:        comment.ID,
		Content:   comment.Content,
		TaskID:    comment.TaskID,
		AuthorID:  comment.AuthorID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	})
}

// ListTaskComments handles listing all comments on a task.
func (ctrl *TaskController) ListTaskComments(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	taskIDStr := c.Params(constants.ParamTaskID)
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidTaskID)
	}

	comments, err := ctrl.commentService.ListTaskComments(uid, taskID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	var res []structs.ResComment
	for _, c := range comments {
		res = append(res, structs.ResComment{
			ID:        c.ID,
			Content:   c.Content,
			TaskID:    c.TaskID,
			AuthorID:  c.AuthorID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Author: structs.ResUser{
				ID:       c.AuthorID,
				Email:    c.AuthorEmail,
				FullName: c.AuthorFullName,
			},
		})
	}

	return utils.JSONSuccess(c, http.StatusOK, res)
}

// DeleteComment handles deleting a comment.
func (ctrl *TaskController) DeleteComment(c *fiber.Ctx) error {
	uidStr := utils.GetString(c.Locals(constants.ContextUID))
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrInvalidUserID)
	}

	commentIDStr := c.Params(constants.ParamCommentID)
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidCommentID)
	}

	err = ctrl.commentService.DeleteComment(uid, commentID)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, fiber.Map{constants.PropMessage: constants.MsgCommentDeleted})
}

func (ctrl *TaskController) mapTaskToRes(t models.Task) structs.ResTask {
	return structs.ResTask{
		ID:          t.ID,
		Number:      t.Number,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Priority:    t.Priority,
		ProjectID:   t.ProjectID,
		AssigneeID:  t.AssigneeID,
		AuthorID:    t.AuthorID,
		DueDate:     t.DueDate,
		CompletedAt: t.CompletedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
