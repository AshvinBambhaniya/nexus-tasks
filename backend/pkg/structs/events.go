package structs

// RealtimeEvent represents a standard structure for WebSocket events
type RealtimeEvent struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Event Types
const (
	EventTypeTaskCreated = "TASK_CREATED"
	EventTypeTaskUpdated = "TASK_UPDATED"
	EventTypeTaskDeleted = "TASK_DELETED"
)
