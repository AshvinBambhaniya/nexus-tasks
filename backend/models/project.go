package models

import (
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

// ProjectTable is the name of the projects table.
const ProjectTable = "projects"

// ProjectMemberTable is the name of the project members table.
const ProjectMemberTable = "project_members"

// ProjectTeamTable is the name of the project teams table.
const ProjectTeamTable = "project_teams"

// ProjectRole defines the role of a user in a project.
type ProjectRole string

const (
	// ProjectRoleAdmin is the admin role in a project.
	ProjectRoleAdmin ProjectRole = "ADMIN"
	// ProjectRoleMember is the member role in a project.
	ProjectRoleMember ProjectRole = "MEMBER"
	// ProjectRoleViewer is the viewer role in a project.
	ProjectRoleViewer ProjectRole = "VIEWER"
)

// Project represents a project in a workspace.
type Project struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	IsArchived  bool      `json:"is_archived" db:"is_archived"`
	WorkspaceID uuid.UUID `json:"workspace_id" db:"workspace_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ProjectMember represents a member of a project.
type ProjectMember struct {
	ProjectID uuid.UUID   `json:"project_id" db:"project_id"`
	UserID    uuid.UUID   `json:"user_id" db:"user_id"`
	Role      ProjectRole `json:"role" db:"role"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
}

// ProjectMemberWithUser represents a project member with user details.
type ProjectMemberWithUser struct {
	ProjectID uuid.UUID   `db:"project_id"`
	UserID    uuid.UUID   `db:"user_id"`
	Role      ProjectRole `db:"role"`
	Email     string      `db:"email"`
	FullName  string      `db:"full_name"`
}

// ProjectTeam represents a team associated with a project.
type ProjectTeam struct {
	ProjectID uuid.UUID `json:"project_id" db:"project_id"`
	TeamID    uuid.UUID `json:"team_id" db:"team_id"`
}

// ProjectTeamWithDetails represents a project team with team details.
type ProjectTeamWithDetails struct {
	ProjectID uuid.UUID `db:"project_id"`
	TeamID    uuid.UUID `db:"team_id"`
	TeamName  string    `db:"name"`
}

// ProjectRepository defines the interface for project data access.
type ProjectRepository interface {
	Create(project Project) (Project, error)
	GetByID(id uuid.UUID) (Project, error)
	ListByWorkspaceID(workspaceID uuid.UUID) ([]Project, error)
	Update(project Project) (Project, error)

	// Members
	AddMember(member ProjectMember) error
	RemoveMember(projectID, userID uuid.UUID) error
	GetMember(projectID, userID uuid.UUID) (ProjectMember, error)
	GetMembers(projectID uuid.UUID) ([]ProjectMemberWithUser, error)

	// Teams
	AddTeam(team ProjectTeam) error
	RemoveTeam(projectID, teamID uuid.UUID) error
	GetTeam(projectID, teamID uuid.UUID) (ProjectTeam, error)
	GetTeams(projectID uuid.UUID) ([]ProjectTeamWithDetails, error)
	ListByTeamID(teamID uuid.UUID) ([]Project, error)
}

// ProjectModel is the implementation of ProjectRepository.
type ProjectModel struct {
	db DbExecutor
}

// InitProjectModel initializes a new ProjectModel.
func InitProjectModel(db DbExecutor) ProjectRepository {
	return &ProjectModel{
		db: db,
	}
}

// Create creates a new project.
func (model *ProjectModel) Create(project Project) (Project, error) {
	var createdProject Project
	found, err := model.db.Insert(ProjectTable).Rows(
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

// GetByID returns a project by its ID.
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

// ListByWorkspaceID lists all projects in a workspace.
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

// Update updates an existing project.
func (model *ProjectModel) Update(project Project) (Project, error) {
	_, err := model.db.Update(ProjectTable).
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

// AddMember adds a member to a project.
func (model *ProjectModel) AddMember(member ProjectMember) error {
	_, err := model.db.Insert(ProjectMemberTable).Rows(
		goqu.Record{
			"project_id": member.ProjectID,
			"user_id":    member.UserID,
			"role":       member.Role,
		},
	).Executor().Exec()
	return err
}

// RemoveMember removes a member from a project.
func (model *ProjectModel) RemoveMember(projectID, userID uuid.UUID) error {
	_, err := model.db.Delete(ProjectMemberTable).
		Where(goqu.Ex{
			"project_id": projectID,
			"user_id":    userID,
		}).Executor().Exec()
	return err
}

// GetMember returns a project member by project ID and user ID.
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

// GetMembers returns all members of a project.
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

// AddTeam adds a team to a project.
func (model *ProjectModel) AddTeam(team ProjectTeam) error {
	_, err := model.db.Insert(ProjectTeamTable).Rows(
		goqu.Record{
			"project_id": team.ProjectID,
			"team_id":    team.TeamID,
		},
	).Executor().Exec()
	return err
}

// RemoveTeam removes a team from a project.
func (model *ProjectModel) RemoveTeam(projectID, teamID uuid.UUID) error {
	_, err := model.db.Delete(ProjectTeamTable).
		Where(goqu.Ex{
			"project_id": projectID,
			"team_id":    teamID,
		}).Executor().Exec()
	return err
}

// GetTeam returns a project team association.
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

// GetTeams returns all teams associated with a project.
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

// ListByTeamID returns all projects associated with a specific team.
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
