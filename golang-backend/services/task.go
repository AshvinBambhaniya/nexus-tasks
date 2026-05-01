package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/realtime"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TaskService struct {
	taskModel      *models.TaskModel
	projectModel   *models.ProjectModel
	workspaceModel *models.WorkspaceModel
	teamModel      *models.TeamModel
	userModel      *models.UserModel
	hub            *realtime.Hub
	db             *goqu.Database
	logger         *zap.Logger
}

func NewTaskService(db *goqu.Database, logger *zap.Logger, taskModel *models.TaskModel, projectModel *models.ProjectModel, workspaceModel *models.WorkspaceModel, teamModel *models.TeamModel, userModel *models.UserModel, hub *realtime.Hub) *TaskService {
	return &TaskService{
		taskModel:      taskModel,
		projectModel:   projectModel,
		workspaceModel: workspaceModel,
		teamModel:      teamModel,
		userModel:      userModel,
		hub:            hub,
		db:             db,
		logger:         logger,
	}
}

func (s *TaskService) CreateTask(userID, projectID uuid.UUID, req structs.ReqCreateTask) (models.Task, error) {
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

	isOk := false
	transaction, err := s.db.Begin()
	if err != nil {
		return models.Task{}, err
	}
	defer func() {
		if isOk {
			transaction.Commit()
		} else {
			transaction.Rollback()
		}
	}()

	// 3. Get Next Task Number
	nextNumber, err := s.taskModel.GetNextTaskNumber(transaction, projectID)
	if err != nil {
		return models.Task{}, err
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

	createdTask, err := s.taskModel.Create(transaction, task)
	if err != nil {
		return models.Task{}, err
	}

	isOk = true

	// Broadcast Event
	if s.hub != nil {
		// Project Channel
		// Payload should match frontend expectation (TaskResponse)
		// Assuming we can serialize struct or map
		s.hub.Broadcast(fmt.Sprintf("project:%d", projectID), map[string]interface{}{
			"type": "TASK_CREATED",
			"task": createdTask,
		})
	}

	return createdTask, nil
}

func (s *TaskService) ListProjectTasks(userID, projectID uuid.UUID, status *models.TaskStatus, assigneeID *uuid.UUID) ([]models.Task, error) {
	err := s.validateProjectAccess(projectID, userID)
	if err != nil {
		return nil, err
	}

	return s.taskModel.ListByProjectID(projectID, status, assigneeID)
}

func (s *TaskService) GetTask(userID, taskID uuid.UUID) (models.Task, error) {
	task, err := s.taskModel.GetByID(taskID)
	if err != nil {
		return models.Task{}, err
	}

	err = s.validateProjectAccess(task.ProjectID, userID)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (s *TaskService) UpdateTask(userID, taskID uuid.UUID, req structs.ReqUpdateTask) (models.Task, error) {
	task, err := s.taskModel.GetByID(taskID)
	if err != nil {
		return models.Task{}, err
	}

	err = s.validateProjectAccess(task.ProjectID, userID)
	if err != nil {
		return models.Task{}, err
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.AssigneeID != nil {
		err := s.validateAssignee(task.ProjectID, *req.AssigneeID)
		if err != nil {
			return models.Task{}, err
		}
		task.AssigneeID = req.AssigneeID
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = &req.DueDate.Time
	}

	// Status Logic
	if req.Status != "" {
		if req.Status == models.TaskStatusDone && task.Status != models.TaskStatusDone {
			now := time.Now()
			task.CompletedAt = &now
		} else if req.Status != models.TaskStatusDone && task.Status == models.TaskStatusDone {
			task.CompletedAt = nil
		}
		task.Status = req.Status
	}

	isOk := false
	transaction, err := s.db.Begin()
	if err != nil {
		return models.Task{}, err
	}
	defer func() {
		if isOk {
			transaction.Commit()
		} else {
			transaction.Rollback()
		}
	}()

	updatedTask, err := s.taskModel.Update(transaction, task)
	if err != nil {
		return models.Task{}, err
	}

	isOk = true

	// Broadcast Event
	if s.hub != nil {
		s.hub.Broadcast(fmt.Sprintf("project:%d", task.ProjectID), map[string]interface{}{
			"type": "TASK_UPDATED",
			"task": updatedTask,
		})
	}

	return updatedTask, nil
}

func (s *TaskService) DeleteTask(userID, taskID uuid.UUID) error {
	task, err := s.taskModel.GetByID(taskID)
	if err != nil {
		return err
	}

	err = s.validateProjectAccess(task.ProjectID, userID)
	if err != nil {
		return err
	}

	isOk := false
	transaction, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if isOk {
			transaction.Commit()
		} else {
			transaction.Rollback()
		}
	}()

	err = s.taskModel.Delete(transaction, taskID)
	if err != nil {
		return err
	}

	isOk = true

	// Broadcast Event
	if s.hub != nil {
		s.hub.Broadcast(fmt.Sprintf("project:%d", task.ProjectID), map[string]interface{}{
			"type":    "TASK_DELETED",
			"task_id": taskID,
		})
	}

	return nil
}

func (s *TaskService) ListMyTasks(userID uuid.UUID) ([]models.Task, error) {
	// Requirement: Get all tasks assigned to current user across all projects
	// In Python: filters by assignee_id
	return s.taskModel.ListByAssigneeID(userID)
}

// Helpers

func (s *TaskService) validateProjectAccess(projectID, userID uuid.UUID) error {
	// Reusing logic similar to ProjectService but simplified or imported
	// Ideally ProjectService exposes `ValidateAccess` public method.
	// Duplicating for now to keep services decoupled or refactor later.

	// 1. Direct Member
	_, err := s.projectModel.GetMember(projectID, userID)
	if err == nil {
		return nil
	}
	// 2. Workspace Admin
	project, err := s.projectModel.GetByID(projectID)
	if err != nil {
		return errors.New("project not found")
	}

	wsMember, err := s.workspaceModel.GetMember(project.WorkspaceID, userID)
	if err == nil && wsMember.Role == models.WorkspaceRoleAdmin {
		return nil
	}

	// 3. Team Member
	projectTeams, err := s.projectModel.GetTeams(projectID)
	if err == nil {
		for _, pt := range projectTeams {
			_, err := s.teamModel.GetMember(pt.TeamID, userID)
			if err == nil {
				return nil
			}
		}
	}

	return errors.New("unauthorized: access denied to project")
}

func (s *TaskService) validateAssignee(projectID, assigneeID uuid.UUID) error {
	// Check Direct Membership
	_, err := s.projectModel.GetMember(projectID, assigneeID)
	if err == nil {
		return nil
	}

	// Check Team Membership
	projectTeams, err := s.projectModel.GetTeams(projectID)
	if err == nil {
		for _, pt := range projectTeams {
			_, err := s.teamModel.GetMember(pt.TeamID, assigneeID)
			if err == nil {
				return nil
			}
		}
	}

	return errors.New("assignee must be a member of the project")
}
