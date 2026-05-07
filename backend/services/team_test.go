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

func setupTeamTest(_ *testing.T) (*teamService, *mockTeamRepository, *mockWorkspaceRepository, *mockStorage) {
	mockTeamRepo := new(mockTeamRepository)
	mockWorkspaceRepo := new(mockWorkspaceRepository)
	mockStor := new(mockStorage)
	logger := zap.NewNop()

	mockStor.On("Teams").Return(mockTeamRepo)
	mockStor.On("Workspaces").Return(mockWorkspaceRepo)

	svc := &teamService{storage: mockStor, logger: logger}

	return svc, mockTeamRepo, mockWorkspaceRepo, mockStor
}

func TestTeamService_CreateTeam(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mt, mw, _ := setupTeamTest(t)
		userID := uuid.New()
		workspaceID := uuid.New()
		teamID := uuid.New()

		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mt.On("CreateTeam", mock.Anything).Return(models.Team{ID: teamID, Name: "Backend"}, nil)
		mt.On("AddMember", mock.Anything).Return(nil)

		res, err := svc.CreateTeam(userID, workspaceID, structs.ReqCreateTeam{Name: "Backend"})
		assert.NoError(t, err)
		assert.Equal(t, "Backend", res.Name)
	})

	t.Run("get workspace member fail", func(t *testing.T) {
		svc, _, mw, _ := setupTeamTest(t)
		mw.On("GetMember", mock.Anything, mock.Anything).Return(models.WorkspaceMember{}, errors.New("db error"))

		_, err := svc.CreateTeam(uuid.New(), uuid.New(), structs.ReqCreateTeam{Name: "Fail"})
		assert.Error(t, err)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, _, mw, _ := setupTeamTest(t)
		mw.On("GetMember", mock.Anything, mock.Anything).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)

		_, err := svc.CreateTeam(uuid.New(), uuid.New(), structs.ReqCreateTeam{Name: "Fail"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("atomic failure - create team", func(t *testing.T) {
		svc, mt, mw, _ := setupTeamTest(t)
		userID := uuid.New()
		workspaceID := uuid.New()

		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mt.On("CreateTeam", mock.Anything).Return(models.Team{}, errors.New("db error"))

		_, err := svc.CreateTeam(userID, workspaceID, structs.ReqCreateTeam{Name: "Fail"})
		assert.Error(t, err)
	})

	t.Run("atomic failure - add member", func(t *testing.T) {
		svc, mt, mw, _ := setupTeamTest(t)
		userID := uuid.New()
		workspaceID := uuid.New()
		teamID := uuid.New()

		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mt.On("CreateTeam", mock.Anything).Return(models.Team{ID: teamID, Name: "Backend"}, nil)
		mt.On("AddMember", mock.Anything).Return(errors.New("db error"))

		_, err := svc.CreateTeam(userID, workspaceID, structs.ReqCreateTeam{Name: "Fail"})
		assert.Error(t, err)
	})
}

func TestTeamService_GetTeam(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mt, _, ms := setupTeamTest(t)
		mp := new(mockProjectRepository)
		ms.On("Projects").Return(mp)

		teamID := uuid.New()
		mt.On("GetByID", teamID).Return(models.Team{ID: teamID, Name: "Team A"}, nil)
		mp.On("ListByTeamID", teamID).Return([]models.Project{{Name: "P1"}}, nil)

		res, err := svc.GetTeam(teamID)
		assert.NoError(t, err)
		assert.Equal(t, "Team A", res.Name)
		assert.Len(t, res.Projects, 1)
	})

	t.Run("team not found", func(t *testing.T) {
		svc, mt, _, _ := setupTeamTest(t)
		mt.On("GetByID", mock.Anything).Return(models.Team{}, errors.New("not found"))

		_, err := svc.GetTeam(uuid.New())
		assert.Error(t, err)
	})

	t.Run("list projects fail", func(t *testing.T) {
		svc, mt, _, ms := setupTeamTest(t)
		mp := new(mockProjectRepository)
		ms.On("Projects").Return(mp)

		teamID := uuid.New()
		mt.On("GetByID", teamID).Return(models.Team{ID: teamID, Name: "Team A"}, nil)
		mp.On("ListByTeamID", teamID).Return([]models.Project{}, errors.New("db error"))

		_, err := svc.GetTeam(teamID)
		assert.Error(t, err)
	})
}

