package models

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
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
	ID        uuid.UUID     `json:"id" db:"id"`
	Name      string        `json:"name" db:"name"`
	Type      WorkspaceType `json:"type" db:"type"`
	OwnerID   uuid.UUID     `json:"owner_id" db:"owner_id"`
	CreatedAt string        `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt string        `json:"updated_at,omitempty" db:"updated_at"`
}

type WorkspaceMember struct {
	WorkspaceID uuid.UUID     `json:"workspace_id" db:"workspace_id"`
	UserID      uuid.UUID     `json:"user_id" db:"user_id"`
	Role        WorkspaceRole `json:"role" db:"role"`
	CreatedAt   string        `json:"created_at" db:"created_at"`
	UpdatedAt   string        `json:"updated_at" db:"updated_at"`
}

type WorkspaceMemberWithUser struct {
	WorkspaceID uuid.UUID     `db:"workspace_id"`
	UserID      uuid.UUID     `db:"user_id"`
	Role        WorkspaceRole `db:"role"`
	Email       string        `db:"email"`
	FullName    string        `db:"full_name"`
}

type WorkspaceModel struct {
	db *goqu.Database
}

func InitWorkspaceModel(goqu *goqu.Database) (WorkspaceModel, error) {
	return WorkspaceModel{
		db: goqu,
	}, nil
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

// GetWorkspacesByUserID returns all workspaces a user is a member of
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

// GetMembers returns all members of a workspace with user details
func (model *WorkspaceModel) ListMembersByWorkspaceId(workspaceID uuid.UUID) ([]WorkspaceMemberWithUser, error) {
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
func (model *WorkspaceModel) CreateWorkspace(transaction *goqu.TxDatabase, ws Workspace) (Workspace, error) {
	var createdWs Workspace

	found, err := transaction.Insert(WorkspaceTable).Rows(
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
func (model *WorkspaceModel) AddMemberTx(transaction *goqu.TxDatabase, member WorkspaceMember) error {
	dataset := transaction.Insert(WorkspaceMemberTable)
	return model.executeAddMember(dataset, member)
}

func (model *WorkspaceModel) AddMember(member WorkspaceMember) error {
	dataset := model.db.Insert(WorkspaceMemberTable)
	return model.executeAddMember(dataset, member)
}

func (model *WorkspaceModel) executeAddMember(dataset *goqu.InsertDataset, member WorkspaceMember) error {
	_, err := dataset.Rows(
		goqu.Record{
			"workspace_id": member.WorkspaceID,
			"user_id":      member.UserID,
			"role":         member.Role,
		},
	).Executor().Exec()

	return err
}

// RemoveMember removes a user from a workspace
func (model *WorkspaceModel) RemoveMember(transaction *goqu.TxDatabase, workspaceID, userID uuid.UUID) error {
	if transaction != nil {
		_, err := transaction.Delete(WorkspaceMemberTable).
			Where(goqu.Ex{
				"workspace_id": workspaceID,
				"user_id":      userID,
			}).Executor().Exec()
		return err
	}

	_, err := model.db.Delete(WorkspaceMemberTable).
		Where(goqu.Ex{
			"workspace_id": workspaceID,
			"user_id":      userID,
		}).Executor().Exec()
	return err
}
