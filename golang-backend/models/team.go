package models

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

const TeamTable = "teams"
const TeamMemberTable = "team_members"

type TeamRole string

const (
	TeamRoleAdmin  TeamRole = "ADMIN"
	TeamRoleMember TeamRole = "MEMBER"
)

type Team struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	WorkspaceID uuid.UUID `json:"workspace_id" db:"workspace_id"`
	CreatedAt   string    `json:"created_at" db:"created_at"`
	UpdatedAt   string    `json:"updated_at" db:"updated_at"`
}

type TeamMember struct {
	TeamID    uuid.UUID `json:"team_id" db:"team_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Role      TeamRole  `json:"role" db:"role"`
	CreatedAt string    `json:"created_at" db:"created_at"`
	UpdatedAt string    `json:"updated_at" db:"updated_at"`
}

type TeamMemberWithUser struct {
	TeamID   uuid.UUID `json:"team_id" db:"team_id"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Role     TeamRole  `json:"role" db:"role"`
	Email    string    `json:"email" db:"email"`
	FullName string    `json:"full_name" db:"full_name"`
}

type TeamRepository interface {
	GetByID(id uuid.UUID) (Team, error)
	ListTeamsByWorkspaceID(workspaceID uuid.UUID) ([]Team, error)
	CreateTeam(team Team) (Team, error)
	AddMember(member TeamMember) error
	ListMembersByTeamId(teamID uuid.UUID) ([]TeamMemberWithUser, error)
	UpdateTeam(team Team) (Team, error)
	DeleteTeam(teamID uuid.UUID) error
	RemoveMember(teamID, userID uuid.UUID) error
	GetMember(teamID, userID uuid.UUID) (TeamMember, error)
}

type TeamModel struct {
	db DbExecutor
}

func InitTeamModel(db DbExecutor) TeamRepository {
	return &TeamModel{
		db: db,
	}
}

func (model *TeamModel) GetByID(id uuid.UUID) (Team, error) {
	team := Team{}
	found, err := model.db.From(TeamTable).Where(goqu.Ex{"id": id}).ScanStruct(&team)
	if err != nil {
		return team, err
	}
	if !found {
		return team, sql.ErrNoRows
	}
	return team, nil
}

func (model *TeamModel) ListTeamsByWorkspaceID(workspaceID uuid.UUID) ([]Team, error) {
	var teams []Team
	err := model.db.From(TeamTable).
		Where(goqu.Ex{"workspace_id": workspaceID}).
		ScanStructs(&teams)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (model *TeamModel) CreateTeam(team Team) (Team, error) {
	var createdTeam Team
	found, err := model.db.Insert(TeamTable).Rows(
		goqu.Record{
			"name":         team.Name,
			"description":  team.Description,
			"workspace_id": team.WorkspaceID,
		},
	).Returning("*").Executor().ScanStruct(&createdTeam)

	if err != nil {
		return team, err
	}
	if !found {
		return team, sql.ErrNoRows
	}
	return createdTeam, nil
}

func (model *TeamModel) AddMember(member TeamMember) error {
	_, err := model.db.Insert(TeamMemberTable).Rows(
		goqu.Record{
			"team_id": member.TeamID,
			"user_id": member.UserID,
			"role":    member.Role,
		},
	).Executor().Exec()
	return err
}

func (model *TeamModel) ListMembersByTeamId(teamID uuid.UUID) ([]TeamMemberWithUser, error) {
	var members []TeamMemberWithUser
	err := model.db.From(TeamMemberTable).
		Join(goqu.T(UserTable), goqu.On(goqu.Ex{TeamMemberTable + ".user_id": goqu.I(UserTable + ".id")})).
		Where(goqu.Ex{TeamMemberTable + ".team_id": teamID}).
		Select(
			TeamMemberTable+".team_id",
			TeamMemberTable+".user_id",
			TeamMemberTable+".role",
			UserTable+".email",
			UserTable+".full_name",
		).ScanStructs(&members)

	if err != nil {
		return nil, err
	}
	return members, nil
}

func (model *TeamModel) UpdateTeam(team Team) (Team, error) {
	_, err := model.db.Update(TeamTable).
		Set(goqu.Record{
			"name":        team.Name,
			"description": team.Description,
		}).
		Where(goqu.Ex{"id": team.ID}).
		Executor().Exec()

	if err != nil {
		return team, err
	}
	return team, nil
}

func (model *TeamModel) DeleteTeam(teamID uuid.UUID) error {
	_, err := model.db.Delete(TeamTable).
		Where(goqu.Ex{"id": teamID}).
		Executor().Exec()
	return err
}

func (model *TeamModel) RemoveMember(teamID, userID uuid.UUID) error {
	_, err := model.db.Delete(TeamMemberTable).
		Where(goqu.Ex{
			"team_id": teamID,
			"user_id": userID,
		}).Executor().Exec()
	return err
}

func (model *TeamModel) GetMember(teamID, userID uuid.UUID) (TeamMember, error) {
	member := TeamMember{}
	found, err := model.db.From(TeamMemberTable).
		Where(goqu.Ex{
			"team_id": teamID,
			"user_id": userID,
		}).ScanStruct(&member)

	if err != nil {
		return member, err
	}
	if !found {
		return member, sql.ErrNoRows
	}
	return member, nil
}