func TestTeamService_ListTeamsByWorkspaceID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mt, _, _ := setupTeamTest(t)
		workspaceID := uuid.New()
		mt.On("ListTeamsByWorkspaceID", workspaceID).Return([]models.Team{{Name: "T1"}}, nil)

		res, err := svc.ListTeamsByWorkspaceID(workspaceID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}

func TestTeamService_UpdateTeam(t *testing.T) {
	t.Run("success as team admin", func(t *testing.T) {
		svc, mt, _, _ := setupTeamTest(t)
		teamID := uuid.New()
		userID := uuid.New()
		workspaceID := uuid.New()

		// Access: Team Admin
		mt.On("GetMember", teamID, userID).Return(models.TeamMember{Role: models.TeamRoleAdmin}, nil)
		mt.On("GetByID", teamID).Return(models.Team{ID: teamID, Name: "Old"}, nil)
		mt.On("UpdateTeam", mock.Anything).Return(models.Team{ID: teamID, Name: "New", Description: "Description"}, nil)

		res, err := svc.UpdateTeam(userID, workspaceID, teamID, structs.ReqUpdateTeam{Name: "New", Description: "Description"})
		assert.NoError(t, err)
		assert.Equal(t, "New", res.Name)
		assert.Equal(t, "Description", res.Description)
	})

	t.Run("success as workspace admin", func(t *testing.T) {
		svc, mt, mw, _ := setupTeamTest(t)
		teamID := uuid.New()
		userID := uuid.New()
		workspaceID := uuid.New()

		// Access: Not Team Admin but Workspace Admin
		mt.On("GetMember", teamID, userID).Return(models.TeamMember{}, errors.New("not found"))
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mt.On("GetByID", teamID).Return(models.Team{ID: teamID, Name: "Old"}, nil)
		mt.On("UpdateTeam", mock.Anything).Return(models.Team{ID: teamID, Name: "New"}, nil)

		res, err := svc.UpdateTeam(userID, workspaceID, teamID, structs.ReqUpdateTeam{Name: "New"})
		assert.NoError(t, err)
		assert.Equal(t, "New", res.Name)
	})

	t.Run("get team member error and not workspace admin", func(t *testing.T) {
		svc, mt, mw, _ := setupTeamTest(t)
		teamID := uuid.New()
		userID := uuid.New()
		workspaceID := uuid.New()

		mt.On("GetMember", teamID, userID).Return(models.TeamMember{}, errors.New("db error"))
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)

		_, err := svc.UpdateTeam(userID, workspaceID, teamID, structs.ReqUpdateTeam{Name: "New"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, mt, mw, _ := setupTeamTest(t)
		tid, rid, wid := uuid.New(), uuid.New(), uuid.New()
		mt.On("GetMember", tid, rid).Return(models.TeamMember{Role: models.TeamRoleMember}, nil)
		mw.On("GetMember", wid, rid).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)
		_, err := svc.UpdateTeam(rid, wid, tid, structs.ReqUpdateTeam{Name: "F"})
		assert.Error(t, err)
	})

	t.Run("team not found", func(t *testing.T) {
		svc, mt, _, _ := setupTeamTest(t)
		teamID := uuid.New()
		userID := uuid.New()
		workspaceID := uuid.New()

		mt.On("GetMember", teamID, userID).Return(models.TeamMember{Role: models.TeamRoleAdmin}, nil)
		mt.On("GetByID", teamID).Return(models.Team{}, errors.New("not found"))

		_, err := svc.UpdateTeam(userID, workspaceID, teamID, structs.ReqUpdateTeam{Name: "New"})
		assert.Error(t, err)
	})

	t.Run("update failure", func(t *testing.T) {
		svc, mt, _, _ := setupTeamTest(t)
		teamID := uuid.New()
		userID := uuid.New()
		workspaceID := uuid.New()

		mt.On("GetMember", teamID, userID).Return(models.TeamMember{Role: models.TeamRoleAdmin}, nil)
		mt.On("GetByID", teamID).Return(models.Team{ID: teamID, Name: "Old"}, nil)
		mt.On("UpdateTeam", mock.Anything).Return(models.Team{}, errors.New("db error"))

		_, err := svc.UpdateTeam(userID, workspaceID, teamID, structs.ReqUpdateTeam{Name: "New"})
		assert.Error(t, err)
	})

}

