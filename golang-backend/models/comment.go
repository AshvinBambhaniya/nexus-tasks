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

type CommentModel struct {
	db *goqu.Database
}

func InitCommentModel(goqu *goqu.Database) (CommentModel, error) {
	return CommentModel{
		db: goqu,
	}, nil
}

func (model *CommentModel) Create(transaction *goqu.TxDatabase, comment Comment) (Comment, error) {
	var createdComment Comment
	found, err := transaction.Insert(CommentTable).Rows(
		goqu.Record{
			"content":   comment.Content,
			"task_id":   comment.TaskID,
			"author_id": comment.AuthorID,
		},
	).Executor().ScanStruct(&createdComment)

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

func (model *CommentModel) Delete(transaction *goqu.TxDatabase, id uuid.UUID) error {
	_, err := transaction.Delete(CommentTable).
		Where(goqu.Ex{"id": id}).
		Executor().Exec()
	return err
}
