package services

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/structs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupProjectTest(_ *testing.T) (*projectService, *mockProjectRepository, *mockWorkspaceRepository, *mockStorage) {
	mockProjectRepo := new(mockProjectRepository)
	mockWorkspaceRepo := new(mockWorkspaceRepository)
	mockStor := new(mockStorage)
	logger := zap.NewNop()

	mockStor.On("Projects").Return(mockProjectRepo)
	mockStor.On("Workspaces").Return(mockWorkspaceRepo)

	svc := &projectService{storage: mockStor, logger: logger}

	return svc, mockProjectRepo, mockWorkspaceRepo, mockStor
}

func TestProjectService_CreateProject(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mp, mw, _ := setupProjectTest(t)
		userID := uuid.New()
		workspaceID := uuid.New()
		projectID := uuid.New()

		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mp.On("Create", mock.Anything).Return(models.Project{ID: projectID, Name: "New Project"}, nil)
		mp.On("AddMember", mock.Anything).Return(nil)

		res, err := svc.CreateProject(userID, workspaceID, structs.ReqCreateProject{Name: "New Project"})

		assert.NoError(t, err)
		assert.Equal(t, "New Project", res.Name)
	})

	t.Run("get member fail", func(t *testing.T) {
		svc, _, mw, _ := setupProjectTest(t)
		userID := uuid.New()
		workspaceID := uuid.New()

		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{}, errors.New("db error"))

		_, err := svc.CreateProject(userID, workspaceID, structs.ReqCreateProject{Name: "New Project"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get workspace member")
	})

	t.Run("fail - project creation", func(t *testing.T) {
		svc, mp, mw, _ := setupProjectTest(t)
		userID := uuid.New()
		workspaceID := uuid.New()

		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mp.On("Create", mock.Anything).Return(models.Project{}, errors.New("db error"))

		_, err := svc.CreateProject(userID, workspaceID, structs.ReqCreateProject{Name: "Fail"})
		assert.Error(t, err)
	})

	t.Run("AddMember failure", func(t *testing.T) {
		svc, mp, mw, _ := setupProjectTest(t)
		userID := uuid.New()
		workspaceID := uuid.New()
		projectID := uuid.New()

		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mp.On("Create", mock.Anything).Return(models.Project{ID: projectID, Name: "New Project"}, nil)
		mp.On("AddMember", mock.Anything).Return(errors.New("db error"))

		_, err := svc.CreateProject(userID, workspaceID, structs.ReqCreateProject{Name: "New Project"})
		assert.Error(t, err)
	})

	t.Run("unauthorized - not workspace admin", func(t *testing.T) {
		svc, _, mw, _ := setupProjectTest(t)
		userID := uuid.New()
		workspaceID := uuid.New()

		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)

		_, err := svc.CreateProject(userID, workspaceID, structs.ReqCreateProject{Name: "Fail"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only workspace admins")
	})

	t.Run("not a member of workspace", func(t *testing.T) {
		svc, _, mw, _ := setupProjectTest(t)
		userID := uuid.New()
		workspaceID := uuid.New()
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{}, sql.ErrNoRows)

		_, err := svc.CreateProject(userID, workspaceID, structs.ReqCreateProject{Name: "Fail"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a member of this workspace")
	})

}

func TestProjectService_GetProject(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mp, _, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()

		// Access validation (Direct member)
		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{Role: models.ProjectRoleMember}, nil)
		mp.On("GetByID", projectID).Return(models.Project{ID: projectID, Name: "Found"}, nil)

		p, err := svc.GetProject(userID, projectID)
		assert.NoError(t, err)
		assert.Equal(t, "Found", p.Name)
	})

	t.Run("project not found", func(t *testing.T) {
		svc, mp, _, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not found"))
		mp.On("GetByID", projectID).Return(models.Project{}, errors.New("not found"))

		_, err := svc.GetProject(userID, projectID)
		assert.Error(t, err)
	})

}

