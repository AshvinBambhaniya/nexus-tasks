package structs

import (
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/google/uuid"
)

// ResUser defines the response payload for a user.
//
// swagger:model ResUser
type ResUser struct {
	// Unique user identifier.
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID uuid.UUID `json:"id"`

	// User's email address.
	// example: user@example.com
	Email string `json:"email"`

	// User's display name.
	// example: Jane Doe
	FullName string `json:"full_name"`

	// Whether the user account is active.
	// example: true
	IsActive bool `json:"is_active"`
}

// ResWorkspace defines the response payload for a workspace.
//
// swagger:model ResWorkspace
type ResWorkspace struct {
	// Unique workspace identifier.
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID uuid.UUID `json:"id"`

	// Workspace display name.
	// example: Acme Corp
	Name string `json:"name"`

	// Workspace type (e.g. PERSONAL, TEAM).
	// example: TEAM
	Type models.WorkspaceType `json:"type"`

	// UUID of the workspace owner.
	// example: 550e8400-e29b-41d4-a716-446655440000
	OwnerID uuid.UUID `json:"owner_id"`
}

// ResWorkspaceMember defines the response payload for a workspace member.
//
// swagger:model ResWorkspaceMember
type ResWorkspaceMember struct {
	// UUID of the workspace.
	// example: 550e8400-e29b-41d4-a716-446655440000
	WorkspaceID uuid.UUID `json:"workspace_id"`

	// UUID of the member user.
	// example: 550e8400-e29b-41d4-a716-446655440001
	UserID uuid.UUID `json:"user_id"`

	// Member's role in the workspace.
	// example: MEMBER
	Role models.WorkspaceRole `json:"role"`

	// Expanded user details.
	User ResUser `json:"user"`
}

// ResTeam defines the response payload for a team.
//
// swagger:model ResTeam
type ResTeam struct {
	// Unique team identifier.
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID uuid.UUID `json:"id"`

	// Team name.
	// example: Backend Squad
	Name string `json:"name"`

	// Team description.
	// example: Responsible for all server-side services.
	Description string `json:"description"`

	// UUID of the workspace this team belongs to.
	// example: 550e8400-e29b-41d4-a716-446655440000
	WorkspaceID uuid.UUID `json:"workspace_id"`
}

// ResTeamWithProjects defines the response payload for a team with its associated projects.
//
// swagger:model ResTeamWithProjects
type ResTeamWithProjects struct {
	ResTeam

	// Projects associated with this team.
	Projects []models.Project `json:"projects"`
}

// ResTeamMember defines the response payload for a team member.
//
// swagger:model ResTeamMember
type ResTeamMember struct {
	// UUID of the team.
	// example: 550e8400-e29b-41d4-a716-446655440000
	TeamID uuid.UUID `json:"team_id"`

	// UUID of the member user.
	// example: 550e8400-e29b-41d4-a716-446655440001
	UserID uuid.UUID `json:"user_id"`

	// Member's role in the team.
	// example: MEMBER
	Role models.TeamRole `json:"role"`

	// Expanded user details.
	User ResUser `json:"user"`
}

// ResProject defines the response payload for a project.
//
// swagger:model ResProject
type ResProject struct {
	// Unique project identifier.
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID uuid.UUID `json:"id"`

	// Project name.
	// example: Website Redesign
	Name string `json:"name"`

	// Project description.
	// example: Full redesign of the marketing website.
	Description string `json:"description"`

	// Whether the project is archived.
	// example: false
	IsArchived bool `json:"is_archived"`

	// UUID of the workspace this project belongs to.
	// example: 550e8400-e29b-41d4-a716-446655440000
	WorkspaceID uuid.UUID `json:"workspace_id"`

	// Timestamp when the project was created.
	// example: 2025-01-15T10:00:00Z
	CreatedAt time.Time `json:"created_at"`
}

// ResProjectMember defines the response payload for a project member.
//
// swagger:model ResProjectMember
type ResProjectMember struct {
	// UUID of the member user.
	// example: 550e8400-e29b-41d4-a716-446655440001
	UserID uuid.UUID `json:"user_id"`

	// Member's full name.
	// example: Jane Doe
	FullName string `json:"full_name"`

	// Member's email address.
	// example: jane@example.com
	Email string `json:"email"`

	// Member's role in the project.
	// example: ADMIN
	Role models.ProjectRole `json:"role"`

	// Whether membership was directly assigned (vs. inherited from a team).
	// example: true
	IsDirect bool `json:"is_direct"`
}

// ResProjectTeam defines the response payload for a project–team association.
//
// swagger:model ResProjectTeam
type ResProjectTeam struct {
	// UUID of the project.
	// example: 550e8400-e29b-41d4-a716-446655440000
	ProjectID uuid.UUID `json:"project_id"`

	// UUID of the team.
	// example: 550e8400-e29b-41d4-a716-446655440002
	TeamID uuid.UUID `json:"team_id"`

	// Name of the team.
	// example: Backend Squad
	TeamName string `json:"team_name"`
}

// ResComment defines the response payload for a task comment.
//
// swagger:model ResComment
type ResComment struct {
	// Unique comment identifier.
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID uuid.UUID `json:"id"`

	// Comment body text.
	// example: LGTM, merging this.
	Content string `json:"content"`

	// UUID of the task this comment belongs to.
	// example: 550e8400-e29b-41d4-a716-446655440001
	TaskID uuid.UUID `json:"task_id"`

	// UUID of the comment author.
	// example: 550e8400-e29b-41d4-a716-446655440002
	AuthorID uuid.UUID `json:"author_id"`

	// Timestamp when the comment was created.
	// example: 2025-06-01T09:00:00Z
	CreatedAt time.Time `json:"created_at"`

	// Timestamp when the comment was last updated.
	// example: 2025-06-01T09:05:00Z
	UpdatedAt time.Time `json:"updated_at"`

	// Expanded author details.
	Author ResUser `json:"author"`
}

