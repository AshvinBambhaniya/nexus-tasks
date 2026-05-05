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

func setupTeamControllerTest() (*fiber.App, *mockTeamService, *TeamController) {
	app := fiber.New()
	mockSvc := new(mockTeamService)
	logger := zap.NewNop()
	ctrl, err := NewTeamController(mockSvc, logger)
	if err != nil {
		panic(err)
	}
	return app, mockSvc, ctrl
}

func TestTeamController_Create(t *testing.T) {
	uid := uuid.New()
	wsID := uuid.New()
	tid := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		wsIDParam      string
		reqBody        interface{}
		setupMocks     func(m *mockTeamService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam: wsID.String(),
			reqBody:   structs.ReqCreateTeam{Name: "T1"},
			setupMocks: func(m *mockTeamService) {
				m.On("CreateTeam", uid, wsID, mock.Anything).Return(models.Team{ID: tid, Name: "T1"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "forbidden",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam: wsID.String(),
			reqBody:   structs.ReqCreateTeam{Name: "T1"},
			setupMocks: func(m *mockTeamService) {
				m.On("CreateTeam", uid, wsID, mock.Anything).Return(models.Team{}, errors.New("unauthorized: only workspace admins can create teams"))
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, "invalid-uuid")
			},
			wsIDParam:      wsID.String(),
			reqBody:        structs.ReqCreateTeam{Name: "T1"},
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_workspace_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      "invalid-uuid",
			reqBody:        structs.ReqCreateTeam{Name: "T1"},
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_request_body",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			reqBody:        "invalid json",
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			reqBody:        structs.ReqCreateTeam{Name: ""},
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam: wsID.String(),
			reqBody:   structs.ReqCreateTeam{Name: "T1"},
			setupMocks: func(m *mockTeamService) {
				m.On("CreateTeam", uid, wsID, mock.Anything).Return(models.Team{}, errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupTeamControllerTest()
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
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTeamController_List(t *testing.T) {
	wsID := uuid.New()

	tests := []struct {
		name           string
		wsIDParam      string
		setupMocks     func(m *mockTeamService)
		expectedStatus int
	}{
		{
			name:      "success",
			wsIDParam: wsID.String(),
			setupMocks: func(m *mockTeamService) {
				m.On("ListTeamsByWorkspaceID", wsID).Return([]models.Team{{ID: uuid.New(), Name: "T"}}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "success_empty_list",
			wsIDParam: wsID.String(),
			setupMocks: func(m *mockTeamService) {
				m.On("ListTeamsByWorkspaceID", wsID).Return([]models.Team{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid_workspace_id_param",
			wsIDParam:      "invalid-uuid",
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "service_error",
			wsIDParam: wsID.String(),
			setupMocks: func(m *mockTeamService) {
				m.On("ListTeamsByWorkspaceID", wsID).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupTeamControllerTest()
			app.Get("/:"+constants.ParamWorkspaceID, ctrl.List)
			tt.setupMocks(mockSvc)

			req := httptest.NewRequest("GET", "/"+tt.wsIDParam, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTeamController_Get(t *testing.T) {
	tid := uuid.New()

	tests := []struct {
		name           string
		tidParam       string
		setupMocks     func(m *mockTeamService)
		expectedStatus int
	}{
		{
			name:     "success",
			tidParam: tid.String(),
			setupMocks: func(m *mockTeamService) {
				m.On("GetTeam", tid).Return(structs.ResTeamWithProjects{
					ResTeam:  structs.ResTeam{ID: tid, Name: "T"},
					Projects: []models.Project{},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid_team_id_param",
			tidParam:       "invalid-uuid",
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "team_not_found",
			tidParam: tid.String(),
			setupMocks: func(m *mockTeamService) {
				m.On("GetTeam", tid).Return(structs.ResTeamWithProjects{}, errors.New("team not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:     "service_error",
			tidParam: tid.String(),
			setupMocks: func(m *mockTeamService) {
				m.On("GetTeam", tid).Return(structs.ResTeamWithProjects{}, errors.New("database error"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupTeamControllerTest()
			app.Get("/:"+constants.ParamTeamID, ctrl.Get)
			tt.setupMocks(mockSvc)

			req := httptest.NewRequest("GET", "/"+tt.tidParam, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTeamController_Update(t *testing.T) {
	uid := uuid.New()
	wsID := uuid.New()
	tid := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		wsIDParam      string
		tidParam       string
		reqBody        interface{}
		setupMocks     func(m *mockTeamService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam: wsID.String(),
			tidParam:  tid.String(),
			reqBody:   structs.ReqUpdateTeam{Name: "Updated"},
			setupMocks: func(m *mockTeamService) {
				m.On("UpdateTeam", uid, wsID, tid, mock.Anything).Return(models.Team{ID: tid, Name: "Updated"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, "invalid-uuid")
			},
			wsIDParam:      wsID.String(),
			tidParam:       tid.String(),
			reqBody:        structs.ReqUpdateTeam{Name: "Updated"},
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_workspace_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      "invalid-uuid",
			tidParam:       tid.String(),
			reqBody:        structs.ReqUpdateTeam{Name: "Updated"},
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_team_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			tidParam:       "invalid-uuid",
			reqBody:        structs.ReqUpdateTeam{Name: "Updated"},
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_request_body",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			tidParam:       tid.String(),
			reqBody:        "invalid json",
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam: wsID.String(),
			tidParam:  tid.String(),
			reqBody:   structs.ReqUpdateTeam{Name: "Updated"},
			setupMocks: func(m *mockTeamService) {
				m.On("UpdateTeam", uid, wsID, tid, mock.Anything).Return(models.Team{}, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupTeamControllerTest()
			app.Patch("/:"+constants.ParamWorkspaceID+"/:"+constants.ParamTeamID, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.Update(c)
			})
			tt.setupMocks(mockSvc)

			body, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)
			path := fmt.Sprintf("/%s/%s", tt.wsIDParam, tt.tidParam)
			req := httptest.NewRequest("PATCH", path, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTeamController_Delete(t *testing.T) {
	uid := uuid.New()
	wsID := uuid.New()
	tid := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		wsIDParam      string
		tidParam       string
		setupMocks     func(m *mockTeamService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam: wsID.String(),
			tidParam:  tid.String(),
			setupMocks: func(m *mockTeamService) {
				m.On("DeleteTeam", uid, wsID, tid).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, "invalid-uuid")
			},
			wsIDParam:      wsID.String(),
			tidParam:       tid.String(),
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_workspace_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      "invalid-uuid",
			tidParam:       tid.String(),
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_team_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			tidParam:       "invalid-uuid",
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam: wsID.String(),
			tidParam:  tid.String(),
			setupMocks: func(m *mockTeamService) {
				m.On("DeleteTeam", uid, wsID, tid).Return(errors.New("unauthorized"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupTeamControllerTest()
			app.Delete("/:"+constants.ParamWorkspaceID+"/:"+constants.ParamTeamID, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.Delete(c)
			})
			tt.setupMocks(mockSvc)

			path := fmt.Sprintf("/%s/%s", tt.wsIDParam, tt.tidParam)
			req := httptest.NewRequest("DELETE", path, nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTeamController_AddMember(t *testing.T) {
	uid := uuid.New()
	wsID := uuid.New()
	tid := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		wsIDParam      string
		tidParam       string
		reqBody        interface{}
		setupMocks     func(m *mockTeamService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam: wsID.String(),
			tidParam:  tid.String(),
			reqBody:   structs.ReqAddTeamMember{Email: "u@t.com", Role: "MEMBER"},
			setupMocks: func(m *mockTeamService) {
				m.On("AddMember", uid, wsID, tid, "u@t.com", "MEMBER").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, "invalid-uuid")
			},
			wsIDParam:      wsID.String(),
			tidParam:       tid.String(),
			reqBody:        structs.ReqAddTeamMember{Email: "u@t.com", Role: "MEMBER"},
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_workspace_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      "invalid-uuid",
			tidParam:       tid.String(),
			reqBody:        structs.ReqAddTeamMember{Email: "u@t.com", Role: "MEMBER"},
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_team_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			tidParam:       "invalid-uuid",
			reqBody:        structs.ReqAddTeamMember{Email: "u@t.com", Role: "MEMBER"},
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_request_body",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			tidParam:       tid.String(),
			reqBody:        "invalid json",
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam: wsID.String(),
			tidParam:  tid.String(),
			reqBody:   structs.ReqAddTeamMember{Email: "u@t.com", Role: "MEMBER"},
			setupMocks: func(m *mockTeamService) {
				m.On("AddMember", uid, wsID, tid, "u@t.com", "MEMBER").Return(errors.New("user not found"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupTeamControllerTest()
			app.Post("/:"+constants.ParamWorkspaceID+"/:"+constants.ParamTeamID+"/members", func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.AddMember(c)
			})
			tt.setupMocks(mockSvc)

			body, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)
			path := fmt.Sprintf("/%s/%s/members", tt.wsIDParam, tt.tidParam)
			req := httptest.NewRequest("POST", path, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTeamController_RemoveMember(t *testing.T) {
	uid := uuid.New()
	wsID := uuid.New()
	tid := uuid.New()
	targetUID := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		wsIDParam      string
		tidParam       string
		targetUIDParam string
		setupMocks     func(m *mockTeamService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			tidParam:       tid.String(),
			targetUIDParam: targetUID.String(),
			setupMocks: func(m *mockTeamService) {
				m.On("RemoveMember", uid, wsID, tid, targetUID).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, "invalid-uuid")
			},
			wsIDParam:      wsID.String(),
			tidParam:       tid.String(),
			targetUIDParam: targetUID.String(),
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_workspace_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      "invalid-uuid",
			tidParam:       tid.String(),
			targetUIDParam: targetUID.String(),
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_team_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			tidParam:       "invalid-uuid",
			targetUIDParam: targetUID.String(),
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_target_user_id_param",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			tidParam:       tid.String(),
			targetUIDParam: "invalid-uuid",
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUid, uid.String())
			},
			wsIDParam:      wsID.String(),
			tidParam:       tid.String(),
			targetUIDParam: targetUID.String(),
			setupMocks: func(m *mockTeamService) {
				m.On("RemoveMember", uid, wsID, tid, targetUID).Return(errors.New("unauthorized"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupTeamControllerTest()
			app.Delete("/:"+constants.ParamWorkspaceID+"/:"+constants.ParamTeamID+"/members/:"+constants.ParamUid, func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.RemoveMember(c)
			})
			tt.setupMocks(mockSvc)

			path := fmt.Sprintf("/%s/%s/members/%s", tt.wsIDParam, tt.tidParam, tt.targetUIDParam)
			req := httptest.NewRequest("DELETE", path, nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTeamController_ListMembers(t *testing.T) {
	tid := uuid.New()

	tests := []struct {
		name           string
		tidParam       string
		setupMocks     func(m *mockTeamService)
		expectedStatus int
	}{
		{
			name:     "success",
			tidParam: tid.String(),
			setupMocks: func(m *mockTeamService) {
				m.On("ListMembersByTeamId", tid).Return([]models.TeamMemberWithUser{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "success_with_members",
			tidParam: tid.String(),
			setupMocks: func(m *mockTeamService) {
				m.On("ListMembersByTeamId", tid).Return([]models.TeamMemberWithUser{
					{UserID: uuid.New()},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid_team_id_param",
			tidParam:       "invalid-uuid",
			setupMocks:     func(m *mockTeamService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "service_error",
			tidParam: tid.String(),
			setupMocks: func(m *mockTeamService) {
				m.On("ListMembersByTeamId", tid).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupTeamControllerTest()
			app.Get("/:"+constants.ParamTeamID+"/members", ctrl.ListMembers)
			tt.setupMocks(mockSvc)

			req := httptest.NewRequest("GET", "/"+tt.tidParam+"/members", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
