package v2

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupAiControllerFullTest() (*fiber.App, *mockAiService, *mockTaskService, *mockCommentService, *AiController) {
	app := fiber.New()
	mockAi := new(mockAiService)
	mockTask := new(mockTaskService)
	mockComment := new(mockCommentService)
	logger := zap.NewNop()

	ctrl, _ := NewAiController(mockAi, mockComment, mockTask, logger)

	app.Post("/api/v2/ai/draft-task", func(c *fiber.Ctx) error {
		return ctrl.DraftTask(c)
	})
	app.Post("/api/v2/workspaces/:wsId/projects/:projectId/tasks/:taskId/ai/summarize-comments", func(c *fiber.Ctx) error {
		c.Locals("userId", c.Get("X-User-Id"))
		return ctrl.SummarizeComments(c)
	})
	app.Get("/api/v2/projects/:projectId/ai/weekly-report", func(c *fiber.Ctx) error {
		c.Locals("userId", c.Get("X-User-Id"))
		return ctrl.GenerateWeeklyReport(c)
	})

	return app, mockAi, mockTask, mockComment, ctrl
}

func TestAiController_DraftTask(t *testing.T) {
	app, ma, dummyTask, dummyComment, dummyCtrl := setupAiControllerFullTest()
	_ = dummyTask
	_ = dummyComment
	_ = dummyCtrl

	t.Run("success", func(t *testing.T) {
		ma.On("DraftTask", "Test title").Return("Draft output", nil).Once()

		reqBody, _ := json.Marshal(structs.ReqDraftTask{Title: "Test title"})
		req := httptest.NewRequest("POST", "/api/v2/ai/draft-task", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("bad json", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v2/ai/draft-task", bytes.NewBufferString("{bad}"))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("validation failure", func(t *testing.T) {
		reqBody, _ := json.Marshal(structs.ReqDraftTask{Title: ""}) // empty title fails validation
		req := httptest.NewRequest("POST", "/api/v2/ai/draft-task", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("service error", func(t *testing.T) {
		ma.On("DraftTask", "Error title").Return("", errors.New("ai error")).Once()

		reqBody, _ := json.Marshal(structs.ReqDraftTask{Title: "Error title"})
		req := httptest.NewRequest("POST", "/api/v2/ai/draft-task", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 500, resp.StatusCode)
	})
}

func TestAiController_SummarizeComments(t *testing.T) {
	app, ma, _, mc, _ := setupAiControllerFullTest()

	t.Run("success", func(t *testing.T) {
		uid := uuid.New()
		taskID := uuid.New()

		mc.On("ListTaskComments", uid, taskID).Return([]models.CommentWithAuthor{
			{Comment: models.Comment{Content: "Hey"}, AuthorFullName: "John"},
		}, nil).Once()

		ma.On("SummarizeComments", mock.Anything).Return("Summary output", nil).Once()

		req := httptest.NewRequest("POST", "/api/v2/workspaces/ws1/projects/proj1/tasks/"+taskID.String()+"/ai/summarize-comments", nil)
		req.Header.Set("X-User-Id", uid.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("invalid user context", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v2/workspaces/ws1/projects/proj1/tasks/"+uuid.New().String()+"/ai/summarize-comments", nil)
		req.Header.Set("X-User-Id", "invalid-uuid")

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 500, resp.StatusCode)
	})

	t.Run("invalid task id", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v2/workspaces/ws1/projects/proj1/tasks/invalid-uuid/ai/summarize-comments", nil)
		req.Header.Set("X-User-Id", uuid.New().String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("service db error", func(t *testing.T) {
		uid := uuid.New()
		taskID := uuid.New()

		mc.On("ListTaskComments", uid, taskID).Return([]models.CommentWithAuthor{}, errors.New("db err")).Once()

		req := httptest.NewRequest("POST", "/api/v2/workspaces/ws1/projects/proj1/tasks/"+taskID.String()+"/ai/summarize-comments", nil)
		req.Header.Set("X-User-Id", uid.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 500, resp.StatusCode)
	})

	t.Run("no comments", func(t *testing.T) {
		uid := uuid.New()
		taskID := uuid.New()

		mc.On("ListTaskComments", uid, taskID).Return([]models.CommentWithAuthor{}, nil).Once()

		req := httptest.NewRequest("POST", "/api/v2/workspaces/ws1/projects/proj1/tasks/"+taskID.String()+"/ai/summarize-comments", nil)
		req.Header.Set("X-User-Id", uid.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("ai error", func(t *testing.T) {
		uid := uuid.New()
		taskID := uuid.New()

		mc.On("ListTaskComments", uid, taskID).Return([]models.CommentWithAuthor{
			{Comment: models.Comment{Content: "Hey"}, AuthorFullName: "John"},
		}, nil).Once()

		ma.On("SummarizeComments", mock.Anything).Return("", errors.New("ai error")).Once()

		req := httptest.NewRequest("POST", "/api/v2/workspaces/ws1/projects/proj1/tasks/"+taskID.String()+"/ai/summarize-comments", nil)
		req.Header.Set("X-User-Id", uid.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 500, resp.StatusCode)
	})
}

func TestAiController_GenerateWeeklyReport(t *testing.T) {
	app, ma, mt, mc, _ := setupAiControllerFullTest()

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		projectID := uuid.New()

		mt.On("ListCompletedTasksInLastDays", userID, projectID, 7).Return([]models.TaskWithAssignee{
			{Task: models.Task{ID: uuid.New(), Title: "Task 1"}},
		}, nil).Once()

		mc.On("ListCommentsForTasks", userID, projectID, mock.Anything).Return([]models.CommentWithAuthor{}, nil).Once()

		ma.On("GenerateWeeklyReport", mock.Anything).Return("Mock AI Report", nil).Once()

		req := httptest.NewRequest("GET", "/api/v2/projects/"+projectID.String()+"/ai/weekly-report", nil)
		req.Header.Set("X-User-Id", userID.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("invalid user context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v2/projects/"+uuid.New().String()+"/ai/weekly-report", nil)
		req.Header.Set("X-User-Id", "bad-uuid")

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 500, resp.StatusCode)
	})

	t.Run("invalid project id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v2/projects/invalid-id/ai/weekly-report", nil)
		req.Header.Set("X-User-Id", uuid.New().String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("task list failure", func(t *testing.T) {
		userID := uuid.New()
		projectID := uuid.New()

		mt.On("ListCompletedTasksInLastDays", userID, projectID, 7).Return([]models.TaskWithAssignee{}, errors.New("db error")).Once()

		req := httptest.NewRequest("GET", "/api/v2/projects/"+projectID.String()+"/ai/weekly-report", nil)
		req.Header.Set("X-User-Id", userID.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 500, resp.StatusCode)
	})

	t.Run("no completed tasks", func(t *testing.T) {
		userID := uuid.New()
		projectID := uuid.New()

		mt.On("ListCompletedTasksInLastDays", userID, projectID, 7).Return([]models.TaskWithAssignee{}, nil).Once()

		req := httptest.NewRequest("GET", "/api/v2/projects/"+projectID.String()+"/ai/weekly-report", nil)
		req.Header.Set("X-User-Id", userID.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("ai error", func(t *testing.T) {
		userID := uuid.New()
		projectID := uuid.New()

		mt.On("ListCompletedTasksInLastDays", userID, projectID, 7).Return([]models.TaskWithAssignee{
			{Task: models.Task{ID: uuid.New(), Title: "Task 1"}},
		}, nil).Once()

		mc.On("ListCommentsForTasks", userID, projectID, mock.Anything).Return([]models.CommentWithAuthor{}, nil).Once()

		ma.On("GenerateWeeklyReport", mock.Anything).Return("", errors.New("ai err")).Once()

		req := httptest.NewRequest("GET", "/api/v2/projects/"+projectID.String()+"/ai/weekly-report", nil)
		req.Header.Set("X-User-Id", userID.String())

		resp, _ := app.Test(req)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, 500, resp.StatusCode)
	})
}
