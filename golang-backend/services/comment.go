package services

import (
	"errors"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CommentService struct {
	commentModel   *models.CommentModel
	taskModel      *models.TaskModel
	projectService *ProjectService
	db             *goqu.Database
	logger         *zap.Logger
}

func NewCommentService(db *goqu.Database, logger *zap.Logger, commentModel *models.CommentModel, taskModel *models.TaskModel, projectService *ProjectService) *CommentService {
	return &CommentService{
		commentModel:   commentModel,
		taskModel:      taskModel,
		projectService: projectService,
		db:             db,
		logger:         logger,
	}
}

func (s *CommentService) CreateComment(userID, taskID uuid.UUID, req structs.ReqCreateComment) (models.Comment, error) {
	// 1. Verify Task Access
	task, err := s.taskModel.GetByID(taskID)
	if err != nil {
		return models.Comment{}, errors.New("task not found")
	}

	// Validate Access
	err = s.projectService.ValidateProjectAccess(task.ProjectID, userID, false)
	if err != nil {
		return models.Comment{}, err
	}

	comment := models.Comment{
		Content:   req.Content,
		TaskID:    taskID,
		AuthorID:  userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	isOk := false
	transaction, err := s.db.Begin()
	if err != nil {
		return models.Comment{}, err
	}
	defer func() {
		if isOk {
			err := transaction.Commit()
			if err != nil {
				s.logger.Error("error during commit in create comment", zap.Error(err))
			}
		} else {
			err := transaction.Rollback()
			if err != nil {
				s.logger.Error("error during rollback in create comment", zap.Error(err))
			}
		}
	}()

	createdComment, err := s.commentModel.Create(transaction, comment)
	if err != nil {
		return models.Comment{}, err
	}

	isOk = true
	return createdComment, nil
}

func (s *CommentService) ListTaskComments(userID, taskID uuid.UUID) ([]models.CommentWithAuthor, error) {
	task, err := s.taskModel.GetByID(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}

	err = s.projectService.ValidateProjectAccess(task.ProjectID, userID, false)
	if err != nil {
		return nil, err
	}

	return s.commentModel.ListByTaskID(taskID)
}

func (s *CommentService) DeleteComment(userID, commentID uuid.UUID) error {
	comment, err := s.commentModel.GetByID(commentID)
	if err != nil {
		return errors.New("comment not found")
	}

	// Check Authorization: Only Author can delete (per Python logic)
	if comment.AuthorID != userID {
		return errors.New("unauthorized: only author can delete comment")
	}

	isOk := false
	transaction, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if isOk {
			err := transaction.Commit()
			if err != nil {
				s.logger.Error("error during commit in delete comment", zap.Error(err))
			}
		} else {
			err := transaction.Rollback()
			if err != nil {
				s.logger.Error("error during rollback in delete comment", zap.Error(err))
			}
		}
	}()

	err = s.commentModel.Delete(transaction, commentID)
	if err != nil {
		return err
	}

	isOk = true
	return nil
}
