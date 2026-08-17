# ⏱️ Nexus Tasks — Time Tracking Feature: Complete Specification

---

## Part 1: Complete User Flow

This section walks through **every screen and interaction** a user encounters, from setting an estimate to viewing analytics.

---

### Flow 1: Setting a Time Estimate on a Task

**Where**: Task creation page or task detail sidebar

```
┌─────────────────────────────────────────────────────────────┐
│  Create New Task                                            │
│─────────────────────────────────────────────────────────────│
│  Title:       [ Fix CORS headers in auth middleware      ]  │
│  Description: [ ████████████████████████████████████████ ]  │
│  Status:      [ TODO          ▾ ]                           │
│  Priority:    [ P1            ▾ ]                           │
│  Assignee:    [ @ashvin       ▾ ]                           │
│  Due Date:    [ 2026-08-20    📅 ]                          │
│                                                             │
│  ⏱️ Estimated Time                                          │
│  ┌──────────┐   ┌──────────┐                                │
│  │  4  hours│   │ 30  mins │                                │
│  └──────────┘   └──────────┘                                │
│                                                             │
│              [ Cancel ]  [ Create Task ✓ ]                  │
└─────────────────────────────────────────────────────────────┘
```

**Steps:**
1. User navigates to `/projects/:projectId/tasks/new` or opens the task creation dialog.
2. Fills in title, description, priority, assignee, due date — all existing fields.
3. **New field**: "Estimated Time" with hours and minutes inputs (stored as `estimated_minutes` integer on the task).
4. Clicks **Create Task**. The API call `POST /tasks` now includes `estimated_minutes: 270`.
5. Estimate can also be added/edited later via the task detail sidebar (`PATCH /tasks/:taskId`).

---

### Flow 2: Starting a Live Timer

**Where**: Task detail page, Kanban card, or task table row

```
Task Detail Page — PROJ-14: Fix CORS headers
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  Status: IN_PROGRESS    Priority: P1    Assignee: @ashvin

  ⏱️ Time Tracking
  ┌───────────────────────────────────────────┐
  │  Estimated: 4h 30m                        │
  │  Logged:    0h 0m                         │
  │  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  0%      │
  │                                           │
  │         [ ▶ Start Timer ]                 │
  └───────────────────────────────────────────┘
```

**Steps:**
1. User clicks **▶ Start Timer** on any task.
2. Frontend calls `POST /api/v2/tasks/:taskId/timer/start`.
3. Backend checks: Does this user already have a running timer on ANY task?
   - **Yes** → Auto-stop the previous timer (duration calculated, saved with empty description), then start the new one.
   - **No** → Create a new `time_entries` row with `start_time = NOW()`, `end_time = NULL`.
4. Backend broadcasts `TIMER_STARTED` event over WebSocket (topic: `workspace:<id>`).
5. Frontend receives the event → activates the **Global Floating Timer Bar**.

---

### Flow 3: The Global Floating Timer Bar

**Where**: Sticky banner at the top of the dashboard layout (visible on every page while a timer is running)

```
┌──────────────────────────────────────────────────────────────────────────────┐
│  ⏱️  PROJ-14: Fix CORS headers    │  01:24:37  │  [ ⏸ Pause ]  [ ⏹ Stop ]  │  ✕ Discard │
└──────────────────────────────────────────────────────────────────────────────┘
```

**Behavior:**
- The bar appears immediately after starting a timer and persists across all page navigations.
- **Live Counter**: A `setInterval` ticks every second, computing elapsed time from `start_time` stored in the Pinia timer store (persisted to localStorage).
- **Task Title is clickable**: Navigates to the task detail page.
- **Stop Button**: Opens the "Log Work" dialog (Flow 4).
- **Discard Button**: Calls `POST /api/v2/tasks/:taskId/timer/discard` — deletes the `time_entries` row without logging any time. Shows a confirmation tooltip first.
- On page refresh or reconnect: The frontend calls `GET /api/v2/timer/active` on mount to restore the running timer state.

---

### Flow 4: Stopping the Timer & Logging Work

**Where**: Modal dialog triggered by clicking "Stop" on the floating bar or task detail

