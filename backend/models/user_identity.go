package models

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

const UserIdentityTable = "user_identities"

type UserIdentity struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	Provider   string    `json:"provider" db:"provider"`
	ProviderID string    `json:"provider_id" db:"provider_id"`
	CreatedAt  string    `json:"created_at,omitempty" db:"created_at"`
}

type UserIdentityRepository interface {
	GetByProvider(provider, providerID string) (UserIdentity, error)
	Create(identity UserIdentity) (UserIdentity, error)
}

type UserIdentityModel struct {
	db DbExecutor
}

func InitUserIdentityModel(db DbExecutor) UserIdentityRepository {
	return &UserIdentityModel{db: db}
}

func (model *UserIdentityModel) GetByProvider(provider, providerID string) (UserIdentity, error) {
	var identity UserIdentity
	found, err := model.db.From(UserIdentityTable).Where(goqu.Ex{
		"provider":    provider,
		"provider_id": providerID,
	}).ScanStruct(&identity)

	if err != nil {
		return identity, err
	}
	if !found {
		return identity, sql.ErrNoRows
	}
	return identity, nil
}

func (model *UserIdentityModel) Create(identity UserIdentity) (UserIdentity, error) {
	var createdIdentity UserIdentity
	found, err := model.db.Insert(UserIdentityTable).Rows(
		goqu.Record{
			"user_id":     identity.UserID,
			"provider":    identity.Provider,
			"provider_id": identity.ProviderID,
		},
	).Returning("*").Executor().ScanStruct(&createdIdentity)

	if err != nil {
		return identity, err
	}
	if !found {
		return identity, sql.ErrNoRows
	}
	return createdIdentity, nil
}
