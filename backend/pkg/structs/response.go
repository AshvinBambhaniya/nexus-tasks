package structs

import (
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/google/uuid"
)

// ResUser defines the response payload for a user.
type ResUser struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	IsActive bool      `json:"is_active"`
}

// ResWorkspace defines the response payload for a workspace.
type ResWorkspace struct {
	ID      uuid.UUID            `json:"id"`
	Name    string               `json:"name"`
	Type    models.WorkspaceType `json:"type"`
	OwnerID uuid.UUID            `json:"owner_id"`
}

// ResWorkspaceMember defines the response payload for a workspace member.
type ResWorkspaceMember struct {
	WorkspaceID uuid.UUID            `json:"workspace_id"`
	UserID      uuid.UUID            `json:"user_id"`
	Role        models.WorkspaceRole `json:"role"`
	User        ResUser              `json:"user"`
}

// ResTeam defines the response payload for a team.
type ResTeam struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
}

// ResTeamWithProjects defines the response payload for a team with its projects.
type ResTeamWithProjects struct {
	ResTeam
	Projects []models.Project `json:"projects"`
}

// ResTeamMember defines the response payload for a team member.
type ResTeamMember struct {
	TeamID uuid.UUID       `json:"team_id"`
	UserID uuid.UUID       `json:"user_id"`
	Role   models.TeamRole `json:"role"`
	User   ResUser         `json:"user"`
}

// ResProject defines the response payload for a project.
type ResProject struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsArchived  bool      `json:"is_archived"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// ResProjectMember defines the response payload for a project member.
type ResProjectMember struct {
	UserID   uuid.UUID          `json:"user_id"`
	FullName string             `json:"full_name"`
	Email    string             `json:"email"`
	Role     models.ProjectRole `json:"role"`
	IsDirect bool               `json:"is_direct"`
}

// ResProjectTeam defines the response payload for a project team.
type ResProjectTeam struct {
	ProjectID uuid.UUID `json:"project_id"`
	TeamID    uuid.UUID `json:"team_id"`
	TeamName  string    `json:"team_name"`
}

// ResComment defines the response payload for a comment.
type ResComment struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	TaskID    uuid.UUID `json:"task_id"`
	AuthorID  uuid.UUID `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Author    ResUser   `json:"author"`
}

// ResTask defines the response payload for a task.
type ResTask struct {
	ID           uuid.UUID           `json:"id"`
	Number       int                 `json:"number"`
	Title        string              `json:"title"`
	Description  string              `json:"description"`
	Status       models.TaskStatus   `json:"status"`
	Priority     models.TaskPriority `json:"priority"`
	ProjectID    uuid.UUID           `json:"project_id"`
	AssigneeID   *uuid.UUID          `json:"assignee_id"`
	AuthorID     *uuid.UUID          `json:"author_id"`
	DueDate      *time.Time          `json:"due_date"`
	CompletedAt  *time.Time          `json:"completed_at"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	CommentCount int                 `json:"comment_count"` // Computed
	Assignee     *ResUser            `json:"assignee"`      // Expanded if needed
	Author       *ResUser            `json:"author"`        // Expanded if needed
}

// ResTaskWithProject defines the response payload for a task with its project.
type ResTaskWithProject struct {
	ResTask
	Project ResProject `json:"project"`
}
