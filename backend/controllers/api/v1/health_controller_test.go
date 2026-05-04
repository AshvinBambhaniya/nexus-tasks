package v1

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

func setupHealthControllerTest() (*fiber.App, *mockHealthService, *HealthController) {
	app := fiber.New()
	mockSvc := new(mockHealthService)
	logger := zap.NewNop()
	ctrl, _ := NewHealthController(mockSvc, logger)
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
			app, mockSvc, ctrl := setupHealthControllerTest()
			app.Get("/healthz", ctrl.Overall)
			tt.setupMocks(mockSvc)

			req := httptest.NewRequest("GET", "/healthz", nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestHealthController_Self(t *testing.T) {
	app, _, ctrl := setupHealthControllerTest()
	app.Get("/healthz/self", ctrl.Self)

	req := httptest.NewRequest("GET", "/healthz/self", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
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
			app, mockSvc, ctrl := setupHealthControllerTest()
			app.Get("/healthz/db", ctrl.Db)
			tt.setupMocks(mockSvc)

			req := httptest.NewRequest("GET", "/healthz/db", nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockSvc.AssertExpectations(t)
		})
	}
}