func TestProjectService_UpdateProject(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mp, _, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil)
		mp.On("GetByID", projectID).Return(models.Project{ID: projectID, Name: "Old"}, nil)
		mp.On("Update", mock.Anything).Return(models.Project{Name: "New", Description: "Description"}, nil)

		res, err := svc.UpdateProject(userID, projectID, structs.ReqUpdateProject{Name: "New", Description: "Description"})
		assert.NoError(t, err)
		assert.Equal(t, "New", res.Name)
		assert.Equal(t, "Description", res.Description)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, mp, mw, _ := setupProjectTest(t)
		uid, pid, wid := uuid.New(), uuid.New(), uuid.New()

		mp.On("GetMember", pid, uid).Return(models.ProjectMember{Role: models.ProjectRoleMember}, nil)
		mp.On("GetByID", pid).Return(models.Project{WorkspaceID: wid}, nil)
		mw.On("GetMember", wid, uid).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)

		_, err := svc.UpdateProject(uid, pid, structs.ReqUpdateProject{Name: "X"})
		assert.Error(t, err)
	})

	t.Run("atomic_failure", func(t *testing.T) {
		svc, mp, _, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil)
		mp.On("GetByID", projectID).Return(models.Project{ID: projectID, Name: "Old"}, nil)
		mp.On("Update", mock.Anything).Return(models.Project{}, errors.New("db error"))

		_, err := svc.UpdateProject(userID, projectID, structs.ReqUpdateProject{Name: "New"})
		assert.Error(t, err)
	})

	t.Run("project not found", func(t *testing.T) {
		svc, mp, _, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not found"))
		mp.On("GetByID", projectID).Return(models.Project{}, errors.New("not found"))

		_, err := svc.UpdateProject(userID, projectID, structs.ReqUpdateProject{Name: "New"})
		assert.Error(t, err)
	})

	t.Run("update_is_archived", func(t *testing.T) {
		svc, mp, _, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()
		isArchived := true

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil)
		mp.On("GetByID", projectID).Return(models.Project{ID: projectID, Name: "Old", IsArchived: false}, nil)
		mp.On("Update", mock.MatchedBy(func(p models.Project) bool {
			return p.IsArchived == true
		})).Return(models.Project{ID: projectID, Name: "Old", IsArchived: true}, nil)

		res, err := svc.UpdateProject(userID, projectID, structs.ReqUpdateProject{IsArchived: &isArchived})
		assert.NoError(t, err)
		assert.True(t, res.IsArchived)
	})

}

