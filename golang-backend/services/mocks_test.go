package services

import (
	"context"

	"github.com/AshvinBambhaniya/nexus-tasks/cli/workers"
	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// --- Storage Mock ---

type mockStorage struct {
	mock.Mock
}

func (m *mockStorage) Users() models.UserRepository {
	args := m.Called()
	return args.Get(0).(models.UserRepository)
}

func (m *mockStorage) Workspaces() models.WorkspaceRepository {
	args := m.Called()
	return args.Get(0).(models.WorkspaceRepository)
}

func (m *mockStorage) Teams() models.TeamRepository {
	args := m.Called()
	return args.Get(0).(models.TeamRepository)
}

func (m *mockStorage) Projects() models.ProjectRepository {
	args := m.Called()
	return args.Get(0).(models.ProjectRepository)
}

func (m *mockStorage) Tasks() models.TaskRepository {
	args := m.Called()
	return args.Get(0).(models.TaskRepository)
}

func (m *mockStorage) Comments() models.CommentRepository {
	args := m.Called()
	return args.Get(0).(models.CommentRepository)
}

func (m *mockStorage) Atomic(ctx context.Context, fn func(models.Storage) error) error {
	// For testing, we usually just execute the function with the mock itself
	return fn(m)
}

func (m *mockStorage) CheckHealth(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// --- Repository Mocks ---

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) GetByEmail(email string) (models.User, error) {
	args := m.Called(email)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockUserRepository) GetByID(id uuid.UUID) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockUserRepository) CreateUser(user models.User) (models.User, error) {
	args := m.Called(user)
	return args.Get(0).(models.User), args.Error(1)
}

type mockWorkspaceRepository struct {
	mock.Mock
}

func (m *mockWorkspaceRepository) GetByID(id uuid.UUID) (models.Workspace, error) {
	args := m.Called(id)
	return args.Get(0).(models.Workspace), args.Error(1)
}

func (m *mockWorkspaceRepository) ListWorkspacesByUserID(userID uuid.UUID) ([]models.Workspace, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Workspace), args.Error(1)
}

func (m *mockWorkspaceRepository) ListMembersByWorkspaceId(workspaceID uuid.UUID) ([]models.WorkspaceMemberWithUser, error) {
	args := m.Called(workspaceID)
	return args.Get(0).([]models.WorkspaceMemberWithUser), args.Error(1)
}

func (m *mockWorkspaceRepository) GetMember(workspaceID, userID uuid.UUID) (models.WorkspaceMember, error) {
	args := m.Called(workspaceID, userID)
	return args.Get(0).(models.WorkspaceMember), args.Error(1)
}

func (m *mockWorkspaceRepository) CreateWorkspace(ws models.Workspace) (models.Workspace, error) {
	args := m.Called(ws)
	return args.Get(0).(models.Workspace), args.Error(1)
}

func (m *mockWorkspaceRepository) AddMember(member models.WorkspaceMember) error {
	args := m.Called(member)
	return args.Error(0)
}

func (m *mockWorkspaceRepository) RemoveMember(workspaceID, userID uuid.UUID) error {
	args := m.Called(workspaceID, userID)
	return args.Error(0)
}

type mockTeamRepository struct {
	mock.Mock
}

func (m *mockTeamRepository) GetByID(id uuid.UUID) (models.Team, error) {
	args := m.Called(id)
	return args.Get(0).(models.Team), args.Error(1)
}

func (m *mockTeamRepository) ListTeamsByWorkspaceID(workspaceID uuid.UUID) ([]models.Team, error) {
	args := m.Called(workspaceID)
	return args.Get(0).([]models.Team), args.Error(1)
}

func (m *mockTeamRepository) CreateTeam(team models.Team) (models.Team, error) {
	args := m.Called(team)
	return args.Get(0).(models.Team), args.Error(1)
}

func (m *mockTeamRepository) AddMember(member models.TeamMember) error {
	args := m.Called(member)
	return args.Error(0)
}

func (m *mockTeamRepository) ListMembersByTeamId(teamID uuid.UUID) ([]models.TeamMemberWithUser, error) {
	args := m.Called(teamID)
	return args.Get(0).([]models.TeamMemberWithUser), args.Error(1)
}

func (m *mockTeamRepository) UpdateTeam(team models.Team) (models.Team, error) {
	args := m.Called(team)
	return args.Get(0).(models.Team), args.Error(1)
}

