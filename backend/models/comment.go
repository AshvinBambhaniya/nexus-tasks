// Package models contains the data models and database repositories for the application.
package models

import (
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

// CommentTable is the name of the table in the database
const CommentTable = "comments"

// Comment represents a comment on a task.
type Comment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	TaskID    uuid.UUID `json:"task_id" db:"task_id"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CommentWithAuthor represents a comment with author details.
type CommentWithAuthor struct {
	Comment
	AuthorEmail    string `db:"author_email"`
	AuthorFullName string `db:"author_full_name"`
}

// CommentRepository defines the interface for comment data access.
type CommentRepository interface {
	Create(comment Comment) (Comment, error)
	GetByID(id uuid.UUID) (Comment, error)
	ListByTaskID(taskID uuid.UUID) ([]CommentWithAuthor, error)
	ListByTaskIDs(taskIDs []uuid.UUID) ([]CommentWithAuthor, error)
	Delete(id uuid.UUID) error
}

// CommentModel is the implementation of CommentRepository.
type CommentModel struct {
	db DbExecutor
}

// InitCommentModel initializes a new CommentModel.
func InitCommentModel(db DbExecutor) CommentRepository {
	return &CommentModel{
		db: db,
	}
}

// Create creates a new comment.
func (model *CommentModel) Create(comment Comment) (Comment, error) {
	var createdComment Comment
	found, err := model.db.Insert(CommentTable).Rows(
		goqu.Record{
			"content":   comment.Content,
			"task_id":   comment.TaskID,
			"author_id": comment.AuthorID,
		},
	).Returning("*").Executor().ScanStruct(&createdComment)

	if err != nil {
		return comment, err
	}
	if !found {
		return comment, sql.ErrNoRows
	}
	return createdComment, nil
}

// GetByID returns a comment by its ID.
func (model *CommentModel) GetByID(id uuid.UUID) (Comment, error) {
	comment := Comment{}
	found, err := model.db.From(CommentTable).Where(goqu.Ex{"id": id}).ScanStruct(&comment)
	if err != nil {
		return comment, err
	}
	if !found {
		return comment, sql.ErrNoRows
	}
	return comment, nil
}

// ListByTaskID returns all comments for a specific task.
func (model *CommentModel) ListByTaskID(taskID uuid.UUID) ([]CommentWithAuthor, error) {
	var comments []CommentWithAuthor
	err := model.db.From(CommentTable).
		LeftJoin(goqu.T(UserTable), goqu.On(goqu.Ex{CommentTable + ".author_id": goqu.I(UserTable + ".id")})).
		Where(goqu.Ex{CommentTable + ".task_id": taskID}).
		Select(
			CommentTable+".*",
			goqu.I(UserTable+".email").As("author_email"),
			goqu.I(UserTable+".full_name").As("author_full_name"),
		).
		Order(goqu.I(CommentTable + ".created_at").Asc()).
		ScanStructs(&comments)

	return comments, err
}

// Delete removes a comment by its ID.
func (model *CommentModel) Delete(id uuid.UUID) error {
	_, err := model.db.Delete(CommentTable).
		Where(goqu.Ex{"id": id}).
		Executor().Exec()
	return err
}

// ListByTaskIDs returns all comments for a specific set of tasks.
func (model *CommentModel) ListByTaskIDs(taskIDs []uuid.UUID) ([]CommentWithAuthor, error) {
	var comments []CommentWithAuthor
	if len(taskIDs) == 0 {
		return comments, nil
	}

	err := model.db.From(CommentTable).
		LeftJoin(goqu.T(UserTable), goqu.On(goqu.Ex{CommentTable + ".author_id": goqu.I(UserTable + ".id")})).
		Where(goqu.Ex{CommentTable + ".task_id": taskIDs}).
		Select(
			CommentTable+".*",
			goqu.I(UserTable+".email").As("author_email"),
			goqu.I(UserTable+".full_name").As("author_full_name"),
		).
		Order(goqu.I(CommentTable + ".created_at").Asc()).
		ScanStructs(&comments)

	return comments, err
}
