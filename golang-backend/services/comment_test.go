package services

import (
	"errors"
	"testing"

	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupCommentTest(t *testing.T) (*commentService, *mockCommentRepository, *mockTaskRepository, *mockProjectService, *mockStorage) {
	mockCommentRepo := new(mockCommentRepository)
	mockTaskRepo := new(mockTaskRepository)
	mockProjSvc := new(mockProjectService)
	mockStor := new(mockStorage)
	logger := zap.NewNop()

	mockStor.On("Comments").Return(mockCommentRepo)
	mockStor.On("Tasks").Return(mockTaskRepo)

	svc := NewCommentService(mockStor, mockProjSvc, logger).(*commentService)

	return svc, mockCommentRepo, mockTaskRepo, mockProjSvc, mockStor
}

func TestCommentService_CreateComment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mc, mt, mps, _ := setupCommentTest(t)
		userID := uuid.New()
		taskID := uuid.New()
		projectID := uuid.New()

		mt.On("GetByID", taskID).Return(models.Task{ProjectID: projectID}, nil)
		mps.On("ValidateProjectAccess", projectID, userID, false).Return(nil)
		mc.On("Create", mock.Anything).Return(models.Comment{Content: "C1"}, nil)

		res, err := svc.CreateComment(userID, taskID, structs.ReqCreateComment{Content: "C1"})
		assert.NoError(t, err)
		assert.Equal(t, "C1", res.Content)
	})

	t.Run("task not found", func(t *testing.T) {
		svc, _, mt, _, _ := setupCommentTest(t)
		mt.On("GetByID", mock.Anything).Return(models.Task{}, errors.New("not found"))

		_, err := svc.CreateComment(uuid.New(), uuid.New(), structs.ReqCreateComment{Content: "C1"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task not found")
	})

	t.Run("access denied", func(t *testing.T) {
		svc, _, mt, mps, _ := setupCommentTest(t)
		taskID := uuid.New()
		projectID := uuid.New()

		mt.On("GetByID", taskID).Return(models.Task{ProjectID: projectID}, nil)
		mps.On("ValidateProjectAccess", projectID, mock.Anything, false).Return(errors.New("unauthorized"))

		_, err := svc.CreateComment(uuid.New(), taskID, structs.ReqCreateComment{Content: "C1"})
		assert.Error(t, err)
	})

	t.Run("atomic failure", func(t *testing.T) {
		svc, mc, mt, mps, _ := setupCommentTest(t)
		userID := uuid.New()
		taskID := uuid.New()
		projectID := uuid.New()

		mt.On("GetByID", taskID).Return(models.Task{ProjectID: projectID}, nil)
		mps.On("ValidateProjectAccess", projectID, userID, false).Return(nil)
		mc.On("Create", mock.Anything).Return(models.Comment{}, errors.New("db error"))

		_, err := svc.CreateComment(userID, taskID, structs.ReqCreateComment{Content: "C1"})
		assert.Error(t, err)
	})

}

func TestCommentService_ListTaskComments(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mc, mt, mps, _ := setupCommentTest(t)
		userID := uuid.New()
		taskID := uuid.New()
		projectID := uuid.New()

		mt.On("GetByID", taskID).Return(models.Task{ID: taskID, ProjectID: projectID}, nil)
		mps.On("ValidateProjectAccess", projectID, userID, false).Return(nil)
		mc.On("ListByTaskID", taskID).Return([]models.CommentWithAuthor{{Comment: models.Comment{Content: "C1"}}}, nil)

		res, err := svc.ListTaskComments(userID, taskID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "C1", res[0].Content)
	})

	t.Run("task not found", func(t *testing.T) {
		svc, _, mt, _, _ := setupCommentTest(t)
		mt.On("GetByID", mock.Anything).Return(models.Task{}, errors.New("not found"))

		_, err := svc.ListTaskComments(uuid.New(), uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task not found")
	})

	t.Run("access denied", func(t *testing.T) {
		svc, _, mt, mps, _ := setupCommentTest(t)
		taskID := uuid.New()
		projectID := uuid.New()

		mt.On("GetByID", taskID).Return(models.Task{ProjectID: projectID}, nil)
		mps.On("ValidateProjectAccess", projectID, mock.Anything, false).Return(errors.New("unauthorized"))

		_, err := svc.ListTaskComments(uuid.New(), taskID)
		assert.Error(t, err)
	})

}

func TestCommentService_DeleteComment(t *testing.T) {
	t.Run("success as author", func(t *testing.T) {
		svc, mc, _, _, _ := setupCommentTest(t)
		userID := uuid.New()
		commentID := uuid.New()

		mc.On("GetByID", commentID).Return(models.Comment{AuthorID: userID}, nil)
		mc.On("Delete", commentID).Return(nil)

		err := svc.DeleteComment(userID, commentID)
		assert.NoError(t, err)
	})

	t.Run("fail - not author", func(t *testing.T) {
		svc, mc, _, _, _ := setupCommentTest(t)
		commentID := uuid.New()

		mc.On("GetByID", commentID).Return(models.Comment{AuthorID: uuid.New()}, nil)

		err := svc.DeleteComment(uuid.New(), commentID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only author")
	})

	t.Run("comment not found", func(t *testing.T) {
		svc, mc, _, _, _ := setupCommentTest(t)
		mc.On("GetByID", mock.Anything).Return(models.Comment{}, errors.New("not found"))

		err := svc.DeleteComment(uuid.New(), uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "comment not found")
	})
}
