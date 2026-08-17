package models

import (
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

// TimeEntryTable is the name of the time_entries table.
const TimeEntryTable = "time_entries"

// TimeEntry represents a time log for a task.
type TimeEntry struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	TaskID          uuid.UUID  `json:"task_id" db:"task_id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	WorkspaceID     uuid.UUID  `json:"workspace_id" db:"workspace_id"`
	Description     string     `json:"description" db:"description"`
	StartTime       time.Time  `json:"start_time" db:"start_time"`
	EndTime         *time.Time `json:"end_time" db:"end_time"`
	DurationMinutes *int       `json:"duration_minutes" db:"duration_minutes"`
	IsManual        bool       `json:"is_manual" db:"is_manual"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// TimeEntryWithDetails includes joined user and task info.
type TimeEntryWithDetails struct {
	TimeEntry
	UserFullName string `db:"user_full_name" json:"user_full_name"`
	UserEmail    string `db:"user_email" json:"user_email"`
	TaskTitle    string `db:"task_title" json:"task_title"`
	TaskNumber   int    `db:"task_number" json:"task_number"`
}

// TaskTimeSummary represents time metrics for a specific task.
type TaskTimeSummary struct {
	TaskID           uuid.UUID `json:"task_id" db:"task_id"`
	TaskNumber       int       `json:"task_number" db:"task_number"`
	TaskTitle        string    `json:"task_title" db:"task_title"`
	EstimatedMinutes *int      `json:"estimated_minutes" db:"estimated_minutes"`
	LoggedMinutes    int       `json:"logged_minutes" db:"logged_minutes"`
	IsOverBudget     bool      `json:"is_over_budget"`
}

// MemberTimeSummary represents time metrics for a specific team member.
type MemberTimeSummary struct {
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	FullName      string    `json:"full_name" db:"full_name"`
	LoggedMinutes int       `json:"logged_minutes" db:"logged_minutes"`
}

// DailyTimeEntry represents total time logged on a specific day.
type DailyTimeEntry struct {
	Date          string `json:"date" db:"date"`
	LoggedMinutes int    `json:"logged_minutes" db:"logged_minutes"`
}

// ProjectTimeAnalytics holds aggregated time metrics for a project.
type ProjectTimeAnalytics struct {
	TotalEstimatedMinutes int                 `json:"total_estimated_minutes"`
	TotalLoggedMinutes    int                 `json:"total_logged_minutes"`
	EstimateAccuracy      int                 `json:"estimate_accuracy_percent"`
	OverBudgetTaskCount   int                 `json:"over_budget_task_count"`
	ByTask                []TaskTimeSummary   `json:"by_task"`
	ByMember              []MemberTimeSummary `json:"by_member"`
	DailyTrend            []DailyTimeEntry    `json:"daily_trend"`
}

// TimeEntryRepository defines the interface for time entry data access.
type TimeEntryRepository interface {
	Create(entry TimeEntry) (TimeEntry, error)
	GetByID(id uuid.UUID) (TimeEntry, error)
	GetActiveByUserID(userID uuid.UUID) (*TimeEntryWithDetails, error)
	ListByTaskID(taskID uuid.UUID) ([]TimeEntryWithDetails, error)
	Update(entry TimeEntry) error
	Delete(id uuid.UUID) error
	SumByTaskID(taskID uuid.UUID) (int, error)
	GetProjectAnalytics(projectID uuid.UUID) (*ProjectTimeAnalytics, error)
	ListByProjectID(projectID uuid.UUID, userID *uuid.UUID, startDate, endDate *time.Time) ([]TimeEntryWithDetails, error)
}

// TimeEntryModel is the implementation of TimeEntryRepository.
type TimeEntryModel struct {
	db DbExecutor
}

// InitTimeEntryModel initializes a new TimeEntryModel.
func InitTimeEntryModel(db DbExecutor) TimeEntryRepository {
	return &TimeEntryModel{
		db: db,
	}
}

// Create inserts a new time entry into the database.
func (model *TimeEntryModel) Create(entry TimeEntry) (TimeEntry, error) {
	var createdEntry TimeEntry
	found, err := model.db.Insert(TimeEntryTable).Rows(
		goqu.Record{
			"task_id":          entry.TaskID,
			"user_id":          entry.UserID,
			"workspace_id":     entry.WorkspaceID,
			"description":      entry.Description,
			"start_time":       entry.StartTime,
			"end_time":         entry.EndTime,
			"duration_minutes": entry.DurationMinutes,
			"is_manual":        entry.IsManual,
		},
	).Returning("*").Executor().ScanStruct(&createdEntry)

	if err != nil {
		return entry, err
	}
	if !found {
		return entry, sql.ErrNoRows
	}
	return createdEntry, nil
}

// GetByID retrieves a time entry by its ID.
func (model *TimeEntryModel) GetByID(id uuid.UUID) (TimeEntry, error) {
	var entry TimeEntry
	found, err := model.db.From(TimeEntryTable).
		Where(goqu.Ex{"id": id}).
		ScanStruct(&entry)

	if err != nil {
		return entry, err
	}
	if !found {
		return entry, sql.ErrNoRows
	}
	return entry, nil
}

// GetActiveByUserID retrieves the active timer for a user.
func (model *TimeEntryModel) GetActiveByUserID(userID uuid.UUID) (*TimeEntryWithDetails, error) {
	var entry TimeEntryWithDetails
	found, err := model.db.From(TimeEntryTable).
		Select(
			"time_entries.*",
			goqu.I("users.full_name").As("user_full_name"),
			goqu.I("users.email").As("user_email"),
			goqu.I("tasks.title").As("task_title"),
			goqu.I("tasks.number").As("task_number"),
		).
		LeftJoin(goqu.T("users"), goqu.On(goqu.Ex{"time_entries.user_id": goqu.I("users.id")})).
		LeftJoin(goqu.T("tasks"), goqu.On(goqu.Ex{"time_entries.task_id": goqu.I("tasks.id")})).
		Where(
			goqu.Ex{"time_entries.user_id": userID},
			goqu.Ex{"time_entries.end_time": nil},
		).
		ScanStruct(&entry)

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &entry, nil
}

// ListByTaskID lists all time entries for a task.
func (model *TimeEntryModel) ListByTaskID(taskID uuid.UUID) ([]TimeEntryWithDetails, error) {
	var entries []TimeEntryWithDetails
	err := model.db.From(TimeEntryTable).
		Select(
			"time_entries.*",
			goqu.I("users.full_name").As("user_full_name"),
			goqu.I("users.email").As("user_email"),
		).
		LeftJoin(goqu.T("users"), goqu.On(goqu.Ex{"time_entries.user_id": goqu.I("users.id")})).
		Where(goqu.Ex{"time_entries.task_id": taskID}).
		Order(goqu.I("time_entries.created_at").Desc()).
		ScanStructs(&entries)

	if err != nil {
		return nil, err
	}
	return entries, nil
}

// Update updates an existing time entry record.
func (model *TimeEntryModel) Update(entry TimeEntry) error {
	record := goqu.Record{
		"description":      entry.Description,
		"end_time":         entry.EndTime,
		"duration_minutes": entry.DurationMinutes,
		"updated_at":       time.Now(),
	}

	_, err := model.db.Update(TimeEntryTable).
		Set(record).
		Where(goqu.Ex{"id": entry.ID}).
		Executor().Exec()

	return err
}

// Delete removes a time entry record by ID.
func (model *TimeEntryModel) Delete(id uuid.UUID) error {
	_, err := model.db.Delete(TimeEntryTable).
		Where(goqu.Ex{"id": id}).
		Executor().Exec()
	return err
}

// SumByTaskID returns the sum of duration_minutes for a task.
func (model *TimeEntryModel) SumByTaskID(taskID uuid.UUID) (int, error) {
	var total struct {
		Sum sql.NullInt64 `db:"sum"`
	}
	found, err := model.db.From(TimeEntryTable).
		Select(goqu.SUM("duration_minutes").As("sum")).
		Where(
			goqu.Ex{"task_id": taskID},
			goqu.I("end_time").IsNotNull(),
		).
		ScanStruct(&total)

	if err != nil {
		return 0, err
	}
	if !found || !total.Sum.Valid {
		return 0, nil
	}
	return int(total.Sum.Int64), nil
}

// GetProjectAnalytics returns aggregated time analytics for a project.
func (model *TimeEntryModel) GetProjectAnalytics(projectID uuid.UUID) (*ProjectTimeAnalytics, error) {
	var analytics ProjectTimeAnalytics
	analytics.ByTask = []TaskTimeSummary{}
	analytics.ByMember = []MemberTimeSummary{}
	analytics.DailyTrend = []DailyTimeEntry{}

	// 1. Task Time Summary
	err := model.db.From(goqu.T("tasks")).
		Select(
			goqu.I("tasks.id").As("task_id"),
			goqu.I("tasks.number").As("task_number"),
			goqu.I("tasks.title").As("task_title"),
			goqu.I("tasks.estimated_minutes").As("estimated_minutes"),
			goqu.COALESCE(goqu.SUM("time_entries.duration_minutes"), 0).As("logged_minutes"),
		).
		LeftJoin(goqu.T("time_entries"), goqu.On(
			goqu.Ex{"tasks.id": goqu.I("time_entries.task_id")},
			goqu.I("time_entries.end_time").IsNotNull(),
		)).
		Where(goqu.Ex{"tasks.project_id": projectID}).
		GroupBy("tasks.id", "tasks.number", "tasks.title", "tasks.estimated_minutes").
		ScanStructs(&analytics.ByTask)

	if err != nil {
		return nil, err
	}

	// Calculate Task metrics
	for i := range analytics.ByTask {
		if analytics.ByTask[i].EstimatedMinutes != nil {
			analytics.TotalEstimatedMinutes += *analytics.ByTask[i].EstimatedMinutes
			if analytics.ByTask[i].LoggedMinutes > *analytics.ByTask[i].EstimatedMinutes {
				analytics.ByTask[i].IsOverBudget = true
				analytics.OverBudgetTaskCount++
			}
		}
		analytics.TotalLoggedMinutes += analytics.ByTask[i].LoggedMinutes
	}

	if analytics.TotalEstimatedMinutes > 0 {
		ratio := float64(analytics.TotalLoggedMinutes) / float64(analytics.TotalEstimatedMinutes)
		if ratio > 1 {
			ratio = 1 / ratio
		}
		analytics.EstimateAccuracy = int(ratio * 100)
	}

	// 2. Member Time Summary
	err = model.db.From(goqu.T("time_entries")).
		Select(
			goqu.I("time_entries.user_id").As("user_id"),
			goqu.I("users.full_name").As("full_name"),
			goqu.SUM("time_entries.duration_minutes").As("logged_minutes"),
		).
		Join(goqu.T("tasks"), goqu.On(goqu.Ex{"time_entries.task_id": goqu.I("tasks.id")})).
		Join(goqu.T("users"), goqu.On(goqu.Ex{"time_entries.user_id": goqu.I("users.id")})).
		Where(
			goqu.Ex{"tasks.project_id": projectID},
			goqu.I("time_entries.end_time").IsNotNull(),
		).
		GroupBy("time_entries.user_id", "users.full_name").
		ScanStructs(&analytics.ByMember)

	if err != nil {
		return nil, err
	}

	// 3. Daily Time Entry
	cutoffDate := time.Now().AddDate(0, 0, -14).Format("2006-01-02")
	err = model.db.From(goqu.T("time_entries")).
		Select(
			goqu.L("TO_CHAR(time_entries.start_time, 'YYYY-MM-DD')").As("date"),
			goqu.SUM("time_entries.duration_minutes").As("logged_minutes"),
		).
		Join(goqu.T("tasks"), goqu.On(goqu.Ex{"time_entries.task_id": goqu.I("tasks.id")})).
		Where(
			goqu.Ex{"tasks.project_id": projectID},
			goqu.I("time_entries.end_time").IsNotNull(),
			goqu.L("time_entries.start_time >= ?", cutoffDate),
		).
		GroupBy(goqu.L("TO_CHAR(time_entries.start_time, 'YYYY-MM-DD')")).
		Order(goqu.L("TO_CHAR(time_entries.start_time, 'YYYY-MM-DD')").Asc()).
		ScanStructs(&analytics.DailyTrend)

	if err != nil {
		return nil, err
	}

	return &analytics, nil
}

// ListByProjectID lists all time entries for a project with optional filters.
func (model *TimeEntryModel) ListByProjectID(projectID uuid.UUID, userID *uuid.UUID, startDate, endDate *time.Time) ([]TimeEntryWithDetails, error) {
	var entries []TimeEntryWithDetails

	q := model.db.From(TimeEntryTable).
		Select(
			"time_entries.*",
			goqu.I("users.full_name").As("user_full_name"),
			goqu.I("users.email").As("user_email"),
			goqu.I("tasks.title").As("task_title"),
			goqu.I("tasks.number").As("task_number"),
		).
		Join(goqu.T("tasks"), goqu.On(goqu.Ex{"time_entries.task_id": goqu.I("tasks.id")})).
		Join(goqu.T("users"), goqu.On(goqu.Ex{"time_entries.user_id": goqu.I("users.id")})).
		Where(
			goqu.Ex{"tasks.project_id": projectID},
			goqu.I("time_entries.end_time").IsNotNull(),
		)

	if userID != nil {
		q = q.Where(goqu.Ex{"time_entries.user_id": *userID})
	}
	if startDate != nil {
		q = q.Where(goqu.I("time_entries.start_time").Gte(*startDate))
	}
	if endDate != nil {
		q = q.Where(goqu.I("time_entries.start_time").Lte(*endDate))
	}

	err := q.Order(goqu.I("time_entries.start_time").Desc()).ScanStructs(&entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
