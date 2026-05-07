package models

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

// WorkspaceTable is the name of the workspaces table
const WorkspaceTable = "workspaces"

// WorkspaceMemberTable is the name of the workspace members table
const WorkspaceMemberTable = "workspace_members"

// WorkspaceType represents the type of a workspace
type WorkspaceType string

const (
	// WorkspaceTypePersonal is a personal workspace
	WorkspaceTypePersonal WorkspaceType = "PERSONAL"
	// WorkspaceTypeTeam is a team workspace
	WorkspaceTypeTeam WorkspaceType = "TEAM"
)

// WorkspaceRole represents the role of a user in a workspace
type WorkspaceRole string

const (
	// WorkspaceRoleAdmin is the administrator role
	WorkspaceRoleAdmin WorkspaceRole = "ADMIN"
	// WorkspaceRoleMember is the standard member role
	WorkspaceRoleMember WorkspaceRole = "MEMBER"
	// WorkspaceRoleViewer is the read-only viewer role
	WorkspaceRoleViewer WorkspaceRole = "VIEWER"
)

// Workspace represents a workspace record
type Workspace struct {
	ID        uuid.UUID     `json:"id" db:"id"`
	Name      string        `json:"name" db:"name"`
	Type      WorkspaceType `json:"type" db:"type"`
	OwnerID   uuid.UUID     `json:"owner_id" db:"owner_id"`
	CreatedAt string        `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt string        `json:"updated_at,omitempty" db:"updated_at"`
}

// WorkspaceMember represents a membership record in a workspace
type WorkspaceMember struct {
	WorkspaceID uuid.UUID     `json:"workspace_id" db:"workspace_id"`
	UserID      uuid.UUID     `json:"user_id" db:"user_id"`
	Role        WorkspaceRole `json:"role" db:"role"`
	CreatedAt   string        `json:"created_at" db:"created_at"`
	UpdatedAt   string        `json:"updated_at" db:"updated_at"`
}

// WorkspaceMemberWithUser represents a workspace member with joined user details
type WorkspaceMemberWithUser struct {
	WorkspaceID uuid.UUID     `db:"workspace_id"`
	UserID      uuid.UUID     `db:"user_id"`
	Role        WorkspaceRole `db:"role"`
	Email       string        `db:"email"`
	FullName    string        `db:"full_name"`
}

// WorkspaceRepository defines the interface for workspace-related database operations
type WorkspaceRepository interface {
	GetByID(id uuid.UUID) (Workspace, error)
	ListWorkspacesByUserID(userID uuid.UUID) ([]Workspace, error)
	ListMembersByWorkspaceID(workspaceID uuid.UUID) ([]WorkspaceMemberWithUser, error)
	GetMember(workspaceID, userID uuid.UUID) (WorkspaceMember, error)
	CreateWorkspace(ws Workspace) (Workspace, error)
	AddMember(member WorkspaceMember) error
	RemoveMember(workspaceID, userID uuid.UUID) error
}

// WorkspaceModel implements the WorkspaceRepository interface
type WorkspaceModel struct {
	db DbExecutor
}

// InitWorkspaceModel initializes the workspace model
func InitWorkspaceModel(db DbExecutor) WorkspaceRepository {
	return &WorkspaceModel{
		db: db,
	}
}

// GetByID gets a workspace by ID
func (model *WorkspaceModel) GetByID(id uuid.UUID) (Workspace, error) {
	ws := Workspace{}
	found, err := model.db.From(WorkspaceTable).Where(goqu.Ex{
		"id": id,
	}).ScanStruct(&ws)

	if err != nil {
		return ws, err
	}
	if !found {
		return ws, sql.ErrNoRows
	}
	return ws, nil
}

// ListWorkspacesByUserID returns all workspaces a user is a member of
func (model *WorkspaceModel) ListWorkspacesByUserID(userID uuid.UUID) ([]Workspace, error) {
	var workspaces []Workspace
	// Join with workspace_members
	err := model.db.From(WorkspaceTable).
		Join(goqu.T(WorkspaceMemberTable), goqu.On(goqu.Ex{WorkspaceTable + ".id": goqu.I(WorkspaceMemberTable + ".workspace_id")})).
		Where(goqu.Ex{WorkspaceMemberTable + ".user_id": userID}).
		Select(WorkspaceTable + ".*").
		ScanStructs(&workspaces)

	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

// ListMembersByWorkspaceID returns all members of a workspace with user details
func (model *WorkspaceModel) ListMembersByWorkspaceID(workspaceID uuid.UUID) ([]WorkspaceMemberWithUser, error) {
	var members []WorkspaceMemberWithUser
	err := model.db.From(WorkspaceMemberTable).
		Join(goqu.T(UserTable), goqu.On(goqu.Ex{WorkspaceMemberTable + ".user_id": goqu.I(UserTable + ".id")})).
		Where(goqu.Ex{WorkspaceMemberTable + ".workspace_id": workspaceID}).
		Select(
			WorkspaceMemberTable+".workspace_id",
			WorkspaceMemberTable+".user_id",
			WorkspaceMemberTable+".role",
			UserTable+".email",
			UserTable+".full_name",
		).
		ScanStructs(&members)

	if err != nil {
		return nil, err
	}
	return members, nil
}

// GetMember returns a specific member record (to check access/role)
func (model *WorkspaceModel) GetMember(workspaceID, userID uuid.UUID) (WorkspaceMember, error) {
	member := WorkspaceMember{}
	found, err := model.db.From(WorkspaceMemberTable).
		Where(goqu.Ex{
			"workspace_id": workspaceID,
			"user_id":      userID,
		}).ScanStruct(&member)

	if err != nil {
		return member, err
	}
	if !found {
		return member, sql.ErrNoRows
	}
	return member, nil
}

// CreateWorkspace inserts a workspace
func (model *WorkspaceModel) CreateWorkspace(ws Workspace) (Workspace, error) {
	var createdWs Workspace

	found, err := model.db.Insert(WorkspaceTable).Rows(
		goqu.Record{
			"name":     ws.Name,
			"type":     ws.Type,
			"owner_id": ws.OwnerID,
		},
	).Returning("*").Executor().ScanStruct(&createdWs)

	if err != nil {
		return ws, err
	}
	if !found {
		return ws, sql.ErrNoRows
	}

	return createdWs, nil
}

// AddMember adds a user to a workspace
func (model *WorkspaceModel) AddMember(member WorkspaceMember) error {
	_, err := model.db.Insert(WorkspaceMemberTable).Rows(
		goqu.Record{
			"workspace_id": member.WorkspaceID,
			"user_id":      member.UserID,
			"role":         member.Role,
		},
	).Executor().Exec()

	return err
}

// RemoveMember removes a user from a workspace
func (model *WorkspaceModel) RemoveMember(workspaceID, userID uuid.UUID) error {
	_, err := model.db.Delete(WorkspaceMemberTable).
		Where(goqu.Ex{
			"workspace_id": workspaceID,
			"user_id":      userID,
		}).Executor().Exec()
	return err
}
