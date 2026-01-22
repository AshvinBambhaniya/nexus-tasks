import enum
from sqlalchemy import Column, Integer, String, ForeignKey, Enum as SqEnum, DateTime, Text
from sqlalchemy.orm import relationship
from core.database import Base
from datetime import datetime

class TaskStatus(str, enum.Enum):
    TODO = "TODO"
    IN_PROGRESS = "IN_PROGRESS"
    DONE = "DONE"
    BACKLOG = "BACKLOG"

class TaskPriority(str, enum.Enum):
    P0 = "P0" # Critical
    P1 = "P1" # High
    P2 = "P2" # Medium
    P3 = "P3" # Low

class Task(Base):
    __tablename__ = "tasks"

    id = Column(Integer, primary_key=True, index=True)
    title = Column(String, nullable=False)
    description = Column(Text, nullable=True)
    status = Column(SqEnum(TaskStatus), default=TaskStatus.TODO, nullable=False)
    priority = Column(SqEnum(TaskPriority), default=TaskPriority.P2, nullable=False)
    
    # Foreign Keys
    workspace_id = Column(Integer, ForeignKey("workspaces.id", ondelete="CASCADE"), nullable=False)
    assignee_id = Column(Integer, ForeignKey("users.id", ondelete="SET NULL"), nullable=True)
    
    # Dates
    due_date = Column(DateTime, nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)

    # Relationships
    workspace = relationship("Workspace", backref="tasks")
    assignee = relationship("User", foreign_keys=[assignee_id])