func TestProjectService_AddMember(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockProjectRepo, mockWorkspaceRepo, mockStor := setupProjectTest(t)
		mockUserRepo := new(mockUserRepository)
		mockStor.On("Users").Return(mockUserRepo)

		requestorID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()
		targetEmail := "user@example.com"
		targetUserID := uuid.New()

		// Access validation (Admin)
		mockProjectRepo.On("GetMember", projectID, requestorID).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil)

		// Find target user
		mockUserRepo.On("GetByEmail", targetEmail).Return(models.User{ID: targetUserID}, nil)

		// Project details for workspace ID
		mockProjectRepo.On("GetByID", projectID).Return(models.Project{ID: projectID, WorkspaceID: workspaceID}, nil)

		// Workspace membership check
		mockWorkspaceRepo.On("GetMember", workspaceID, targetUserID).Return(models.WorkspaceMember{}, nil)

		// Check if already in project
		mockProjectRepo.On("GetMember", projectID, targetUserID).Return(models.ProjectMember{}, errors.New("not found"))

		// Add Member
		mockProjectRepo.On("AddMember", mock.MatchedBy(func(m models.ProjectMember) bool {
			return m.ProjectID == projectID && m.UserID == targetUserID
		})).Return(nil)

		member, err := svc.AddMember(requestorID, projectID, structs.ReqAddProjectMember{Email: targetEmail})

		assert.NoError(t, err)
		assert.Equal(t, targetUserID, member.UserID)
	})

	t.Run("user not found", func(t *testing.T) {
		svc, mp, _, ms := setupProjectTest(t)
		mu := new(mockUserRepository)

		ms.On("Users").Return(mu)
		uid, pid := uuid.New(), uuid.New()
		mp.On("GetMember", pid, uid).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil)
		mu.On("GetByEmail", mock.Anything).Return(models.User{}, errors.New("not found"))

		_, err := svc.AddMember(uid, pid, structs.ReqAddProjectMember{Email: "none@test.com"})
		assert.Error(t, err)
	})

	t.Run("user not in workspace", func(t *testing.T) {
		svc, mp, mw, ms := setupProjectTest(t)
		mu := new(mockUserRepository)

		ms.On("Users").Return(mu)
		uid, pid, wid, tid := uuid.New(), uuid.New(), uuid.New(), uuid.New()
		mp.On("GetMember", pid, uid).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil)
		mu.On("GetByEmail", mock.Anything).Return(models.User{ID: tid}, nil)
		mp.On("GetByID", pid).Return(models.Project{WorkspaceID: wid}, nil)
		mw.On("GetMember", wid, tid).Return(models.WorkspaceMember{}, errors.New("not in ws"))

		_, err := svc.AddMember(uid, pid, structs.ReqAddProjectMember{Email: "user@test.com"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a member of the workspace")
	})

	t.Run("user already in project", func(t *testing.T) {
		svc, mp, mw, ms := setupProjectTest(t)
		mu := new(mockUserRepository)

		ms.On("Users").Return(mu)
		uid, pid, wid, tid := uuid.New(), uuid.New(), uuid.New(), uuid.New()
		mp.On("GetMember", pid, uid).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil).Once()
		mu.On("GetByEmail", mock.Anything).Return(models.User{ID: tid}, nil)
		mp.On("GetByID", pid).Return(models.Project{WorkspaceID: wid}, nil)
		mw.On("GetMember", wid, tid).Return(models.WorkspaceMember{}, nil)
		mp.On("GetMember", pid, tid).Return(models.ProjectMember{}, nil) // Already in project

		_, err := svc.AddMember(uid, pid, structs.ReqAddProjectMember{Email: "user@test.com"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user already in project")
	})

	t.Run("unauthorized - not project admin", func(t *testing.T) {
		svc, mp, mw, _ := setupProjectTest(t)
		uid, pid, wid := uuid.New(), uuid.New(), uuid.New()
		mp.On("GetMember", pid, uid).Return(models.ProjectMember{Role: models.ProjectRoleMember}, nil)
		// Since it's not a project admin, it checks if it's a workspace admin
		mp.On("GetByID", pid).Return(models.Project{WorkspaceID: wid}, nil)
		mw.On("GetMember", wid, uid).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)

		_, err := svc.AddMember(uid, pid, structs.ReqAddProjectMember{Email: "user@test.com"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("atomic failure", func(t *testing.T) {
		svc, mp, mw, ms := setupProjectTest(t)
		mu := new(mockUserRepository)

		ms.On("Users").Return(mu)
		uid, pid, wid, tid := uuid.New(), uuid.New(), uuid.New(), uuid.New()
		mp.On("GetMember", pid, uid).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil)
		mu.On("GetByEmail", mock.Anything).Return(models.User{ID: tid}, nil)
		mp.On("GetByID", pid).Return(models.Project{WorkspaceID: wid}, nil)
		mw.On("GetMember", wid, tid).Return(models.WorkspaceMember{}, nil)
		mp.On("GetMember", pid, tid).Return(models.ProjectMember{}, errors.New("not found"))
		mp.On("AddMember", mock.Anything).Return(errors.New("db error"))

		_, err := svc.AddMember(uid, pid, structs.ReqAddProjectMember{Email: "user@test.com"})
		assert.Error(t, err)
	})

}

func TestProjectService_RemoveMember(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mp, _, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()
		targetID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil)
		mp.On("RemoveMember", projectID, targetID).Return(nil)

		err := svc.RemoveMember(userID, projectID, targetID)
		assert.NoError(t, err)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, mp, mw, _ := setupProjectTest(t)
		uid, pid, wid := uuid.New(), uuid.New(), uuid.New()
		mp.On("GetMember", pid, uid).Return(models.ProjectMember{Role: models.ProjectRoleMember}, nil)
		mp.On("GetByID", pid).Return(models.Project{WorkspaceID: wid}, nil)
		mw.On("GetMember", wid, uid).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)

		err := svc.RemoveMember(uid, pid, uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})
}

