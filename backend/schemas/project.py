from pydantic import BaseModel
from typing import Optional, List
from datetime import datetime
from models.project import ProjectRole

class ProjectBase(BaseModel):
    name: str
    description: Optional[str] = None

class ProjectCreate(ProjectBase):
    pass

class ProjectUpdate(BaseModel):
    name: Optional[str] = None
    description: Optional[str] = None
    is_archived: Optional[bool] = None

class ProjectResponse(ProjectBase):
    id: int
    workspace_id: int
    created_at: datetime
    is_archived: bool

    class Config:
        from_attributes = True

class ProjectMemberAdd(BaseModel):
    email: str
    role: ProjectRole = ProjectRole.MEMBER

class ProjectMemberResponse(BaseModel):
    user_id: int
    email: str
    role: ProjectRole
    is_direct: bool = True # Default to True for backward compatibility/simplicity

    class Config:
        from_attributes = True

class ProjectTeamAdd(BaseModel):
    team_id: int

class ProjectTeamResponse(BaseModel):
    project_id: int
    team_id: int
    team_name: str

    class Config:
        from_attributes = True
