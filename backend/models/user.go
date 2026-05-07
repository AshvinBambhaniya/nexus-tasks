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

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	GetByEmail(email string) (User, error)
	GetByID(id uuid.UUID) (User, error)
	CreateUser(user User) (User, error)
}

// UserModel implements user related database operations
type UserModel struct {
	db DbExecutor
}

// InitUserModel initializes the user model
func InitUserModel(db DbExecutor) UserRepository {
	return &UserModel{
		db: db,
	}
}

// GetByEmail retrieves a user by their email address
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

// GetByID retrieves a user by their unique ID
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

// CreateUser inserts a new user into the database
func (model *UserModel) CreateUser(user User) (User, error) {
	var createdUser User

	found, err := model.db.Insert(UserTable).Rows(
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
