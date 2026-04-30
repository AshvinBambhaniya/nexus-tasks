package structs

import (
	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/google/uuid"
)

type ReqRegisterUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
}

type ReqLoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ReqCreateWorkspace struct {
	Name string `json:"name" validate:"required"`
}

type ReqCreateTeam struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type ReqInviteWorkspaceMember struct {
	Email string `json:"email" validate:"required,email"`
}

type ReqUpdateTeam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ReqAddTeamMember struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required"` // Should be validated against ENUM
}

type ReqCreateProject struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type ReqUpdateProject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsArchived  *bool  `json:"is_archived"`
}

type ReqAddProjectMember struct {
	Email string             `json:"email" validate:"required,email"`
	Role  models.ProjectRole `json:"role"` // Optional, default MEMBER
}

type ReqAddProjectTeam struct {
	TeamID uuid.UUID `json:"team_id" validate:"required"`
}

type ReqCreateComment struct {
	Content string `json:"content" validate:"required"`
}

type ReqCreateTask struct {
	Title       string              `json:"title" validate:"required"`
	Description string              `json:"description"`
	Status      models.TaskStatus   `json:"status"`   // Optional, default TODO
	Priority    models.TaskPriority `json:"priority"` // Optional, default P2
	DueDate     *CustomTime         `json:"due_date"`
	AssigneeID  *uuid.UUID          `json:"assignee_id"`
}

type ReqUpdateTask struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Status      models.TaskStatus   `json:"status"`
	Priority    models.TaskPriority `json:"priority"`
	DueDate     *CustomTime         `json:"due_date"`
	AssigneeID  *uuid.UUID          `json:"assignee_id"`
}
