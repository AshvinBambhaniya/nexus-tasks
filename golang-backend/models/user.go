package models

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

// UserTable represent table name
const UserTable = "users"

// User model
type User struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Email          string    `json:"email" db:"email"`
	FullName       string    `json:"full_name" db:"full_name"`
	HashedPassword string    `json:"-" db:"hashed_password"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      string    `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt      string    `json:"updated_at,omitempty" db:"updated_at"`
}

// UserModel implements user related database operations
type UserModel struct {
	db *goqu.Database
}

// InitUserModel Init model
func InitUserModel(goqu *goqu.Database) (UserModel, error) {
	return UserModel{
		db: goqu,
	}, nil
}

// GetUserByEmail get user by email
func (model *UserModel) GetByEmail(email string) (User, error) {
	user := User{}
	found, err := model.db.From(UserTable).Where(goqu.Ex{
		"email": email,
	}).ScanStruct(&user)

	if err != nil {
		return user, err
	}

	if !found {
		return user, sql.ErrNoRows
	}

	return user, nil
}

// GetByID get user by id
func (model *UserModel) GetByID(id uuid.UUID) (User, error) {
	user := User{}
	found, err := model.db.From(UserTable).Where(goqu.Ex{
		"id": id,
	}).ScanStruct(&user)

	if err != nil {
		return user, err
	}

	if !found {
		return user, sql.ErrNoRows
	}

	return user, nil
}

// CreateUser inserts a new user
func (model *UserModel) CreateUser(transaction *goqu.TxDatabase, user User) (User, error) {
	var createdUser User

	found, err := transaction.Insert(UserTable).Rows(
		goqu.Record{
			"email":           user.Email,
			"full_name":       user.FullName,
			"hashed_password": user.HashedPassword,
			"is_active":       user.IsActive,
		},
	).Returning("*").Executor().ScanStruct(&createdUser)

	if err != nil {
		return user, err
	}
	if !found {
		return user, sql.ErrNoRows
	}
	return createdUser, nil
}
