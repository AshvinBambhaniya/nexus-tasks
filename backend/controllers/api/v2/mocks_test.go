package v2

import (
	"context"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/structs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func getMockArg[T any](args mock.Arguments, index int) T {
	val := args.Get(index)
	switch v := val.(type) {
	case T:
		return v
	default:
		var zero T
		return zero
	}
}

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) Register(email, password, fullName string) (models.User, string, error) {
	args := m.Called(email, password, fullName)
	return getMockArg[models.User](args, 0), args.String(1), args.Error(2)
}

func (m *mockUserService) Authenticate(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func (m *mockUserService) GetByID(userID uuid.UUID) (models.User, error) {
	args := m.Called(userID)
	return getMockArg[models.User](args, 0), args.Error(1)
}

type mockWorkspaceService struct {
	mock.Mock
}

func (m *mockWorkspaceService) CreateWorkspace(ownerID uuid.UUID, req structs.ReqCreateWorkspace) (models.Workspace, error) {
	args := m.Called(ownerID, req)
	return getMockArg[models.Workspace](args, 0), args.Error(1)
}

func (m *mockWorkspaceService) ListWorkspacesByUserID(userID uuid.UUID) ([]models.Workspace, error) {
	args := m.Called(userID)
	return getMockArg[[]models.Workspace](args, 0), args.Error(1)
}

func (m *mockWorkspaceService) ListMembersByWorkspaceID(workspaceID uuid.UUID) ([]models.WorkspaceMemberWithUser, error) {
	args := m.Called(workspaceID)
	return getMockArg[[]models.WorkspaceMemberWithUser](args, 0), args.Error(1)
}

func (m *mockWorkspaceService) InviteMember(requestorID, workspaceID uuid.UUID, email string) error {
	args := m.Called(requestorID, workspaceID, email)
	return args.Error(0)
}

func (m *mockWorkspaceService) RemoveMember(requestorID, workspaceID, userID uuid.UUID) error {
	args := m.Called(requestorID, workspaceID, userID)
	return args.Error(0)
}

func (m *mockWorkspaceService) ValidateAccess(userID, workspaceID uuid.UUID) (models.WorkspaceMember, error) {
	args := m.Called(userID, workspaceID)
	return getMockArg[models.WorkspaceMember](args, 0), args.Error(1)
}

type mockProjectService struct {
	mock.Mock
}

func (m *mockProjectService) CreateProject(userID, workspaceID uuid.UUID, req structs.ReqCreateProject) (models.Project, error) {
	args := m.Called(userID, workspaceID, req)
	return getMockArg[models.Project](args, 0), args.Error(1)
}

func (m *mockProjectService) GetProject(userID, projectID uuid.UUID) (models.Project, error) {
	args := m.Called(userID, projectID)
	return getMockArg[models.Project](args, 0), args.Error(1)
}

func (m *mockProjectService) UpdateProject(userID, projectID uuid.UUID, req structs.ReqUpdateProject) (models.Project, error) {
	args := m.Called(userID, projectID, req)
	return getMockArg[models.Project](args, 0), args.Error(1)
}

func (m *mockProjectService) AddMember(userID, projectID uuid.UUID, req structs.ReqAddProjectMember) (models.ProjectMember, error) {
	args := m.Called(userID, projectID, req)
	return getMockArg[models.ProjectMember](args, 0), args.Error(1)
}

func (m *mockProjectService) RemoveMember(userID, projectID, targetUserID uuid.UUID) error {
	args := m.Called(userID, projectID, targetUserID)
	return args.Error(0)
}

func (m *mockProjectService) ListMembers(userID, projectID uuid.UUID) ([]structs.ResProjectMember, error) {
	args := m.Called(userID, projectID)
	return getMockArg[[]structs.ResProjectMember](args, 0), args.Error(1)
}

func (m *mockProjectService) ListByWorkspaceID(workspaceID uuid.UUID) ([]models.Project, error) {
	args := m.Called(workspaceID)
	return getMockArg[[]models.Project](args, 0), args.Error(1)
}

func (m *mockProjectService) AddTeam(userID, projectID, teamID uuid.UUID) (structs.ResProjectTeam, error) {
	args := m.Called(userID, projectID, teamID)
	return getMockArg[structs.ResProjectTeam](args, 0), args.Error(1)
}

func (m *mockProjectService) RemoveTeam(userID, projectID, teamID uuid.UUID) error {
	args := m.Called(userID, projectID, teamID)
	return args.Error(0)
}

func (m *mockProjectService) ListTeams(userID, projectID uuid.UUID) ([]structs.ResProjectTeam, error) {
	args := m.Called(userID, projectID)
	return getMockArg[[]structs.ResProjectTeam](args, 0), args.Error(1)
}

func (m *mockProjectService) ValidateProjectAccess(projectID, userID uuid.UUID, requireAdmin bool) error {
	args := m.Called(projectID, userID, requireAdmin)
	return args.Error(0)
}

type mockTeamService struct {
	mock.Mock
}

func (m *mockTeamService) CreateTeam(userID uuid.UUID, workspaceID uuid.UUID, req structs.ReqCreateTeam) (models.Team, error) {
	args := m.Called(userID, workspaceID, req)
	return getMockArg[models.Team](args, 0), args.Error(1)
}

func (m *mockTeamService) GetTeam(teamID uuid.UUID) (structs.ResTeamWithProjects, error) {
	args := m.Called(teamID)
	return getMockArg[structs.ResTeamWithProjects](args, 0), args.Error(1)
}

func (m *mockTeamService) ListTeamsByWorkspaceID(workspaceID uuid.UUID) ([]models.Team, error) {
	args := m.Called(workspaceID)
	return getMockArg[[]models.Team](args, 0), args.Error(1)
}

func (m *mockTeamService) UpdateTeam(requestorID, workspaceID, teamID uuid.UUID, req structs.ReqUpdateTeam) (models.Team, error) {
	args := m.Called(requestorID, workspaceID, teamID, req)
	return getMockArg[models.Team](args, 0), args.Error(1)
}

func (m *mockTeamService) DeleteTeam(requestorID, workspaceID, teamID uuid.UUID) error {
	args := m.Called(requestorID, workspaceID, teamID)
	return args.Error(0)
}

func (m *mockTeamService) AddMember(requestorID, workspaceID, teamID uuid.UUID, email, role string) error {
	args := m.Called(requestorID, workspaceID, teamID, email, role)
	return args.Error(0)
}

func (m *mockTeamService) RemoveMember(requestorID, workspaceID, teamID, userID uuid.UUID) error {
	args := m.Called(requestorID, workspaceID, teamID, userID)
	return args.Error(0)
}

func (m *mockTeamService) ListMembersByTeamID(teamID uuid.UUID) ([]models.TeamMemberWithUser, error) {
	args := m.Called(teamID)
	return getMockArg[[]models.TeamMemberWithUser](args, 0), args.Error(1)
}

type mockTaskService struct {
	mock.Mock
}

func (m *mockTaskService) CreateTask(userID, projectID uuid.UUID, req structs.ReqCreateTask) (models.Task, error) {
	args := m.Called(userID, projectID, req)
	return getMockArg[models.Task](args, 0), args.Error(1)
}

func (m *mockTaskService) ListProjectTasks(userID, projectID uuid.UUID, status *models.TaskStatus, assigneeID *uuid.UUID) ([]models.TaskWithAssignee, error) {
	args := m.Called(userID, projectID, status, assigneeID)
	return getMockArg[[]models.TaskWithAssignee](args, 0), args.Error(1)
}

func (m *mockTaskService) GetTask(userID, taskID uuid.UUID) (models.TaskWithAssignee, error) {
	args := m.Called(userID, taskID)
	return getMockArg[models.TaskWithAssignee](args, 0), args.Error(1)
}

func (m *mockTaskService) UpdateTask(userID, taskID uuid.UUID, req structs.ReqUpdateTask) (models.Task, error) {
	args := m.Called(userID, taskID, req)
	return getMockArg[models.Task](args, 0), args.Error(1)
}

func (m *mockTaskService) DeleteTask(userID, taskID uuid.UUID) error {
	args := m.Called(userID, taskID)
	return args.Error(0)
}

func (m *mockTaskService) ListMyTasks(userID uuid.UUID) ([]models.TaskWithAssignee, error) {
	args := m.Called(userID)
	return getMockArg[[]models.TaskWithAssignee](args, 0), args.Error(1)
}

func (m *mockTaskService) ListCompletedTasksInLastDays(userID, projectID uuid.UUID, days int) ([]models.TaskWithAssignee, error) {
	args := m.Called(userID, projectID, days)
	return getMockArg[[]models.TaskWithAssignee](args, 0), args.Error(1)
}

type mockCommentService struct {
	mock.Mock
}

func (m *mockCommentService) CreateComment(userID, taskID uuid.UUID, req structs.ReqCreateComment) (models.Comment, error) {
	args := m.Called(userID, taskID, req)
	return getMockArg[models.Comment](args, 0), args.Error(1)
}

func (m *mockCommentService) ListTaskComments(userID, taskID uuid.UUID) ([]models.CommentWithAuthor, error) {
	args := m.Called(userID, taskID)
	return getMockArg[[]models.CommentWithAuthor](args, 0), args.Error(1)
}

func (m *mockCommentService) DeleteComment(userID, commentID uuid.UUID) error {
	args := m.Called(userID, commentID)
	return args.Error(0)
}

func (m *mockCommentService) ListCommentsForTasks(userID, projectID uuid.UUID, taskIDs []uuid.UUID) ([]models.CommentWithAuthor, error) {
	args := m.Called(userID, projectID, taskIDs)
	return getMockArg[[]models.CommentWithAuthor](args, 0), args.Error(1)
}

type mockHealthService struct {
	mock.Mock
}

func (m *mockHealthService) CheckDatabaseHealth(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
