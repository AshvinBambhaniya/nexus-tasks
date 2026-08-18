package structs

import (
	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/google/uuid"
)

// ReqRegisterUser defines the request payload for user registration.
//
// swagger:model ReqRegisterUser
type ReqRegisterUser struct {
	// User's email address.
	// Required: true
	// example: user@example.com
	Email string `json:"email" validate:"required,email"`

	// User's password. Minimum 6 characters.
	// Required: true
	// min length: 6
	// example: s3cr3tPassw0rd
	Password string `json:"password" validate:"required,min=6"`

	// User's display name.
	// Required: true
	// example: Jane Doe
	FullName string `json:"full_name" validate:"required"`
}

// ReqLoginUser defines the request payload for user login.
//
// swagger:model ReqLoginUser
type ReqLoginUser struct {
	// User's email address.
	// Required: true
	// example: user@example.com
	Email string `json:"email" validate:"required,email"`

	// User's account password.
	// Required: true
	// example: s3cr3tPassw0rd
	Password string `json:"password" validate:"required"`
}

// ReqCreateWorkspace defines the request payload for creating a workspace.
//
// swagger:model ReqCreateWorkspace
type ReqCreateWorkspace struct {
	// Workspace display name.
	// Required: true
	// example: Acme Corp
	Name string `json:"name" validate:"required"`
}

// ReqCreateTeam defines the request payload for creating a team.
//
// swagger:model ReqCreateTeam
type ReqCreateTeam struct {
	// Team name.
	// Required: true
	// example: Backend Squad
	Name string `json:"name" validate:"required"`

	// Optional team description.
	// example: Responsible for all server-side services.
	Description string `json:"description"`
}

// ReqInviteWorkspaceMember defines the request payload for inviting a member to a workspace.
//
// swagger:model ReqInviteWorkspaceMember
type ReqInviteWorkspaceMember struct {
	// Email of the user to invite.
	// Required: true
	// example: colleague@example.com
	Email string `json:"email" validate:"required,email"`
}

// ReqUpdateTeam defines the request payload for updating a team.
//
// swagger:model ReqUpdateTeam
type ReqUpdateTeam struct {
	// New team name.
	// example: Platform Team
	Name string `json:"name"`

	// New team description.
	// example: Handles platform infrastructure.
	Description string `json:"description"`
}

// ReqAddTeamMember defines the request payload for adding a member to a team.
//
// swagger:model ReqAddTeamMember
type ReqAddTeamMember struct {
	// Email of the user to add.
	// Required: true
	// example: dev@example.com
	Email string `json:"email" validate:"required,email"`

	// Role to assign to the member. Must be a valid team role enum value.
	// Required: true
	// example: MEMBER
	Role string `json:"role" validate:"required"`
}

// ReqCreateProject defines the request payload for creating a project.
//
// swagger:model ReqCreateProject
type ReqCreateProject struct {
	// Project name.
	// Required: true
	// example: Website Redesign
	Name string `json:"name" validate:"required"`

	// Optional project description.
	// example: Full redesign of the marketing website.
	Description string `json:"description"`
}

// ReqUpdateProject defines the request payload for updating a project.
//
// swagger:model ReqUpdateProject
type ReqUpdateProject struct {
	// New project name.
	// example: Website Redesign v2
	Name string `json:"name"`

	// New project description.
	// example: Updated scope including mobile.
	Description string `json:"description"`

	// Whether to archive this project.
	// example: false
	IsArchived *bool `json:"is_archived"`
}

// ReqAddProjectMember defines the request payload for adding a member to a project.
//
// swagger:model ReqAddProjectMember
type ReqAddProjectMember struct {
	// Email of the user to add.
	// Required: true
	// example: dev@example.com
	Email string `json:"email" validate:"required,email"`

	// Role to assign. Defaults to MEMBER if omitted.
	// example: MEMBER
	Role models.ProjectRole `json:"role"`
}

