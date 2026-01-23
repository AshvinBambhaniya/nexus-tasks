from typing import List
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from core.database import get_db
from models.user import User
from models.workspace import Workspace, WorkspaceMember, WorkspaceType, WorkspaceRole
from schemas.workspace import WorkspaceCreate, WorkspaceResponse, WorkspaceMemberInvite, WorkspaceMemberResponse
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

@router.get("/{workspace_id}/members", response_model=List[WorkspaceMemberResponse])
def list_workspace_members(
    workspace_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    List all members of a workspace.
    """
    validate_workspace_access(workspace_id, db, current_user.id)
    
    members = db.query(WorkspaceMember).filter(
        WorkspaceMember.workspace_id == workspace_id
    ).join(WorkspaceMember.user).all()
    
    return members

@router.post("/{workspace_id}/members", response_model=WorkspaceMemberResponse)
def invite_member(
    workspace_id: int,
    invite: WorkspaceMemberInvite,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Invite a user to the workspace. Only ADMINs can invite.
    """
    # 1. Validate Admin Access
    member = validate_workspace_access(workspace_id, db, current_user.id)
    if member.role != WorkspaceRole.ADMIN:
         raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN, 
            detail="Only Admins can invite members"
        )

    # 2. Check if user exists
    user_to_add = db.query(User).filter(User.email == invite.email).first()
    if not user_to_add:
        raise HTTPException(status_code=404, detail="User with this email not found")

    # 3. Check if already member
    existing_member = db.query(WorkspaceMember).filter(
        WorkspaceMember.workspace_id == workspace_id,
        WorkspaceMember.user_id == user_to_add.id
    ).first()
    if existing_member:
        raise HTTPException(status_code=400, detail="User is already a member")

    # 4. Add Member
    new_member = WorkspaceMember(
        workspace_id=workspace_id,
        user_id=user_to_add.id,
        role=WorkspaceRole.MEMBER
    )
    db.add(new_member)
    db.commit()
    db.refresh(new_member)
    
    return new_member

@router.delete("/{workspace_id}/members/{user_id}")
def remove_member(
    workspace_id: int,
    user_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Remove a member from the workspace. Only ADMINs can remove.
    """
    # 1. Validate Admin Access
    requester = validate_workspace_access(workspace_id, db, current_user.id)
    if requester.role != WorkspaceRole.ADMIN:
         raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN, 
            detail="Only Admins can remove members"
        )
    
    # 2. Get Member to remove
    member_to_remove = db.query(WorkspaceMember).filter(
        WorkspaceMember.workspace_id == workspace_id,
        WorkspaceMember.user_id == user_id
    ).first()
    
    if not member_to_remove:
        raise HTTPException(status_code=404, detail="Member not found")
        
    # Prevent removing self if it leaves no admins? (Simple check: can remove self)
    # But usually we prevent removing the last admin. 
    # For now, let's just allow it, but maybe prevent removing self via this specific endpoint if we want explicit 'leave' logic separately.
    # Requirement says "Remove a member". 
    
    db.delete(member_to_remove)
    db.commit()
    
    return {"status": "success", "message": "Member removed"}
