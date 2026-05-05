package services

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupWorkspaceTest(t *testing.T) (*workspaceService, *mockWorkspaceRepository, *mockUserRepository, *mockStorage, *mockPublisher) {
	mockWorkspaceRepo := new(mockWorkspaceRepository)
	mockUserRepo := new(mockUserRepository)
	mockStor := new(mockStorage)
	mockPub := new(mockPublisher)
	logger := zap.NewNop()

	mockStor.On("Workspaces").Return(mockWorkspaceRepo)
	mockStor.On("Users").Return(mockUserRepo)

	svc := &workspaceService{storage: mockStor, logger: logger, publisher: mockPub}

	return svc, mockWorkspaceRepo, mockUserRepo, mockStor, mockPub
}

func TestWorkspaceService_CreateWorkspace(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mw, _, _, _ := setupWorkspaceTest(t)
		ownerID, workspaceID := uuid.New(), uuid.New()

		mw.On("CreateWorkspace", mock.Anything).Return(models.Workspace{ID: workspaceID, Name: "Test"}, nil)
		mw.On("AddMember", mock.Anything).Return(nil)

		ws, err := svc.CreateWorkspace(ownerID, structs.ReqCreateWorkspace{Name: "Test"})
		assert.NoError(t, err)
		assert.Equal(t, "Test", ws.Name)
	})

	t.Run("create fail", func(t *testing.T) {
		svc, mw, _, _, _ := setupWorkspaceTest(t)
		mw.On("CreateWorkspace", mock.Anything).Return(models.Workspace{}, errors.New("fail"))
		_, err := svc.CreateWorkspace(uuid.New(), structs.ReqCreateWorkspace{Name: "X"})
		assert.Error(t, err)
	})

	t.Run("add member fail", func(t *testing.T) {
		svc, mw, _, _, _ := setupWorkspaceTest(t)
		mw.On("CreateWorkspace", mock.Anything).Return(models.Workspace{ID: uuid.New()}, nil)
		mw.On("AddMember", mock.Anything).Return(errors.New("fail"))
		_, err := svc.CreateWorkspace(uuid.New(), structs.ReqCreateWorkspace{Name: "X"})
		assert.Error(t, err)
	})
}

func TestWorkspaceService_ListWorkspacesByUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mw, _, _, _ := setupWorkspaceTest(t)
		userID := uuid.New()
		mw.On("ListWorkspacesByUserID", userID).Return([]models.Workspace{{Name: "W1"}}, nil)

		res, err := svc.ListWorkspacesByUserID(userID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}

func TestWorkspaceService_ListMembersByWorkspaceId(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mw, _, _, _ := setupWorkspaceTest(t)
		workspaceID := uuid.New()
		mw.On("ListMembersByWorkspaceId", workspaceID).Return([]models.WorkspaceMemberWithUser{{Email: "test@test.com"}}, nil)

		res, err := svc.ListMembersByWorkspaceId(workspaceID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "test@test.com", res[0].Email)
	})
}

