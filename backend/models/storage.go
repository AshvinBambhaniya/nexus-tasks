package models

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
)

// DbExecutor is an interface that both goqu.Database and goqu.TxDatabase implement.
type DbExecutor interface {
	From(from ...interface{}) *goqu.SelectDataset
	Insert(table interface{}) *goqu.InsertDataset
	Update(table interface{}) *goqu.UpdateDataset
	Delete(table interface{}) *goqu.DeleteDataset
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// Storage defines the interface for the entire data access layer.
type Storage interface {
	Users() UserRepository
	Workspaces() WorkspaceRepository
	Teams() TeamRepository
	Projects() ProjectRepository
	Tasks() TaskRepository
	Comments() CommentRepository
	Notifications() NotificationRepository

	Atomic(ctx context.Context, fn func(Storage) error) error
	CheckHealth(ctx context.Context) error
}

type storage struct {
	db DbExecutor
}

// NewStorage creates a new storage instance
func NewStorage(db *goqu.Database) Storage {
	return &storage{db: db}
}

func (s *storage) Users() UserRepository {
	return InitUserModel(s.db)
}

func (s *storage) Workspaces() WorkspaceRepository {
	return InitWorkspaceModel(s.db)
}

func (s *storage) Teams() TeamRepository {
	return InitTeamModel(s.db)
}

func (s *storage) Projects() ProjectRepository {
	return InitProjectModel(s.db)
}

func (s *storage) Tasks() TaskRepository {
	return InitTaskModel(s.db)
}

func (s *storage) Comments() CommentRepository {
	return InitCommentModel(s.db)
}

func (s *storage) Notifications() NotificationRepository {
	return InitNotificationModel(s.db)
}

func (s *storage) Atomic(ctx context.Context, fn func(Storage) error) error {
	var tx *goqu.TxDatabase
	var err error

	if _, ok := s.db.(*goqu.TxDatabase); ok {
		// Already in a transaction, just use it
		return fn(s)
	}

	if d, ok := s.db.(*goqu.Database); ok {
		tx, err = d.BeginTx(ctx, nil)
	} else {
		return fmt.Errorf("unsupported database executor type")
	}

	if err != nil {
		return err
	}

	txStorage := &storage{db: tx}
	if err := fn(txStorage); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction rollback failed: %v, original error: %w", rbErr, err)
		}
		return err
	}

	return tx.Commit()
}

func (s *storage) CheckHealth(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "SELECT 1")
	return err
}
