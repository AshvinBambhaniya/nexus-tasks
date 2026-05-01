package services

import (
	"errors"
	"testing"

	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupWebsocketTest(t *testing.T) (*websocketService, *mockWorkspaceService, *mockProjectService) {
	mockWsSvc := new(mockWorkspaceService)
	mockProjSvc := new(mockProjectService)
	logger := zap.NewNop()

	svc := NewWebsocketService(mockWsSvc, mockProjSvc, logger).(*websocketService)

	return svc, mockWsSvc, mockProjSvc
}

func TestWebsocketService_GetConnectionTopics(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockWsSvc, mockProjSvc := setupWebsocketTest(t)

		userID := uuid.New()
		workspaceID := uuid.New()
		projectID1 := uuid.New()
		projectID2 := uuid.New()

		// Validate Access
		mockWsSvc.On("ValidateAccess", userID, workspaceID).Return(models.WorkspaceMember{}, nil)

		// List Projects
		mockProjSvc.On("ListByWorkspaceID", workspaceID).Return([]models.Project{
			{ID: projectID1},
			{ID: projectID2},
		}, nil)

		// Check Project Access for each
		mockProjSvc.On("ValidateProjectAccess", projectID1, userID, false).Return(nil)
		mockProjSvc.On("ValidateProjectAccess", projectID2, userID, false).Return(nil)

		topics, err := svc.GetConnectionTopics(userID, workspaceID)

		assert.NoError(t, err)
		assert.Len(t, topics, 3) // 1 workspace + 2 projects
		assert.Contains(t, topics, "workspace:"+workspaceID.String())
		assert.Contains(t, topics, "project:"+projectID1.String())
		assert.Contains(t, topics, "project:"+projectID2.String())
	})

	t.Run("workspace access denied", func(t *testing.T) {
		svc, mockWsSvc, _ := setupWebsocketTest(t)
		mockWsSvc.On("ValidateAccess", mock.Anything, mock.Anything).Return(models.WorkspaceMember{}, errors.New("denied"))

		_, err := svc.GetConnectionTopics(uuid.New(), uuid.New())
		assert.Error(t, err)
	})

	t.Run("partial project access", func(t *testing.T) {
		svc, mockWsSvc, mockProjSvc := setupWebsocketTest(t)
		uid, wid, pid1, pid2 := uuid.New(), uuid.New(), uuid.New(), uuid.New()

		mockWsSvc.On("ValidateAccess", uid, wid).Return(models.WorkspaceMember{}, nil)
		mockProjSvc.On("ListByWorkspaceID", wid).Return([]models.Project{{ID: pid1}, {ID: pid2}}, nil)

		// Access to pid1, but not pid2
		mockProjSvc.On("ValidateProjectAccess", pid1, uid, false).Return(nil)
		mockProjSvc.On("ValidateProjectAccess", pid2, uid, false).Return(errors.New("denied"))

		topics, _ := svc.GetConnectionTopics(uid, wid)
		assert.Len(t, topics, 2) // WS + Project1
		assert.Contains(t, topics, "project:"+pid1.String())
	})
}
