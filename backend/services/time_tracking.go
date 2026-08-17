package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/realtime"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/structs"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TimeTrackingService defines the interface for time tracking business logic
type TimeTrackingService interface {
	StartTimer(userID, taskID uuid.UUID) (*models.TimeEntryWithDetails, error)
	StopTimer(userID, taskID uuid.UUID, req structs.ReqStopTimer) (*models.TimeEntry, error)
	DiscardTimer(userID uuid.UUID) error
	GetActiveTimer(userID uuid.UUID) (*models.TimeEntryWithDetails, error)
	LogManualTime(userID, taskID uuid.UUID, req structs.ReqLogManualTime) (*models.TimeEntry, error)
	ListTaskTimeEntries(userID, taskID uuid.UUID) ([]models.TimeEntryWithDetails, int, *int, error)
	DeleteTimeEntry(userID uuid.UUID, entryID uuid.UUID) error
	GetProjectAnalytics(userID, projectID uuid.UUID) (*models.ProjectTimeAnalytics, error)
	ListProjectTimeEntries(userID, projectID uuid.UUID, targetUserID *uuid.UUID, startDate, endDate *time.Time) ([]models.TimeEntryWithDetails, error)
}

type timeTrackingService struct {
	storage models.Storage
	hub     realtime.IHub
	logger  *zap.Logger
}

// NewTimeTrackingService creates a new time tracking service instance
func NewTimeTrackingService(storage models.Storage, logger *zap.Logger, hub realtime.IHub) TimeTrackingService {
	return &timeTrackingService{
		storage: storage,
		logger:  logger,
		hub:     hub,
	}
}

func (s *timeTrackingService) validateProjectAccess(projectID, userID uuid.UUID) error {
	// 1. Direct Member
	_, err := s.storage.Projects().GetMember(projectID, userID)
	if err == nil {
		return nil
	}
	// 2. Workspace Admin
	project, err := s.storage.Projects().GetByID(projectID)
	if err != nil {
		return errors.New("project not found")
	}

	member, err := s.storage.Workspaces().GetMember(project.WorkspaceID, userID)
	if err == nil && member.Role == models.WorkspaceRoleAdmin {
		return nil
	}

	// 3. Team Member
	teams, err := s.storage.Projects().GetTeams(projectID)
	if err != nil {
		return errors.New("access denied")
	}

	for _, team := range teams {
		_, err := s.storage.Teams().GetMember(team.TeamID, userID)
		if err == nil {
			return nil
		}
	}

	return errors.New("access denied")
}

func (s *timeTrackingService) validateTaskAccess(taskID, userID uuid.UUID) (*models.TaskWithAssignee, error) {
	task, err := s.storage.Tasks().GetByID(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}

	err = s.validateProjectAccess(task.ProjectID, userID)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *timeTrackingService) StartTimer(userID, taskID uuid.UUID) (*models.TimeEntryWithDetails, error) {
	task, err := s.validateTaskAccess(taskID, userID)
	if err != nil {
		return nil, err
	}

	project, err := s.storage.Projects().GetByID(task.ProjectID)
	if err != nil {
		return nil, err
	}

	err = s.storage.Atomic(context.Background(), func(tx models.Storage) error {
		active, err := tx.TimeEntries().GetActiveByUserID(userID)
		if err != nil {
			return err
		}

		if active != nil {
			if active.TaskID == taskID {
				return errors.New("timer is already running for this task")
			}

			// Stop the current timer
			now := time.Now()
			dur := int(math.Round(now.Sub(active.StartTime).Minutes()))
			if dur < 1 {
				dur = 1
			}

			active.EndTime = &now
			active.DurationMinutes = &dur
			err = tx.TimeEntries().Update(active.TimeEntry)
			if err != nil {
				return err
			}
		}

		entry := models.TimeEntry{
			TaskID:      taskID,
			UserID:      userID,
			WorkspaceID: project.WorkspaceID,
			StartTime:   time.Now(),
			IsManual:    false,
		}

		_, err = tx.TimeEntries().Create(entry)
		return err
	})

	if err != nil {
		return nil, err
	}

	result, err := s.storage.TimeEntries().GetActiveByUserID(userID)
	if err != nil {
		return nil, err
	}

	if s.hub != nil {
		s.hub.Broadcast(fmt.Sprintf("user:%s", userID.String()), structs.RealtimeEvent{
			Type: structs.EventTypeTimerStarted,
			Payload: structs.ResActiveTimer{
				ID:         result.ID,
				TaskID:     result.TaskID,
				TaskTitle:  result.TaskTitle,
				TaskNumber: result.TaskNumber,
				StartTime:  result.StartTime,
			},
		})
	}

	return result, nil
}

