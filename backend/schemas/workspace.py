from pydantic import BaseModel, EmailStr
from typing import Optional
from models.workspace import WorkspaceType, WorkspaceRole

class WorkspaceBase(BaseModel):
    name: str

class WorkspaceCreate(WorkspaceBase):
    pass

class WorkspaceResponse(WorkspaceBase):
    id: int
    type: WorkspaceType
    owner_id: int

    class Config:
        from_attributes = True

class WorkspaceMemberInvite(BaseModel):
    email: EmailStr

class UserInfo(BaseModel):
    id: int
    email: str
    
    class Config:
        from_attributes = True

class WorkspaceMemberResponse(BaseModel):
    workspace_id: int
    user_id: int
    role: WorkspaceRole
    user: UserInfo # Include user details
    
    class Config:
        from_attributes = True
