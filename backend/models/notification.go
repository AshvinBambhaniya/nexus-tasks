package models

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

// NotificationEntityType represents the type of entity a notification refers to.
type NotificationEntityType string

const (
	// EntityTypeTask represents a notification tied to a task.
	EntityTypeTask NotificationEntityType = "TASK"
	// EntityTypeComment represents a notification tied to a comment.
	EntityTypeComment NotificationEntityType = "COMMENT"
	// EntityTypeProject represents a notification tied to a project.
	EntityTypeProject NotificationEntityType = "PROJECT"
)

// NotificationType represents the type of action that triggered the notification.
type NotificationType string

const (
	// NotificationTypeAssigned is triggered when a user is assigned to a task.
	NotificationTypeAssigned NotificationType = "ASSIGNED"
	// NotificationTypeMentioned is triggered when a user is mentioned in a comment or task.
	NotificationTypeMentioned NotificationType = "MENTIONED"
	// NotificationTypeStatusChanged is triggered when a task's status is changed.
	NotificationTypeStatusChanged NotificationType = "STATUS_CHANGED"
	// NotificationTypeCommentAdded is triggered when a new comment is added.
	NotificationTypeCommentAdded NotificationType = "COMMENT_ADDED"
)

// Notification represents a single inbox notification event.
type Notification struct {
	ID         uuid.UUID              `json:"id" db:"id" goqu:"skipinsert"`
	UserID     uuid.UUID              `json:"user_id" db:"user_id"`
	ActorID    uuid.UUID              `json:"actor_id" db:"actor_id"`
	EntityID   uuid.UUID              `json:"entity_id" db:"entity_id"`
	EntityType NotificationEntityType `json:"entity_type" db:"entity_type"`
	Type       NotificationType       `json:"type" db:"type"`
	Title      string                 `json:"title" db:"title"`
	Body       *string                `json:"body" db:"body"`
	IsRead     bool                   `json:"is_read" db:"is_read"`
	IsCleared  bool                   `json:"is_cleared" db:"is_cleared"`
	CreatedAt  time.Time              `json:"created_at" db:"created_at" goqu:"skipinsert"`
	UpdatedAt  time.Time              `json:"updated_at" db:"updated_at" goqu:"skipinsert"`
}

// NotificationRepository provides operations to manage notifications.
type NotificationRepository interface {
	Create(n *Notification) error
	GetInbox(userID uuid.UUID) ([]Notification, error)
	MarkAsRead(id uuid.UUID, userID uuid.UUID) error
	MarkAsCleared(id uuid.UUID, userID uuid.UUID) error
	ClearAll(userID uuid.UUID) error
}

type notificationModel struct {
	db DbExecutor
}

// InitNotificationModel initializes and returns a new NotificationRepository.
func InitNotificationModel(db DbExecutor) NotificationRepository {
	return &notificationModel{db: db}
}

func (m *notificationModel) Create(n *Notification) error {
	_, err := m.db.Insert("notifications").Rows(n).Executor().Exec()
	return err
}

func (m *notificationModel) GetInbox(userID uuid.UUID) ([]Notification, error) {
	var notifications []Notification
	err := m.db.From("notifications").
		Where(goqu.Ex{"user_id": userID, "is_cleared": false}).
		Order(goqu.I("created_at").Desc()).
		ScanStructs(&notifications)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (m *notificationModel) MarkAsRead(id uuid.UUID, userID uuid.UUID) error {
	_, err := m.db.Update("notifications").
		Set(goqu.Record{"is_read": true, "updated_at": goqu.L("CURRENT_TIMESTAMP")}).
		Where(goqu.Ex{"id": id, "user_id": userID}).
		Executor().Exec()
	return err
}

func (m *notificationModel) MarkAsCleared(id uuid.UUID, userID uuid.UUID) error {
	_, err := m.db.Update("notifications").
		Set(goqu.Record{"is_cleared": true, "updated_at": goqu.L("CURRENT_TIMESTAMP")}).
		Where(goqu.Ex{"id": id, "user_id": userID}).
		Executor().Exec()
	return err
}

func (m *notificationModel) ClearAll(userID uuid.UUID) error {
	_, err := m.db.Update("notifications").
		Set(goqu.Record{"is_cleared": true, "updated_at": goqu.L("CURRENT_TIMESTAMP")}).
		Where(goqu.Ex{"user_id": userID, "is_read": true, "is_cleared": false}).
		Executor().Exec()
	return err
}
