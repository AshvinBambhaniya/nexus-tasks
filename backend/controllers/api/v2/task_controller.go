package v2

import (
	"net/http"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/structs"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/services"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/utils"
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
//
/*
 swagger:operation POST /workspaces/{workspaceId}/projects/{projectId}/tasks tasks createTask

 # Create a task

 Creates a new task within the specified project. The authenticated user becomes the task author.

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
   - name: projectId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the project.
   - name: body
     in: body
     required: true
     schema:
        "$ref": "#/definitions/ReqCreateTask"

 responses:

	201:
	  description: Task created successfully.
	  schema:
      "$ref": "#/definitions/ResTask"
	400:
	  description: Validation error or invalid path parameter.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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

	return utils.JSONSuccess(c, http.StatusCreated, ctrl.mapTaskToRes(models.TaskWithAssignee{Task: task}))
}

// ListProjectTasks handles listing all tasks in a project.
//
/*
 swagger:operation GET /workspaces/{workspaceId}/projects/{projectId}/tasks tasks listProjectTasks

 # List tasks in a project

 Returns all tasks belonging to the specified project, with optional filtering by status or assignee.

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
   - name: projectId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the project.
   - name: status
     in: query
     type: string
     description: "Filter tasks by status (e.g. TODO, IN_PROGRESS, DONE)."
   - name: assignee_id
     in: query
     type: string
     format: uuid
     description: Filter tasks by assignee UUID.

 responses:

	200:
	  description: List of tasks.
	  schema:
	    type: array
	    items:
	      "$ref": "#/definitions/ResTask"
	400:
	  description: Invalid path parameter.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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
//
/*
 swagger:operation GET /workspaces/{workspaceId}/projects/{projectId}/tasks/{taskId} tasks getTask

 # Get a task by ID

 Returns the full details of a single task, including its assignee and author.

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
   - name: projectId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the project.
   - name: taskId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the task.

 responses:

	200:
	  description: Task details.
	  schema:
      "$ref": "#/definitions/ResTask"
	400:
	  description: Invalid task ID.
	401:
	  description: Not authenticated.
	404:
	  description: Task not found.
	500:
	  description: Internal server error.
*/
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
//
/*
 swagger:operation PATCH /workspaces/{workspaceId}/projects/{projectId}/tasks/{taskId} tasks updateTask

 # Update a task

 Partially updates the specified task. Only non-zero fields in the request body are applied.

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
   - name: projectId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the project.
   - name: taskId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the task.
   - name: body
     in: body
     required: true
     schema:
        "$ref": "#/definitions/ReqUpdateTask"

 responses:

	200:
	  description: Task updated.
	  schema:
      "$ref": "#/definitions/ResTask"
	400:
	  description: Validation error or invalid task ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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

	return utils.JSONSuccess(c, http.StatusOK, ctrl.mapTaskToRes(models.TaskWithAssignee{Task: task}))
}

// DeleteTask handles deleting a task.
//
/*
 swagger:operation DELETE /workspaces/{workspaceId}/projects/{projectId}/tasks/{taskId} tasks deleteTask

 # Delete a task

 Permanently removes a task and all its associated comments.

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
   - name: projectId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the project.
   - name: taskId
     in: path
     required: true
     type: string
     format: uuid
     description: UUID of the task to delete.

 responses:

	200:
	  description: Task deleted.
	400:
	  description: Invalid task ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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
//
/*
 swagger:operation GET /tasks/me tasks listMyTasks

 # List tasks assigned to me

 Returns all tasks currently assigned to the authenticated user across all projects.

 ---
 produces:
 - application/json
 security:
 - cookieAuth: []
 - apiKeyAuth: []
 responses:

	200:
	  description: List of tasks assigned to the current user.
	  schema:
	    type: array
	    items:
	      "$ref": "#/definitions/ResTask"
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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
//
/*
 swagger:operation POST /workspaces/{workspaceId}/projects/{projectId}/tasks/{taskId}/comments comments createComment

 # Post a comment on a task

 Adds a new comment to the specified task. Mentioned users will receive notifications.

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
   - name: projectId
     in: path
     required: true
     type: string
     format: uuid
   - name: taskId
     in: path
     required: true
     type: string
     format: uuid
   - name: body
     in: body
     required: true
     schema:
        "$ref": "#/definitions/ReqCreateComment"

 responses:

	201:
	  description: Comment created.
	  schema:
      "$ref": "#/definitions/ResComment"
	400:
	  description: Validation error or invalid task ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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
//
/*
 swagger:operation GET /workspaces/{workspaceId}/projects/{projectId}/tasks/{taskId}/comments comments listTaskComments

 # List comments on a task

 Returns all comments on the specified task, ordered by creation time ascending.

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
   - name: projectId
     in: path
     required: true
     type: string
     format: uuid
   - name: taskId
     in: path
     required: true
     type: string
     format: uuid

 responses:

	200:
	  description: List of comments.
	  schema:
	    type: array
	    items:
	      "$ref": "#/definitions/ResComment"
	400:
	  description: Invalid task ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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
//
/*
 swagger:operation DELETE /workspaces/{workspaceId}/projects/{projectId}/tasks/{taskId}/comments/{commentId} comments deleteComment

 # Delete a comment

 Permanently removes a comment. Only the comment author may delete their own comment.

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
   - name: projectId
     in: path
     required: true
     type: string
     format: uuid
   - name: taskId
     in: path
     required: true
     type: string
     format: uuid
   - name: commentId
     in: path
     required: true
     type: string
     format: uuid

 responses:

	200:
	  description: Comment deleted.
	400:
	  description: Invalid comment ID.
	401:
	  description: Not authenticated.
	500:
	  description: Internal server error.
*/
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

func (ctrl *TaskController) mapTaskToRes(t models.TaskWithAssignee) structs.ResTask {
	res := structs.ResTask{
		ID:               t.ID,
		Number:           t.Number,
		Title:            t.Title,
		Description:      t.Description,
		Status:           t.Status,
		Priority:         t.Priority,
		ProjectID:        t.ProjectID,
		AssigneeID:       t.AssigneeID,
		AuthorID:         t.AuthorID,
		DueDate:          t.DueDate,
		EstimatedMinutes: t.EstimatedMinutes,
		CompletedAt:      t.CompletedAt,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}

	if t.AssigneeID != nil && t.AssigneeEmail != nil {
		fullName := ""
		if t.AssigneeFullName != nil {
			fullName = *t.AssigneeFullName
		}
		res.Assignee = &structs.ResUser{
			ID:       *t.AssigneeID,
			Email:    *t.AssigneeEmail,
			FullName: fullName,
		}
	}

	if t.AuthorID != nil && t.AuthorEmail != nil {
		fullName := ""
		if t.AuthorFullName != nil {
			fullName = *t.AuthorFullName
		}
		res.Author = &structs.ResUser{
			ID:       *t.AuthorID,
			Email:    *t.AuthorEmail,
			FullName: fullName,
		}
	}

	return res
}