```
┌─────────────────────────────────────────────────┐
│  Log Work — PROJ-14: Fix CORS headers           │
│─────────────────────────────────────────────────│
│                                                 │
│  Duration:  01h 24m 37s                         │
│             (auto-calculated, editable)         │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐    │
│  │  1  hours│   │ 25  mins │   │ 00  secs │    │
│  └──────────┘   └──────────┘   └──────────┘    │
│                                                 │
│  Work Description:                              │
│  ┌─────────────────────────────────────────┐    │
│  │ Debugged CORS preflight handling.       │    │
│  │ Added wildcard origin support for dev.  │    │
│  └─────────────────────────────────────────┘    │
│                                                 │
│           [ Cancel ]   [ Log Time ✓ ]           │
└─────────────────────────────────────────────────┘
```

**Steps:**
1. User clicks **⏹ Stop** on the floating bar or the task detail page.
2. The "Log Work" dialog opens, pre-filled with the elapsed duration (editable in case the user wants to round).
3. User types a work description (what they accomplished during this session).
4. Clicks **Log Time ✓**.
5. Frontend calls `POST /api/v2/tasks/:taskId/timer/stop` with `{ description, duration_minutes }`.
6. Backend sets `end_time = NOW()`, stores `duration_minutes` and `description`.
7. Backend broadcasts `TIMER_STOPPED` event via WebSocket.
8. The floating timer bar disappears. The task's time tracking widget updates.

---

### Flow 5: Manual Time Logging (Without Timer)

**Where**: Task detail page → "Log Time" button

```
Task Detail — Time Tracking Section
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ⏱️ Time Tracking
  ┌───────────────────────────────────────────┐
  │  Estimated: 4h 30m                        │
  │  Logged:    1h 25m                        │
  │  ████████░░░░░░░░░░░░░░░░░░░░░░  31%     │
  │                                           │
  │   [ ▶ Start Timer ]  [ + Log Time ]       │
  └───────────────────────────────────────────┘
```

**Steps:**
1. User clicks **+ Log Time** (useful for retroactive logging, e.g., "I spent 2 hours on this yesterday").
2. A dialog opens:
   ```
   ┌────────────────────────────────────────────┐
   │  Log Time Manually                         │
   │──────────────────────────────────────────── │
   │  Date:      [ 2026-08-16  📅 ]             │
   │  Duration:  [ 2 ] hours  [ 0 ] mins        │
   │  Description:                               │
   │  [ Researched CORS spec, wrote unit tests ] │
   │                                             │
   │          [ Cancel ]  [ Log Time ✓ ]         │
   └────────────────────────────────────────────┘
   ```
3. Frontend calls `POST /api/v2/tasks/:taskId/time-entries` with `{ duration_minutes: 120, description, date, is_manual: true }`.
4. The work log timeline updates to show the new entry with a "Manual" badge.

---

### Flow 6: Viewing Work Log History on a Task

**Where**: Task detail page → "Time Log" section below the time tracking widget

```
  📋 Work Log
  ┌───────────────────────────────────────────────────────────────┐
  │  @ashvin  •  Today, 2:30 PM  •  1h 25m  •  ⏱️ Timer         │
  │  "Debugged CORS preflight handling. Added wildcard origin."   │
  │                                                        [ 🗑 ] │
  ├───────────────────────────────────────────────────────────────┤
  │  @ashvin  •  Yesterday, 4:00 PM  •  2h 0m  •  ✏️ Manual      │
  │  "Researched CORS spec, wrote unit tests"                     │
  │                                                        [ 🗑 ] │
  ├───────────────────────────────────────────────────────────────┤
  │  @priya   •  Aug 14, 11:00 AM  •  0h 45m  •  ⏱️ Timer        │
  │  "Code review and suggested fixes"                            │
  │                                                               │
  └───────────────────────────────────────────────────────────────┘

  Total Logged: 4h 10m / 4h 30m estimated  ━━━━━━━━━━━━━━━━░░ 93%
```

**Rules:**
- Entries are ordered newest first.
- Each entry shows: author avatar, timestamp, duration, source badge (Timer / Manual), and description.
- Delete button (🗑) visible only to the entry owner or a workspace/project admin.
- The summary bar at the bottom shows cumulative time vs estimate with a color-coded progress bar:
  - **Blue/Green**: 0–80% of estimate
  - **Amber**: 80–100% of estimate
  - **Red**: Over 100% of estimate (with `+Xh Ym over` label)

