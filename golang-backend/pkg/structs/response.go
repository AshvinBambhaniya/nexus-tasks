package structs

import (
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/google/uuid"
)

type ResUser struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	IsActive bool      `json:"is_active"`
}

type ResWorkspace struct {
	ID      uuid.UUID            `json:"id"`
	Name    string               `json:"name"`
	Type    models.WorkspaceType `json:"type"`
	OwnerID uuid.UUID            `json:"owner_id"`
}

type ResWorkspaceMember struct {
	WorkspaceID uuid.UUID            `json:"workspace_id"`
	UserID      uuid.UUID            `json:"user_id"`
	Role        models.WorkspaceRole `json:"role"`
	User        ResUser              `json:"user"`
}

type ResTeam struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
}

type ResTeamMember struct {
	TeamID uuid.UUID       `json:"team_id"`
	UserID uuid.UUID       `json:"user_id"`
	Role   models.TeamRole `json:"role"`
	User   ResUser         `json:"user"`
}

type ResProject struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsArchived  bool      `json:"is_archived"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type ResProjectMember struct {
	UserID   uuid.UUID          `json:"user_id"`
	Email    string             `json:"email"`
	Role     models.ProjectRole `json:"role"`
	IsDirect bool               `json:"is_direct"`
}

type ResProjectTeam struct {
	ProjectID uuid.UUID `json:"project_id"`
	TeamID    uuid.UUID `json:"team_id"`
	TeamName  string    `json:"team_name"`
}
