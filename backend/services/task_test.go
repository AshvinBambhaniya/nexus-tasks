package services

import (
	"errors"
	"testing"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupTaskTest(_ *testing.T) (*taskService, *mockTaskRepository, *mockProjectRepository, *mockWorkspaceRepository, *mockStorage, *mockHub) {
	mockTaskRepo := new(mockTaskRepository)
	mockProjectRepo := new(mockProjectRepository)
	mockWorkspaceRepo := new(mockWorkspaceRepository)
	mockStor := new(mockStorage)
	mockH := new(mockHub)
	logger := zap.NewNop()

	mockStor.On("Tasks").Return(mockTaskRepo)
	mockStor.On("Projects").Return(mockProjectRepo)
	mockStor.On("Workspaces").Return(mockWorkspaceRepo)

	svc := &taskService{
		storage: mockStor,
		logger:  logger,
		hub:     mockH,
	}

	return svc, mockTaskRepo, mockProjectRepo, mockWorkspaceRepo, mockStor, mockH
}

func TestTaskService_CreateTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mt, mp, _, _, mh := setupTaskTest(t)
		userID := uuid.New()
		projectID := uuid.New()
		taskID := uuid.New()
		dueDate := time.Now().Add(24 * time.Hour)

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)
		mt.On("Create", mock.Anything).Return(models.Task{ID: taskID, Title: "Bug", DueDate: &dueDate}, nil)
		mh.On("Broadcast", mock.Anything, mock.Anything).Return()

		res, err := svc.CreateTask(userID, projectID, structs.ReqCreateTask{Title: "Bug", DueDate: &structs.CustomTime{Time: dueDate}})
		assert.NoError(t, err)
		assert.Equal(t, taskID, res.ID)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, _, mp, mw, ms, _ := setupTaskTest(t)
		mt := new(mockTeamRepository)
		ms.On("Teams").Return(mt)

		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()

		// 1. Direct member check fails
		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not member"))
		// 2. Project details
		mp.On("GetByID", projectID).Return(models.Project{WorkspaceID: workspaceID}, nil)
		// 3. Workspace admin check fails
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)
		// 4. Team member check fails
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, nil)

		_, err := svc.CreateTask(userID, projectID, structs.ReqCreateTask{Title: "Fail"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("invalid assignee", func(t *testing.T) {
		svc, _, mp, _, _, _ := setupTaskTest(t)
		userID := uuid.New()
		projectID := uuid.New()
		assigneeID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)
		// Assignee is not in project
		mp.On("GetMember", projectID, assigneeID).Return(models.ProjectMember{}, errors.New("not found"))
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, nil)

		_, err := svc.CreateTask(userID, projectID, structs.ReqCreateTask{Title: "T", AssigneeID: &assigneeID})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "assignee must be a member")
	})

	t.Run("atomic_failure", func(t *testing.T) {
		svc, mt, mp, _, _, _ := setupTaskTest(t)
		userID := uuid.New()
		projectID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)
		mt.On("Create", mock.Anything).Return(models.Task{}, errors.New("db error"))

		_, err := svc.CreateTask(userID, projectID, structs.ReqCreateTask{Title: "Fail"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestTaskService_ListProjectTasks(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mt, mp, _, _, _ := setupTaskTest(t)
		projectID := uuid.New()
		userID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)
		mt.On("ListByProjectID", projectID, mock.Anything, mock.Anything).Return([]models.Task{{Title: "T1"}}, nil)

		res, err := svc.ListProjectTasks(userID, projectID, nil, nil)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, _, mp, mw, ms, _ := setupTaskTest(t)
		tr := new(mockTeamRepository)
		ms.On("Teams").Return(tr)
		projectID := uuid.New()
		userID := uuid.New()
		workspaceID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not member"))
		mp.On("GetByID", projectID).Return(models.Project{WorkspaceID: workspaceID}, nil)
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, nil)

		_, err := svc.ListProjectTasks(userID, projectID, nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})
}

func TestTaskService_GetTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mt, mp, _, _, _ := setupTaskTest(t)
		taskID := uuid.New()
		userID := uuid.New()
		projectID := uuid.New()

		mt.On("GetByID", taskID).Return(models.Task{ID: taskID, ProjectID: projectID}, nil)
		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)

		res, err := svc.GetTask(userID, taskID)
		assert.NoError(t, err)
		assert.Equal(t, taskID, res.ID)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, mt, mp, mw, ms, _ := setupTaskTest(t)

		tr := new(mockTeamRepository)
		ms.On("Teams").Return(tr)
		taskID := uuid.New()
		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()

		mt.On("GetByID", taskID).Return(models.Task{ID: taskID, ProjectID: projectID}, nil)
		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not member"))
		mp.On("GetByID", projectID).Return(models.Project{WorkspaceID: workspaceID}, nil)
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, nil)

		_, err := svc.GetTask(userID, taskID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("task not found", func(t *testing.T) {
		svc, mt, _, _, _, _ := setupTaskTest(t)
		mt.On("GetByID", mock.Anything).Return(models.Task{}, errors.New("not found"))

		_, err := svc.GetTask(uuid.New(), uuid.New())
		assert.Error(t, err)
	})
}

