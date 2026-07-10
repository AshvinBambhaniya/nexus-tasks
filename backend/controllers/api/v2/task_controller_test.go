package v2

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupTaskControllerTest() (*fiber.App, *mockTaskService, *mockCommentService, *TaskController) {
	app := fiber.New()
	mockTaskSvc := new(mockTaskService)
	mockCommentSvc := new(mockCommentService)
	logger := zap.NewNop()
	ctrl, err := NewTaskController(mockTaskSvc, mockCommentSvc, logger)
	if err != nil {
		panic(err)
	}
	return app, mockTaskSvc, mockCommentSvc, ctrl
}

func TestTaskController_CreateTask(t *testing.T) {
	uid := uuid.New()
	projectID := uuid.New()
	taskID := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		projectIDParam string
		reqBody        interface{}
		setupMocks     func(mt *mockTaskService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			projectIDParam: projectID.String(),
			reqBody:        structs.ReqCreateTask{Title: "Task 1"},
			setupMocks: func(mt *mockTaskService) {
				mt.On("CreateTask", uid, projectID, mock.Anything).Return(models.Task{ID: taskID, Title: "Task 1"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			projectIDParam: projectID.String(),
			reqBody:        structs.ReqCreateTask{Title: "Task 1"},
			setupMocks:     func(_ *mockTaskService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_project_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			projectIDParam: "invalid-uuid",
			reqBody:        structs.ReqCreateTask{Title: "Task 1"},
			setupMocks:     func(_ *mockTaskService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_request_body",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			projectIDParam: projectID.String(),
			reqBody:        "invalid json",
			setupMocks:     func(_ *mockTaskService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			projectIDParam: projectID.String(),
			reqBody:        structs.ReqCreateTask{Title: ""},
			setupMocks:     func(_ *mockTaskService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			projectIDParam: projectID.String(),
			reqBody:        structs.ReqCreateTask{Title: "Task 1"},
			setupMocks: func(mt *mockTaskService) {
				mt.On("CreateTask", uid, projectID, mock.Anything).Return(models.Task{}, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockTaskSvc, _, ctrl := setupTaskControllerTest()
			app.Post("/:"+constants.ParamProjectID, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.CreateTask(c)
			})
			tt.setupMocks(mockTaskSvc)

			body, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)
			req := httptest.NewRequest("POST", "/"+tt.projectIDParam, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTaskController_ListProjectTasks(t *testing.T) {
	uid := uuid.New()
	projectID := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		projectIDParam string
		queryParams    string
		setupMocks     func(mt *mockTaskService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			projectIDParam: projectID.String(),
			setupMocks: func(mt *mockTaskService) {
				mt.On("ListProjectTasks", uid, projectID, mock.Anything, mock.Anything).Return([]models.Task{{ID: uuid.New()}}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "success_with_filters",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			projectIDParam: projectID.String(),
			queryParams:    "?status=TODO&assignee_id=" + uuid.New().String(),
			setupMocks: func(mt *mockTaskService) {
				mt.On("ListProjectTasks", uid, projectID, mock.Anything, mock.Anything).Return([]models.Task{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			projectIDParam: projectID.String(),
			setupMocks:     func(_ *mockTaskService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_project_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			projectIDParam: "invalid-uuid",
			setupMocks:     func(_ *mockTaskService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			projectIDParam: projectID.String(),
			setupMocks: func(mt *mockTaskService) {
				mt.On("ListProjectTasks", uid, projectID, mock.Anything, mock.Anything).Return(nil, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockTaskSvc, _, ctrl := setupTaskControllerTest()
			app.Get("/:"+constants.ParamProjectID, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.ListProjectTasks(c)
			})
			tt.setupMocks(mockTaskSvc)

			req := httptest.NewRequest("GET", "/"+tt.projectIDParam+tt.queryParams, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTaskController_GetTask(t *testing.T) {
	uid := uuid.New()
	taskID := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		taskIDParam    string
		setupMocks     func(mt *mockTaskService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam: taskID.String(),
			setupMocks: func(mt *mockTaskService) {
				mt.On("GetTask", uid, taskID).Return(models.TaskWithAssignee{Task: models.Task{ID: taskID}}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			taskIDParam:    taskID.String(),
			setupMocks:     func(_ *mockTaskService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_task_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam:    "invalid-uuid",
			setupMocks:     func(_ *mockTaskService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "task_not_found",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam: taskID.String(),
			setupMocks: func(mt *mockTaskService) {
				mt.On("GetTask", uid, taskID).Return(models.TaskWithAssignee{}, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockTaskSvc, _, ctrl := setupTaskControllerTest()
			app.Get("/:"+constants.ParamTaskID, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.GetTask(c)
			})
			tt.setupMocks(mockTaskSvc)

			req := httptest.NewRequest("GET", "/"+tt.taskIDParam, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTaskController_UpdateTask(t *testing.T) {
	uid := uuid.New()
	taskID := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		taskIDParam    string
		reqBody        interface{}
		setupMocks     func(mt *mockTaskService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam: taskID.String(),
			reqBody:     structs.ReqUpdateTask{Title: "Updated"},
			setupMocks: func(mt *mockTaskService) {
				mt.On("UpdateTask", uid, taskID, mock.Anything).Return(models.Task{ID: taskID, Title: "Updated"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			taskIDParam:    taskID.String(),
			reqBody:        structs.ReqUpdateTask{Title: "Updated"},
			setupMocks:     func(_ *mockTaskService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_task_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam:    "invalid-uuid",
			reqBody:        structs.ReqUpdateTask{Title: "Updated"},
			setupMocks:     func(_ *mockTaskService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_request_body",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam:    taskID.String(),
			reqBody:        "invalid json",
			setupMocks:     func(_ *mockTaskService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam: taskID.String(),
			reqBody:     structs.ReqUpdateTask{Title: "Updated"},
			setupMocks: func(mt *mockTaskService) {
				mt.On("UpdateTask", uid, taskID, mock.Anything).Return(models.Task{}, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockTaskSvc, _, ctrl := setupTaskControllerTest()
			app.Patch("/:"+constants.ParamTaskID, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.UpdateTask(c)
			})
			tt.setupMocks(mockTaskSvc)

			body, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)
			req := httptest.NewRequest("PATCH", "/"+tt.taskIDParam, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTaskController_DeleteTask(t *testing.T) {
	uid := uuid.New()
	taskID := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		taskIDParam    string
		setupMocks     func(mt *mockTaskService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam: taskID.String(),
			setupMocks: func(mt *mockTaskService) {
				mt.On("DeleteTask", uid, taskID).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			taskIDParam:    taskID.String(),
			setupMocks:     func(_ *mockTaskService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_task_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam:    "invalid-uuid",
			setupMocks:     func(_ *mockTaskService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam: taskID.String(),
			setupMocks: func(mt *mockTaskService) {
				mt.On("DeleteTask", uid, taskID).Return(errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockTaskSvc, _, ctrl := setupTaskControllerTest()
			app.Delete("/:"+constants.ParamTaskID, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.DeleteTask(c)
			})
			tt.setupMocks(mockTaskSvc)

			req := httptest.NewRequest("DELETE", "/"+tt.taskIDParam, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTaskController_ListMyTasks(t *testing.T) {
	uid := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		setupMocks     func(mt *mockTaskService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			setupMocks: func(mt *mockTaskService) {
				mt.On("ListMyTasks", uid).Return([]models.Task{{ID: uuid.New()}}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			setupMocks:     func(_ *mockTaskService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			setupMocks: func(mt *mockTaskService) {
				mt.On("ListMyTasks", uid).Return(nil, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockTaskSvc, _, ctrl := setupTaskControllerTest()
			app.Get("/", func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.ListMyTasks(c)
			})
			tt.setupMocks(mockTaskSvc)

			req := httptest.NewRequest("GET", "/", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTaskController_CreateComment(t *testing.T) {
	uid := uuid.New()
	taskID := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		taskIDParam    string
		reqBody        interface{}
		setupMocks     func(mc *mockCommentService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam: taskID.String(),
			reqBody:     structs.ReqCreateComment{Content: "Nice task"},
			setupMocks: func(mc *mockCommentService) {
				mc.On("CreateComment", uid, taskID, mock.Anything).Return(models.Comment{ID: uuid.New(), Content: "Nice task"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			taskIDParam:    taskID.String(),
			reqBody:        structs.ReqCreateComment{Content: "Nice task"},
			setupMocks:     func(_ *mockCommentService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_task_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam:    "invalid-uuid",
			reqBody:        structs.ReqCreateComment{Content: "Nice task"},
			setupMocks:     func(_ *mockCommentService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_request_body",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam:    taskID.String(),
			reqBody:        "invalid json",
			setupMocks:     func(_ *mockCommentService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam:    taskID.String(),
			reqBody:        structs.ReqCreateComment{Content: ""},
			setupMocks:     func(_ *mockCommentService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam: taskID.String(),
			reqBody:     structs.ReqCreateComment{Content: "Nice task"},
			setupMocks: func(mc *mockCommentService) {
				mc.On("CreateComment", uid, taskID, mock.Anything).Return(models.Comment{}, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, _, mockCommentSvc, ctrl := setupTaskControllerTest()
			app.Post("/:"+constants.ParamTaskID, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.CreateComment(c)
			})
			tt.setupMocks(mockCommentSvc)

			body, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)
			req := httptest.NewRequest("POST", "/"+tt.taskIDParam, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTaskController_ListTaskComments(t *testing.T) {
	uid := uuid.New()
	taskID := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		taskIDParam    string
		setupMocks     func(mc *mockCommentService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam: taskID.String(),
			setupMocks: func(mc *mockCommentService) {
				mc.On("ListTaskComments", uid, taskID).Return([]models.CommentWithAuthor{{Comment: models.Comment{ID: uuid.New()}}}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			taskIDParam:    taskID.String(),
			setupMocks:     func(_ *mockCommentService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_task_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam:    "invalid-uuid",
			setupMocks:     func(_ *mockCommentService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			taskIDParam: taskID.String(),
			setupMocks: func(mc *mockCommentService) {
				mc.On("ListTaskComments", uid, taskID).Return(nil, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, _, mockCommentSvc, ctrl := setupTaskControllerTest()
			app.Get("/:"+constants.ParamTaskID, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.ListTaskComments(c)
			})
			tt.setupMocks(mockCommentSvc)

			req := httptest.NewRequest("GET", "/"+tt.taskIDParam, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTaskController_DeleteComment(t *testing.T) {
	uid := uuid.New()
	commentID := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		commentIDParam string
		setupMocks     func(mc *mockCommentService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			commentIDParam: commentID.String(),
			setupMocks: func(mc *mockCommentService) {
				mc.On("DeleteComment", uid, commentID).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			commentIDParam: commentID.String(),
			setupMocks:     func(_ *mockCommentService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_comment_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			commentIDParam: "invalid-uuid",
			setupMocks:     func(_ *mockCommentService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			commentIDParam: commentID.String(),
			setupMocks: func(mc *mockCommentService) {
				mc.On("DeleteComment", uid, commentID).Return(errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, _, mockCommentSvc, ctrl := setupTaskControllerTest()
			app.Delete("/:"+constants.ParamCommentID, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.DeleteComment(c)
			})
			tt.setupMocks(mockCommentSvc)

			req := httptest.NewRequest("DELETE", "/"+tt.commentIDParam, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
