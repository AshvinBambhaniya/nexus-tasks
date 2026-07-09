package models

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

// APIKeyTable represents the table name
const APIKeyTable = "api_keys"

// APIKey represents a personal access token for external integrations.
type APIKey struct {
	ID          uuid.UUID  `json:"id" db:"id" goqu:"skipinsert"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	Name        string     `json:"name" db:"name"`
	TokenPrefix string     `json:"token_prefix" db:"token_prefix"`
	TokenHash   string     `json:"-" db:"token_hash"`
	LastUsedAt  *time.Time `json:"last_used_at" db:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at" goqu:"skipinsert"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at" goqu:"skipinsert"`
}

// APIKeyRepository defines the interface for api key database operations.
type APIKeyRepository interface {
	Create(apiKey APIKey) (APIKey, error)
	GetByTokenHash(hash string) (APIKey, error)
	ListByUserID(userID uuid.UUID) ([]APIKey, error)
	Delete(id uuid.UUID, userID uuid.UUID) error
	UpdateLastUsed(id uuid.UUID) error
}

type apiKeyModel struct {
	db DbExecutor
}

// InitAPIKeyModel initializes the api key model.
func InitAPIKeyModel(db DbExecutor) APIKeyRepository {
	return &apiKeyModel{db: db}
}

// Create inserts a new API key record into the database.
func (m *apiKeyModel) Create(apiKey APIKey) (APIKey, error) {
	var created APIKey

	found, err := m.db.Insert(APIKeyTable).Rows(
		goqu.Record{
			"user_id":      apiKey.UserID,
			"name":         apiKey.Name,
			"token_prefix": apiKey.TokenPrefix,
			"token_hash":   apiKey.TokenHash,
			"expires_at":   apiKey.ExpiresAt,
		},
	).Returning("*").Executor().ScanStruct(&created)

	if err != nil {
		return apiKey, err
	}
	if !found {
		return apiKey, err
	}

	return created, nil
}

// GetByTokenHash retrieves an API key by its SHA-256 token hash.
func (m *apiKeyModel) GetByTokenHash(hash string) (APIKey, error) {
	var apiKey APIKey

	found, err := m.db.From(APIKeyTable).Where(goqu.Ex{
		"token_hash": hash,
	}).ScanStruct(&apiKey)

	if err != nil {
		return apiKey, err
	}
	if !found {
		return apiKey, err
	}

	return apiKey, nil
}

// ListByUserID retrieves all API keys for a given user, ordered by creation date.
func (m *apiKeyModel) ListByUserID(userID uuid.UUID) ([]APIKey, error) {
	var keys []APIKey

	err := m.db.From(APIKeyTable).
		Where(goqu.Ex{"user_id": userID}).
		Order(goqu.I("created_at").Desc()).
		ScanStructs(&keys)

	if err != nil {
		return nil, err
	}

	return keys, nil
}

// Delete removes an API key by ID, scoped to the owning user.
func (m *apiKeyModel) Delete(id uuid.UUID, userID uuid.UUID) error {
	_, err := m.db.Delete(APIKeyTable).
		Where(goqu.Ex{"id": id, "user_id": userID}).
		Executor().Exec()

	return err
}

// UpdateLastUsed sets the last_used_at timestamp to the current time.
func (m *apiKeyModel) UpdateLastUsed(id uuid.UUID) error {
	_, err := m.db.Update(APIKeyTable).
		Set(goqu.Record{"last_used_at": goqu.L("CURRENT_TIMESTAMP"), "updated_at": goqu.L("CURRENT_TIMESTAMP")}).
		Where(goqu.Ex{"id": id}).
		Executor().Exec()

	return err
}
