from typing import List, Optional

from fastapi import APIRouter, BackgroundTasks, Depends, HTTPException, Query
from sqlalchemy.orm import Session

from api.v1.auth import get_current_user
from core.database import get_db
from core.websocket import manager
from models.comment import Comment
from models.project import Project, ProjectMember, ProjectTeam
from models.task import Task, TaskStatus
from models.team import TeamMember
from models.user import User
from models.workspace import WorkspaceMember, WorkspaceRole
from schemas.comment import CommentCreate, CommentResponse
from schemas.task import TaskCreate, TaskResponse, TaskUpdate, TaskWithProject

router = APIRouter()


def validate_assignee(project_id: int, assignee_id: Optional[int], db: Session):
    if not assignee_id:
        return

    # Check Direct Membership
    member = (
        db.query(ProjectMember)
        .filter(
            ProjectMember.project_id == project_id, ProjectMember.user_id == assignee_id
        )
        .first()
    )
    if member:
        return

    # Check Team Membership
    team_member = (
        db.query(ProjectTeam)
        .join(TeamMember, ProjectTeam.team_id == TeamMember.team_id)
        .filter(ProjectTeam.project_id == project_id, TeamMember.user_id == assignee_id)
        .first()
    )
    if team_member:
        return

    raise HTTPException(
        status_code=400,
        detail="Assignee must be a member of the project (directly or via team)",
    )


def validate_project_access(project_id: int, db: Session, user_id: int) -> bool:
    # 1. Check direct project membership
    member = (
        db.query(ProjectMember)
        .filter(
            ProjectMember.project_id == project_id, ProjectMember.user_id == user_id
        )
        .first()
    )
    if member:
        return True

    # 2. Check Workspace Admin (Master Access)
    project = db.query(Project).filter(Project.id == project_id).first()
    if not project:
        raise HTTPException(status_code=404, detail="Project not found")

    ws_member = (
        db.query(WorkspaceMember)
        .filter(
            WorkspaceMember.workspace_id == project.workspace_id,
            WorkspaceMember.user_id == user_id,
        )
        .first()
    )

    if ws_member and ws_member.role == WorkspaceRole.ADMIN:
        return True

    raise HTTPException(status_code=403, detail="Not a member of this project")


@router.get("/tasks/{task_id}", response_model=TaskResponse)
def get_task(
    task_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    task = db.query(Task).filter(Task.id == task_id).first()
    if not task:
        raise HTTPException(status_code=404, detail="Task not found")

    validate_project_access(task.project_id, db, current_user.id)
    return task


@router.get("/tasks/{task_id}/comments", response_model=List[CommentResponse])
def list_task_comments(
    task_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    task = db.query(Task).filter(Task.id == task_id).first()
    if not task:
        raise HTTPException(status_code=404, detail="Task not found")

    validate_project_access(task.project_id, db, current_user.id)

    comments = (
        db.query(Comment)
        .filter(Comment.task_id == task_id)
        .order_by(Comment.created_at.asc())
        .all()
    )
    return comments


@router.post("/tasks/{task_id}/comments", response_model=CommentResponse)
def create_comment(
    task_id: int,
    comment: CommentCreate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    task = db.query(Task).filter(Task.id == task_id).first()
    if not task:
        raise HTTPException(status_code=404, detail="Task not found")

    validate_project_access(task.project_id, db, current_user.id)

    new_comment = Comment(
        content=comment.content, task_id=task_id, author_id=current_user.id
    )
    db.add(new_comment)
    db.commit()
    db.refresh(new_comment)
    return new_comment


@router.get("/tasks/me", response_model=List[TaskWithProject])
def get_my_tasks(
    current_user: User = Depends(get_current_user), db: Session = Depends(get_db)
):
    """
    Get all tasks assigned to the current user across all projects.
    """
    tasks = (
        db.query(Task)
        .join(Task.project)
        .filter(Task.assignee_id == current_user.id)
        .order_by(Task.priority.asc(), Task.due_date.asc())
        .all()
    )

    return tasks


@router.post("/projects/{project_id}/tasks", response_model=TaskResponse)
def create_task(
    project_id: int,
    task: TaskCreate,
    background_tasks: BackgroundTasks,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """
    Create a new task in a specific project.
    """
    validate_project_access(project_id, db, current_user.id)

    # Force project_id from URL to match body or override
    task_data = task.model_dump()
    task_data["project_id"] = project_id

    validate_assignee(project_id, task_data.get("assignee_id"), db)

    new_task = Task(**task_data)
    new_task.author_id = current_user.id  # Set author
    db.add(new_task)
    db.commit()
    db.refresh(new_task)

    # Broadcast Event
    # We broadcast to the PROJECT channel now, not workspace
    # Need to update WebSocket manager to handle Project channels or generic channels
    # For now, let's use project_id as the channel ID
    background_tasks.add_task(
        manager.broadcast,
        project_id,
        {
            "type": "TASK_CREATED",
            "task": TaskResponse.model_validate(new_task).model_dump(mode="json"),
        },
    )

    return new_task


@router.get("/projects/{project_id}/tasks", response_model=List[TaskResponse])
def list_project_tasks(
    project_id: int,
    status: Optional[TaskStatus] = Query(None),
    assignee_id: Optional[int] = Query(None),
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    validate_project_access(project_id, db, current_user.id)

    query = db.query(Task).filter(Task.project_id == project_id)

    if status:
        query = query.filter(Task.status == status)
    if assignee_id:
        query = query.filter(Task.assignee_id == assignee_id)

    return query.all()


@router.delete("/comments/{comment_id}")
def delete_comment(
    comment_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    comment = db.query(Comment).filter(Comment.id == comment_id).first()
    if not comment:
        raise HTTPException(status_code=404, detail="Comment not found")

    # Check permissions: Author or Admin?
    # For now, allow author to delete.
    if comment.author_id != current_user.id:
        # Optionally check if project admin
        # Simple check: Only author can delete for now to be safe
        raise HTTPException(
            status_code=403, detail="Not authorized to delete this comment"
        )

    db.delete(comment)
    db.commit()

    return {"status": "success", "message": "Comment deleted"}


@router.patch("/tasks/{task_id}", response_model=TaskResponse)
def update_task(
    task_id: int,
    task_update: TaskUpdate,
    background_tasks: BackgroundTasks,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    task = db.query(Task).filter(Task.id == task_id).first()
    if not task:
        raise HTTPException(status_code=404, detail="Task not found")

    validate_project_access(task.project_id, db, current_user.id)

    update_data = task_update.model_dump(exclude_unset=True)

    if "assignee_id" in update_data:
        validate_assignee(task.project_id, update_data["assignee_id"], db)

    for key, value in update_data.items():
        setattr(task, key, value)

    db.commit()
    db.refresh(task)

    background_tasks.add_task(
        manager.broadcast,
        task.project_id,
        {
            "type": "TASK_UPDATED",
            "task": TaskResponse.model_validate(task).model_dump(mode="json"),
        },
    )

    return task


@router.delete("/tasks/{task_id}")
def delete_task(
    task_id: int,
    background_tasks: BackgroundTasks,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    task = db.query(Task).filter(Task.id == task_id).first()
    if not task:
        raise HTTPException(status_code=404, detail="Task not found")

    project_id = task.project_id
    validate_project_access(project_id, db, current_user.id)

    db.delete(task)
    db.commit()

    background_tasks.add_task(
        manager.broadcast, project_id, {"type": "TASK_DELETED", "task_id": task_id}
    )

    return {"status": "success", "message": "Task deleted"}
