package v2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupProjectControllerTest(t *testing.T) (*fiber.App, *mockProjectService, *ProjectController) {
	app := fiber.New()
	mockSvc := new(mockProjectService)
	logger := zap.NewNop()
	ctrl, err := NewProjectController(mockSvc, logger)
	assert.NoError(t, err)
	return app, mockSvc, ctrl
}

func TestProjectController_Create(t *testing.T) {
	uid := uuid.New()
	wsID := uuid.New()
	pid := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		wsIDParam      string
		reqBody        interface{}
		setupMocks     func(m *mockProjectService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			wsIDParam: wsID.String(),
			reqBody:   structs.ReqCreateProject{Name: "P1"},
			setupMocks: func(m *mockProjectService) {
				m.On("CreateProject", uid, wsID, mock.Anything).Return(models.Project{ID: pid, Name: "P1"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "forbidden",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			wsIDParam: wsID.String(),
			reqBody:   structs.ReqCreateProject{Name: "P1"},
			setupMocks: func(m *mockProjectService) {
				m.On("CreateProject", uid, wsID, mock.Anything).Return(models.Project{}, errors.New("unauthorized: only workspace admins can create projects"))
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			wsIDParam:      wsID.String(),
			reqBody:        structs.ReqCreateProject{Name: "P1"},
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_workspace_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			wsIDParam:      "invalid-uuid",
			reqBody:        structs.ReqCreateProject{Name: "P1"},
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_request_body",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			wsIDParam:      wsID.String(),
			reqBody:        "invalid json",
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			wsIDParam:      wsID.String(),
			reqBody:        structs.ReqCreateProject{Name: ""},
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			wsIDParam: wsID.String(),
			reqBody:   structs.ReqCreateProject{Name: "P1"},
			setupMocks: func(m *mockProjectService) {
				m.On("CreateProject", uid, wsID, mock.Anything).Return(models.Project{}, errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupProjectControllerTest(t)
			app.Post("/:"+constants.ParamWorkspaceID, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.Create(c)
			})
			tt.setupMocks(mockSvc)

			body, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)
			req := httptest.NewRequest("POST", "/"+tt.wsIDParam, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestProjectController_List(t *testing.T) {
	wsID := uuid.New()

	tests := []struct {
		name           string
		wsIDParam      string
		setupMocks     func(m *mockProjectService)
		expectedStatus int
	}{
		{
			name:      "success",
			wsIDParam: wsID.String(),
			setupMocks: func(m *mockProjectService) {
				m.On("ListByWorkspaceID", wsID).Return([]models.Project{{ID: uuid.New(), Name: "P"}}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "success_empty_list",
			wsIDParam: wsID.String(),
			setupMocks: func(m *mockProjectService) {
				m.On("ListByWorkspaceID", wsID).Return([]models.Project{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid_workspace_id_param",
			wsIDParam:      "invalid-uuid",
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "service_error",
			wsIDParam: wsID.String(),
			setupMocks: func(m *mockProjectService) {
				m.On("ListByWorkspaceID", wsID).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupProjectControllerTest(t)
			app.Get("/:"+constants.ParamWorkspaceID, ctrl.List)
			tt.setupMocks(mockSvc)

			req := httptest.NewRequest("GET", "/"+tt.wsIDParam, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestProjectController_Get(t *testing.T) {
	uid := uuid.New()
	pid := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		pidParam       string
		setupMocks     func(m *mockProjectService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			setupMocks: func(m *mockProjectService) {
				m.On("GetProject", uid, pid).Return(models.Project{ID: pid, Name: "P"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			pidParam:       pid.String(),
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_project_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       "invalid-uuid",
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "project_not_found",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			setupMocks: func(m *mockProjectService) {
				m.On("GetProject", uid, pid).Return(models.Project{}, errors.New("project not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupProjectControllerTest(t)
			app.Get("/:"+constants.ParamProjectID, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.Get(c)
			})
			tt.setupMocks(mockSvc)

			req := httptest.NewRequest("GET", "/"+tt.pidParam, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestProjectController_Update(t *testing.T) {
	uid := uuid.New()
	pid := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		pidParam       string
		reqBody        interface{}
		setupMocks     func(m *mockProjectService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			reqBody:  structs.ReqUpdateProject{Name: "Updated"},
			setupMocks: func(m *mockProjectService) {
				m.On("UpdateProject", uid, pid, mock.Anything).Return(models.Project{ID: pid, Name: "Updated"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			pidParam:       pid.String(),
			reqBody:        structs.ReqUpdateProject{Name: "Updated"},
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_project_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       "invalid-uuid",
			reqBody:        structs.ReqUpdateProject{Name: "Updated"},
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_request_body",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       pid.String(),
			reqBody:        "invalid json",
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			reqBody:  structs.ReqUpdateProject{Name: "Updated"},
			setupMocks: func(m *mockProjectService) {
				m.On("UpdateProject", uid, pid, mock.Anything).Return(models.Project{}, errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupProjectControllerTest(t)
			app.Patch("/:"+constants.ParamProjectID, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.Update(c)
			})
			tt.setupMocks(mockSvc)

			body, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)
			req := httptest.NewRequest("PATCH", "/"+tt.pidParam, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestProjectController_AddMember(t *testing.T) {
	uid := uuid.New()
	pid := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		pidParam       string
		reqBody        interface{}
		setupMocks     func(m *mockProjectService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			reqBody:  structs.ReqAddProjectMember{Email: "u@t.com", Role: "MEMBER"},
			setupMocks: func(m *mockProjectService) {
				m.On("AddMember", uid, pid, mock.Anything).Return(models.ProjectMember{UserID: uuid.New()}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			pidParam:       pid.String(),
			reqBody:        structs.ReqAddProjectMember{Email: "u@t.com", Role: "MEMBER"},
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_project_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       "invalid-uuid",
			reqBody:        structs.ReqAddProjectMember{Email: "u@t.com", Role: "MEMBER"},
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_request_body",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       pid.String(),
			reqBody:        "invalid json",
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation_error_missing_email",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       pid.String(),
			reqBody:        structs.ReqAddProjectMember{Email: "", Role: "MEMBER"},
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			reqBody:  structs.ReqAddProjectMember{Email: "u@t.com", Role: "MEMBER"},
			setupMocks: func(m *mockProjectService) {
				m.On("AddMember", uid, pid, mock.Anything).Return(models.ProjectMember{}, errors.New("user not found"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupProjectControllerTest(t)
			app.Post("/:"+constants.ParamProjectID+"/members", func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.AddMember(c)
			})
			tt.setupMocks(mockSvc)

			body, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)
			req := httptest.NewRequest("POST", "/"+tt.pidParam+"/members", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestProjectController_RemoveMember(t *testing.T) {
	uid := uuid.New()
	pid := uuid.New()
	targetUID := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		pidParam       string
		targetUIDParam string
		setupMocks     func(m *mockProjectService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       pid.String(),
			targetUIDParam: targetUID.String(),
			setupMocks: func(m *mockProjectService) {
				m.On("RemoveMember", uid, pid, targetUID).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			pidParam:       pid.String(),
			targetUIDParam: targetUID.String(),
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_project_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       "invalid-uuid",
			targetUIDParam: targetUID.String(),
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_target_user_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       pid.String(),
			targetUIDParam: "invalid-uuid",
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       pid.String(),
			targetUIDParam: targetUID.String(),
			setupMocks: func(m *mockProjectService) {
				m.On("RemoveMember", uid, pid, targetUID).Return(errors.New("unauthorized"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupProjectControllerTest(t)
			app.Delete("/:"+constants.ParamProjectID+"/members/:"+constants.ParamUID, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.RemoveMember(c)
			})
			tt.setupMocks(mockSvc)

			path := fmt.Sprintf("/%s/members/%s", tt.pidParam, tt.targetUIDParam)
			req := httptest.NewRequest("DELETE", path, nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestProjectController_ListMembers(t *testing.T) {
	uid := uuid.New()
	pid := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		pidParam       string
		setupMocks     func(m *mockProjectService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			setupMocks: func(m *mockProjectService) {
				m.On("ListMembers", uid, pid).Return([]structs.ResProjectMember{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "success_with_members",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			setupMocks: func(m *mockProjectService) {
				m.On("ListMembers", uid, pid).Return([]structs.ResProjectMember{
					{UserID: uuid.New(), Email: "m@t.com"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			pidParam:       pid.String(),
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_project_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       "invalid-uuid",
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			setupMocks: func(m *mockProjectService) {
				m.On("ListMembers", uid, pid).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupProjectControllerTest(t)
			app.Get("/:"+constants.ParamProjectID+"/members", func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.ListMembers(c)
			})
			tt.setupMocks(mockSvc)

			req := httptest.NewRequest("GET", "/"+tt.pidParam+"/members", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestProjectController_AddTeam(t *testing.T) {
	uid := uuid.New()
	pid := uuid.New()
	tid := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		pidParam       string
		reqBody        interface{}
		setupMocks     func(m *mockProjectService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			reqBody:  structs.ReqAddProjectTeam{TeamID: tid},
			setupMocks: func(m *mockProjectService) {
				m.On("AddTeam", uid, pid, tid).Return(structs.ResProjectTeam{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			pidParam:       pid.String(),
			reqBody:        structs.ReqAddProjectTeam{TeamID: tid},
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_project_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       "invalid-uuid",
			reqBody:        structs.ReqAddProjectTeam{TeamID: tid},
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_request_body",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       pid.String(),
			reqBody:        "invalid json",
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation_error_missing_team_id",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       pid.String(),
			reqBody:        structs.ReqAddProjectTeam{TeamID: uuid.UUID{}},
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			reqBody:  structs.ReqAddProjectTeam{TeamID: tid},
			setupMocks: func(m *mockProjectService) {
				m.On("AddTeam", uid, pid, tid).Return(structs.ResProjectTeam{}, errors.New("team not found"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupProjectControllerTest(t)
			app.Post("/:"+constants.ParamProjectID+"/teams", func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.AddTeam(c)
			})
			tt.setupMocks(mockSvc)

			body, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)
			req := httptest.NewRequest("POST", "/"+tt.pidParam+"/teams", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestProjectController_RemoveTeam(t *testing.T) {
	uid := uuid.New()
	pid := uuid.New()
	tid := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		pidParam       string
		tidParam       string
		setupMocks     func(m *mockProjectService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			tidParam: tid.String(),
			setupMocks: func(m *mockProjectService) {
				m.On("RemoveTeam", uid, pid, tid).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			pidParam:       pid.String(),
			tidParam:       tid.String(),
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_project_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       "invalid-uuid",
			tidParam:       tid.String(),
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_team_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       pid.String(),
			tidParam:       "invalid-uuid",
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			tidParam: tid.String(),
			setupMocks: func(m *mockProjectService) {
				m.On("RemoveTeam", uid, pid, tid).Return(errors.New("unauthorized"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupProjectControllerTest(t)
			app.Delete("/:"+constants.ParamProjectID+"/teams/:"+constants.ParamTeamID, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.RemoveTeam(c)
			})
			tt.setupMocks(mockSvc)

			path := fmt.Sprintf("/%s/teams/%s", tt.pidParam, tt.tidParam)
			req := httptest.NewRequest("DELETE", path, nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestProjectController_ListTeams(t *testing.T) {
	uid := uuid.New()
	pid := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		pidParam       string
		setupMocks     func(m *mockProjectService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			setupMocks: func(m *mockProjectService) {
				m.On("ListTeams", uid, pid).Return([]structs.ResProjectTeam{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "success_with_teams",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			setupMocks: func(m *mockProjectService) {
				m.On("ListTeams", uid, pid).Return([]structs.ResProjectTeam{
					{TeamID: uuid.New(), TeamName: "Team A"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			pidParam:       pid.String(),
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_project_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam:       "invalid-uuid",
			setupMocks:     func(_ *mockProjectService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, uid.String())
			},
			pidParam: pid.String(),
			setupMocks: func(m *mockProjectService) {
				m.On("ListTeams", uid, pid).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupProjectControllerTest(t)
			app.Get("/:"+constants.ParamProjectID+"/teams", func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.ListTeams(c)
			})
			tt.setupMocks(mockSvc)

			req := httptest.NewRequest("GET", "/"+tt.pidParam+"/teams", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
