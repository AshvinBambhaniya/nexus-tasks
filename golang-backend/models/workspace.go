package models

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

const WorkspaceTable = "workspaces"
const WorkspaceMemberTable = "workspace_members"

type WorkspaceType string

const (
	WorkspaceTypePersonal WorkspaceType = "PERSONAL"
	WorkspaceTypeTeam     WorkspaceType = "TEAM"
)

type WorkspaceRole string

const (
	WorkspaceRoleAdmin  WorkspaceRole = "ADMIN"
	WorkspaceRoleMember WorkspaceRole = "MEMBER"
	WorkspaceRoleViewer WorkspaceRole = "VIEWER"
)

type Workspace struct {
	ID      int           `json:"id" db:"id"`
	Name    string        `json:"name" db:"name"`
	Type    WorkspaceType `json:"type" db:"type"`
	OwnerID int           `json:"owner_id" db:"owner_id"`
}

type WorkspaceMember struct {
	WorkspaceID int           `json:"workspace_id" db:"workspace_id"`
	UserID      int           `json:"user_id" db:"user_id"`
	Role        WorkspaceRole `json:"role" db:"role"`
}

type WorkspaceModel struct {
	db *goqu.Database
}

func InitWorkspaceModel(goqu *goqu.Database) (WorkspaceModel, error) {
	return WorkspaceModel{
		db: goqu,
	}, nil
}

// CreateWorkspace inserts a workspace
func (model *WorkspaceModel) CreateWorkspace(transaction *goqu.TxDatabase, ws Workspace) (Workspace, error) {

	var createdWs Workspace
	found, err := transaction.Insert(WorkspaceTable).Rows(
		goqu.Record{
			"name":     ws.Name,
			"type":     ws.Type,
			"owner_id": ws.OwnerID,
		},
	).Executor().ScanStruct(&createdWs)

	if err != nil {
		return ws, err
	}
	if !found {
		return ws, sql.ErrNoRows
	}

	return createdWs, nil
}

// AddMember adds a user to a workspace
func (model *WorkspaceModel) AddMember(transaction *goqu.TxDatabase, member WorkspaceMember) error {

	_, err := transaction.Insert(WorkspaceMemberTable).Rows(
		goqu.Record{
			"workspace_id": member.WorkspaceID,
			"user_id":      member.UserID,
			"role":         member.Role,
		},
	).Executor().Exec()

	return err
}