func (m *mockTeamRepository) DeleteTeam(teamID uuid.UUID) error {
	args := m.Called(teamID)
	return args.Error(0)
}

func (m *mockTeamRepository) RemoveMember(teamID, userID uuid.UUID) error {
	args := m.Called(teamID, userID)
	return args.Error(0)
}

func (m *mockTeamRepository) GetMember(teamID, userID uuid.UUID) (models.TeamMember, error) {
	args := m.Called(teamID, userID)
	return args.Get(0).(models.TeamMember), args.Error(1)
}

type mockProjectRepository struct {
	mock.Mock
}

func (m *mockProjectRepository) Create(project models.Project) (models.Project, error) {
	args := m.Called(project)
	return args.Get(0).(models.Project), args.Error(1)
}

func (m *mockProjectRepository) GetByID(id uuid.UUID) (models.Project, error) {
	args := m.Called(id)
	return args.Get(0).(models.Project), args.Error(1)
}

func (m *mockProjectRepository) ListByWorkspaceID(workspaceID uuid.UUID) ([]models.Project, error) {
	args := m.Called(workspaceID)
	return args.Get(0).([]models.Project), args.Error(1)
}

func (m *mockProjectRepository) Update(project models.Project) (models.Project, error) {
	args := m.Called(project)
	return args.Get(0).(models.Project), args.Error(1)
}

func (m *mockProjectRepository) AddMember(member models.ProjectMember) error {
	args := m.Called(member)
	return args.Error(0)
}

func (m *mockProjectRepository) RemoveMember(projectID, userID uuid.UUID) error {
	args := m.Called(projectID, userID)
	return args.Error(0)
}

func (m *mockProjectRepository) GetMember(projectID, userID uuid.UUID) (models.ProjectMember, error) {
	args := m.Called(projectID, userID)
	return args.Get(0).(models.ProjectMember), args.Error(1)
}

func (m *mockProjectRepository) GetMembers(projectID uuid.UUID) ([]models.ProjectMemberWithUser, error) {
	args := m.Called(projectID)
	return args.Get(0).([]models.ProjectMemberWithUser), args.Error(1)
}

func (m *mockProjectRepository) AddTeam(team models.ProjectTeam) error {
	args := m.Called(team)
	return args.Error(0)
}

func (m *mockProjectRepository) RemoveTeam(projectID, teamID uuid.UUID) error {
	args := m.Called(projectID, teamID)
	return args.Error(0)
}

func (m *mockProjectRepository) GetTeam(projectID, teamID uuid.UUID) (models.ProjectTeam, error) {
	args := m.Called(projectID, teamID)
	return args.Get(0).(models.ProjectTeam), args.Error(1)
}

func (m *mockProjectRepository) GetTeams(projectID uuid.UUID) ([]models.ProjectTeamWithDetails, error) {
	args := m.Called(projectID)
	return args.Get(0).([]models.ProjectTeamWithDetails), args.Error(1)
}

func (m *mockProjectRepository) ListByTeamID(teamID uuid.UUID) ([]models.Project, error) {
	args := m.Called(teamID)
	return args.Get(0).([]models.Project), args.Error(1)
}

type mockTaskRepository struct {
	mock.Mock
}

func (m *mockTaskRepository) Create(task models.Task) (models.Task, error) {
	args := m.Called(task)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *mockTaskRepository) GetByID(id uuid.UUID) (models.Task, error) {
	args := m.Called(id)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *mockTaskRepository) Update(task models.Task) (models.Task, error) {
	args := m.Called(task)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *mockTaskRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockTaskRepository) ListByProjectID(projectID uuid.UUID, status *models.TaskStatus, assigneeID *uuid.UUID) ([]models.Task, error) {
	args := m.Called(projectID, status, assigneeID)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *mockTaskRepository) ListByAssigneeID(assigneeID uuid.UUID) ([]models.Task, error) {
	args := m.Called(assigneeID)
	return args.Get(0).([]models.Task), args.Error(1)
}

type mockCommentRepository struct {
	mock.Mock
}

func (m *mockCommentRepository) Create(comment models.Comment) (models.Comment, error) {
	args := m.Called(comment)
	return args.Get(0).(models.Comment), args.Error(1)
}

func (m *mockCommentRepository) GetByID(id uuid.UUID) (models.Comment, error) {
	args := m.Called(id)
	return args.Get(0).(models.Comment), args.Error(1)
}

