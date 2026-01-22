import enum
from sqlalchemy import Column, Integer, String, ForeignKey, Enum as SqEnum
from sqlalchemy.orm import relationship
from core.database import Base

class WorkspaceType(str, enum.Enum):
    PERSONAL = "PERSONAL"
    TEAM = "TEAM"

class WorkspaceRole(str, enum.Enum):
    ADMIN = "ADMIN"
    MEMBER = "MEMBER"
    VIEWER = "VIEWER"

class WorkspaceMember(Base):
    __tablename__ = "workspace_members"

    workspace_id = Column(Integer, ForeignKey("workspaces.id"), primary_key=True)
    user_id = Column(Integer, ForeignKey("users.id"), primary_key=True)
    role = Column(SqEnum(WorkspaceRole), default=WorkspaceRole.MEMBER, nullable=False)

    # Relationships
    workspace = relationship("Workspace", back_populates="members")
    user = relationship("User", back_populates="workspaces")

class Workspace(Base):
    __tablename__ = "workspaces"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, nullable=False)
    type = Column(SqEnum(WorkspaceType), default=WorkspaceType.PERSONAL, nullable=False)
    owner_id = Column(Integer, ForeignKey("users.id"), nullable=False)

    # Relationships
    members = relationship("WorkspaceMember", back_populates="workspace", cascade="all, delete-orphan")
    owner = relationship("User", foreign_keys=[owner_id])
