package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/config"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/pkg/jwt"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UserService defines the interface for user-related business logic
type UserService interface {
	Register(email, password, fullName string) (models.User, string, error)
	Authenticate(email, password string) (string, error)
	GetByID(userID uuid.UUID) (models.User, error)
	ProvisionSSOUser(email, fullName, provider, providerId string) (models.User, string, error)
}

type userService struct {
	storage models.Storage
	logger  *zap.Logger
	config  *config.AppConfig
}

// NewUserService creates a new user service instance
func NewUserService(storage models.Storage, logger *zap.Logger, cfg *config.AppConfig) UserService {
	return &userService{
		storage: storage,
		logger:  logger,
		config:  cfg,
	}
}

func (s *userService) Register(email, password, fullName string) (models.User, string, error) {

	// Check if user exists
	_, err := s.storage.Users().GetByEmail(email)
	if err == nil {
		return models.User{}, "", errors.New("email already registered")
	}

	//  Hash Password
	hashedPwd, err := utils.PasswordHash(password)
	if err != nil {
		return models.User{}, "", err
	}

	var createdUser models.User
	var token string

	err = s.storage.Atomic(context.Background(), func(txStorage models.Storage) error {
		user := models.User{
			Email:          email,
			FullName:       fullName,
			HashedPassword: hashedPwd,
			IsActive:       true,
		}

		// A. Create User
		var err error
		createdUser, err = txStorage.Users().CreateUser(user)
		if err != nil {
			s.logger.Error("failed to create user", zap.Error(err))
			return err
		}

		// B. Create Personal Workspace
		personalWs := models.Workspace{
			Name:    constants.WorkspaceTypePersonal,
			Type:    models.WorkspaceTypePersonal,
			OwnerID: createdUser.ID,
		}
		createdWs, err := txStorage.Workspaces().CreateWorkspace(personalWs)
		if err != nil {
			s.logger.Error("failed to create personal workspace", zap.Error(err))
			return err
		}

		// C. Add Member (Admin)
		err = txStorage.Workspaces().AddMember(models.WorkspaceMember{
			WorkspaceID: createdWs.ID,
			UserID:      createdUser.ID,
			Role:        models.WorkspaceRoleAdmin,
		})
		if err != nil {
			s.logger.Error("failed to add user to personal workspace", zap.Error(err))
			return err
		}

		// Generate Token
		token, err = jwt.CreateToken(s.config.Secret, createdUser.ID.String(), time.Now().Add(time.Duration(s.config.JwtExpirationHours)*time.Hour))
		if err != nil {
			s.logger.Error("failed to generate token", zap.Error(err))
			return err
		}

		return nil
	})

	if err != nil {
		return models.User{}, "", err
	}

	return createdUser, token, nil
}

func (s *userService) Authenticate(email, password string) (string, error) {
	user, err := s.storage.Users().GetByEmail(email)
	if err != nil {
		s.logger.Debug("failed to get user by email", zap.String("email", email), zap.Error(err))
		return "", err
	}

	if !utils.CheckPasswordHash(password, user.HashedPassword) {
		return "", sql.ErrNoRows // Treat as not found/invalid
	}

	// Generate Token
	// Using ID (UUID) as Subject.
	token, err := jwt.CreateToken(s.config.Secret, user.ID.String(), time.Now().Add(time.Duration(s.config.JwtExpirationHours)*time.Hour))
	if err != nil {
		s.logger.Error("failed to generate token", zap.Error(err))
		return "", err
	}

	return token, nil
}

func (s *userService) GetByID(userID uuid.UUID) (models.User, error) {
	return s.storage.Users().GetByID(userID)
}

func (s *userService) ProvisionSSOUser(email, fullName, provider, providerId string) (models.User, string, error) {
	var user models.User

	// 1. Check for existing linked identity
	identity, err := s.storage.UserIdentities().GetByProvider(provider, providerId)
	if err == nil {
		user, _ = s.storage.Users().GetByID(identity.UserID)
		token, err := jwt.CreateToken(s.config.Secret, user.ID.String(), time.Now().Add(time.Duration(s.config.JwtExpirationHours)*time.Hour))
		if err != nil {
			return models.User{}, "", err
		}
		return user, token, nil
	}

	err = s.storage.Atomic(context.Background(), func(tx models.Storage) error {
		// 2. Fallback: Check if local user exists by email for first-time link
		existingUser, err := tx.Users().GetByEmail(email)

		if err == nil {
			// 3. Link account: User exists locally, link this new SSO identity
			user = existingUser
		} else {
			// 4. Provision: User does not exist, create them
			newUser := models.User{
				Email:    email,
				FullName: fullName,
				IsActive: true,
			}
			user, err = tx.Users().CreateUser(newUser)
			if err != nil {
				return err
			}

			// Also create a default personal workspace for the new user
			personalWs := models.Workspace{
				Name:    string(models.WorkspaceTypePersonal),
				Type:    models.WorkspaceTypePersonal,
				OwnerID: user.ID,
			}
			createdWs, err := tx.Workspaces().CreateWorkspace(personalWs)
			if err != nil {
				return err
			}
			err = tx.Workspaces().AddMember(models.WorkspaceMember{
				WorkspaceID: createdWs.ID,
				UserID:      user.ID,
				Role:        models.WorkspaceRoleAdmin,
			})
			if err != nil {
				return err
			}
		}

		// Insert the identity record mapping the SSO provider to the local User ID
		userIdentity := models.UserIdentity{
			UserID:     user.ID,
			Provider:   provider,
			ProviderID: providerId,
		}
		_, err = tx.UserIdentities().Create(userIdentity)
		return err
	})

	if err != nil {
		s.logger.Error("failed to provision sso user", zap.Error(err))
		return models.User{}, "", err
	}

	token, err := jwt.CreateToken(s.config.Secret, user.ID.String(), time.Now().Add(time.Duration(s.config.JwtExpirationHours)*time.Hour))
	if err != nil {
		return models.User{}, "", err
	}

	return user, token, nil
}
