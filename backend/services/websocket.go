package services

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// WebsocketService defines the interface for websocket-related business logic
type WebsocketService interface {
	GetConnectionTopics(userID, workspaceID uuid.UUID) ([]string, error)
}

type websocketService struct {
	workspaceService WorkspaceService
	projectService   ProjectService
	logger           *zap.Logger
}

// NewWebsocketService creates a new websocket service instance
func NewWebsocketService(
	workspaceService WorkspaceService,
	projectService ProjectService,
	logger *zap.Logger,
) WebsocketService {
	return &websocketService{
		workspaceService: workspaceService,
		projectService:   projectService,
		logger:           logger,
	}
}

func (s *websocketService) GetConnectionTopics(userID, workspaceID uuid.UUID) ([]string, error) {
	// 1. Verify Workspace Access
	_, err := s.workspaceService.ValidateAccess(userID, workspaceID)
	if err != nil {
		return nil, err
	}

	var topics []string
	// 2. Add Workspace Topic
	topics = append(topics, fmt.Sprintf("workspace:%s", workspaceID.String()))

	// 3. Get all projects in workspace and check access
	projects, err := s.projectService.ListByWorkspaceID(workspaceID)
	if err == nil {
		for _, p := range projects {
			if s.checkProjectAccess(p.ID, userID) == nil {
				topics = append(topics, fmt.Sprintf("project:%s", p.ID.String()))
			}
		}
	}

	return topics, nil
}

func (s *websocketService) checkProjectAccess(projectID, userID uuid.UUID) error {
	// Reuse logic from projectService if possible, or keep it here if it's specific to websocket sub logic
	// Actually, ProjectService has ValidateProjectAccess
	return s.projectService.ValidateProjectAccess(projectID, userID, false)
}
