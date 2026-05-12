package v2

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupHealthControllerTest(t *testing.T) (*fiber.App, *mockHealthService, *HealthController) {
	app := fiber.New()
	mockSvc := new(mockHealthService)
	logger := zap.NewNop()
	ctrl, err := NewHealthController(mockSvc, logger)
	assert.NoError(t, err)
	return app, mockSvc, ctrl
}

func TestHealthController_Overall(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(m *mockHealthService)
		expectedStatus int
	}{
		{
			name: "success",
			setupMocks: func(m *mockHealthService) {
				m.On("CheckDatabaseHealth", mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "failure",
			setupMocks: func(m *mockHealthService) {
				m.On("CheckDatabaseHealth", mock.Anything).Return(errors.New("db connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupHealthControllerTest(t)
			app.Get("/healthz", ctrl.Overall)
			tt.setupMocks(mockSvc)

			req := httptest.NewRequest("GET", "/healthz", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestHealthController_Self(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		app, _, ctrl := setupHealthControllerTest(t)
		app.Get("/healthz/self", ctrl.Self)

		req := httptest.NewRequest("GET", "/healthz/self", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestHealthController_Db(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(m *mockHealthService)
		expectedStatus int
	}{
		{
			name: "success",
			setupMocks: func(m *mockHealthService) {
				m.On("CheckDatabaseHealth", mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "failure",
			setupMocks: func(m *mockHealthService) {
				m.On("CheckDatabaseHealth", mock.Anything).Return(errors.New("db connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockSvc, ctrl := setupHealthControllerTest(t)
			app.Get("/healthz/db", ctrl.Db)
			tt.setupMocks(mockSvc)

			req := httptest.NewRequest("GET", "/healthz/db", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockSvc.AssertExpectations(t)
		})
	}
}
