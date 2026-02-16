# Implementation Plan: Task Completion Date

## Objective

Add a `completed_at` timestamp to tasks to track when they are moved to the "DONE" status. This enables better analytics (lead time, cycle time) and historical tracking.

## Phase 1: Backend Implementation

1.  **Update SQLAlchemy Model**

    - Modify `backend/models/task.py`:
      - Add `completed_at = Column(DateTime, nullable=True)` to the `Task` class.

2.  **Generate Database Migration**

    - Run the autogenerate command:
      ```bash
      alembic revision --autogenerate -m "add_completed_at_to_tasks"
      ```
    - Verify the generated migration script in `backend/alembic/versions/` to ensure it correctly adds the column.
    - Apply the migration: `alembic upgrade head`.

3.  **Update Pydantic Schemas**

    - Modify `backend/schemas/task.py`:
      - Add `completed_at: Optional[datetime] = None` to the `TaskResponse` schema.

4.  **Update API Logic (Status Change Handler)**
    - Modify `backend/api/v1/tasks.py` (specifically the update endpoint):
      - When updating a task, check if the `status` field is changing.
      - If `new_status == TaskStatus.DONE` and `old_status != TaskStatus.DONE`:
        - Set `completed_at = datetime.utcnow()`.
      - If `new_status != TaskStatus.DONE` and `old_status == TaskStatus.DONE` (task is re-opened):
        - Set `completed_at = None`.

## Phase 2: Frontend Implementation

1.  **Update Type Definitions**

    - Modify `frontend/src/types/index.ts`:
      - Add `completed_at?: string;` to the `Task` interface.

2.  **UI Updates**
    - **Task Detail Page**:
      - File: `frontend/src/app/(dashboard)/projects/[projectId]/tasks/[taskId]/page.tsx`
      - Logic: If `task.status === "DONE"`, display "Completed on {formatted_date}" near the status badge or in the metadata section.
    - **Task Cards (Board & List views)**:
      - Files: `frontend/src/components/tasks/task-card.tsx` and `board-card.tsx`.
      - Logic: Optionally show a completion indicator or tooltip with the date.

## Phase 3: Verification

1.  **Backend Verification**:

    - Create a new task -> status starts as TODO -> `completed_at` is None.
    - Update status to DONE -> `completed_at` is set to current time.
    - Update status back to IN_PROGRESS -> `completed_at` becomes None.

2.  **Frontend Verification**:
    - Verify the UI correctly displays the completion date for completed tasks.
    - Ensure no errors occur for tasks without a completion date.

## Questions/Notes

- **Backfilling**: Existing tasks that are already "DONE" will have `completed_at = NULL` initially. We will leave them as is for now unless a backfill script is requested (e.g., setting them to their `updated_at` time).
