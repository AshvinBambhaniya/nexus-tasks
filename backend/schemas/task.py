from datetime import datetime
from typing import Optional

from pydantic import BaseModel

from models.task import TaskPriority, TaskStatus


class TaskBase(BaseModel):
    title: str
    description: Optional[str] = None
    status: Optional[TaskStatus] = TaskStatus.TODO
    priority: Optional[TaskPriority] = TaskPriority.P2
    due_date: Optional[datetime] = None
    assignee_id: Optional[int] = None


class TaskCreate(TaskBase):
    pass


class TaskUpdate(BaseModel):
    title: Optional[str] = None
    description: Optional[str] = None
    status: Optional[TaskStatus] = None
    priority: Optional[TaskPriority] = None
    due_date: Optional[datetime] = None
    assignee_id: Optional[int] = None


class UserInfo(BaseModel):
    id: int
    email: str

    class Config:
        from_attributes = True


class TaskResponse(TaskBase):
    id: int
    project_id: int
    created_at: datetime
    updated_at: datetime
    assignee: Optional[UserInfo] = None
    author: Optional[UserInfo] = None  # Added author
    comment_count: int = 0

    class Config:
        from_attributes = True


class ProjectInfo(BaseModel):
    id: int
    name: str

    class Config:
        from_attributes = True


class TaskWithProject(TaskResponse):
    project: ProjectInfo
