package models

import (
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

const TaskTable = "tasks"

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "TODO"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusDone       TaskStatus = "DONE"
	TaskStatusBacklog    TaskStatus = "BACKLOG"
)

type TaskPriority string

const (
	TaskPriorityP0 TaskPriority = "P0" // Critical
	TaskPriorityP1 TaskPriority = "P1" // High
	TaskPriorityP2 TaskPriority = "P2" // Medium
	TaskPriorityP3 TaskPriority = "P3" // Low
)

type Task struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	Title       string       `json:"title" db:"title"`
	Description string       `json:"description" db:"description"`
	Status      TaskStatus   `json:"status" db:"status"`
	Priority    TaskPriority `json:"priority" db:"priority"`
	ProjectID   uuid.UUID    `json:"project_id" db:"project_id"`
	AssigneeID  *uuid.UUID   `json:"assignee_id" db:"assignee_id"` // Nullable
	AuthorID    *uuid.UUID   `json:"author_id" db:"author_id"`     // Nullable
	DueDate     *time.Time   `json:"due_date" db:"due_date"`       // Nullable
	CompletedAt *time.Time   `json:"completed_at" db:"completed_at"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}

// TaskWithDetails includes joined user info if needed,
// but for now, we'll keep it simple and load relations separately or via simple joins if critical.
// The Python version loads relationships. We can do a struct with embedded user info.

type TaskWithAssignee struct {
	Task
	AssigneeEmail    *string `db:"assignee_email"`
	AssigneeFullName *string `db:"assignee_full_name"`
}

type TaskModel struct {
	db *goqu.Database
}

func InitTaskModel(goqu *goqu.Database) (TaskModel, error) {
	return TaskModel{
		db: goqu,
	}, nil
}

func (model *TaskModel) Create(transaction *goqu.TxDatabase, task Task) (Task, error) {
	var createdTask Task
	found, err := transaction.Insert(TaskTable).Rows(
		goqu.Record{
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"priority":    task.Priority,
			"project_id":  task.ProjectID,
			"assignee_id": task.AssigneeID,
			"author_id":   task.AuthorID,
			"due_date":    task.DueDate,
		},
	).Returning("*").Executor().ScanStruct(&createdTask)

	if err != nil {
		return task, err
	}
	if !found {
		return task, sql.ErrNoRows
	}
	return createdTask, nil
}

func (model *TaskModel) GetByID(id uuid.UUID) (Task, error) {
	task := Task{}
	found, err := model.db.From(TaskTable).Where(goqu.Ex{"id": id}).ScanStruct(&task)
	if err != nil {
		return task, err
	}
	if !found {
		return task, sql.ErrNoRows
	}
	return task, nil
}

func (model *TaskModel) Update(transaction *goqu.TxDatabase, task Task) (Task, error) {
	// Only update fields that are usually mutable via this method
	// Or pass the full struct.
	// For partial updates, the Service layer usually merges changes.
	// Here we update everything passed in the struct except immutable IDs.

	record := goqu.Record{
		"title":        task.Title,
		"description":  task.Description,
		"status":       task.Status,
		"priority":     task.Priority,
		"assignee_id":  task.AssigneeID,
		"due_date":     task.DueDate,
		"completed_at": task.CompletedAt,
		"updated_at":   time.Now(),
	}

	_, err := transaction.Update(TaskTable).
		Set(record).
		Where(goqu.Ex{"id": task.ID}).
		Executor().Exec()

	if err != nil {
		return task, err
	}
	return task, nil
}

func (model *TaskModel) Delete(transaction *goqu.TxDatabase, id uuid.UUID) error {
	_, err := transaction.Delete(TaskTable).
		Where(goqu.Ex{"id": id}).
		Executor().Exec()
	return err
}

func (model *TaskModel) ListByProjectID(projectID uuid.UUID, status *TaskStatus, assigneeID *uuid.UUID) ([]Task, error) {
	query := model.db.From(TaskTable).Where(goqu.Ex{"project_id": projectID})

	if status != nil {
		query = query.Where(goqu.Ex{"status": *status})
	}
	if assigneeID != nil {
		query = query.Where(goqu.Ex{"assignee_id": *assigneeID})
	}

	var tasks []Task
	err := query.ScanStructs(&tasks)
	return tasks, err
}

func (model *TaskModel) ListByAssigneeID(assigneeID uuid.UUID) ([]Task, error) {
	var tasks []Task
	err := model.db.From(TaskTable).
		Where(goqu.Ex{"assignee_id": assigneeID}).
		Order(
			goqu.I("priority").Asc(),
			goqu.I("due_date").Asc(),
		).
		ScanStructs(&tasks)
	return tasks, err
}
