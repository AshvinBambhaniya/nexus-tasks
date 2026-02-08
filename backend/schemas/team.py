from pydantic import BaseModel
from typing import Optional, List
from datetime import datetime
from models.team import TeamRole

class TeamBase(BaseModel):
    name: str
    description: Optional[str] = None

class TeamCreate(TeamBase):
    pass

class TeamResponse(TeamBase):
    id: int
    workspace_id: int
    created_at: datetime

    class Config:
        from_attributes = True

class TeamMemberAdd(BaseModel):
    email: str
    role: TeamRole = TeamRole.MEMBER

class TeamMemberResponse(BaseModel):
    user_id: int
    email: str
    role: TeamRole

    class Config:
        from_attributes = True