// ResTask defines the response payload for a task.
//
// swagger:model ResTask
type ResTask struct {
	// Unique task identifier.
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID uuid.UUID `json:"id"`

	// Human-readable task number within the project.
	// example: 42
	Number int `json:"number"`

	// Task title.
	// example: Implement OAuth2 login
	Title string `json:"title"`

	// Task description.
	// example: Integrate OAuth2 with the existing auth service.
	Description string `json:"description"`

	// Current task status.
	// example: IN_PROGRESS
	Status models.TaskStatus `json:"status"`

	// Task priority.
	// example: P1
	Priority models.TaskPriority `json:"priority"`

	// UUID of the project this task belongs to.
	// example: 550e8400-e29b-41d4-a716-446655440001
	ProjectID uuid.UUID `json:"project_id"`

	// UUID of the assignee. Null if unassigned.
	// example: 550e8400-e29b-41d4-a716-446655440002
	AssigneeID *uuid.UUID `json:"assignee_id"`

	// UUID of the task author. Null if unknown.
	// example: 550e8400-e29b-41d4-a716-446655440003
	AuthorID *uuid.UUID `json:"author_id"`

	// Task due date. Null if not set.
	// example: 2025-12-31T23:59:59Z
	DueDate *time.Time `json:"due_date"`

	// Estimated effort in minutes. Null if not set.
	// example: 120
	EstimatedMinutes *int `json:"estimated_minutes"`

	// Timestamp when the task was completed. Null if not yet completed.
	// example: 2025-11-30T17:00:00Z
	CompletedAt *time.Time `json:"completed_at"`

	// Timestamp when the task was created.
	// example: 2025-10-01T08:00:00Z
	CreatedAt time.Time `json:"created_at"`

	// Timestamp when the task was last updated.
	// example: 2025-11-01T12:00:00Z
	UpdatedAt time.Time `json:"updated_at"`

	// Total number of comments on this task.
	// example: 5
	CommentCount int `json:"comment_count"`

	// Expanded assignee details. Null if unassigned.
	Assignee *ResUser `json:"assignee"`

	// Expanded author details. Null if unknown.
	Author *ResUser `json:"author"`
}

// ResTaskWithProject defines the response payload for a task with its parent project.
//
// swagger:model ResTaskWithProject
type ResTaskWithProject struct {
	ResTask

	// The project this task belongs to.
	Project ResProject `json:"project"`
}

// ResTimeEntry defines the response payload for a time entry.
//
// swagger:model ResTimeEntry
type ResTimeEntry struct {
	// Unique time entry identifier.
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID uuid.UUID `json:"id"`

	// UUID of the task this entry is associated with.
	// example: 550e8400-e29b-41d4-a716-446655440001
	TaskID uuid.UUID `json:"task_id"`

	// UUID of the user who logged the time.
	// example: 550e8400-e29b-41d4-a716-446655440002
	UserID uuid.UUID `json:"user_id"`

	// Full name of the user who logged the time.
	// example: Jane Doe
	UserFullName string `json:"user_full_name"`

	// Description of work performed.
	// example: Reviewed PR and added test cases.
	Description string `json:"description"`

	// Timer start time (RFC3339).
	// example: 2025-11-01T09:00:00Z
	StartTime time.Time `json:"start_time"`

	// Timer end time. Null if the timer is still running.
	// example: 2025-11-01T10:30:00Z
	EndTime *time.Time `json:"end_time"`

	// Total duration in minutes. Null if the timer is still running.
	// example: 90
	DurationMinutes *int `json:"duration_minutes"`

	// Whether this entry was manually logged (vs. started/stopped via timer).
	// example: false
	IsManual bool `json:"is_manual"`

	// Timestamp when the entry record was created.
	// example: 2025-11-01T10:30:05Z
	CreatedAt time.Time `json:"created_at"`
}

// ResTimeEntryWithTask defines the response payload for a time entry with task context.
//
// swagger:model ResTimeEntryWithTask
type ResTimeEntryWithTask struct {
	ResTimeEntry

	// Title of the associated task.
	// example: Implement OAuth2 login
	TaskTitle string `json:"task_title"`

	// Human-readable number of the associated task.
	// example: 42
	TaskNumber int `json:"task_number"`
}

// ResActiveTimer defines the response for the currently running timer.
//
// swagger:model ResActiveTimer
type ResActiveTimer struct {
	// Unique time entry identifier for the active timer.
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID uuid.UUID `json:"id"`

	// UUID of the task being timed.
	// example: 550e8400-e29b-41d4-a716-446655440001
	TaskID uuid.UUID `json:"task_id"`

	// Title of the task being timed.
	// example: Implement OAuth2 login
	TaskTitle string `json:"task_title"`

	// Human-readable number of the task being timed.
	// example: 42
	TaskNumber int `json:"task_number"`

	// Time when the timer was started (RFC3339).
	// example: 2025-11-01T09:00:00Z
	StartTime time.Time `json:"start_time"`
}

// ResTaskTimeEntries defines the aggregated response for time entries on a task.
//
// swagger:model ResTaskTimeEntries
type ResTaskTimeEntries struct {
	// List of individual time entries.
	Entries []ResTimeEntry `json:"entries"`

	// Total logged time in minutes across all entries.
	// example: 240
	TotalLoggedMinutes int `json:"total_logged_minutes"`

	// Task's estimated effort in minutes. Null if not set.
	// example: 300
	EstimatedMinutes *int `json:"estimated_minutes"`
}
