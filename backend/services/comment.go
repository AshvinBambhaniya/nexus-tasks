// Package services provides business logic for the application.
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/structs"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CommentService defines the interface for comment-related business logic
type CommentService interface {
	CreateComment(userID, taskID uuid.UUID, req structs.ReqCreateComment) (models.Comment, error)
	ListTaskComments(userID, taskID uuid.UUID) ([]models.CommentWithAuthor, error)
	ListCommentsForTasks(userID, projectID uuid.UUID, taskIDs []uuid.UUID) ([]models.CommentWithAuthor, error)
	DeleteComment(userID, commentID uuid.UUID) error
}

type commentService struct {
	storage        models.Storage
	projectService ProjectService
	logger         *zap.Logger
}

// NewCommentService creates a new comment service instance
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
		if err != nil {
			return err
		}

		mentionedMap := make(map[uuid.UUID]bool)
		for _, mentionedID := range req.MentionedUserIDs {
			if mentionedID == userID {
				continue // Don't notify the user who made the comment
			}
			// Validate if the mentioned user is part of the project
			if err := s.projectService.ValidateProjectAccess(task.ProjectID, mentionedID, false); err == nil {
				mentionedMap[mentionedID] = true
				body := "You were mentioned in a comment."
				_ = txStorage.Notifications().Create(&models.Notification{
					UserID:     mentionedID,
					ActorID:    userID,
					EntityID:   taskID,
					EntityType: models.EntityTypeTask,
					Type:       models.NotificationTypeMentioned,
					Title:      fmt.Sprintf("Mentioned on %s", task.Title),
					Body:       &body,
				})
			}
		}

		// Notify Author if not mentioned
		if task.AuthorID != nil && *task.AuthorID != userID && !mentionedMap[*task.AuthorID] {
			body := "A new comment was added"
			_ = txStorage.Notifications().Create(&models.Notification{
				UserID:     *task.AuthorID,
				ActorID:    userID,
				EntityID:   taskID,
				EntityType: models.EntityTypeTask,
				Type:       models.NotificationTypeCommentAdded,
				Title:      fmt.Sprintf("New comment on %s", task.Title),
				Body:       &body,
			})
		}

		// Notify Assignee if not mentioned
		if task.AssigneeID != nil && *task.AssigneeID != userID && (task.AuthorID == nil || *task.AssigneeID != *task.AuthorID) && !mentionedMap[*task.AssigneeID] {
			body := "A new comment was added"
			_ = txStorage.Notifications().Create(&models.Notification{
				UserID:     *task.AssigneeID,
				ActorID:    userID,
				EntityID:   taskID,
				EntityType: models.EntityTypeTask,
				Type:       models.NotificationTypeCommentAdded,
				Title:      fmt.Sprintf("New comment on %s", task.Title),
				Body:       &body,
			})
		}

		return nil
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

func (s *commentService) ListCommentsForTasks(userID, projectID uuid.UUID, taskIDs []uuid.UUID) ([]models.CommentWithAuthor, error) {
	err := s.projectService.ValidateProjectAccess(projectID, userID, false)
	if err != nil {
		return nil, err
	}

	return s.storage.Comments().ListByTaskIDs(taskIDs)
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
