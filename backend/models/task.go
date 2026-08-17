package models

import (
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

// TaskTable is the name of the tasks table.
const TaskTable = "tasks"

// TaskStatus defines the status of a task.
type TaskStatus string

const (
	// TaskStatusTodo represents a task that is yet to be started.
	TaskStatusTodo TaskStatus = "TODO"
	// TaskStatusInProgress represents a task that is currently being worked on.
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	// TaskStatusDone represents a task that has been completed.
	TaskStatusDone TaskStatus = "DONE"
	// TaskStatusBacklog represents a task that is in the backlog.
	TaskStatusBacklog TaskStatus = "BACKLOG"
)

// TaskPriority defines the priority of a task.
type TaskPriority string

const (
	// TaskPriorityP0 represents critical priority.
	TaskPriorityP0 TaskPriority = "P0" // Critical
	// TaskPriorityP1 represents high priority.
	TaskPriorityP1 TaskPriority = "P1" // High
	// TaskPriorityP2 represents medium priority.
	TaskPriorityP2 TaskPriority = "P2" // Medium
	// TaskPriorityP3 represents low priority.
	TaskPriorityP3 TaskPriority = "P3" // Low
)

// Task represents a task in a project.
type Task struct {
	ID               uuid.UUID    `json:"id" db:"id"`
	Number           int          `json:"number" db:"number"`
	Title            string       `json:"title" db:"title"`
	Description      string       `json:"description" db:"description"`
	Status           TaskStatus   `json:"status" db:"status"`
	Priority         TaskPriority `json:"priority" db:"priority"`
	ProjectID        uuid.UUID    `json:"project_id" db:"project_id"`
	AssigneeID       *uuid.UUID   `json:"assignee_id" db:"assignee_id"` // Nullable
	AuthorID         *uuid.UUID   `json:"author_id" db:"author_id"`     // Nullable
	DueDate          *time.Time   `json:"due_date" db:"due_date"`       // Nullable
	EstimatedMinutes *int         `json:"estimated_minutes" db:"estimated_minutes"`
	CompletedAt      *time.Time   `json:"completed_at" db:"completed_at"`
	CreatedAt        time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at" db:"updated_at"`
}

// TaskWithDetails includes joined user info if needed,
// but for now, we'll keep it simple and load relations separately or via simple joins if critical.
// The Python version loads relationships. We can do a struct with embedded user info.

// TaskWithAssignee represents a task with assignee details.
type TaskWithAssignee struct {
	Task
	AssigneeEmail    *string `db:"assignee_email"`
	AssigneeFullName *string `db:"assignee_full_name"`
	AuthorEmail      *string `db:"author_email"`
	AuthorFullName   *string `db:"author_full_name"`
}

// TaskModel is the implementation of TaskRepository.
type TaskModel struct {
	db DbExecutor
}

// TaskRepository defines the interface for task data access.
type TaskRepository interface {
	Create(task Task) (Task, error)
	GetByID(id uuid.UUID) (TaskWithAssignee, error)
	Update(task Task) (Task, error)
	Delete(id uuid.UUID) error
	ListByProjectID(projectID uuid.UUID, status *TaskStatus, assigneeID *uuid.UUID) ([]TaskWithAssignee, error)
	ListByAssigneeID(assigneeID uuid.UUID) ([]TaskWithAssignee, error)
	GetNextTaskNumber(projectID uuid.UUID) (int, error)
	ListCompletedTasksInLastDays(projectID uuid.UUID, days int) ([]TaskWithAssignee, error)
}

// InitTaskModel initializes a new TaskModel.
func InitTaskModel(db DbExecutor) TaskRepository {
	return &TaskModel{
		db: db,
	}
}

// Create inserts a new task into the database.
func (model *TaskModel) Create(task Task) (Task, error) {
	var createdTask Task
	found, err := model.db.Insert(TaskTable).Rows(
		goqu.Record{
			"number":            task.Number,
			"title":             task.Title,
			"description":       task.Description,
			"status":            task.Status,
			"priority":          task.Priority,
			"project_id":        task.ProjectID,
			"assignee_id":       task.AssigneeID,
			"author_id":         task.AuthorID,
			"due_date":          task.DueDate,
			"estimated_minutes": task.EstimatedMinutes,
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

// GetNextTaskNumber returns the next available task number for a project.
func (model *TaskModel) GetNextTaskNumber(projectID uuid.UUID) (int, error) {
	var maxNum struct {
		Max sql.NullInt64 `db:"max"`
	}
	found, err := model.db.From(TaskTable).
		Select(goqu.MAX("number").As("max")).
		Where(goqu.Ex{"project_id": projectID}).
		ScanStruct(&maxNum)

	if err != nil {
		return 0, err
	}
	if !found || !maxNum.Max.Valid {
		return 1, nil
	}
	return int(maxNum.Max.Int64) + 1, nil
}

// GetByID retrieves a task by its ID.
func (model *TaskModel) GetByID(id uuid.UUID) (TaskWithAssignee, error) {
	var task TaskWithAssignee
	found, err := model.db.From(TaskTable).
		Select(
			"tasks.*",
			goqu.I("assignee.email").As("assignee_email"),
			goqu.I("assignee.full_name").As("assignee_full_name"),
			goqu.I("author.email").As("author_email"),
			goqu.I("author.full_name").As("author_full_name"),
		).
		LeftJoin(goqu.T("users").As("assignee"), goqu.On(goqu.Ex{"tasks.assignee_id": goqu.I("assignee.id")})).
		LeftJoin(goqu.T("users").As("author"), goqu.On(goqu.Ex{"tasks.author_id": goqu.I("author.id")})).
		Where(goqu.Ex{"tasks.id": id}).
		ScanStruct(&task)

	if err != nil {
		return task, err
	}
	if !found {
		return task, sql.ErrNoRows
	}
	return task, nil
}

// Update updates an existing task record.
func (model *TaskModel) Update(task Task) (Task, error) {
	// Only update fields that are usually mutable via this method
	// Or pass the full struct.
	// For partial updates, the Service layer usually merges changes.
	// Here we update everything passed in the struct except immutable IDs.

	record := goqu.Record{
		"title":             task.Title,
		"description":       task.Description,
		"status":            task.Status,
		"priority":          task.Priority,
		"assignee_id":       task.AssigneeID,
		"due_date":          task.DueDate,
		"estimated_minutes": task.EstimatedMinutes,
		"completed_at":      task.CompletedAt,
		"updated_at":        time.Now(),
	}

	_, err := model.db.Update(TaskTable).
		Set(record).
		Where(goqu.Ex{"id": task.ID}).
		Executor().Exec()

	if err != nil {
		return task, err
	}
	return task, nil
}

// Delete removes a task record by ID.
func (model *TaskModel) Delete(id uuid.UUID) error {
	_, err := model.db.Delete(TaskTable).
		Where(goqu.Ex{"id": id}).
		Executor().Exec()
	return err
}

// ListByProjectID lists all tasks in a project, optionally filtered by status and assignee.
func (model *TaskModel) ListByProjectID(projectID uuid.UUID, status *TaskStatus, assigneeID *uuid.UUID) ([]TaskWithAssignee, error) {
	query := model.db.From(TaskTable).
		Select(
			"tasks.*",
			goqu.I("assignee.email").As("assignee_email"),
			goqu.I("assignee.full_name").As("assignee_full_name"),
			goqu.I("author.email").As("author_email"),
			goqu.I("author.full_name").As("author_full_name"),
		).
		LeftJoin(goqu.T("users").As("assignee"), goqu.On(goqu.Ex{"tasks.assignee_id": goqu.I("assignee.id")})).
		LeftJoin(goqu.T("users").As("author"), goqu.On(goqu.Ex{"tasks.author_id": goqu.I("author.id")})).
		Where(goqu.Ex{"tasks.project_id": projectID})

	if status != nil {
		query = query.Where(goqu.Ex{"tasks.status": *status})
	}
	if assigneeID != nil {
		query = query.Where(goqu.Ex{"tasks.assignee_id": *assigneeID})
	}

	var tasks []TaskWithAssignee
	err := query.ScanStructs(&tasks)
	return tasks, err
}

// ListByAssigneeID lists all tasks assigned to a specific user.
func (model *TaskModel) ListByAssigneeID(assigneeID uuid.UUID) ([]TaskWithAssignee, error) {
	var tasks []TaskWithAssignee
	err := model.db.From(TaskTable).
		Select(
			"tasks.*",
			goqu.I("assignee.email").As("assignee_email"),
			goqu.I("assignee.full_name").As("assignee_full_name"),
			goqu.I("author.email").As("author_email"),
			goqu.I("author.full_name").As("author_full_name"),
		).
		LeftJoin(goqu.T("users").As("assignee"), goqu.On(goqu.Ex{"tasks.assignee_id": goqu.I("assignee.id")})).
		LeftJoin(goqu.T("users").As("author"), goqu.On(goqu.Ex{"tasks.author_id": goqu.I("author.id")})).
		Where(goqu.Ex{"tasks.assignee_id": assigneeID}).
		Order(
			goqu.I("tasks.priority").Asc(),
			goqu.I("tasks.due_date").Asc(),
		).
		ScanStructs(&tasks)
	return tasks, err
}

// ListCompletedTasksInLastDays returns tasks that were completed in the last N days.
func (model *TaskModel) ListCompletedTasksInLastDays(projectID uuid.UUID, days int) ([]TaskWithAssignee, error) {
	var tasks []TaskWithAssignee
	cutoffDate := time.Now().AddDate(0, 0, -days)

	err := model.db.From(TaskTable).
		Select(
			"tasks.*",
			goqu.I("assignee.email").As("assignee_email"),
			goqu.I("assignee.full_name").As("assignee_full_name"),
			goqu.I("author.email").As("author_email"),
			goqu.I("author.full_name").As("author_full_name"),
		).
		LeftJoin(goqu.T("users").As("assignee"), goqu.On(goqu.Ex{"tasks.assignee_id": goqu.I("assignee.id")})).
		LeftJoin(goqu.T("users").As("author"), goqu.On(goqu.Ex{"tasks.author_id": goqu.I("author.id")})).
		Where(
			goqu.Ex{
				"tasks.project_id": projectID,
				"tasks.status":     TaskStatusDone,
			},
			goqu.I("tasks.completed_at").Gte(cutoffDate),
		).
		Order(goqu.I("tasks.completed_at").Desc()).
		ScanStructs(&tasks)

	return tasks, err
}