func TestWorkspaceService_InviteMember(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mw, mu, _, mp := setupWorkspaceTest(t)
		requestorID := uuid.New()
		workspaceID := uuid.New()
		targetUserID := uuid.New()
		email := "target@test.com"

		mw.On("GetMember", workspaceID, requestorID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mu.On("GetByEmail", email).Return(models.User{ID: targetUserID}, nil)
		mw.On("GetMember", workspaceID, targetUserID).Return(models.WorkspaceMember{}, sql.ErrNoRows)
		mw.On("AddMember", mock.Anything).Return(nil)
		mw.On("GetByID", workspaceID).Return(models.Workspace{Name: "Test WS"}, nil)

		// Expect notification
		mp.On("Publish", mock.Anything, mock.Anything).Return(nil)

		err := svc.InviteMember(requestorID, workspaceID, email)
		assert.NoError(t, err)
		mp.AssertExpectations(t)
	})

	t.Run("success even if publish fails", func(t *testing.T) {
		svc, mw, mu, _, mp := setupWorkspaceTest(t)
		requestorID := uuid.New()
		workspaceID := uuid.New()
		targetUserID := uuid.New()
		email := "target@test.com"

		mw.On("GetMember", workspaceID, requestorID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mu.On("GetByEmail", email).Return(models.User{ID: targetUserID}, nil)
		mw.On("GetMember", workspaceID, targetUserID).Return(models.WorkspaceMember{}, sql.ErrNoRows)
		mw.On("AddMember", mock.Anything).Return(nil)
		mw.On("GetByID", workspaceID).Return(models.Workspace{Name: "Test WS"}, nil)

		// Expect notification failure
		mp.On("Publish", mock.Anything, mock.Anything).Return(errors.New("mq error"))

		err := svc.InviteMember(requestorID, workspaceID, email)
		assert.NoError(t, err) // Should still succeed
		mp.AssertExpectations(t)
	})

	t.Run("user already member", func(t *testing.T) {
		svc, mw, mu, _, _ := setupWorkspaceTest(t)
		workspaceID := uuid.New()
		requestorID := uuid.New()
		targetUserID := uuid.New()
		email := "already@test.com"

		mw.On("GetMember", workspaceID, requestorID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mu.On("GetByEmail", email).Return(models.User{ID: targetUserID}, nil)
		mw.On("GetMember", workspaceID, targetUserID).Return(models.WorkspaceMember{}, nil)

		err := svc.InviteMember(requestorID, workspaceID, email)
		assert.Error(t, err)
		assert.Equal(t, "user is already a member", err.Error())
	})

	t.Run("unauthorized - not admin", func(t *testing.T) {
		svc, mw, _, _, _ := setupWorkspaceTest(t)
		workspaceID := uuid.New()
		requestorID := uuid.New()

		mw.On("GetMember", workspaceID, requestorID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)

		err := svc.InviteMember(requestorID, workspaceID, "any@test.com")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only admins")
	})

	t.Run("unauthorized - not member", func(t *testing.T) {
		svc, mw, _, _, _ := setupWorkspaceTest(t)
		workspaceID := uuid.New()
		requestorID := uuid.New()

		mw.On("GetMember", workspaceID, requestorID).Return(models.WorkspaceMember{}, sql.ErrNoRows)

		err := svc.InviteMember(requestorID, workspaceID, "any@test.com")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("user not found", func(t *testing.T) {
		svc, mw, mu, _, _ := setupWorkspaceTest(t)
		wid, rid := uuid.New(), uuid.New()
		mw.On("GetMember", wid, rid).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mu.On("GetByEmail", mock.Anything).Return(models.User{}, sql.ErrNoRows)
		err := svc.InviteMember(rid, wid, "none@t.com")
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("get by email other error", func(t *testing.T) {
		svc, mw, mu, _, _ := setupWorkspaceTest(t)
		wid, rid := uuid.New(), uuid.New()
		mw.On("GetMember", wid, rid).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mu.On("GetByEmail", mock.Anything).Return(models.User{}, errors.New("db error"))
		err := svc.InviteMember(rid, wid, "none@t.com")
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("add member fail", func(t *testing.T) {
		svc, mw, mu, _, _ := setupWorkspaceTest(t)
		wid, rid, tid := uuid.New(), uuid.New(), uuid.New()
		email := "test@test.com"

		mw.On("GetMember", wid, rid).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mu.On("GetByEmail", email).Return(models.User{ID: tid}, nil)
		mw.On("GetMember", wid, tid).Return(models.WorkspaceMember{}, sql.ErrNoRows)
		mw.On("AddMember", mock.Anything).Return(errors.New("db error"))

		err := svc.InviteMember(rid, wid, email)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestWorkspaceService_RemoveMember(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mw, _, _, _ := setupWorkspaceTest(t)
		requestorID := uuid.New()
		workspaceID := uuid.New()
		targetUserID := uuid.New()

		mw.On("GetMember", workspaceID, requestorID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mw.On("RemoveMember", workspaceID, targetUserID).Return(nil)

		err := svc.RemoveMember(requestorID, workspaceID, targetUserID)
		assert.NoError(t, err)
	})

	t.Run("unauthorized - not admin", func(t *testing.T) {
		svc, mw, _, _, _ := setupWorkspaceTest(t)
		requestorID := uuid.New()
		workspaceID := uuid.New()

		mw.On("GetMember", workspaceID, requestorID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)

		err := svc.RemoveMember(requestorID, workspaceID, uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only admins")
	})

	t.Run("unauthorized - not member", func(t *testing.T) {
		svc, mw, _, _, _ := setupWorkspaceTest(t)
		wid, rid := uuid.New(), uuid.New()

		mw.On("GetMember", wid, rid).Return(models.WorkspaceMember{}, sql.ErrNoRows)

		err := svc.RemoveMember(rid, wid, uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("remove member fail", func(t *testing.T) {
		svc, mw, _, _, _ := setupWorkspaceTest(t)
		wid, rid, tid := uuid.New(), uuid.New(), uuid.New()

		mw.On("GetMember", wid, rid).Return(models.WorkspaceMember{Role: models.WorkspaceRoleAdmin}, nil)
		mw.On("RemoveMember", wid, tid).Return(errors.New("db error"))

		err := svc.RemoveMember(rid, wid, tid)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestWorkspaceService_ValidateAccess(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mw, _, _, _ := setupWorkspaceTest(t)
		userID := uuid.New()
		workspaceID := uuid.New()

		mw.On("GetMember", workspaceID, userID).Return(models.WorkspaceMember{Role: models.WorkspaceRoleMember}, nil)

		member, err := svc.ValidateAccess(userID, workspaceID)
		assert.NoError(t, err)
		assert.Equal(t, models.WorkspaceRoleMember, member.Role)
	})

	t.Run("no access", func(t *testing.T) {
		svc, mw, _, _, _ := setupWorkspaceTest(t)
		mw.On("GetMember", mock.Anything, mock.Anything).Return(models.WorkspaceMember{}, sql.ErrNoRows)

		_, err := svc.ValidateAccess(uuid.New(), uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})
}
