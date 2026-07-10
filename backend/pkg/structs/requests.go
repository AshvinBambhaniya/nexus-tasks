package structs

import (
	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/google/uuid"
)

// ReqRegisterUser defines the request payload for user registration.
type ReqRegisterUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
}

// ReqLoginUser defines the request payload for user login.
type ReqLoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// ReqCreateWorkspace defines the request payload for creating a workspace.
type ReqCreateWorkspace struct {
	Name string `json:"name" validate:"required"`
}

// ReqCreateTeam defines the request payload for creating a team.
type ReqCreateTeam struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// ReqInviteWorkspaceMember defines the request payload for inviting a member to a workspace.
type ReqInviteWorkspaceMember struct {
	Email string `json:"email" validate:"required,email"`
}

// ReqUpdateTeam defines the request payload for updating a team.
type ReqUpdateTeam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ReqAddTeamMember defines the request payload for adding a member to a team.
type ReqAddTeamMember struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required"` // Should be validated against ENUM
}

// ReqCreateProject defines the request payload for creating a project.
type ReqCreateProject struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// ReqUpdateProject defines the request payload for updating a project.
type ReqUpdateProject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsArchived  *bool  `json:"is_archived"`
}

// ReqAddProjectMember defines the request payload for adding a member to a project.
type ReqAddProjectMember struct {
	Email string             `json:"email" validate:"required,email"`
	Role  models.ProjectRole `json:"role"` // Optional, default MEMBER
}

// ReqAddProjectTeam defines the request payload for adding a team to a project.
type ReqAddProjectTeam struct {
	TeamID uuid.UUID `json:"team_id" validate:"required"`
}

// ReqCreateComment defines the request payload for creating a comment.
type ReqCreateComment struct {
	Content          string      `json:"content" validate:"required"`
	MentionedUserIDs []uuid.UUID `json:"mentioned_user_ids,omitempty"`
}

// ReqCreateTask defines the request payload for creating a task.
type ReqCreateTask struct {
	Title       string              `json:"title" validate:"required"`
	Description string              `json:"description"`
	Status      models.TaskStatus   `json:"status"`   // Optional, default TODO
	Priority    models.TaskPriority `json:"priority"` // Optional, default P2
	DueDate     *CustomTime         `json:"due_date"`
	AssigneeID  *uuid.UUID          `json:"assignee_id"`
}

// ReqUpdateTask defines the request payload for updating a task.
type ReqUpdateTask struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Status      models.TaskStatus   `json:"status"`
	Priority    models.TaskPriority `json:"priority"`
	DueDate     *CustomTime         `json:"due_date"`
	AssigneeID  *uuid.UUID          `json:"assignee_id"`
}

// ReqDraftTask defines the request payload for drafting a task via AI.
type ReqDraftTask struct {
	Title string `json:"title" validate:"required"`
}
