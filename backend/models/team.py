import enum
from datetime import datetime

from sqlalchemy import Column, DateTime, ForeignKey, Integer, String, Text
from sqlalchemy import Enum as SqEnum
from sqlalchemy.orm import relationship

from core.database import Base


class TeamRole(str, enum.Enum):
    ADMIN = "ADMIN"
    MEMBER = "MEMBER"


class TeamMember(Base):
    __tablename__ = "team_members"

    team_id = Column(
        Integer, ForeignKey("teams.id", ondelete="CASCADE"), primary_key=True
    )
    user_id = Column(
        Integer, ForeignKey("users.id", ondelete="CASCADE"), primary_key=True
    )
    role = Column(SqEnum(TeamRole), default=TeamRole.MEMBER, nullable=False)

    # Relationships
    team = relationship("Team", back_populates="members")
    user = relationship("User")


class Team(Base):
    __tablename__ = "teams"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, nullable=False)
    description = Column(Text, nullable=True)
    workspace_id = Column(
        Integer, ForeignKey("workspaces.id", ondelete="CASCADE"), nullable=False
    )
    created_at = Column(DateTime, default=datetime.utcnow)

    # Relationships
    workspace = relationship("Workspace", back_populates="teams")
    members = relationship(
        "TeamMember", back_populates="team", cascade="all, delete-orphan"
    )
