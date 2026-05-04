package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupWorkspaceControllerTest() (*fiber.App, *mockWorkspaceService, *WorkspaceController) {
	app := fiber.New()
	mockSvc := new(mockWorkspaceService)
	logger := zap.NewNop()
	// publisher is ignored by the actual constructor implementation
	ctrl, _ := NewWorkspaceController(mockSvc, logger, nil)
	return app, mockSvc, ctrl
}

func TestWorkspaceController_Create(t *testing.T) {
	uid := uuid.New()
	wsID := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		reqBody        interface{}
		setupMocks     func(m *mockWorkspaceService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			reqBody: structs.ReqCreateWorkspace{Name: "New Workspace"},
			setupMocks: func(m *mockWorkspaceService) {
				m.On("CreateWorkspace", uid, mock.MatchedBy(func(req structs.ReqCreateWorkspace) bool {
					return req.Name == "New Workspace"
				})).Return(models.Workspace{ID: wsID, Name: "New Workspace"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, "invalid-uuid")
			},
			reqBody:        structs.ReqCreateWorkspace{Name: "X"},
			setupMocks:     func(m *mockWorkspaceService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_request_body",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			reqBody:        "invalid json",
			setupMocks:     func(m *mockWorkspaceService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation_failure",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			reqBody:        structs.ReqCreateWorkspace{Name: ""}, // empty name
			setupMocks:     func(m *mockWorkspaceService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			reqBody: structs.ReqCreateWorkspace{Name: "New Workspace"},
			setupMocks: func(m *mockWorkspaceService) {
				m.On("CreateWorkspace", mock.Anything, mock.Anything).Return(models.Workspace{}, errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupWorkspaceControllerTest()
			app.Post("/", func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.Create(c)
			})
			tt.setupMocks(mockSvc)

			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestWorkspaceController_List(t *testing.T) {
	uid := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		setupMocks     func(m *mockWorkspaceService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			setupMocks: func(m *mockWorkspaceService) {
				m.On("ListWorkspacesByUserID", uid).Return([]models.Workspace{{ID: uuid.New(), Name: "W1"}}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, "invalid-uuid")
			},
			setupMocks:     func(m *mockWorkspaceService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			setupMocks: func(m *mockWorkspaceService) {
				m.On("ListWorkspacesByUserID", uid).Return(nil, errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupWorkspaceControllerTest()
			app.Get("/", func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.List(c)
			})
			tt.setupMocks(mockSvc)

			req := httptest.NewRequest("GET", "/", nil)
			resp, _ := app.Test(req)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestWorkspaceController_ListMembers(t *testing.T) {
	wsID := uuid.New()

	tests := []struct {
		name           string
		wsIDParam      string
		setupMocks     func(m *mockWorkspaceService)
		expectedStatus int
	}{
		{
			name:      "success",
			wsIDParam: wsID.String(),
			setupMocks: func(m *mockWorkspaceService) {
				m.On("ListMembersByWorkspaceId", wsID).Return([]models.WorkspaceMemberWithUser{{UserID: uuid.New(), Email: "u@t.com"}}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid_workspace_id",
			wsIDParam:      "invalid",
			setupMocks:     func(m *mockWorkspaceService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "service_error",
			wsIDParam: wsID.String(),
			setupMocks: func(m *mockWorkspaceService) {
				m.On("ListMembersByWorkspaceId", wsID).Return(nil, errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupWorkspaceControllerTest()
			app.Get("/:"+constants.ParamWorkspaceID+"/members", ctrl.ListMembers)
			tt.setupMocks(mockSvc)

			req := httptest.NewRequest("GET", "/"+tt.wsIDParam+"/members", nil)
			resp, _ := app.Test(req)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestWorkspaceController_InviteMember(t *testing.T) {
	uid := uuid.New()
	wsID := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		wsIDParam      string
		reqBody        interface{}
		setupMocks     func(m *mockWorkspaceService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam: wsID.String(),
			reqBody:   structs.ReqInviteWorkspaceMember{Email: "invite@test.com"},
			setupMocks: func(m *mockWorkspaceService) {
				m.On("InviteMember", uid, wsID, "invite@test.com").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "unauthorized_user_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, "bad")
			},
			wsIDParam:      wsID.String(),
			reqBody:        structs.ReqInviteWorkspaceMember{Email: "x@t.com"},
			setupMocks:     func(m *mockWorkspaceService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_workspace_id",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      "invalid",
			reqBody:        structs.ReqInviteWorkspaceMember{Email: "invite@test.com"},
			setupMocks:     func(m *mockWorkspaceService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_body",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			reqBody:        "invalid",
			setupMocks:     func(m *mockWorkspaceService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation_failure - invalid email",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			reqBody:        structs.ReqInviteWorkspaceMember{Email: "invalid-email"},
			setupMocks:     func(m *mockWorkspaceService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam: wsID.String(),
			reqBody:   structs.ReqInviteWorkspaceMember{Email: "invite@test.com"},
			setupMocks: func(m *mockWorkspaceService) {
				m.On("InviteMember", uid, wsID, "invite@test.com").Return(errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupWorkspaceControllerTest()
			app.Post("/:"+constants.ParamWorkspaceID+"/invite", func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.InviteMember(c)
			})
			tt.setupMocks(mockSvc)

			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest("POST", "/"+tt.wsIDParam+"/invite", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestWorkspaceController_RemoveMember(t *testing.T) {
	uid := uuid.New()
	wsID := uuid.New()
	targetUID := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		wsIDParam      string
		targetUIDParam string
		setupMocks     func(m *mockWorkspaceService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			targetUIDParam: targetUID.String(),
			setupMocks: func(m *mockWorkspaceService) {
				m.On("RemoveMember", uid, wsID, targetUID).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "unauthorized_user_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, "bad")
			},
			wsIDParam:      wsID.String(),
			targetUIDParam: targetUID.String(),
			setupMocks:     func(m *mockWorkspaceService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_workspace_id",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      "invalid",
			targetUIDParam: targetUID.String(),
			setupMocks:     func(m *mockWorkspaceService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_target_user_id",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			targetUIDParam: "invalid",
			setupMocks:     func(m *mockWorkspaceService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			targetUIDParam: targetUID.String(),
			setupMocks: func(m *mockWorkspaceService) {
				m.On("RemoveMember", uid, wsID, targetUID).Return(errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupWorkspaceControllerTest()
			app.Delete("/:"+constants.ParamWorkspaceID+"/members/:"+constants.ParamUid, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.RemoveMember(c)
			})
			tt.setupMocks(mockSvc)

			path := fmt.Sprintf("/%s/members/%s", tt.wsIDParam, tt.targetUIDParam)
			req := httptest.NewRequest("DELETE", path, nil)

			resp, _ := app.Test(req)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