---

### Flow 7: Project Time Analytics Dashboard

**Where**: Project detail page → New "Analytics" tab (alongside Tasks, Board, Members, Settings)

```
Project: nexus-tasks-backend  →  [ Tasks ] [ Board ] [ Analytics ] [ Members ] [ Settings ]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐
  │ Total        │  │ Total        │  │ Estimate     │  │ Over-Budget      │
  │ Estimated    │  │ Logged       │  │ Accuracy     │  │ Tasks            │
  │              │  │              │  │              │  │                  │
  │   120h 30m   │  │   98h 45m    │  │    82%       │  │      3           │
  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────────┘

  ┌─── Estimated vs Logged (by Task) ────────────────────────────────────────┐
  │                                                                          │
  │  PROJ-14  ████████████████████  4.5h est                                 │
  │           ████████████████████████  5.2h logged  (+0.7h over)            │
  │                                                                          │
  │  PROJ-15  ██████████████████████████  8.0h est                           │
  │           ████████████████████  6.5h logged                              │
  │                                                                          │
  │  PROJ-16  ████████  2.0h est                                             │
  │           ████████  1.8h logged                                          │
  └──────────────────────────────────────────────────────────────────────────┘

  ┌─── Hours by Team Member ──────┐  ┌─── Daily Logged Hours (Last 14d) ────┐
  │                               │  │                                       │
  │       ┌────┐                  │  │  8h ┤         ╭─╮                     │
  │       │████│ @ashvin  42.5h   │  │  6h ┤    ╭─╮  │ │  ╭─╮               │
  │       │████│ @priya   31.0h   │  │  4h ┤ ╭─╮│ │╭─╮ │  │ │╭─╮           │
  │       │████│ @dev     25.3h   │  │  2h ┤ │ ││ ││ │ │  │ ││ │           │
  │       └────┘                  │  │  0h ┼─┴─┴┴─┴┴─┴─┴──┴─┴┴─┴──         │
  │                               │  │      Mon Tue Wed Thu Fri Sat Sun      │
  └───────────────────────────────┘  └───────────────────────────────────────┘

                            [ ⬇ Export CSV ]  [ ⬇ Export JSON ]
```

**Data Source**: `GET /api/v2/workspaces/:wid/projects/:pid/time-analytics`

**Response Shape:**
```json
{
  "summary": {
    "total_estimated_minutes": 7230,
    "total_logged_minutes": 5925,
    "estimate_accuracy_percent": 82,
    "over_budget_task_count": 3,
    "total_tasks_with_estimates": 15,
    "total_tasks_with_logs": 12
  },
  "by_task": [
    {
      "task_id": "...", "task_number": 14, "task_title": "Fix CORS headers",
      "estimated_minutes": 270, "logged_minutes": 312, "is_over_budget": true
    }
  ],
  "by_member": [
    { "user_id": "...", "full_name": "Ashvin", "logged_minutes": 2550 }
  ],
  "daily_trend": [
    { "date": "2026-08-10", "logged_minutes": 360 },
    { "date": "2026-08-11", "logged_minutes": 480 }
  ]
}
```

---

### Flow 8: AI Integration — Enhanced Weekly Report

**Where**: Project detail page → "Generate AI Weekly Report" button (already exists)

The AI prompt for `GenerateWeeklyReport` is enhanced to include time data:

**Before (existing):**
> *"12 tasks completed this sprint..."*

**After (with time tracking):**
> *"12 tasks completed this sprint. The team logged **98h 45m** against a total estimate of **120h 30m** (82% accuracy). 3 tasks exceeded their estimates — notably PROJ-14 (Fix CORS headers) went 1.2h over due to unexpected preflight handling complexity. @priya was the most efficient contributor at 94% estimate accuracy."*

---

### Flow 9: MCP Agent Interaction

An AI agent (e.g. Claude Desktop) connected via the MCP server can:

```
Agent: "I've been working on PROJ-14 for the last 2 hours debugging CORS issues. Log that time."

→ MCP Tool Call: log_task_time(task_id: "...", duration_minutes: 120, description: "Debugging CORS preflight issues")

Agent: "Start tracking time on PROJ-15 now."

→ MCP Tool Call: start_task_timer(task_id: "...")

Agent: "How is the project doing on time?"

→ MCP Tool Call: get_project_time_analytics(project_id: "...")
→ Response: "Team has logged 98h of 120h estimated (82% accuracy). 3 tasks over budget."
```

