package structs

import (
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
