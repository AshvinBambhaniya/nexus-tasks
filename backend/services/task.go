package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/realtime"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/structs"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TaskService defines the interface for task-related business logic
type TaskService interface {
	CreateTask(userID, projectID uuid.UUID, req structs.ReqCreateTask) (models.Task, error)
	ListProjectTasks(userID, projectID uuid.UUID, status *models.TaskStatus, assigneeID *uuid.UUID) ([]models.TaskWithAssignee, error)
	GetTask(userID, taskID uuid.UUID) (models.TaskWithAssignee, error)
	UpdateTask(userID, taskID uuid.UUID, req structs.ReqUpdateTask) (models.Task, error)
	DeleteTask(userID, taskID uuid.UUID) error
	ListMyTasks(userID uuid.UUID) ([]models.TaskWithAssignee, error)
	ListCompletedTasksInLastDays(userID, projectID uuid.UUID, days int) ([]models.TaskWithAssignee, error)
}

type taskService struct {
	storage models.Storage
	hub     realtime.IHub
	logger  *zap.Logger
}

// NewTaskService creates a new task service instance
func NewTaskService(storage models.Storage, logger *zap.Logger, hub realtime.IHub) TaskService {
	return &taskService{
		storage: storage,
		logger:  logger,
		hub:     hub,
	}
}

func (s *taskService) CreateTask(userID, projectID uuid.UUID, req structs.ReqCreateTask) (models.Task, error) {
	// 1. Verify Project Access
	err := s.validateProjectAccess(projectID, userID)
	if err != nil {
		return models.Task{}, err
	}

	// 2. Validate Assignee
	if req.AssigneeID != nil {
		err := s.validateAssignee(projectID, *req.AssigneeID)
		if err != nil {
			return models.Task{}, err
		}
	}

	var createdTask models.Task
	err = s.storage.Atomic(context.Background(), func(txStorage models.Storage) error {
		status := req.Status
		if status == "" {
			status = models.TaskStatusTodo
		}
		priority := req.Priority
		if priority == "" {
			priority = models.TaskPriorityP2
		}

		var dueDate *time.Time
		if req.DueDate != nil {
			dueDate = &req.DueDate.Time
		}

		nextNumber, nextNumbererr := txStorage.Tasks().GetNextTaskNumber(projectID)
		if nextNumbererr != nil {
			return nextNumbererr
		}

		task := models.Task{
			Number:      nextNumber,
			Title:       req.Title,
			Description: req.Description,
			Status:      status,
			Priority:    priority,
			ProjectID:   projectID,
			AssigneeID:  req.AssigneeID,
			AuthorID:    &userID,
			DueDate:     dueDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		var err error
		createdTask, err = txStorage.Tasks().Create(task)
		if err != nil {
			return err
		}

		if req.AssigneeID != nil && *req.AssigneeID != userID {
			s.logger.Info("Creating notification for new task assignment", zap.String("assigneeID", req.AssigneeID.String()), zap.String("taskID", createdTask.ID.String()))
			body := "You have been assigned to a new task."
			_ = txStorage.Notifications().Create(&models.Notification{
				UserID:     *req.AssigneeID,
				ActorID:    userID,
				EntityID:   createdTask.ID,
				EntityType: models.EntityTypeTask,
				Type:       models.NotificationTypeAssigned,
				Title:      fmt.Sprintf("You were assigned to %s", createdTask.Title),
				Body:       &body,
			})
		}

		return nil
	})

	if err != nil {
		return models.Task{}, err
	}

	// Broadcast Event (After commit)
	if s.hub != nil {
		s.hub.Broadcast(fmt.Sprintf("project:%s", projectID.String()), structs.RealtimeEvent{
			Type:    structs.EventTypeTaskCreated,
			Payload: createdTask,
		})
	}

	return createdTask, nil
}

func (s *taskService) ListProjectTasks(userID, projectID uuid.UUID, status *models.TaskStatus, assigneeID *uuid.UUID) ([]models.TaskWithAssignee, error) {
	err := s.validateProjectAccess(projectID, userID)
	if err != nil {
		return nil, err
	}

	return s.storage.Tasks().ListByProjectID(projectID, status, assigneeID)
}

func (s *taskService) GetTask(userID, taskID uuid.UUID) (models.TaskWithAssignee, error) {
	task, err := s.storage.Tasks().GetByID(taskID)
	if err != nil {
		return models.TaskWithAssignee{}, err
	}

	err = s.validateProjectAccess(task.ProjectID, userID)
	if err != nil {
		return models.TaskWithAssignee{}, err
	}

	return task, nil
}

func (s *taskService) UpdateTask(userID, taskID uuid.UUID, req structs.ReqUpdateTask) (models.Task, error) {
	var updatedTask models.Task
	err := s.storage.Atomic(context.Background(), func(txStorage models.Storage) error {
		task, err := txStorage.Tasks().GetByID(taskID)
		if err != nil {
			return err
		}

		if err := s.internalValidateProjectAccess(txStorage, task.ProjectID, userID); err != nil {
			return err
		}

		s.applyBasicUpdates(&task.Task, req)

		statusChanged := req.Status != "" && req.Status != task.Status
		assigneeChanged := req.AssigneeID != nil && (task.AssigneeID == nil || *req.AssigneeID != *task.AssigneeID)

		if assigneeChanged {
			if err := s.internalValidateAssignee(txStorage, task.ProjectID, *req.AssigneeID); err != nil {
				return err
			}
			task.AssigneeID = req.AssigneeID
		}

		s.handleStatusUpdate(&task.Task, req.Status)

		updatedTask, err = txStorage.Tasks().Update(task.Task)
		if err != nil {
			return err
		}

		s.dispatchUpdateNotifications(txStorage, userID, task, updatedTask, assigneeChanged, statusChanged, req)
		return nil
	})

	if err != nil {
		return models.Task{}, err
	}

	if s.hub != nil {
		s.hub.Broadcast(fmt.Sprintf("project:%s", updatedTask.ProjectID.String()), structs.RealtimeEvent{
			Type:    structs.EventTypeTaskUpdated,
			Payload: updatedTask,
		})
	}

	return updatedTask, nil
}

func (s *taskService) dispatchUpdateNotifications(txStorage models.Storage, userID uuid.UUID, task models.TaskWithAssignee, updatedTask models.Task, assigneeChanged, statusChanged bool, req structs.ReqUpdateTask) {
	if assigneeChanged && *req.AssigneeID != userID {
		body := "You have been assigned to a task."
		_ = txStorage.Notifications().Create(&models.Notification{
			UserID:     *req.AssigneeID,
			ActorID:    userID,
			EntityID:   updatedTask.ID,
			EntityType: models.EntityTypeTask,
			Type:       models.NotificationTypeAssigned,
			Title:      fmt.Sprintf("You were assigned to %s", updatedTask.Title),
			Body:       &body,
		})
	}

	if statusChanged && req.Status == models.TaskStatusDone && task.AuthorID != nil && *task.AuthorID != userID {
		body := "Task moved to DONE"
		_ = txStorage.Notifications().Create(&models.Notification{
			UserID:     *task.AuthorID,
			ActorID:    userID,
			EntityID:   updatedTask.ID,
			EntityType: models.EntityTypeTask,
			Type:       models.NotificationTypeStatusChanged,
			Title:      fmt.Sprintf("%s was marked as Done", updatedTask.Title),
			Body:       &body,
		})
	}
}

func (s *taskService) applyBasicUpdates(task *models.Task, req structs.ReqUpdateTask) {
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = &req.DueDate.Time
	}
}

