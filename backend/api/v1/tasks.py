from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status, Query
from sqlalchemy.orm import Session

from core.database import get_db
from models.user import User
from models.task import Task, TaskStatus
from models.workspace import WorkspaceMember
from schemas.task import TaskCreate, TaskResponse, TaskUpdate
from api.v1.auth import get_current_user
from api.v1.workspaces import validate_workspace_access

router = APIRouter()

@router.post("/workspaces/{workspace_id}/tasks", response_model=TaskResponse)
def create_task(
    workspace_id: int,
    task: TaskCreate,
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
    return task

@router.delete("/tasks/{task_id}")
def delete_task(
    task_id: int,
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

    # 2. Validate Access
    validate_workspace_access(task.workspace_id, db, current_user.id)

    # 3. Delete
    db.delete(task)
    db.commit()
    return {"status": "success", "message": "Task deleted"}
