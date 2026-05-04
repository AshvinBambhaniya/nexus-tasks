package models

import (
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

const CommentTable = "comments"

type Comment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	TaskID    uuid.UUID `json:"task_id" db:"task_id"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CommentWithAuthor struct {
	Comment
	AuthorEmail    string `db:"author_email"`
	AuthorFullName string `db:"author_full_name"`
}

type CommentRepository interface {
	Create(comment Comment) (Comment, error)
	GetByID(id uuid.UUID) (Comment, error)
	ListByTaskID(taskID uuid.UUID) ([]CommentWithAuthor, error)
	Delete(id uuid.UUID) error
}

type CommentModel struct {
	db DbExecutor
}

func InitCommentModel(db DbExecutor) CommentRepository {
	return &CommentModel{
		db: db,
	}
}

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

func (model *CommentModel) Delete(id uuid.UUID) error {
	_, err := model.db.Delete(CommentTable).
		Where(goqu.Ex{"id": id}).
		Executor().Exec()
	return err
}
