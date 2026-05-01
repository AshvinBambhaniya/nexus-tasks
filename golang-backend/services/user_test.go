package services

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/AshvinBambhaniya/nexus-tasks/models"
	"github.com/AshvinBambhaniya/nexus-tasks/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupUserTest(t *testing.T) (*userService, *mockUserRepository, *mockWorkspaceRepository, *mockStorage) {
	mockUserRepo := new(mockUserRepository)
	mockWorkspaceRepo := new(mockWorkspaceRepository)
	mockStor := new(mockStorage)

	logger := zap.NewNop()
	cfg := &config.AppConfig{
		Secret:             "testsecret",
		JwtExpirationHours: 1,
	}

	mockStor.On("Users").Return(mockUserRepo)
	mockStor.On("Workspaces").Return(mockWorkspaceRepo)

	svc := NewUserService(mockStor, logger, cfg).(*userService)

	return svc, mockUserRepo, mockWorkspaceRepo, mockStor
}

func TestUserService_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockUserRepo, mockWorkspaceRepo, _ := setupUserTest(t)

		email := "test@example.com"
		userID := uuid.New()
		workspaceID := uuid.New()

		mockUserRepo.On("GetByEmail", email).Return(models.User{}, sql.ErrNoRows)
		mockUserRepo.On("CreateUser", mock.Anything).Return(models.User{ID: userID, Email: email}, nil)
		mockWorkspaceRepo.On("CreateWorkspace", mock.Anything).Return(models.Workspace{ID: workspaceID}, nil)
		mockWorkspaceRepo.On("AddMember", mock.Anything).Return(nil)

		user, token, err := svc.Register(email, "password123", "Full Name")

		assert.NoError(t, err)
		assert.Equal(t, email, user.Email)
		assert.NotEmpty(t, token)
	})

	t.Run("email already registered", func(t *testing.T) {
		svc, mockUserRepo, _, _ := setupUserTest(t)

		email := "existing@example.com"
		mockUserRepo.On("GetByEmail", email).Return(models.User{Email: email}, nil)

		_, _, err := svc.Register(email, "password123", "Name")

		assert.Error(t, err)
		assert.Equal(t, "email already registered", err.Error())
	})

	t.Run("user creation failure", func(t *testing.T) {
		svc, mockUserRepo, _, _ := setupUserTest(t)

		email := "test@example.com"
		mockUserRepo.On("GetByEmail", email).Return(models.User{}, sql.ErrNoRows)
		mockUserRepo.On("CreateUser", mock.Anything).Return(models.User{}, errors.New("db error"))

		_, _, err := svc.Register(email, "password123", "Name")

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())

	})

	t.Run("workspace creation failure", func(t *testing.T) {
		svc, mockUserRepo, mockWorkspaceRepo, _ := setupUserTest(t)
		mockUserRepo.On("GetByEmail", mock.Anything).Return(models.User{}, sql.ErrNoRows)
		mockUserRepo.On("CreateUser", mock.Anything).Return(models.User{ID: uuid.New()}, nil)
		mockWorkspaceRepo.On("CreateWorkspace", mock.Anything).Return(models.Workspace{}, errors.New("ws error"))
		_, _, err := svc.Register("ws-fail@test.com", "pwd", "Name")
		assert.Error(t, err)
	})

	t.Run("add member failure", func(t *testing.T) {
		svc, mockUserRepo, mockWorkspaceRepo, _ := setupUserTest(t)
		userID := uuid.New()
		mockUserRepo.On("GetByEmail", mock.Anything).Return(models.User{}, sql.ErrNoRows)
		mockUserRepo.On("CreateUser", mock.Anything).Return(models.User{ID: userID}, nil)
		mockWorkspaceRepo.On("CreateWorkspace", mock.Anything).Return(models.Workspace{ID: uuid.New()}, nil)
		mockWorkspaceRepo.On("AddMember", mock.Anything).Return(errors.New("mem error"))
		_, _, err := svc.Register("mem-fail@test.com", "pwd", "Name")
		assert.Error(t, err)
	})
}

func TestUserService_Authenticate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockUserRepo, _, _ := setupUserTest(t)

		email := "test@example.com"
		password := "password123"
		hashedPassword, _ := utils.PasswordHash(password)
		userID := uuid.New()

		mockUserRepo.On("GetByEmail", email).Return(models.User{ID: userID, HashedPassword: hashedPassword}, nil)

		token, err := svc.Authenticate(email, password)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("user not found", func(t *testing.T) {
		svc, mockUserRepo, _, _ := setupUserTest(t)
		mockUserRepo.On("GetByEmail", mock.Anything).Return(models.User{}, sql.ErrNoRows)
		_, err := svc.Authenticate("none@test.com", "pwd")
		assert.Error(t, err)
	})

	t.Run("invalid password", func(t *testing.T) {
		svc, mockUserRepo, _, _ := setupUserTest(t)

		email := "test@example.com"
		hashedPassword, _ := utils.PasswordHash("correct")

		mockUserRepo.On("GetByEmail", email).Return(models.User{HashedPassword: hashedPassword}, nil)

		token, err := svc.Authenticate(email, "wrong")

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func TestUserService_GetByID(t *testing.T) {
	svc, mockUserRepo, _, _ := setupUserTest(t)
	userID := uuid.New()

	mockUserRepo.On("GetByID", userID).Return(models.User{ID: userID, Email: "test@example.com"}, nil)

	user, err := svc.GetByID(userID)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
}
