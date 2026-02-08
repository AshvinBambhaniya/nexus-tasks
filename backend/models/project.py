import enum
from sqlalchemy import Column, Integer, String, ForeignKey, Enum as SqEnum, DateTime, Text, Boolean
from sqlalchemy.orm import relationship
from core.database import Base
from datetime import datetime

class ProjectRole(str, enum.Enum):
    ADMIN = "ADMIN"
    MEMBER = "MEMBER"
    VIEWER = "VIEWER"

class ProjectMember(Base):
    __tablename__ = "project_members"

    project_id = Column(Integer, ForeignKey("projects.id", ondelete="CASCADE"), primary_key=True)
    user_id = Column(Integer, ForeignKey("users.id", ondelete="CASCADE"), primary_key=True)
    role = Column(SqEnum(ProjectRole), default=ProjectRole.MEMBER, nullable=False)

    # Relationships
    project = relationship("Project", back_populates="members")
    user = relationship("User")

class ProjectTeam(Base):
    __tablename__ = "project_teams"

    project_id = Column(Integer, ForeignKey("projects.id", ondelete="CASCADE"), primary_key=True)
    team_id = Column(Integer, ForeignKey("teams.id", ondelete="CASCADE"), primary_key=True)

    # Relationships
    project = relationship("Project", back_populates="teams")
    team = relationship("Team")

class Project(Base):
    __tablename__ = "projects"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, nullable=False)
    description = Column(Text, nullable=True)
    is_archived = Column(Boolean, default=False)
    workspace_id = Column(Integer, ForeignKey("workspaces.id", ondelete="CASCADE"), nullable=False)
    created_at = Column(DateTime, default=datetime.utcnow)

    # Relationships
    workspace = relationship("Workspace", back_populates="projects")
    members = relationship("ProjectMember", back_populates="project", cascade="all, delete-orphan")
    teams = relationship("ProjectTeam", back_populates="project", cascade="all, delete-orphan")
    tasks = relationship("Task", back_populates="project", cascade="all, delete-orphan")
