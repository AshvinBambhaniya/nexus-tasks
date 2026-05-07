package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupAuthTest(t *testing.T) (*fiber.App, *mockUserService, *AuthController) {
	app := fiber.New()
	mockSvc := new(mockUserService)
	logger := zap.NewNop()
	cfg := &config.AppConfig{JwtExpirationHours: 24}
	ctrl, err := NewAuthController(mockSvc, logger, cfg)
	assert.NoError(t, err)
	return app, mockSvc, ctrl
}

func TestAuthController_Register(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		reqBody        interface{}
		setupMocks     func(m *mockUserService)
		expectedStatus int
	}{
		{
			name: "success",
			reqBody: structs.ReqRegisterUser{
				Email:    "test@example.com",
				Password: "password123",
				FullName: "Test User",
			},
			setupMocks: func(m *mockUserService) {
				m.On("Register", "test@example.com", "password123", "Test User").
					Return(models.User{ID: userID, Email: "test@example.com"}, "token123", nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid_request_body",
			reqBody:        "invalid json",
			setupMocks:     func(_ *mockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation_failure_missing_email",
			reqBody: structs.ReqRegisterUser{
				Password: "password123",
				FullName: "Test User",
			},
			setupMocks:     func(_ *mockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			reqBody: structs.ReqRegisterUser{
				Email:    "fail@test.com",
				Password: "password123",
				FullName: "Test User",
			},
			setupMocks: func(m *mockUserService) {
				m.On("Register", mock.Anything, mock.Anything, mock.Anything).
					Return(models.User{}, "", errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupAuthTest(t)
			app.Post("/register", ctrl.Register)
			tt.setupMocks(mockSvc)

			body, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestAuthController_Login(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        interface{}
		setupMocks     func(m *mockUserService)
		expectedStatus int
	}{
		{
			name:    "success",
			reqBody: structs.ReqLoginUser{Email: "test@test.com", Password: "pwd"},
			setupMocks: func(m *mockUserService) {
				m.On("Authenticate", "test@test.com", "pwd").Return("token-abc", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "invalid_credentials",
			reqBody: structs.ReqLoginUser{Email: "test@test.com", Password: "wrong"},
			setupMocks: func(m *mockUserService) {
				m.On("Authenticate", mock.Anything, mock.Anything).Return("", errors.New("auth fail"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid_request_body",
			reqBody:        "invalid json",
			setupMocks:     func(_ *mockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "validation_failure_missing_email",
			reqBody:        structs.ReqLoginUser{Password: "pwd"},
			setupMocks:     func(_ *mockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupAuthTest(t)
			app.Post("/login", ctrl.Login)
			tt.setupMocks(mockSvc)

			body, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestAuthController_Logout(t *testing.T) {
	app, _, ctrl := setupAuthTest(t)
	app.Post("/logout", ctrl.Logout)

	r := httptest.NewRequest("POST", "/logout", nil)
	resp, err := app.Test(r)
	assert.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthController_Me(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		setupContext   func(c *fiber.Ctx)
		setupMocks     func(m *mockUserService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, userID.String())
			},
			setupMocks: func(m *mockUserService) {
				m.On("GetByID", userID).Return(models.User{ID: userID, Email: "me@test.com"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_user_id_in_context",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, "invalid-uuid")
			},
			setupMocks: func(m *mockUserService) {
				m.On("GetByID", userID).Return(models.User{}, nil)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "user_not_found",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(constants.ContextUID, userID.String())
			},
			setupMocks: func(m *mockUserService) {
				m.On("GetByID", userID).Return(models.User{}, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupAuthTest(t)
			app.Get("/me", func(c *fiber.Ctx) error {
				tt.setupContext(c)
				return ctrl.Me(c)
			})
			tt.setupMocks(mockSvc)

			req := httptest.NewRequest("GET", "/me", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