func (s *timeTrackingService) StopTimer(userID, taskID uuid.UUID, req structs.ReqStopTimer) (*models.TimeEntry, error) {
	_, err := s.validateTaskAccess(taskID, userID)
	if err != nil {
		return nil, err
	}

	active, err := s.storage.TimeEntries().GetActiveByUserID(userID)
	if err != nil {
		return nil, err
	}

	if active == nil || active.TaskID != taskID {
		return nil, errors.New("no active timer found for this task")
	}

	now := time.Now()
	var dur int
	if req.DurationMinutes != nil {
		dur = *req.DurationMinutes
	} else {
		dur = int(math.Round(now.Sub(active.StartTime).Minutes()))
		if dur < 1 {
			dur = 1
		}
	}

	active.EndTime = &now
	active.DurationMinutes = &dur
	active.Description = req.Description

	err = s.storage.TimeEntries().Update(active.TimeEntry)
	if err != nil {
		return nil, err
	}

	if s.hub != nil {
		s.hub.Broadcast(fmt.Sprintf("user:%s", userID.String()), structs.RealtimeEvent{
			Type:    structs.EventTypeTimerStopped,
			Payload: active.TimeEntry,
		})
	}

	return &active.TimeEntry, nil
}

func (s *timeTrackingService) DiscardTimer(userID uuid.UUID) error {
	active, err := s.storage.TimeEntries().GetActiveByUserID(userID)
	if err != nil {
		return err
	}

	if active == nil {
		return errors.New("no active timer found")
	}

	err = s.storage.TimeEntries().Delete(active.ID)
	if err != nil {
		return err
	}

	if s.hub != nil {
		s.hub.Broadcast(fmt.Sprintf("user:%s", userID.String()), structs.RealtimeEvent{
			Type:    structs.EventTypeTimerDiscarded,
			Payload: nil,
		})
	}

	return nil
}

func (s *timeTrackingService) GetActiveTimer(userID uuid.UUID) (*models.TimeEntryWithDetails, error) {
	return s.storage.TimeEntries().GetActiveByUserID(userID)
}

func (s *timeTrackingService) LogManualTime(userID, taskID uuid.UUID, req structs.ReqLogManualTime) (*models.TimeEntry, error) {
	task, err := s.validateTaskAccess(taskID, userID)
	if err != nil {
		return nil, err
	}

	project, err := s.storage.Projects().GetByID(task.ProjectID)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	if req.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err == nil {
			startTime = parsedDate
		}
	}

	endTime := startTime.Add(time.Duration(req.DurationMinutes) * time.Minute)

	entry := models.TimeEntry{
		TaskID:          taskID,
		UserID:          userID,
		WorkspaceID:     project.WorkspaceID,
		Description:     req.Description,
		StartTime:       startTime,
		EndTime:         &endTime,
		DurationMinutes: &req.DurationMinutes,
		IsManual:        true,
	}

	created, err := s.storage.TimeEntries().Create(entry)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (s *timeTrackingService) ListTaskTimeEntries(userID, taskID uuid.UUID) ([]models.TimeEntryWithDetails, int, *int, error) {
	task, err := s.validateTaskAccess(taskID, userID)
	if err != nil {
		return nil, 0, nil, err
	}

	entries, err := s.storage.TimeEntries().ListByTaskID(taskID)
	if err != nil {
		return nil, 0, nil, err
	}

	total, err := s.storage.TimeEntries().SumByTaskID(taskID)
	if err != nil {
		return nil, 0, nil, err
	}

	return entries, total, task.EstimatedMinutes, nil
}

func (s *timeTrackingService) DeleteTimeEntry(userID uuid.UUID, entryID uuid.UUID) error {
	entry, err := s.storage.TimeEntries().GetByID(entryID)
	if err != nil {
		return errors.New("time entry not found")
	}

	if entry.UserID == userID {
		return s.storage.TimeEntries().Delete(entryID)
	}

	task, err := s.validateTaskAccess(entry.TaskID, userID)
	if err != nil {
		return err
	}

	// Verify project admin
	project, err := s.storage.Projects().GetByID(task.ProjectID)
	if err != nil {
		return errors.New("project not found")
	}

	wsMember, err := s.storage.Workspaces().GetMember(project.WorkspaceID, userID)
	if err == nil && wsMember.Role == models.WorkspaceRoleAdmin {
		return s.storage.TimeEntries().Delete(entryID)
	}

	member, err := s.storage.Projects().GetMember(task.ProjectID, userID)
	if err == nil && member.Role == models.ProjectRoleAdmin {
		return s.storage.TimeEntries().Delete(entryID)
	}

	return errors.New("unauthorized to delete this time entry")
}

func (s *timeTrackingService) GetProjectAnalytics(userID, projectID uuid.UUID) (*models.ProjectTimeAnalytics, error) {
	err := s.validateProjectAccess(projectID, userID)
	if err != nil {
		return nil, err
	}

	return s.storage.TimeEntries().GetProjectAnalytics(projectID)
}

func (s *timeTrackingService) ListProjectTimeEntries(userID, projectID uuid.UUID, targetUserID *uuid.UUID, startDate, endDate *time.Time) ([]models.TimeEntryWithDetails, error) {
	err := s.validateProjectAccess(projectID, userID)
	if err != nil {
		return nil, err
	}
	return s.storage.TimeEntries().ListByProjectID(projectID, targetUserID, startDate, endDate)
}
