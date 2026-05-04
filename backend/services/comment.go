package services

import (
	"context"
	"errors"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CommentService interface {
	CreateComment(userID, taskID uuid.UUID, req structs.ReqCreateComment) (models.Comment, error)
	ListTaskComments(userID, taskID uuid.UUID) ([]models.CommentWithAuthor, error)
	DeleteComment(userID, commentID uuid.UUID) error
}

type commentService struct {
	storage        models.Storage
	projectService ProjectService
	logger         *zap.Logger
}

func NewCommentService(storage models.Storage, projectService ProjectService, logger *zap.Logger) CommentService {
	return &commentService{
		storage:        storage,
		projectService: projectService,
		logger:         logger,
	}
}

func (s *commentService) CreateComment(userID, taskID uuid.UUID, req structs.ReqCreateComment) (models.Comment, error) {
	// 1. Verify Task Access
	task, err := s.storage.Tasks().GetByID(taskID)
	if err != nil {
		return models.Comment{}, errors.New("task not found")
	}

	// Validate Access
	err = s.projectService.ValidateProjectAccess(task.ProjectID, userID, false)
	if err != nil {
		return models.Comment{}, err
	}

	var createdComment models.Comment
	err = s.storage.Atomic(context.Background(), func(txStorage models.Storage) error {
		comment := models.Comment{
			Content:   req.Content,
			TaskID:    taskID,
			AuthorID:  userID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		var err error
		createdComment, err = txStorage.Comments().Create(comment)
		return err
	})

	if err != nil {
		return models.Comment{}, err
	}

	return createdComment, nil
}

func (s *commentService) ListTaskComments(userID, taskID uuid.UUID) ([]models.CommentWithAuthor, error) {
	task, err := s.storage.Tasks().GetByID(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}

	err = s.projectService.ValidateProjectAccess(task.ProjectID, userID, false)
	if err != nil {
		return nil, err
	}

	return s.storage.Comments().ListByTaskID(taskID)
}

func (s *commentService) DeleteComment(userID, commentID uuid.UUID) error {
	comment, err := s.storage.Comments().GetByID(commentID)
	if err != nil {
		return errors.New("comment not found")
	}

	// Check Authorization: Only Author can delete (per Python logic)
	if comment.AuthorID != userID {
		return errors.New("unauthorized: only author can delete comment")
	}

	return s.storage.Atomic(context.Background(), func(txStorage models.Storage) error {
		return txStorage.Comments().Delete(commentID)
	})
}
