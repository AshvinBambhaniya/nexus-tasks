from typing import List
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from core.database import get_db
from models.user import User
from models.workspace import Workspace, WorkspaceMember, WorkspaceType, WorkspaceRole
from schemas.workspace import WorkspaceCreate, WorkspaceResponse
from api.v1.auth import get_current_user

router = APIRouter()

@router.post("/", response_model=WorkspaceResponse)
def create_workspace(
    workspace: WorkspaceCreate, 
    current_user: User = Depends(get_current_user), 
    db: Session = Depends(get_db)
):
    """
    Create a new TEAM workspace. 
    The creator becomes the ADMIN.
    """
    new_ws = Workspace(
        name=workspace.name,
        type=WorkspaceType.TEAM, # Explicitly creating a TEAM workspace
        owner_id=current_user.id
    )
    db.add(new_ws)
    db.commit()
    db.refresh(new_ws)

    # Add creator as ADMIN
    member = WorkspaceMember(
        workspace_id=new_ws.id,
        user_id=current_user.id,
        role=WorkspaceRole.ADMIN
    )
    db.add(member)
    db.commit()
    
    return new_ws

@router.get("/", response_model=List[WorkspaceResponse])
def list_my_workspaces(
    current_user: User = Depends(get_current_user), 
    db: Session = Depends(get_db)
):
    """
    List all workspaces the current user is a member of.
    """
    # Join WorkspaceMember to find workspaces for this user
    workspaces = (
        db.query(Workspace)
        .join(WorkspaceMember, Workspace.id == WorkspaceMember.workspace_id)
        .filter(WorkspaceMember.user_id == current_user.id)
        .all()
    )
    return workspaces

# Reusable Dependency for future endpoints (e.g. Tasks)
def validate_workspace_access(workspace_id: int, db: Session, user_id: int) -> WorkspaceMember:
    member = db.query(WorkspaceMember).filter(
        WorkspaceMember.workspace_id == workspace_id,
        WorkspaceMember.user_id == user_id
    ).first()
    
    if not member:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN, 
            detail="Not a member of this workspace"
        )
    return member