func TestProjectService_ListMembers(t *testing.T) {
	t.Run("success with team members", func(t *testing.T) {
		svc, mp, _, ms := setupProjectTest(t)
		mt := new(mockTeamRepository)
		ms.On("Teams").Return(mt)

		userID := uuid.New()
		projectID := uuid.New()
		teamID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)
		mp.On("GetMembers", projectID).Return([]models.ProjectMemberWithUser{
			{UserID: userID, Email: "admin@test.com", Role: models.ProjectRoleAdmin},
		}, nil)
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{{TeamID: teamID}}, nil)
		mt.On("ListMembersByTeamID", teamID).Return([]models.TeamMemberWithUser{
			{UserID: uuid.New(), Email: "team@test.com"},
		}, nil)

		members, err := svc.ListMembers(userID, projectID)
		assert.NoError(t, err)
		assert.Len(t, members, 2)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, mp, mw, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not found"))
		mp.On("GetByID", projectID).Return(models.Project{WorkspaceID: workspaceID}, nil)
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, nil)

		_, err := svc.ListMembers(userID, projectID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("get direct members fail", func(t *testing.T) {
		svc, mp, _, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)
		mp.On("GetMembers", projectID).Return([]models.ProjectMemberWithUser{}, errors.New("db error"))

		_, err := svc.ListMembers(userID, projectID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("get teams fail", func(t *testing.T) {
		svc, mp, _, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)
		mp.On("GetMembers", projectID).Return([]models.ProjectMemberWithUser{}, nil)
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, errors.New("db error"))

		_, err := svc.ListMembers(userID, projectID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestProjectService_ListByWorkspaceID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mp, _, _ := setupProjectTest(t)
		workspaceID := uuid.New()
		mp.On("ListByWorkspaceID", workspaceID).Return([]models.Project{{Name: "P1"}}, nil)

		res, err := svc.ListByWorkspaceID(workspaceID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}

func TestProjectService_AddTeam(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mp, _, ms := setupProjectTest(t)
		mt := new(mockTeamRepository)
		ms.On("Teams").Return(mt)

		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()
		teamID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil)
		mt.On("GetByID", teamID).Return(models.Team{ID: teamID, WorkspaceID: workspaceID, Name: "Team A"}, nil)
		mp.On("GetByID", projectID).Return(models.Project{ID: projectID, WorkspaceID: workspaceID}, nil)
		mp.On("GetTeam", projectID, teamID).Return(models.ProjectTeam{}, errors.New("not found"))
		mp.On("AddTeam", mock.Anything).Return(nil)

		res, err := svc.AddTeam(userID, projectID, teamID)
		assert.NoError(t, err)
		assert.Equal(t, teamID, res.TeamID)
	})

	// unauthorized
	t.Run("unauthorized", func(t *testing.T) {
		svc, mp, mw, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not found"))
		mp.On("GetByID", projectID).Return(models.Project{WorkspaceID: workspaceID}, nil)
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, nil)

		_, err := svc.AddTeam(userID, projectID, uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("workspace mismatch", func(t *testing.T) {
		svc, mp, _, ms := setupProjectTest(t)
		mt := new(mockTeamRepository)

		ms.On("Teams").Return(mt)
		uid, pid, wid := uuid.New(), uuid.New(), uuid.New()
		mp.On("GetMember", pid, uid).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil)
		mt.On("GetByID", mock.Anything).Return(models.Team{WorkspaceID: uuid.New()}, nil) // Different WS
		mp.On("GetByID", pid).Return(models.Project{WorkspaceID: wid}, nil)

		_, err := svc.AddTeam(uid, pid, uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "same workspace")
	})

	t.Run("team not exists", func(t *testing.T) {
		svc, mp, _, ms := setupProjectTest(t)
		mt := new(mockTeamRepository)

		ms.On("Teams").Return(mt)
		uid, pid := uuid.New(), uuid.New()
		mp.On("GetMember", pid, uid).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil)
		mt.On("GetByID", mock.Anything).Return(models.Team{}, errors.New("not found"))

		_, err := svc.AddTeam(uid, pid, uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "team not found")
	})

	t.Run("project not exist", func(t *testing.T) {
		svc, mp, _, ms := setupProjectTest(t)
		mt := new(mockTeamRepository)

		ms.On("Teams").Return(mt)
		uid, pid := uuid.New(), uuid.New()
		mp.On("GetMember", pid, uid).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil)
		mt.On("GetByID", mock.Anything).Return(models.Team{WorkspaceID: uuid.New()}, nil)
		mp.On("GetByID", pid).Return(models.Project{}, errors.New("not found"))

		_, err := svc.AddTeam(uid, pid, uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "project not found")
	})

	t.Run("team already added", func(t *testing.T) {
		svc, mp, _, ms := setupProjectTest(t)
		mt := new(mockTeamRepository)

		ms.On("Teams").Return(mt)
		uid, pid, wid := uuid.New(), uuid.New(), uuid.New()
		mp.On("GetMember", pid, uid).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil)
		mt.On("GetByID", mock.Anything).Return(models.Team{WorkspaceID: wid}, nil)
		mp.On("GetByID", pid).Return(models.Project{WorkspaceID: wid}, nil)
		mp.On("GetTeam", pid, mock.Anything).Return(models.ProjectTeam{}, nil) // Already added

		_, err := svc.AddTeam(uid, pid, uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "team already assigned to project")
	})

	// atomic failure
	t.Run("atomic failure", func(t *testing.T) {
		svc, mp, _, ms := setupProjectTest(t)
		mt := new(mockTeamRepository)
		ms.On("Teams").Return(mt)

		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()
		teamID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil)
		mt.On("GetByID", teamID).Return(models.Team{ID: teamID, WorkspaceID: workspaceID, Name: "Team A"}, nil)
		mp.On("GetByID", projectID).Return(models.Project{ID: projectID, WorkspaceID: workspaceID}, nil)
		mp.On("GetTeam", projectID, teamID).Return(models.ProjectTeam{}, errors.New("not found"))
		mp.On("AddTeam", mock.Anything).Return(errors.New("db error"))

		_, err := svc.AddTeam(userID, projectID, teamID)
		assert.Error(t, err)

	})

}