// ReqAddProjectTeam defines the request payload for adding a team to a project.
//
// swagger:model ReqAddProjectTeam
type ReqAddProjectTeam struct {
	// UUID of the team to associate with the project.
	// Required: true
	// example: 550e8400-e29b-41d4-a716-446655440000
	TeamID uuid.UUID `json:"team_id" validate:"required"`
}

// ReqCreateComment defines the request payload for creating a comment.
//
// swagger:model ReqCreateComment
type ReqCreateComment struct {
	// Comment body text.
	// Required: true
	// example: LGTM, merging this.
	Content string `json:"content" validate:"required"`

	// List of UUIDs of users mentioned in the comment.
	// example: ["550e8400-e29b-41d4-a716-446655440000"]
	MentionedUserIDs []uuid.UUID `json:"mentioned_user_ids,omitempty"`
}

// ReqCreateTask defines the request payload for creating a task.
//
// swagger:model ReqCreateTask
type ReqCreateTask struct {
	// Task title.
	// Required: true
	// example: Implement OAuth2 login
	Title string `json:"title" validate:"required"`

	// Optional task description.
	// example: Integrate OAuth2 with the existing auth service.
	Description string `json:"description"`

	// Task status. Defaults to TODO if omitted.
	// example: TODO
	Status models.TaskStatus `json:"status"`

	// Task priority. Defaults to P2 if omitted.
	// example: P1
	Priority models.TaskPriority `json:"priority"`

	// Optional due date in RFC3339 format.
	// example: "2025-12-31T23:59:59Z"
	DueDate *CustomTime `json:"due_date"`

	// Optional UUID of the assignee.
	// example: 550e8400-e29b-41d4-a716-446655440000
	AssigneeID *uuid.UUID `json:"assignee_id"`

	// Estimated effort in minutes.
	// example: 120
	EstimatedMinutes *int `json:"estimated_minutes"`
}

// ReqUpdateTask defines the request payload for updating a task.
//
// swagger:model ReqUpdateTask
type ReqUpdateTask struct {
	// Updated task title.
	// example: Implement OAuth2 login (revised)
	Title string `json:"title"`

	// Updated task description.
	// example: Added PKCE support.
	Description string `json:"description"`

	// Updated task status.
	// example: IN_PROGRESS
	Status models.TaskStatus `json:"status"`

	// Updated task priority.
	// example: P2
	Priority models.TaskPriority `json:"priority"`

	// Updated due date in RFC3339 format.
	// example: "2026-01-15T23:59:59Z"
	DueDate *CustomTime `json:"due_date"`

	// Updated assignee UUID.
	// example: 550e8400-e29b-41d4-a716-446655440000
	AssigneeID *uuid.UUID `json:"assignee_id"`

	// Updated estimated effort in minutes.
	// example: 90
	EstimatedMinutes *int `json:"estimated_minutes"`
}

// ReqDraftTask defines the request payload for drafting a task description via AI.
//
// swagger:model ReqDraftTask
type ReqDraftTask struct {
	// Task title used as the AI prompt seed.
	// Required: true
	// example: Set up CI/CD pipeline with GitHub Actions
	Title string `json:"title" validate:"required"`
}

// ReqStopTimer defines the request payload for stopping a running timer.
//
// swagger:model ReqStopTimer
type ReqStopTimer struct {
	// Optional description for the time entry.
	// example: Worked on authentication module.
	Description string `json:"description"`

	// Optional override for the recorded duration in minutes.
	// example: 45
	DurationMinutes *int `json:"duration_minutes"`
}

// ReqLogManualTime defines the request payload for manually logging time on a task.
//
// swagger:model ReqLogManualTime
type ReqLogManualTime struct {
	// Duration of work in minutes.
	// Required: true
	// minimum: 1
	// example: 60
	DurationMinutes int `json:"duration_minutes" validate:"required,min=1"`

	// Optional description of work done.
	// example: Code review and feedback.
	Description string `json:"description"`

	// Date the work was performed (YYYY-MM-DD). Defaults to today if omitted.
	// example: "2025-11-01"
	Date string `json:"date"`
}