---

## Part 2: Complete Implementation Plan

This section maps every code change to the existing codebase, following established project patterns.

---

### Phase 1: Backend — Database & Models

#### Step 1.1: Create Migration

**New file**: `backend/database/migrations/000013_add_time_tracking.up.sql`

```sql
ALTER TABLE tasks ADD COLUMN estimated_minutes INT DEFAULT NULL;

CREATE TABLE time_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    description TEXT DEFAULT '',
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    duration_minutes INT DEFAULT NULL,
    is_manual BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_time_entries_task_id ON time_entries(task_id);
CREATE INDEX idx_time_entries_user_id ON time_entries(user_id);
CREATE INDEX idx_time_entries_workspace_id ON time_entries(workspace_id);
CREATE UNIQUE INDEX idx_active_user_timer ON time_entries(user_id) WHERE end_time IS NULL;
```

**New file**: `backend/database/migrations/000013_add_time_tracking.down.sql`

```sql
DROP TABLE IF EXISTS time_entries;
ALTER TABLE tasks DROP COLUMN IF EXISTS estimated_minutes;
```

#### Step 1.2: Define Go Models

**New file**: `backend/models/time_entry.go`

Following the same pattern as existing models like [`models/task.go`](file:///mnt/c/BACKUP/workspace/nexus-tasks/backend/models/task.go):

```go
package models

import "time"

type TimeEntry struct {
    ID              string     `db:"id" goqu:"skipinsert,skipupdate"`
    TaskID          string     `db:"task_id"`
    UserID          string     `db:"user_id"`
    WorkspaceID     string     `db:"workspace_id"`
    Description     string     `db:"description"`
    StartTime       time.Time  `db:"start_time"`
    EndTime         *time.Time `db:"end_time"`
    DurationMinutes *int       `db:"duration_minutes"`
    IsManual        bool       `db:"is_manual"`
    CreatedAt       time.Time  `db:"created_at" goqu:"skipinsert,skipupdate"`
    UpdatedAt       time.Time  `db:"updated_at" goqu:"skipinsert,skipupdate"`
    // Joined fields (not persisted)
    UserFullName    string     `db:"user_full_name" goqu:"skipinsert,skipupdate"`
}
```

#### Step 1.3: Define Repository Interface

**Edit**: [`backend/models/storage.go`](file:///mnt/c/BACKUP/workspace/nexus-tasks/backend/models/storage.go)

Add `TimeEntryRepository` to the `Storage` interface, following the existing pattern:

```go
type TimeEntryRepository interface {
    Create(ctx context.Context, entry *TimeEntry) error
    GetByID(ctx context.Context, id string) (*TimeEntry, error)
    GetActiveByUserID(ctx context.Context, userID string) (*TimeEntry, error)
    ListByTaskID(ctx context.Context, taskID string) ([]TimeEntry, error)
    Update(ctx context.Context, entry *TimeEntry) error
    Delete(ctx context.Context, id string) error
    // Analytics queries
    GetTotalLoggedByTaskID(ctx context.Context, taskID string) (int, error)
    GetProjectAnalytics(ctx context.Context, projectID string, days int) (*TimeAnalytics, error)
}
```

#### Step 1.4: Update Task Model

**Edit**: [`backend/models/task.go`](file:///mnt/c/BACKUP/workspace/nexus-tasks/backend/models/task.go)

Add `EstimatedMinutes *int` field to the `Task` struct and to `CreateTaskRequest` / `UpdateTaskRequest`.

#### Step 1.5: Implement Repository (Goqu)

**New file**: `backend/database/time_entry_repository.go`

Following the same Goqu query builder + scan pattern used in [`database/task_repository.go`](file:///mnt/c/BACKUP/workspace/nexus-tasks/backend/database/task_repository.go):
- `GetActiveByUserID`: `SELECT ... FROM time_entries WHERE user_id = ? AND end_time IS NULL`
- `GetProjectAnalytics`: Aggregation queries joining `time_entries` ↔ `tasks` ↔ `users` with GROUP BY for by-task, by-member, and daily-trend breakdowns.

---

### Phase 2: Backend — Service Layer

#### Step 2.1: Create TimeTrackingService

**New file**: `backend/services/time_tracking.go`

Following the same constructor pattern as [`services/task.go`](file:///mnt/c/BACKUP/workspace/nexus-tasks/backend/services/task.go):

```go
type TimeTrackingService struct {
    storage models.Storage
    hub     *realtime.Hub
    logger  *slog.Logger
}

func NewTimeTrackingService(storage models.Storage, hub *realtime.Hub, logger *slog.Logger) *TimeTrackingService
```

**Methods:**

| Method | Logic |
|---|---|
| `StartTimer(ctx, userID, taskID, workspaceID)` | 1. Check for active timer via `GetActiveByUserID`. 2. If exists, auto-stop it (calculate duration, set `end_time`). 3. Create new `time_entries` row with `start_time=now`, `end_time=NULL`. 4. Broadcast `TIMER_STARTED` via WebSocket hub. |
| `StopTimer(ctx, userID, taskID, description, durationOverride)` | 1. Get active timer, validate it belongs to `taskID`. 2. Set `end_time=now`, compute `duration_minutes` (or use override). 3. Set description. 4. Broadcast `TIMER_STOPPED`. |
| `DiscardTimer(ctx, userID)` | 1. Get active timer. 2. Delete the row. 3. Broadcast `TIMER_DISCARDED`. |
| `LogManualTime(ctx, userID, taskID, workspaceID, durationMinutes, description, date)` | 1. Validate project access. 2. Create `time_entries` with `is_manual=true`, `start_time=date`, `end_time=date+duration`, `duration_minutes`. |
| `GetActiveTimer(ctx, userID)` | Return the active timer entry (where `end_time IS NULL`), joined with task title. |
| `ListTaskTimeEntries(ctx, taskID)` | Return all entries for a task + total logged summary. |
| `DeleteTimeEntry(ctx, userID, entryID)` | Validate ownership or admin role, then delete. |
| `GetProjectAnalytics(ctx, projectID, days)` | Call repository aggregation, return structured analytics response. |

#### Step 2.2: Add WebSocket Events

**Edit**: [`backend/pkg/structs/events.go`](file:///mnt/c/BACKUP/workspace/nexus-tasks/backend/pkg/structs/events.go)

Add new event constants:
```go
const (
    EventTimerStarted  = "TIMER_STARTED"
    EventTimerStopped  = "TIMER_STOPPED"
    EventTimerDiscarded = "TIMER_DISCARDED"
)
```

#### Step 2.3: Enhance AI Service

**Edit**: [`backend/services/ai_service.go`](file:///mnt/c/BACKUP/workspace/nexus-tasks/backend/services/ai_service.go)

Update `GenerateWeeklyReport` to:
1. Accept time analytics data (total estimated, total logged, over-budget tasks).
2. Include a "Time & Estimation Analysis" section in the AI prompt.

---

### Phase 3: Backend — Controllers & Routes

#### Step 3.1: Create TimeTrackingController

**New file**: `backend/controllers/api/v2/time_tracking_controller.go`

Following the pattern in [`controllers/api/v2/task_controller.go`](file:///mnt/c/BACKUP/workspace/nexus-tasks/backend/controllers/api/v2/task_controller.go):

```go
type TimeTrackingController struct {
    service *services.TimeTrackingService
    logger  *slog.Logger
}
```

**Handler methods:**
- `GetActiveTimer(c *fiber.Ctx) error` — `GET /timer/active`
- `StartTimer(c *fiber.Ctx) error` — `POST /tasks/:taskId/timer/start`
- `StopTimer(c *fiber.Ctx) error` — `POST /tasks/:taskId/timer/stop`
- `DiscardTimer(c *fiber.Ctx) error` — `POST /tasks/:taskId/timer/discard`
- `LogManualTime(c *fiber.Ctx) error` — `POST /tasks/:taskId/time-entries`
- `ListTaskTimeEntries(c *fiber.Ctx) error` — `GET /tasks/:taskId/time-entries`
- `DeleteTimeEntry(c *fiber.Ctx) error` — `DELETE /time-entries/:entryId`
- `GetProjectAnalytics(c *fiber.Ctx) error` — `GET /.../projects/:projectId/time-analytics`

#### Step 3.2: Register Routes

**Edit**: [`backend/routes/main.go`](file:///mnt/c/BACKUP/workspace/nexus-tasks/backend/routes/main.go)

Add new route group after the existing task routes:

```go
// Time Tracking
timer := api.Group("/timer", middlewares.Authenticated(storage, apiKeyService))
timer.Get("/active", timeTrackingController.GetActiveTimer)

// Task-scoped time tracking (nested under existing task routes)
taskGroup.Post("/:taskId/timer/start", timeTrackingController.StartTimer)
taskGroup.Post("/:taskId/timer/stop", timeTrackingController.StopTimer)
taskGroup.Post("/:taskId/timer/discard", timeTrackingController.DiscardTimer)
taskGroup.Post("/:taskId/time-entries", timeTrackingController.LogManualTime)
taskGroup.Get("/:taskId/time-entries", timeTrackingController.ListTaskTimeEntries)

// Time entry management
api.Delete("/time-entries/:entryId", middlewares.Authenticated(storage, apiKeyService), timeTrackingController.DeleteTimeEntry)

// Project analytics
projectGroup.Get("/:projectId/time-analytics", timeTrackingController.GetProjectAnalytics)
```

---

### Phase 4: Backend — MCP Server Integration

#### Step 4.1: Add MCP Tools

**New file**: `backend/mcp/tools/time_tracking.go`

Following the pattern in [`mcp/tools/tasks.go`](file:///mnt/c/BACKUP/workspace/nexus-tasks/backend/mcp/tools/tasks.go):

```go
func RegisterTimeTrackingTools(server *mcp_server.MCPServer, service *services.TimeTrackingService) {
    // start_task_timer
    server.AddTool(mcp.NewTool("start_task_timer",
        mcp.WithDescription("Start a live timer on a task"),
        mcp.WithString("task_id", mcp.Required(), mcp.Description("Task UUID")),
    ), handleStartTimer(service))

    // stop_task_timer
    server.AddTool(mcp.NewTool("stop_task_timer",
        mcp.WithDescription("Stop the currently running timer and log work"),
        mcp.WithString("task_id", mcp.Required()),
        mcp.WithString("description", mcp.Description("Work summary")),
    ), handleStopTimer(service))

    // log_task_time
    server.AddTool(mcp.NewTool("log_task_time",
        mcp.WithDescription("Manually log time spent on a task"),
        mcp.WithString("task_id", mcp.Required()),
        mcp.WithNumber("duration_minutes", mcp.Required()),
        mcp.WithString("description"),
    ), handleLogTime(service))

    // get_project_time_analytics
    server.AddTool(mcp.NewTool("get_project_time_analytics",
        mcp.WithDescription("Get time tracking analytics for a project"),
        mcp.WithString("project_id", mcp.Required()),
    ), handleGetAnalytics(service))
}
```

#### Step 4.2: Register in MCP Server Bootstrap

**Edit**: [`backend/cli/mcp.go`](file:///mnt/c/BACKUP/workspace/nexus-tasks/backend/cli/mcp.go)

Add `tools.RegisterTimeTrackingTools(mcpServer, timeTrackingService)` alongside the existing tool registrations.

---

### Phase 5: Frontend — State & Composables

#### Step 5.1: Create Timer Store

**New file**: `frontend/app/stores/timer.ts`

```typescript
// Pinia store with persist: true (survives page refresh)
export const useTimerStore = defineStore('timer', () => {
  const activeTimer = ref<ActiveTimer | null>(null)
  const elapsedSeconds = ref(0)
  let intervalId: ReturnType<typeof setInterval> | null = null

  function startTicking() {
    stopTicking()
    intervalId = setInterval(() => {
      if (activeTimer.value) {
        const start = new Date(activeTimer.value.start_time).getTime()
        elapsedSeconds.value = Math.floor((Date.now() - start) / 1000)
      }
    }, 1000)
  }

  function stopTicking() {
    if (intervalId) clearInterval(intervalId)
    intervalId = null
    elapsedSeconds.value = 0
  }

  // ... setActiveTimer, clearActiveTimer, formatted elapsed time computed
  return { activeTimer, elapsedSeconds, startTicking, stopTicking }
}, { persist: true })
```

#### Step 5.2: Create Time Tracking Composable

**New file**: `frontend/app/composables/useTimeTracking.ts`

Following the pattern in [`composables/useTasks.ts`](file:///mnt/c/BACKUP/workspace/nexus-tasks/frontend/app/composables/useTasks.ts):

```typescript
export function useTimeTracking() {
  const { useMutation } = useApi()

  const fetchActiveTimer = () => useApi<ActiveTimer>('/timer/active')

  const startTimer = (taskId: string) =>
    useMutation<ActiveTimer>(`/tasks/${taskId}/timer/start`, { method: 'POST' })

  const stopTimer = (taskId: string, body: { description: string; duration_minutes?: number }) =>
    useMutation(`/tasks/${taskId}/timer/stop`, { method: 'POST', body })

  const discardTimer = (taskId: string) =>
    useMutation(`/tasks/${taskId}/timer/discard`, { method: 'POST' })

  const logManualTime = (taskId: string, body: ManualTimeLog) =>
    useMutation(`/tasks/${taskId}/time-entries`, { method: 'POST', body })

  const fetchTaskTimeEntries = (taskId: string) =>
    useApi<TimeEntriesResponse>(`/tasks/${taskId}/time-entries`)

  const deleteTimeEntry = (entryId: string) =>
    useMutation(`/time-entries/${entryId}`, { method: 'DELETE' })

  const fetchProjectAnalytics = (projectId: string) =>
    useApi<TimeAnalytics>(`/projects/${projectId}/time-analytics`)

  return { fetchActiveTimer, startTimer, stopTimer, discardTimer,
           logManualTime, fetchTaskTimeEntries, deleteTimeEntry, fetchProjectAnalytics }
}
```

#### Step 5.3: Add TypeScript Types

**Edit**: [`frontend/app/types/`](file:///mnt/c/BACKUP/workspace/nexus-tasks/frontend/app/types)

**New file**: `frontend/app/types/time_tracking.ts`

```typescript
export interface TimeEntry {
  id: string
  task_id: string
  user_id: string
  user_full_name: string
  description: string
  start_time: string
  end_time: string | null
  duration_minutes: number | null
  is_manual: boolean
  created_at: string
}

export interface ActiveTimer {
  id: string
  task_id: string
  task_title: string
  task_number: number
  project_id: string
  project_name: string
  start_time: string
}

export interface TimeAnalytics {
  summary: {
    total_estimated_minutes: number
    total_logged_minutes: number
    estimate_accuracy_percent: number
    over_budget_task_count: number
  }
  by_task: TaskTimeSummary[]
  by_member: MemberTimeSummary[]
  daily_trend: DailyTimeEntry[]
}
```

#### Step 5.4: Extend WebSocket Handler

**Edit**: [`frontend/app/composables/useSocket.ts`](file:///mnt/c/BACKUP/workspace/nexus-tasks/frontend/app/composables/useSocket.ts)

Add cases for `TIMER_STARTED`, `TIMER_STOPPED`, `TIMER_DISCARDED` events to invalidate time tracking caches and update the timer store.

---

### Phase 6: Frontend — UI Components

#### Step 6.1: Global Floating Timer Bar

**New file**: `frontend/app/components/layout/ActiveTimerBar.vue`

- Renders conditionally when `timerStore.activeTimer` is not null.
- Shows task title (linked), live counter, Stop and Discard buttons.
- Mounted in [`layouts/dashboard.vue`](file:///mnt/c/BACKUP/workspace/nexus-tasks/frontend/app/layouts/dashboard.vue) at the top of the layout.

#### Step 6.2: Time Tracking Widget (Task Detail Sidebar)

**New file**: `frontend/app/components/tasks/TimeTrackingWidget.vue`

- Shows estimated time, logged time, progress bar with color coding.
- Start Timer / Log Time buttons.
- Placed inside the task detail page ([`pages/projects/[projectId]/tasks/[taskId].vue`](file:///mnt/c/BACKUP/workspace/nexus-tasks/frontend/app/pages/projects/%5BprojectId%5D/tasks/%5BtaskId%5D.vue)).

#### Step 6.3: Work Log Timeline

**New file**: `frontend/app/components/tasks/TimeLogList.vue`

- Lists all `time_entries` for the current task.
- Each entry: avatar, name, timestamp, duration, manual/timer badge, description, delete button.

#### Step 6.4: Log Work Dialog

**New file**: `frontend/app/components/tasks/LogTimeDialog.vue`

- Uses `BaseDialog` from [`components/ui/BaseDialog.vue`](file:///mnt/c/BACKUP/workspace/nexus-tasks/frontend/app/components/ui/BaseDialog.vue).
- Two modes: "Stop Timer" (pre-filled duration) and "Manual Log" (empty form with date picker).

#### Step 6.5: Estimate Input

**Edit**: Task creation page ([`pages/projects/[projectId]/tasks/new.vue`](file:///mnt/c/BACKUP/workspace/nexus-tasks/frontend/app/pages/projects/%5BprojectId%5D/tasks/new.vue)) and task detail sidebar.

Add hours/minutes input fields bound to `estimated_minutes`.

#### Step 6.6: Project Analytics Tab

**New file**: `frontend/app/components/project/tabs/ProjectTabsAnalytics.vue`

- KPI stat cards (reuse [`dashboard/StatsCard.vue`](file:///mnt/c/BACKUP/workspace/nexus-tasks/frontend/app/components/dashboard/StatsCard.vue) pattern).
- Charts rendered with lightweight inline SVG bars or a minimal chart library.
- Export buttons for CSV/JSON download.
- Registered as a new tab in the project detail page ([`pages/projects/[projectId]/index.vue`](file:///mnt/c/BACKUP/workspace/nexus-tasks/frontend/app/pages/projects/%5BprojectId%5D/index.vue)).

#### Step 6.7: Restore Timer on App Load

**Edit**: [`layouts/dashboard.vue`](file:///mnt/c/BACKUP/workspace/nexus-tasks/frontend/app/layouts/dashboard.vue)

On mount, call `GET /api/v2/timer/active`. If a timer is running, populate `timerStore` and start ticking.

---

### Phase 7: Testing & Edge Cases

| Edge Case | Handling |
|---|---|
| User starts timer, closes browser, reopens | `GET /timer/active` on dashboard mount restores the timer; elapsed time computed from `start_time` |
| User starts timer on Task A, then starts on Task B | Auto-stop Task A's timer (empty description, computed duration), start Task B |
| User runs timer for 8+ hours (forgot to stop) | Duration calculated accurately; user can edit duration in the stop dialog before saving |
| Task deleted while timer is running | `ON DELETE CASCADE` removes the `time_entries` row; frontend WebSocket `TASK_DELETED` event clears the timer store |
| Two tabs open simultaneously | Pinia `persist: true` + WebSocket events keep both tabs in sync |
| Concurrent team members timing same task | Allowed — each user has their own timer; the unique index is per-user (`idx_active_user_timer` on `user_id WHERE end_time IS NULL`) |
| Manual log with 0 minutes | Validation rejects `duration_minutes < 1` |
| Analytics on project with no time data | Return zeroed summary, empty arrays — UI shows "No time data yet" empty state |

---

### Implementation Order Summary

```
Phase 1 — Backend Database & Models
  ├── 1.1  Migration 000013
  ├── 1.2  TimeEntry model
  ├── 1.3  Storage interface update
  ├── 1.4  Task model (add estimated_minutes)
  └── 1.5  Goqu repository implementation

Phase 2 — Backend Services
  ├── 2.1  TimeTrackingService
  ├── 2.2  WebSocket event constants
  └── 2.3  AI service enhancement

Phase 3 — Backend Controllers & Routes
  ├── 3.1  TimeTrackingController
  └── 3.2  Route registration

Phase 4 — MCP Integration
  ├── 4.1  MCP time tracking tools
  └── 4.2  MCP server bootstrap

Phase 5 — Frontend State & API
  ├── 5.1  Timer Pinia store
  ├── 5.2  useTimeTracking composable
  ├── 5.3  TypeScript types
  └── 5.4  WebSocket event handling

Phase 6 — Frontend UI
  ├── 6.1  ActiveTimerBar (global)
  ├── 6.2  TimeTrackingWidget (task detail)
  ├── 6.3  TimeLogList (task detail)
  ├── 6.4  LogTimeDialog (modal)
  ├── 6.5  Estimate input fields
  ├── 6.6  ProjectTabsAnalytics
  └── 6.7  Timer restore on load

Phase 7 — Testing & Polish
```
