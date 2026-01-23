from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status, Query, BackgroundTasks
from sqlalchemy.orm import Session

from core.database import get_db
from core.websocket import manager
from models.user import User
from models.task import Task, TaskStatus
from models.workspace import WorkspaceMember
from schemas.task import TaskCreate, TaskResponse, TaskUpdate, TaskWithWorkspace
from api.v1.auth import get_current_user
from api.v1.workspaces import validate_workspace_access

router = APIRouter()

@router.get("/tasks/me", response_model=List[TaskWithWorkspace])
def get_my_tasks(
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Get all tasks assigned to the current user across all workspaces.
    Sorted by Priority (P0 first) and Due Date (earliest first).
    """
    tasks = db.query(Task).join(Task.workspace).filter(
        Task.assignee_id == current_user.id
    ).order_by(
        Task.priority.asc(), # P0 comes before P1 (lexicographically P0 < P1) if enum is string, let's verify Enum order or definition
        Task.due_date.asc()
    ).all()
    
    return tasks

@router.post("/workspaces/{workspace_id}/tasks", response_model=TaskResponse)
def create_task(
    workspace_id: int,
    task: TaskCreate,
    background_tasks: BackgroundTasks,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Create a new task in a specific workspace.
    """
    # 1. Validate Access
    validate_workspace_access(workspace_id, db, current_user.id)

    # 2. Create Task
    new_task = Task(
        **task.model_dump(),
        workspace_id=workspace_id
    )
    db.add(new_task)
    db.commit()
    db.refresh(new_task)

    # 3. Broadcast Event
    background_tasks.add_task(
        manager.broadcast, 
        workspace_id, 
        {"type": "TASK_CREATED", "task": TaskResponse.model_validate(new_task).model_dump(mode='json')}
    )

    return new_task

@router.get("/workspaces/{workspace_id}/tasks", response_model=List[TaskResponse])
def list_workspace_tasks(
    workspace_id: int,
    status: Optional[TaskStatus] = Query(None),
    assignee_id: Optional[int] = Query(None),
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    List tasks in a workspace with optional filters.
    """
    # 1. Validate Access
    validate_workspace_access(workspace_id, db, current_user.id)

    # 2. Query Tasks
    query = db.query(Task).filter(Task.workspace_id == workspace_id)
    
    if status:
        query = query.filter(Task.status == status)
    if assignee_id:
        query = query.filter(Task.assignee_id == assignee_id)
        
    return query.all()

@router.patch("/tasks/{task_id}", response_model=TaskResponse)
def update_task(
    task_id: int,
    task_update: TaskUpdate,
    background_tasks: BackgroundTasks,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Update a task. Verifies user is member of the task's workspace.
    """
    # 1. Get Task
    task = db.query(Task).filter(Task.id == task_id).first()
    if not task:
        raise HTTPException(status_code=404, detail="Task not found")

    # 2. Validate Access (using task.workspace_id)
    validate_workspace_access(task.workspace_id, db, current_user.id)

    # 3. Update Fields
    update_data = task_update.model_dump(exclude_unset=True)
    for key, value in update_data.items():
        setattr(task, key, value)

    db.commit()
    db.refresh(task)

    # 4. Broadcast Event
    background_tasks.add_task(
        manager.broadcast, 
        task.workspace_id, 
        {"type": "TASK_UPDATED", "task": TaskResponse.model_validate(task).model_dump(mode='json')}
    )

    return task

@router.delete("/tasks/{task_id}")
def delete_task(
    task_id: int,
    background_tasks: BackgroundTasks,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Delete a task.
    """
    # 1. Get Task
    task = db.query(Task).filter(Task.id == task_id).first()
    if not task:
        raise HTTPException(status_code=404, detail="Task not found")

    workspace_id = task.workspace_id # Store before delete

    # 2. Validate Access
    validate_workspace_access(workspace_id, db, current_user.id)

    # 3. Delete
    db.delete(task)
    db.commit()

    # 4. Broadcast Event
    background_tasks.add_task(
        manager.broadcast, 
        workspace_id, 
        {"type": "TASK_DELETED", "task_id": task_id}
    )

    return {"status": "success", "message": "Task deleted"}