func TestTaskService_UpdateTask(t *testing.T) {
	t.Run("success with status transition", func(t *testing.T) {
		svc, mt, mp, _, _, mh := setupTaskTest(t)
		taskID := uuid.New()
		userID := uuid.New()
		projectID := uuid.New()
		assigneeID := uuid.New()
		dueDate := time.Now().Add(48 * time.Hour)

		mt.On("GetByID", taskID).Return(models.Task{ID: taskID, ProjectID: projectID, Status: "TODO"}, nil)
		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)

		// Validate assignee
		mp.On("GetMember", projectID, assigneeID).Return(models.ProjectMember{}, nil)

		// Update to DONE
		mt.On("Update", mock.MatchedBy(func(t models.Task) bool {
			return t.Status == "DONE" && t.CompletedAt != nil && *t.AssigneeID == assigneeID && t.DueDate.Equal(dueDate)
		})).Return(models.Task{ID: taskID, Status: "DONE", ProjectID: projectID}, nil)

		mh.On("Broadcast", mock.Anything, mock.Anything).Return()

		res, err := svc.UpdateTask(userID, taskID, structs.ReqUpdateTask{
			Title:       "Title",
			Description: "Description",
			AssigneeID:  &assigneeID,
			Priority:    models.TaskPriorityP0,
			Status:      "DONE",
			DueDate:     &structs.CustomTime{Time: dueDate},
		})
		assert.NoError(t, err)
		assert.Equal(t, models.TaskStatus("DONE"), res.Status)
	})

	t.Run("task not found", func(t *testing.T) {
		svc, mt, _, _, _, _ := setupTaskTest(t)
		taskID := uuid.New()
		userID := uuid.New()

		mt.On("GetByID", taskID).Return(models.Task{}, errors.New("not found"))

		_, err := svc.UpdateTask(userID, taskID, structs.ReqUpdateTask{Title: "X"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, mt, mp, mw, ms, _ := setupTaskTest(t)
		tr := new(mockTeamRepository)
		ms.On("Teams").Return(tr)

		taskID := uuid.New()
		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()

		mt.On("GetByID", taskID).Return(models.Task{ID: taskID, ProjectID: projectID}, nil)
		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not member"))
		mp.On("GetByID", projectID).Return(models.Project{WorkspaceID: workspaceID}, nil)
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, nil)

		_, err := svc.UpdateTask(userID, taskID, structs.ReqUpdateTask{Title: "X"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("invalid assignee", func(t *testing.T) {
		svc, mt, mp, _, _, _ := setupTaskTest(t)
		taskID := uuid.New()
		userID := uuid.New()
		projectID := uuid.New()
		assigneeID := uuid.New()

		mt.On("GetByID", taskID).Return(models.Task{ID: taskID, ProjectID: projectID}, nil)
		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)
		// Assignee check fails
		mp.On("GetMember", projectID, assigneeID).Return(models.ProjectMember{}, errors.New("not found"))
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, nil)

		_, err := svc.UpdateTask(userID, taskID, structs.ReqUpdateTask{AssigneeID: &assigneeID})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "assignee must be a member")
	})

	t.Run("update failure", func(t *testing.T) {
		svc, mt, mp, _, _, _ := setupTaskTest(t)
		taskID := uuid.New()
		userID := uuid.New()
		projectID := uuid.New()

		mt.On("GetByID", taskID).Return(models.Task{ID: taskID, ProjectID: projectID}, nil)
		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)
		mt.On("Update", mock.Anything).Return(models.Task{}, errors.New("db error"))

		_, err := svc.UpdateTask(userID, taskID, structs.ReqUpdateTask{Title: "X"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("reopen task clears CompletedAt", func(t *testing.T) {
		svc, mt, mp, _, _, mh := setupTaskTest(t)
		taskID := uuid.New()
		userID := uuid.New()
		projectID := uuid.New()
		completedAt := time.Now()

		mt.On("GetByID", taskID).Return(models.Task{ID: taskID, ProjectID: projectID, Status: "DONE", CompletedAt: &completedAt}, nil)
		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)
		mt.On("Update", mock.MatchedBy(func(t models.Task) bool {
			return t.Status == "TODO" && t.CompletedAt == nil
		})).Return(models.Task{ID: taskID, Status: "TODO", ProjectID: projectID}, nil)
		mh.On("Broadcast", mock.Anything, mock.Anything).Return()

		res, err := svc.UpdateTask(userID, taskID, structs.ReqUpdateTask{Status: "TODO"})
		assert.NoError(t, err)
		assert.Equal(t, models.TaskStatus("TODO"), res.Status)
		assert.Nil(t, res.CompletedAt)
	})
}

func TestTaskService_DeleteTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mt, mp, _, _, mh := setupTaskTest(t)
		taskID := uuid.New()
		projectID := uuid.New()
		userID := uuid.New()

		mt.On("GetByID", taskID).Return(models.Task{ID: taskID, ProjectID: projectID}, nil)
		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)
		mt.On("Delete", taskID).Return(nil)
		mh.On("Broadcast", mock.Anything, mock.Anything).Return()

		err := svc.DeleteTask(userID, taskID)
		assert.NoError(t, err)
	})

	t.Run("task not found", func(t *testing.T) {
		svc, mt, _, _, _, _ := setupTaskTest(t)
		taskID := uuid.New()
		userID := uuid.New()

		mt.On("GetByID", taskID).Return(models.Task{}, errors.New("not found"))

		err := svc.DeleteTask(userID, taskID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, mt, mp, mw, ms, _ := setupTaskTest(t)
		tr := new(mockTeamRepository)
		ms.On("Teams").Return(tr)

		taskID := uuid.New()
		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()

		mt.On("GetByID", taskID).Return(models.Task{ID: taskID, ProjectID: projectID}, nil)
		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not member"))
		mp.On("GetByID", projectID).Return(models.Project{WorkspaceID: workspaceID}, nil)
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, nil)

		err := svc.DeleteTask(userID, taskID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("delete failure", func(t *testing.T) {
		svc, mt, mp, _, _, _ := setupTaskTest(t)
		taskID := uuid.New()
		userID := uuid.New()
		projectID := uuid.New()

		mt.On("GetByID", taskID).Return(models.Task{ID: taskID, ProjectID: projectID}, nil)
		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)
		mt.On("Delete", taskID).Return(errors.New("db error"))

		err := svc.DeleteTask(userID, taskID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestTaskService_ListMyTasks(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mt, _, _, _, _ := setupTaskTest(t)
		userID := uuid.New()

		mt.On("ListByAssigneeID", userID).Return([]models.Task{{Title: "My Task"}}, nil)

		res, err := svc.ListMyTasks(userID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}

func TestTaskService_InternalHelpers(t *testing.T) {
	svc, _, mp, mw, ms, _ := setupTaskTest(t)
	tr := new(mockTeamRepository)
	ms.On("Teams").Return(tr)

	t.Run("internalValidateProjectAccess - direct member", func(t *testing.T) {
		userID := uuid.New()
		projectID := uuid.New()
		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil).Once()

		err := svc.internalValidateProjectAccess(ms, projectID, userID)
		assert.NoError(t, err)
	})

	t.Run("internalValidateProjectAccess - workspace admin", func(t *testing.T) {
		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not found")).Once()
		mp.On("GetByID", projectID).Return(models.Project{WorkspaceID: workspaceID}, nil).Once()
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil).Once()

		err := svc.internalValidateProjectAccess(ms, projectID, userID)
		assert.NoError(t, err)
	})

	t.Run("internalValidateProjectAccess - team member", func(t *testing.T) {
		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()
		teamID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not found")).Once()
		mp.On("GetByID", projectID).Return(models.Project{WorkspaceID: workspaceID}, nil).Once()
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil).Once()
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{{TeamID: teamID}}, nil).Once()
		tr.On("GetMember", teamID, userID).Return(models.TeamMember{}, nil).Once()

		err := svc.internalValidateProjectAccess(ms, projectID, userID)
		assert.NoError(t, err)
	})

	t.Run("internalValidateProjectAccess - unauthorized", func(t *testing.T) {
		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not found")).Once()
		mp.On("GetByID", projectID).Return(models.Project{WorkspaceID: workspaceID}, nil).Once()
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil).Once()
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, nil).Once()

		err := svc.internalValidateProjectAccess(ms, projectID, userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("internalValidateAssignee - direct member", func(t *testing.T) {
		assigneeID := uuid.New()
		projectID := uuid.New()
		mp.On("GetMember", projectID, assigneeID).Return(models.ProjectMember{}, nil).Once()

		err := svc.internalValidateAssignee(ms, projectID, assigneeID)
		assert.NoError(t, err)
	})

	t.Run("internalValidateAssignee - team member", func(t *testing.T) {
		assigneeID := uuid.New()
		projectID := uuid.New()
		teamID := uuid.New()

		mp.On("GetMember", projectID, assigneeID).Return(models.ProjectMember{}, errors.New("not found")).Once()
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{{TeamID: teamID}}, nil).Once()
		tr.On("GetMember", teamID, assigneeID).Return(models.TeamMember{}, nil).Once()

		err := svc.internalValidateAssignee(ms, projectID, assigneeID)
		assert.NoError(t, err)
	})

	t.Run("internalValidateAssignee - failure", func(t *testing.T) {
		assigneeID := uuid.New()
		projectID := uuid.New()

		mp.On("GetMember", projectID, assigneeID).Return(models.ProjectMember{}, errors.New("not found")).Once()
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, nil).Once()

		err := svc.internalValidateAssignee(ms, projectID, assigneeID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "assignee must be a member")
	})
}
