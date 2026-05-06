package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupHealthTest(_ *testing.T) (*healthService, *mockStorage) {
	mockStor := new(mockStorage)
	logger := zap.NewNop()

	svc := &healthService{storage: mockStor, logger: logger}

	return svc, mockStor
}

func TestHealthService_CheckDatabaseHealth(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockStor := setupHealthTest(t)

		mockStor.On("CheckHealth", mock.Anything).Return(nil)

		err := svc.CheckDatabaseHealth(context.Background())

		assert.NoError(t, err)
		mockStor.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		svc, mockStor := setupHealthTest(t)

		mockStor.On("CheckHealth", mock.Anything).Return(errors.New("db connection failed"))

		err := svc.CheckDatabaseHealth(context.Background())

		assert.Error(t, err)
		assert.Equal(t, "db connection failed", err.Error())
		mockStor.AssertExpectations(t)
	})
}