func (m *mockCommentRepository) ListByTaskID(taskID uuid.UUID) ([]models.CommentWithAuthor, error) {
	args := m.Called(taskID)
	return args.Get(0).([]models.CommentWithAuthor), args.Error(1)
}

func (m *mockCommentRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// --- Hub Mock ---

type mockHub struct {
	mock.Mock
}

func (m *mockHub) Run() {
	m.Called()
}

func (m *mockHub) Subscribe(topic string, conn *websocket.Conn) {
	m.Called(topic, conn)
}

func (m *mockHub) Unsubscribe(topic string, conn *websocket.Conn) {
	m.Called(topic, conn)
}

func (m *mockHub) Broadcast(topic string, payload interface{}) {
	m.Called(topic, payload)
}

// --- Service Mocks (Cross-service dependencies) ---

type mockProjectService struct {
	mock.Mock
}

func (m *mockProjectService) CreateProject(userID, workspaceID uuid.UUID, req structs.ReqCreateProject) (models.Project, error) {
	args := m.Called(userID, workspaceID, req)
	return args.Get(0).(models.Project), args.Error(1)
}

func (m *mockProjectService) GetProject(userID, projectID uuid.UUID) (models.Project, error) {
	args := m.Called(userID, projectID)
	return args.Get(0).(models.Project), args.Error(1)
}

func (m *mockProjectService) UpdateProject(userID, projectID uuid.UUID, req structs.ReqUpdateProject) (models.Project, error) {
	args := m.Called(userID, projectID, req)
	return args.Get(0).(models.Project), args.Error(1)
}

func (m *mockProjectService) AddMember(userID, projectID uuid.UUID, req structs.ReqAddProjectMember) (models.ProjectMember, error) {
	args := m.Called(userID, projectID, req)
	return args.Get(0).(models.ProjectMember), args.Error(1)
}

func (m *mockProjectService) RemoveMember(userID, projectID, targetUserID uuid.UUID) error {
	args := m.Called(userID, projectID, targetUserID)
	return args.Error(0)
}

func (m *mockProjectService) ListMembers(userID, projectID uuid.UUID) ([]structs.ResProjectMember, error) {
	args := m.Called(userID, projectID)
	return args.Get(0).([]structs.ResProjectMember), args.Error(1)
}

func (m *mockProjectService) ListByWorkspaceID(workspaceID uuid.UUID) ([]models.Project, error) {
	args := m.Called(workspaceID)
	return args.Get(0).([]models.Project), args.Error(1)
}

func (m *mockProjectService) AddTeam(userID, projectID, teamID uuid.UUID) (structs.ResProjectTeam, error) {
	args := m.Called(userID, projectID, teamID)
	return args.Get(0).(structs.ResProjectTeam), args.Error(1)
}

func (m *mockProjectService) RemoveTeam(userID, projectID, teamID uuid.UUID) error {
	args := m.Called(userID, projectID, teamID)
	return args.Error(0)
}

func (m *mockProjectService) ListTeams(userID, projectID uuid.UUID) ([]structs.ResProjectTeam, error) {
	args := m.Called(userID, projectID)
	return args.Get(0).([]structs.ResProjectTeam), args.Error(1)
}

func (m *mockProjectService) ValidateProjectAccess(projectID, userID uuid.UUID, requireAdmin bool) error {
	args := m.Called(projectID, userID, requireAdmin)
	return args.Error(0)
}

type mockWorkspaceService struct {
	mock.Mock
}

func (m *mockWorkspaceService) CreateWorkspace(ownerID uuid.UUID, req structs.ReqCreateWorkspace) (models.Workspace, error) {
	args := m.Called(ownerID, req)
	return args.Get(0).(models.Workspace), args.Error(1)
}

func (m *mockWorkspaceService) ListWorkspacesByUserID(userID uuid.UUID) ([]models.Workspace, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Workspace), args.Error(1)
}

func (m *mockWorkspaceService) ListMembersByWorkspaceId(workspaceID uuid.UUID) ([]models.WorkspaceMemberWithUser, error) {
	args := m.Called(workspaceID)
	return args.Get(0).([]models.WorkspaceMemberWithUser), args.Error(1)
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
	return args.Get(0).(models.WorkspaceMember), args.Error(1)
}

type mockPublisher struct {
	mock.Mock
}

func (m *mockPublisher) Publish(topic string, handle workers.Handler) error {
	args := m.Called(topic, handle)
	return args.Error(0)
}