func (s *taskService) handleStatusUpdate(task *models.Task, newStatus models.TaskStatus) {
	if newStatus == "" || newStatus == task.Status {
		return
	}

	if newStatus == models.TaskStatusDone && task.Status != models.TaskStatusDone {
		now := time.Now()
		task.CompletedAt = &now
	} else if newStatus != models.TaskStatusDone && task.Status == models.TaskStatusDone {
		task.CompletedAt = nil
	}
	task.Status = newStatus
}

func (s *taskService) DeleteTask(userID, taskID uuid.UUID) error {
	var projectID uuid.UUID
	err := s.storage.Atomic(context.Background(), func(txStorage models.Storage) error {
		task, err := txStorage.Tasks().GetByID(taskID)
		if err != nil {
			return err
		}
		projectID = task.ProjectID

		err = s.internalValidateProjectAccess(txStorage, task.ProjectID, userID)
		if err != nil {
			return err
		}

		return txStorage.Tasks().Delete(taskID)
	})

	if err != nil {
		return err
	}

	// Broadcast Event
	if s.hub != nil {
		s.hub.Broadcast(fmt.Sprintf("project:%s", projectID.String()), structs.RealtimeEvent{
			Type:    structs.EventTypeTaskDeleted,
			Payload: taskID,
		})
	}

	return nil
}

func (s *taskService) ListMyTasks(userID uuid.UUID) ([]models.TaskWithAssignee, error) {
	return s.storage.Tasks().ListByAssigneeID(userID)
}

func (s *taskService) ListCompletedTasksInLastDays(userID, projectID uuid.UUID, days int) ([]models.TaskWithAssignee, error) {
	err := s.validateProjectAccess(projectID, userID)
	if err != nil {
		return nil, err
	}

	return s.storage.Tasks().ListCompletedTasksInLastDays(projectID, days)
}

// Helpers

func (s *taskService) validateProjectAccess(projectID, userID uuid.UUID) error {
	return s.internalValidateProjectAccess(s.storage, projectID, userID)
}

func (s *taskService) internalValidateProjectAccess(st models.Storage, projectID, userID uuid.UUID) error {
	// 1. Direct Member
	_, err := st.Projects().GetMember(projectID, userID)
	if err == nil {
		return nil
	}
	// 2. Workspace Admin
	project, err := st.Projects().GetByID(projectID)
	if err != nil {
		return errors.New("project not found")
	}

	wsMember, err := st.Workspaces().GetMember(project.WorkspaceID, userID)
	if err == nil && wsMember.Role == models.WorkspaceRoleAdmin {
		return nil
	}

	// 3. Team Member
	projectTeams, err := st.Projects().GetTeams(projectID)
	if err == nil {
		for _, pt := range projectTeams {
			_, err := st.Teams().GetMember(pt.TeamID, userID)
			if err == nil {
				return nil
			}
		}
	}

	return errors.New("unauthorized: access denied to project")
}

func (s *taskService) validateAssignee(projectID, assigneeID uuid.UUID) error {
	return s.internalValidateAssignee(s.storage, projectID, assigneeID)
}

func (s *taskService) internalValidateAssignee(st models.Storage, projectID, assigneeID uuid.UUID) error {
	// Check Direct Membership
	_, err := st.Projects().GetMember(projectID, assigneeID)
	if err == nil {
		return nil
	}

	// Check Team Membership
	projectTeams, err := st.Projects().GetTeams(projectID)
	if err == nil {
		for _, pt := range projectTeams {
			_, err := st.Teams().GetMember(pt.TeamID, assigneeID)
			if err == nil {
				return nil
			}
		}
	}

	return errors.New("assignee must be a member of the project")
}