func TestTeamService_DeleteTeam(t *testing.T) {
	t.Run("success as team admin", func(t *testing.T) {
		svc, mt, _, _ := setupTeamTest(t)
		teamID := uuid.New()
		userID := uuid.New()
		workspaceID := uuid.New()

		mt.On("GetMember", teamID, userID).Return(models.TeamMember{Role: models.TeamRoleAdmin}, nil)
		mt.On("DeleteTeam", teamID).Return(nil)

		err := svc.DeleteTeam(userID, workspaceID, teamID)
		assert.NoError(t, err)
	})

	t.Run("success as workspace admin", func(t *testing.T) {
		svc, mt, mw, _ := setupTeamTest(t)
		teamID := uuid.New()
		userID := uuid.New()
		workspaceID := uuid.New()

		// Not team admin
		mt.On("GetMember", teamID, userID).Return(models.TeamMember{}, errors.New("not found"))
		// But is Workspace Admin
		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mt.On("DeleteTeam", teamID).Return(nil)

		err := svc.DeleteTeam(userID, workspaceID, teamID)
		assert.NoError(t, err)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, mt, mw, _ := setupTeamTest(t)
		tid, rid, wid := uuid.New(), uuid.New(), uuid.New()
		mt.On("GetMember", tid, rid).Return(models.TeamMember{Role: models.TeamRoleMember}, nil)
		mw.On("GetMember", wid, rid).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)
		err := svc.DeleteTeam(rid, wid, tid)
		assert.Error(t, err)
	})

	t.Run("delete failure", func(t *testing.T) {
		svc, mt, _, _ := setupTeamTest(t)
		teamID := uuid.New()
		userID := uuid.New()
		workspaceID := uuid.New()

		mt.On("GetMember", teamID, userID).Return(models.TeamMember{Role: models.TeamRoleAdmin}, nil)
		mt.On("DeleteTeam", teamID).Return(errors.New("db error"))

		err := svc.DeleteTeam(userID, workspaceID, teamID)
		assert.Error(t, err)
	})
}