func TestProjectService_RemoveTeam(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mp, _, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()
		teamID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{Role: models.ProjectRoleAdmin}, nil)
		mp.On("RemoveTeam", projectID, teamID).Return(nil)

		err := svc.RemoveTeam(userID, projectID, teamID)
		assert.NoError(t, err)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, mp, mw, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not found"))
		mp.On("GetByID", projectID).Return(models.Project{WorkspaceID: workspaceID}, nil)
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, nil)

		err := svc.RemoveTeam(userID, projectID, uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})
}

func TestProjectService_ListTeams(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mp, _, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{{TeamName: "T1"}}, nil)

		res, err := svc.ListTeams(userID, projectID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, mp, mw, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not found"))
		mp.On("GetByID", projectID).Return(models.Project{WorkspaceID: workspaceID}, nil)
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, nil)

		_, err := svc.ListTeams(userID, projectID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("GetTeams_failure", func(t *testing.T) {
		svc, mp, _, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, nil)
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, errors.New("db error"))

		_, err := svc.ListTeams(userID, projectID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

}

func TestProjectService_ValidateProjectAccess(t *testing.T) {
	t.Run("direct member", func(t *testing.T) {
		svc, mp, _, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{Role: models.ProjectRoleMember}, nil)

		err := svc.ValidateProjectAccess(projectID, userID, false)
		assert.NoError(t, err)
	})

	t.Run("workspace admin", func(t *testing.T) {
		svc, mp, mw, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not found"))
		mp.On("GetByID", projectID).Return(models.Project{WorkspaceID: workspaceID}, nil)
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)

		err := svc.ValidateProjectAccess(projectID, userID, true)
		assert.NoError(t, err)
	})

	t.Run("team member access", func(t *testing.T) {
		svc, mp, mw, ms := setupProjectTest(t)
		mt := new(mockTeamRepository)
		ms.On("Teams").Return(mt)

		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()
		teamID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not found"))
		mp.On("GetByID", projectID).Return(models.Project{WorkspaceID: workspaceID}, nil)
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{{TeamID: teamID}}, nil)
		mt.On("GetMember", teamID, userID).Return(models.TeamMember{}, nil)

		err := svc.ValidateProjectAccess(projectID, userID, false)
		assert.NoError(t, err)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, mp, mw, _ := setupProjectTest(t)
		userID := uuid.New()
		projectID := uuid.New()
		workspaceID := uuid.New()

		mp.On("GetMember", projectID, userID).Return(models.ProjectMember{}, errors.New("not found"))
		mp.On("GetByID", projectID).Return(models.Project{WorkspaceID: workspaceID}, nil)
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)
		mp.On("GetTeams", projectID).Return([]models.ProjectTeamWithDetails{}, nil)

		err := svc.ValidateProjectAccess(projectID, userID, false)
		assert.Error(t, err)
		assert.Equal(t, "unauthorized", err.Error())
	})
}
