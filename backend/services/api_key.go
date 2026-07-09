package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/constants"
	"github.com/AshvinBambhaniya/nexus-tasks/v2/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// APIKeyService provides methods for managing personal access tokens.
type APIKeyService interface {
	GenerateKey(userID uuid.UUID, name string) (rawToken string, key models.APIKey, err error)
	ValidateToken(rawToken string) (models.User, uuid.UUID, error)
	ListKeys(userID uuid.UUID) ([]models.APIKey, error)
	RevokeKey(id uuid.UUID, userID uuid.UUID) error
}

type apiKeyService struct {
	storage models.Storage
	logger  *zap.Logger
}

// NewAPIKeyService returns a new instance of APIKeyService.
func NewAPIKeyService(storage models.Storage, logger *zap.Logger) APIKeyService {
	return &apiKeyService{
		storage: storage,
		logger:  logger,
	}
}

// GenerateKey creates a new personal access token for the given user.
// The raw token is returned once and never stored; only its SHA-256 hash is persisted.
func (s *apiKeyService) GenerateKey(userID uuid.UUID, name string) (string, models.APIKey, error) {
	// Generate 20 cryptographically random bytes → 40 hex characters
	randomBytes := make([]byte, 20)
	if _, err := rand.Read(randomBytes); err != nil {
		s.logger.Error("failed to generate random bytes for api key", zap.Error(err))
		return "", models.APIKey{}, fmt.Errorf("failed to generate token: %w", err)
	}

	rawToken := constants.APIKeyPrefix + hex.EncodeToString(randomBytes)
	tokenHash := hashToken(rawToken)
	tokenPrefix := rawToken[:8]

	apiKey := models.APIKey{
		UserID:      userID,
		Name:        name,
		TokenPrefix: tokenPrefix,
		TokenHash:   tokenHash,
	}

	created, err := s.storage.APIKeys().Create(apiKey)
	if err != nil {
		s.logger.Error("failed to create api key", zap.Error(err))
		return "", models.APIKey{}, fmt.Errorf("failed to create api key: %w", err)
	}

	return rawToken, created, nil
}

// ValidateToken verifies a raw token by hashing it and looking up the hash.
// Returns the associated user and the api key ID if valid.
func (s *apiKeyService) ValidateToken(rawToken string) (models.User, uuid.UUID, error) {
	tokenHash := hashToken(rawToken)

	apiKey, err := s.storage.APIKeys().GetByTokenHash(tokenHash)
	if err != nil {
		return models.User{}, uuid.Nil, errors.New("invalid api key")
	}

	// Check if the key is zero-value (not found)
	if apiKey.ID == uuid.Nil {
		return models.User{}, uuid.Nil, errors.New("invalid api key")
	}

	// Check expiry
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return models.User{}, uuid.Nil, errors.New("api key expired")
	}

	// Update last used timestamp (fire-and-forget)
	go func() {
		if err := s.storage.APIKeys().UpdateLastUsed(apiKey.ID); err != nil {
			s.logger.Warn("failed to update api key last_used_at", zap.Error(err))
		}
	}()

	// Fetch the user associated with this key
	user, err := s.storage.Users().GetByID(apiKey.UserID)
	if err != nil {
		s.logger.Error("failed to fetch user for api key", zap.Error(err))
		return models.User{}, uuid.Nil, errors.New("user not found for api key")
	}

	return user, apiKey.ID, nil
}

// ListKeys retrieves all API keys for a given user.
func (s *apiKeyService) ListKeys(userID uuid.UUID) ([]models.APIKey, error) {
	keys, err := s.storage.APIKeys().ListByUserID(userID)
	if err != nil {
		s.logger.Error("failed to list api keys", zap.Error(err))
		return nil, err
	}

	return keys, nil
}

// RevokeKey deletes an API key, scoped to the owning user.
func (s *apiKeyService) RevokeKey(id uuid.UUID, userID uuid.UUID) error {
	err := s.storage.APIKeys().Delete(id, userID)
	if err != nil {
		s.logger.Error("failed to revoke api key", zap.Error(err))
		return err
	}

	return nil
}

// hashToken computes the SHA-256 hash of a raw token string and returns it as a hex string.
func hashToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(hash[:])
}
