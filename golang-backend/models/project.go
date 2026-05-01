package models

import (
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

const ProjectTable = "projects"
const ProjectMemberTable = "project_members"
const ProjectTeamTable = "project_teams"

type ProjectRole string

const (
	ProjectRoleAdmin  ProjectRole = "ADMIN"
	ProjectRoleMember ProjectRole = "MEMBER"
	ProjectRoleViewer ProjectRole = "VIEWER"
)

type Project struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	IsArchived  bool      `json:"is_archived" db:"is_archived"`
	WorkspaceID uuid.UUID `json:"workspace_id" db:"workspace_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ProjectMember struct {
	ProjectID uuid.UUID   `json:"project_id" db:"project_id"`
	UserID    uuid.UUID   `json:"user_id" db:"user_id"`
	Role      ProjectRole `json:"role" db:"role"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
}

type ProjectMemberWithUser struct {
	ProjectID uuid.UUID   `db:"project_id"`
	UserID    uuid.UUID   `db:"user_id"`
	Role      ProjectRole `db:"role"`
	Email     string      `db:"email"`
	FullName  string      `db:"full_name"`
}

type ProjectTeam struct {
	ProjectID uuid.UUID `json:"project_id" db:"project_id"`
	TeamID    uuid.UUID `json:"team_id" db:"team_id"`
}

type ProjectTeamWithDetails struct {
	ProjectID uuid.UUID `db:"project_id"`
	TeamID    uuid.UUID `db:"team_id"`
	TeamName  string    `db:"name"`
}

type ProjectModel struct {
	db *goqu.Database
}

func InitProjectModel(goqu *goqu.Database) (ProjectModel, error) {
	return ProjectModel{
		db: goqu,
	}, nil
}

func (model *ProjectModel) Create(transaction *goqu.TxDatabase, project Project) (Project, error) {
	var createdProject Project
	found, err := transaction.Insert(ProjectTable).Rows(
		goqu.Record{
			"name":         project.Name,
			"description":  project.Description,
			"workspace_id": project.WorkspaceID,
			"is_archived":  project.IsArchived,
		},
	).Returning("*").Executor().ScanStruct(&createdProject)

	if err != nil {
		return project, err
	}
	if !found {
		return project, sql.ErrNoRows
	}
	return createdProject, nil
}

func (model *ProjectModel) GetByID(id uuid.UUID) (Project, error) {
	project := Project{}
	found, err := model.db.From(ProjectTable).Where(goqu.Ex{"id": id}).ScanStruct(&project)
	if err != nil {
		return project, err
	}
	if !found {
		return project, sql.ErrNoRows
	}
	return project, nil
}

func (model *ProjectModel) ListByWorkspaceID(workspaceID uuid.UUID) ([]Project, error) {
	var projects []Project
	err := model.db.From(ProjectTable).
		Where(goqu.Ex{"workspace_id": workspaceID, "is_archived": false}).
		ScanStructs(&projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (model *ProjectModel) Update(transaction *goqu.TxDatabase, project Project) (Project, error) {
	_, err := transaction.Update(ProjectTable).
		Set(goqu.Record{
			"name":        project.Name,
			"description": project.Description,
			"is_archived": project.IsArchived,
		}).
		Where(goqu.Ex{"id": project.ID}).
		Executor().Exec()

	if err != nil {
		return project, err
	}
	return project, nil
}

// Members
func (model *ProjectModel) AddMember(transaction *goqu.TxDatabase, member ProjectMember) error {
	_, err := transaction.Insert(ProjectMemberTable).Rows(
		goqu.Record{
			"project_id": member.ProjectID,
			"user_id":    member.UserID,
			"role":       member.Role,
		},
	).Executor().Exec()
	return err
}

func (model *ProjectModel) RemoveMember(transaction *goqu.TxDatabase, projectID, userID uuid.UUID) error {
	_, err := transaction.Delete(ProjectMemberTable).
		Where(goqu.Ex{
			"project_id": projectID,
			"user_id":    userID,
		}).Executor().Exec()
	return err
}

func (model *ProjectModel) GetMember(projectID, userID uuid.UUID) (ProjectMember, error) {
	member := ProjectMember{}
	found, err := model.db.From(ProjectMemberTable).
		Where(goqu.Ex{
			"project_id": projectID,
			"user_id":    userID,
		}).ScanStruct(&member)

	if err != nil {
		return member, err
	}
	if !found {
		return member, sql.ErrNoRows
	}
	return member, nil
}

func (model *ProjectModel) GetMembers(projectID uuid.UUID) ([]ProjectMemberWithUser, error) {
	var members []ProjectMemberWithUser
	err := model.db.From(ProjectMemberTable).
		Join(goqu.T(UserTable), goqu.On(goqu.Ex{ProjectMemberTable + ".user_id": goqu.I(UserTable + ".id")})).
		Where(goqu.Ex{ProjectMemberTable + ".project_id": projectID}).
		Select(
			ProjectMemberTable+".project_id",
			ProjectMemberTable+".user_id",
			ProjectMemberTable+".role",
			UserTable+".email",
			UserTable+".full_name",
		).ScanStructs(&members)

	if err != nil {
		return nil, err
	}
	return members, nil
}

// Teams
func (model *ProjectModel) AddTeam(transaction *goqu.TxDatabase, team ProjectTeam) error {
	_, err := transaction.Insert(ProjectTeamTable).Rows(
		goqu.Record{
			"project_id": team.ProjectID,
			"team_id":    team.TeamID,
		},
	).Executor().Exec()
	return err
}

func (model *ProjectModel) RemoveTeam(transaction *goqu.TxDatabase, projectID, teamID uuid.UUID) error {
	_, err := transaction.Delete(ProjectTeamTable).
		Where(goqu.Ex{
			"project_id": projectID,
			"team_id":    teamID,
		}).Executor().Exec()
	return err
}

func (model *ProjectModel) GetTeam(projectID, teamID uuid.UUID) (ProjectTeam, error) {
	team := ProjectTeam{}
	found, err := model.db.From(ProjectTeamTable).
		Where(goqu.Ex{
			"project_id": projectID,
			"team_id":    teamID,
		}).ScanStruct(&team)

	if err != nil {
		return team, err
	}
	if !found {
		return team, sql.ErrNoRows
	}
	return team, nil
}

func (model *ProjectModel) GetTeams(projectID uuid.UUID) ([]ProjectTeamWithDetails, error) {
	var teams []ProjectTeamWithDetails
	err := model.db.From(ProjectTeamTable).
		Join(goqu.T(TeamTable), goqu.On(goqu.Ex{ProjectTeamTable + ".team_id": goqu.I(TeamTable + ".id")})).
		Where(goqu.Ex{ProjectTeamTable + ".project_id": projectID}).
		Select(
			ProjectTeamTable+".project_id",
			ProjectTeamTable+".team_id",
			TeamTable+".name",
		).ScanStructs(&teams)

	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (model *ProjectModel) ListByTeamID(teamID uuid.UUID) ([]Project, error) {
	var projects []Project
	err := model.db.From(ProjectTable).
		Join(goqu.T(ProjectTeamTable), goqu.On(goqu.Ex{ProjectTable + ".id": goqu.I(ProjectTeamTable + ".project_id")})).
		Select(ProjectTable + ".*").
		Where(goqu.Ex{
			ProjectTeamTable + ".team_id": teamID,
			ProjectTable + ".is_archived": false,
		}).
		ScanStructs(&projects)

	if err != nil {
		return nil, err
	}
	return projects, nil
}