func TestTeamService_AddMember(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mt, _, ms := setupTeamTest(t)
		mu := new(mockUserRepository)
		ms.On("Users").Return(mu)

		teamID := uuid.New()
		requestorID := uuid.New()
		targetUserID := uuid.New()
		email := "target@test.com"

		mt.On("GetMember", teamID, requestorID).Return(models.TeamMember{Role: models.TeamRoleAdmin}, nil)
		mu.On("GetByEmail", email).Return(models.User{ID: targetUserID}, nil)
		mt.On("GetMember", teamID, targetUserID).Return(models.TeamMember{}, errors.New("not found"))
		mt.On("AddMember", mock.Anything).Return(nil)

		err := svc.AddMember(requestorID, uuid.New(), teamID, email, "MEMBER")
		assert.NoError(t, err)
	})

	t.Run("success as workspace admin", func(t *testing.T) {
		svc, mt, mw, ms := setupTeamTest(t)
		mu := new(mockUserRepository)
		ms.On("Users").Return(mu)

		teamID := uuid.New()
		workspaceID := uuid.New()
		requestorID := uuid.New()
		targetUserID := uuid.New()
		email := "target@test.com"

		// 1. Not team admin
		mt.On("GetMember", teamID, requestorID).Return(models.TeamMember{}, errors.New("not found"))
		// 2. Is Workspace Admin
		mw.On("GetMember", workspaceID, requestorID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)

		mu.On("GetByEmail", email).Return(models.User{ID: targetUserID}, nil)
		mt.On("GetMember", teamID, targetUserID).Return(models.TeamMember{}, errors.New("not found"))
		mt.On("AddMember", mock.Anything).Return(nil)

		err := svc.AddMember(requestorID, workspaceID, teamID, email, "MEMBER")
		assert.NoError(t, err)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, mt, mw, ms := setupTeamTest(t)
		mu := new(mockUserRepository)
		ms.On("Users").Return(mu)
		tid, rid, wid := uuid.New(), uuid.New(), uuid.New()
		mt.On("GetMember", tid, rid).Return(models.TeamMember{Role: models.TeamRoleMember}, nil)
		mw.On("GetMember", wid, rid).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)
		err := svc.AddMember(rid, wid, tid, "target@test.com", "MEMBER")
		assert.Error(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		svc, mt, _, ms := setupTeamTest(t)
		mu := new(mockUserRepository)
		ms.On("Users").Return(mu)

		teamID := uuid.New()
		mt.On("GetMember", teamID, mock.Anything).Return(models.TeamMember{Role: models.TeamRoleAdmin}, nil)
		mu.On("GetByEmail", mock.Anything).Return(models.User{}, errors.New("not found"))

		err := svc.AddMember(uuid.New(), uuid.New(), teamID, "missing@test.com", "MEMBER")
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("already in team", func(t *testing.T) {
		svc, mt, _, ms := setupTeamTest(t)
		mu := new(mockUserRepository)
		ms.On("Users").Return(mu)
		tid, rid, tid2 := uuid.New(), uuid.New(), uuid.New()
		mt.On("GetMember", tid, rid).Return(models.TeamMember{Role: models.TeamRoleAdmin}, nil)
		mu.On("GetByEmail", mock.Anything).Return(models.User{ID: tid2}, nil)
		mt.On("GetMember", tid, tid2).Return(models.TeamMember{}, nil)
		err := svc.AddMember(rid, uuid.New(), tid, "in@t.com", "MEMBER")
		assert.Error(t, err)
	})

	t.Run("invalid role", func(t *testing.T) {
		svc, mt, _, ms := setupTeamTest(t)
		mu := new(mockUserRepository)
		ms.On("Users").Return(mu)
		tid, rid, tid2 := uuid.New(), uuid.New(), uuid.New()
		mt.On("GetMember", tid, rid).Return(models.TeamMember{Role: models.TeamRoleAdmin}, nil)
		mu.On("GetByEmail", mock.Anything).Return(models.User{ID: tid2}, nil)
		mt.On("GetMember", tid, tid2).Return(models.TeamMember{}, errors.New("none"))
		err := svc.AddMember(rid, uuid.New(), tid, "x@t.com", "SUPER_ADMIN")
		assert.Error(t, err)
	})

}

func TestTeamService_RemoveMember(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mt, _, _ := setupTeamTest(t)
		teamID := uuid.New()
		requestorID := uuid.New()
		targetUserID := uuid.New()

		mt.On("GetMember", teamID, requestorID).Return(models.TeamMember{Role: models.TeamRoleAdmin}, nil)
		mt.On("RemoveMember", teamID, targetUserID).Return(nil)

		err := svc.RemoveMember(requestorID, uuid.New(), teamID, targetUserID)
		assert.NoError(t, err)
	})

	t.Run("success as workspace admin", func(t *testing.T) {
		svc, mt, mw, _ := setupTeamTest(t)
		teamID := uuid.New()
		workspaceID := uuid.New()
		requestorID := uuid.New()
		targetUserID := uuid.New()

		// 1. Not team admin
		mt.On("GetMember", teamID, requestorID).Return(models.TeamMember{}, errors.New("not found"))
		// 2. Is Workspace Admin
		mw.On("GetMember", workspaceID, requestorID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mt.On("RemoveMember", teamID, targetUserID).Return(nil)

		err := svc.RemoveMember(requestorID, workspaceID, teamID, targetUserID)
		assert.NoError(t, err)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc, mt, mw, _ := setupTeamTest(t)
		tid, rid, wid, uid := uuid.New(), uuid.New(), uuid.New(), uuid.New()
		mt.On("GetMember", tid, rid).Return(models.TeamMember{Role: models.TeamRoleMember}, nil)
		mw.On("GetMember", wid, rid).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)
		err := svc.RemoveMember(rid, wid, tid, uid)
		assert.Error(t, err)
	})

	t.Run("remove member failure", func(t *testing.T) {
		svc, mt, _, _ := setupTeamTest(t)
		teamID := uuid.New()
		requestorID := uuid.New()
		targetUserID := uuid.New()

		mt.On("GetMember", teamID, requestorID).Return(models.TeamMember{Role: models.TeamRoleAdmin}, nil)
		mt.On("RemoveMember", teamID, targetUserID).Return(errors.New("db error"))

		err := svc.RemoveMember(requestorID, uuid.New(), teamID, targetUserID)
		assert.Error(t, err)
	})

}

func TestTeamService_ListMembersByTeamID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mt, _, _ := setupTeamTest(t)
		teamID := uuid.New()
		mt.On("ListMembersByTeamID", teamID).Return([]models.TeamMemberWithUser{{Email: "u1@test.com"}}, nil)

		res, err := svc.ListMembersByTeamID(teamID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}
