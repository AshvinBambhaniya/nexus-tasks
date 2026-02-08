from typing import List
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from core.database import get_db
from models.user import User
from models.team import Team, TeamMember, TeamRole
from models.workspace import WorkspaceMember, WorkspaceRole
from schemas.team import TeamCreate, TeamResponse, TeamMemberAdd, TeamMemberResponse
from api.v1.auth import get_current_user

router = APIRouter()

def validate_team_access(team_id: int, db: Session, user_id: int, require_admin: bool = False) -> TeamMember:
    """
    Check if user is a member of the team.
    Workspace Admins also implicitly have access (managed via separate check usually, but for simplicity here we check direct membership).
    To strictly follow RBAC: Workspace Admin should have access even if not in TeamMember table. 
    But typically we add them to the team or check workspace role too.
    """
    # Check direct team membership
    member = db.query(TeamMember).filter(
        TeamMember.team_id == team_id,
        TeamMember.user_id == user_id
    ).first()

    if member:
        if require_admin and member.role != TeamRole.ADMIN:
             # Check if they are workspace admin?
             # Let's fallback to checking workspace admin if not team admin
             pass
        else:
            return member

    # If not in team, check if Workspace Admin
    team = db.query(Team).filter(Team.id == team_id).first()
    if not team:
        raise HTTPException(status_code=404, detail="Team not found")
        
    ws_member = db.query(WorkspaceMember).filter(
        WorkspaceMember.workspace_id == team.workspace_id,
        WorkspaceMember.user_id == user_id
    ).first()

    if ws_member and ws_member.role == WorkspaceRole.ADMIN:
        return TeamMember(team_id=team_id, user_id=user_id, role=TeamRole.ADMIN) # Mock object for auth success

    if member:
         # Was member but failed admin check above
         raise HTTPException(status_code=403, detail="Team Admin privileges required")

    raise HTTPException(status_code=403, detail="Not a member of this team")

@router.post("/workspaces/{workspace_id}/teams", response_model=TeamResponse)
def create_team(
    workspace_id: int,
    team_data: TeamCreate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Create a new Team in a Workspace.
    Only Workspace Admins can create teams.
    """
    # Verify Workspace Admin
    ws_member = db.query(WorkspaceMember).filter(
        WorkspaceMember.workspace_id == workspace_id,
        WorkspaceMember.user_id == current_user.id
    ).first()

    if not ws_member or ws_member.role != WorkspaceRole.ADMIN:
        raise HTTPException(status_code=403, detail="Only Workspace Admins can create teams")

    new_team = Team(
        name=team_data.name,
        description=team_data.description,
        workspace_id=workspace_id
    )
    db.add(new_team)
    db.commit()
    db.refresh(new_team)

    # Add Creator as Team Admin (Optional, but good for UX)
    member = TeamMember(
        team_id=new_team.id,
        user_id=current_user.id,
        role=TeamRole.ADMIN
    )
    db.add(member)
    db.commit()
    
    return new_team

@router.get("/workspaces/{workspace_id}/teams", response_model=List[TeamResponse])
def list_workspace_teams(
    workspace_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    List all teams in a workspace.
    User must be a member of the workspace.
    """
    # Check basic workspace membership
    ws_member = db.query(WorkspaceMember).filter(
        WorkspaceMember.workspace_id == workspace_id,
        WorkspaceMember.user_id == current_user.id
    ).first()
    
    if not ws_member:
        raise HTTPException(status_code=403, detail="Not a member of this workspace")

    # Return all teams (open visibility) or only teams user is in?
    # Requirement: "Teams (Independent Group of Users)". Usually you can see what teams exist.
    teams = db.query(Team).filter(Team.workspace_id == workspace_id).all()
    return teams

@router.post("/teams/{team_id}/members", response_model=TeamMemberResponse)
def add_team_member(
    team_id: int,
    member_data: TeamMemberAdd,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Add a user to a team.
    Requires Team Admin or Workspace Admin.
    """
    validate_team_access(team_id, db, current_user.id, require_admin=True)

    # Find user
    user_to_add = db.query(User).filter(User.email == member_data.email).first()
    if not user_to_add:
        raise HTTPException(status_code=404, detail="User not found")

    # Check existence
    exists = db.query(TeamMember).filter(
        TeamMember.team_id == team_id,
        TeamMember.user_id == user_to_add.id
    ).first()
    if exists:
        raise HTTPException(status_code=400, detail="User already in team")

    new_member = TeamMember(
        team_id=team_id,
        user_id=user_to_add.id,
        role=member_data.role
    )
    db.add(new_member)
    db.commit()
    db.refresh(new_member)

    return {
        "user_id": new_member.user_id,
        "role": new_member.role,
        "email": user_to_add.email
    }

@router.get("/teams/{team_id}/members", response_model=List[TeamMemberResponse])
def list_team_members(
    team_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    validate_team_access(team_id, db, current_user.id)
    
    members = db.query(TeamMember).filter(TeamMember.team_id == team_id).join(TeamMember.user).all()
    # Pydantic will extract email from the joined User relationship if configured, 
    # but schema expects flat structure. 
    # We might need to manually construct response or rely on ORM mapping if schema matches.
    # Schema: user_id, email, role.
    # TeamMember has user_id, role. 'email' is on .user.email.
    
    # Let's return list of dicts to be safe
    return [
        {
            "user_id": m.user_id,
            "role": m.role,
            "email": m.user.email
        }
        for m in members
    ]
