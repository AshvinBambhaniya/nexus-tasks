from typing import List
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from core.database import get_db
from models.user import User
from models.project import Project, ProjectMember, ProjectRole, ProjectTeam
from models.team import TeamMember, Team
from models.workspace import WorkspaceMember, WorkspaceRole
from schemas.project import ProjectCreate, ProjectResponse, ProjectMemberAdd, ProjectMemberResponse, ProjectTeamAdd, ProjectTeamResponse, ProjectUpdate
from api.v1.auth import get_current_user

router = APIRouter()

def validate_project_access(project_id: int, db: Session, user_id: int, require_admin: bool = False) -> ProjectMember:
    # Check direct project membership
    member = db.query(ProjectMember).filter(
        ProjectMember.project_id == project_id,
        ProjectMember.user_id == user_id
    ).first()

    if member:
        if require_admin and member.role != ProjectRole.ADMIN:
             pass # Fallback to workspace check
        else:
            return member

    # Check Team Access (Implicit Membership)
    # If a user is a member of a team that is assigned to this project
    team_access = db.query(ProjectTeam).join(
        TeamMember, ProjectTeam.team_id == TeamMember.team_id
    ).filter(
        ProjectTeam.project_id == project_id,
        TeamMember.user_id == user_id
    ).first()

    if team_access:
        # If require_admin is True, Team access (usually MEMBER role) is likely insufficient 
        # unless we define Team Roles in Projects (future). 
        # For now, Team access grants standard MEMBER access.
        if not require_admin:
             return ProjectMember(project_id=project_id, user_id=user_id, role=ProjectRole.MEMBER)

    # Check Workspace Admin
    project = db.query(Project).filter(Project.id == project_id).first()

    if not project:
        raise HTTPException(status_code=404, detail="Project not found")
        
    ws_member = db.query(WorkspaceMember).filter(
        WorkspaceMember.workspace_id == project.workspace_id,
        WorkspaceMember.user_id == user_id
    ).first()

    if ws_member and ws_member.role == WorkspaceRole.ADMIN:
        return ProjectMember(project_id=project_id, user_id=user_id, role=ProjectRole.ADMIN)

    if member:
         raise HTTPException(status_code=403, detail="Project Admin privileges required")

    raise HTTPException(status_code=403, detail="Not a member of this project")

@router.post("/workspaces/{workspace_id}/projects", response_model=ProjectResponse)
def create_project(
    workspace_id: int,
    project_data: ProjectCreate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Create a new Project in a Workspace.
    Workspace Admins can create projects.
    """
    ws_member = db.query(WorkspaceMember).filter(
        WorkspaceMember.workspace_id == workspace_id,
        WorkspaceMember.user_id == current_user.id
    ).first()

    if not ws_member or ws_member.role != WorkspaceRole.ADMIN:
        raise HTTPException(status_code=403, detail="Only Workspace Admins can create projects")

    new_project = Project(
        name=project_data.name,
        description=project_data.description,
        workspace_id=workspace_id
    )
    db.add(new_project)
    db.commit()
    db.refresh(new_project)

    # Add Creator as Admin
    member = ProjectMember(
        project_id=new_project.id,
        user_id=current_user.id,
        role=ProjectRole.ADMIN
    )
    db.add(member)
    db.commit()
    
    return new_project

@router.get("/workspaces/{workspace_id}/projects", response_model=List[ProjectResponse])
def list_workspace_projects(
    workspace_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    ws_member = db.query(WorkspaceMember).filter(
        WorkspaceMember.workspace_id == workspace_id,
        WorkspaceMember.user_id == current_user.id
    ).first()
    
    if not ws_member:
        raise HTTPException(status_code=403, detail="Not a member of this workspace")

    # Return all active projects in workspace
    projects = db.query(Project).filter(
        Project.workspace_id == workspace_id,
        Project.is_archived == False
    ).all()
    return projects

@router.post("/projects/{project_id}/members", response_model=ProjectMemberResponse)
def add_project_member(
    project_id: int,
    member_data: ProjectMemberAdd,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    validate_project_access(project_id, db, current_user.id, require_admin=True)

    user_to_add = db.query(User).filter(User.email == member_data.email).first()
    if not user_to_add:
        raise HTTPException(status_code=404, detail="User not found")

    project = db.query(Project).filter(Project.id == project_id).first()
    if not project:
        raise HTTPException(status_code=404, detail="Project not found")

    # Check if user is in Workspace
    ws_member = db.query(WorkspaceMember).filter(
        WorkspaceMember.workspace_id == project.workspace_id,
        WorkspaceMember.user_id == user_to_add.id
    ).first()
    if not ws_member:
        raise HTTPException(status_code=400, detail="User must be a member of the workspace first")

    exists = db.query(ProjectMember).filter(
        ProjectMember.project_id == project_id,
        ProjectMember.user_id == user_to_add.id
    ).first()
    if exists:
        raise HTTPException(status_code=400, detail="User already in project")

    new_member = ProjectMember(
        project_id=project_id,
        user_id=user_to_add.id,
        role=member_data.role
    )
    db.add(new_member)
    db.commit()
    db.refresh(new_member)

    return new_member

@router.get("/projects/{project_id}/members", response_model=List[ProjectMemberResponse])
def list_project_members(
    project_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    validate_project_access(project_id, db, current_user.id)
    
    # 1. Direct Members
    direct_members = db.query(ProjectMember).filter(ProjectMember.project_id == project_id).join(ProjectMember.user).all()
    
    # 2. Team Members (Implicit Access)
    team_members = db.query(TeamMember).join(
        ProjectTeam, ProjectTeam.team_id == TeamMember.team_id
    ).filter(
        ProjectTeam.project_id == project_id
    ).join(TeamMember.user).all()

    # 3. Combine and Deduplicate
    all_members = {}
    
    # Add direct members
    for m in direct_members:
        all_members[m.user_id] = {
            "user_id": m.user_id,
            "role": m.role,
            "email": m.user.email,
            "is_direct": True
        }
        
    # Add team members if not already present
    for m in team_members:
        if m.user_id not in all_members:
            all_members[m.user_id] = {
                "user_id": m.user_id,
                "role": ProjectRole.MEMBER, # Implicit members get MEMBER role
                "email": m.user.email,
                "is_direct": False
            }
    
    return list(all_members.values())

@router.delete("/projects/{project_id}/members/{user_id}")
def remove_project_member(
    project_id: int,
    user_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    validate_project_access(project_id, db, current_user.id, require_admin=True)

    member = db.query(ProjectMember).filter(
        ProjectMember.project_id == project_id,
        ProjectMember.user_id == user_id
    ).first()
    
    if not member:
        raise HTTPException(status_code=404, detail="Member not found")

    db.delete(member)
    db.commit()
    
    return {"status": "success", "message": "Member removed"}

@router.post("/projects/{project_id}/teams", response_model=ProjectTeamResponse)
def add_project_team(
    project_id: int,
    team_data: ProjectTeamAdd,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    validate_project_access(project_id, db, current_user.id, require_admin=True)

    # Validate Team
    team = db.query(Team).filter(Team.id == team_data.team_id).first()
    if not team:
        raise HTTPException(status_code=404, detail="Team not found")

    # Validate Project
    project = db.query(Project).filter(Project.id == project_id).first()
    if not project:
        raise HTTPException(status_code=404, detail="Project not found")

    # Ensure Team is in same Workspace
    if team.workspace_id != project.workspace_id:
        raise HTTPException(status_code=400, detail="Team must belong to the same workspace")

    exists = db.query(ProjectTeam).filter(
        ProjectTeam.project_id == project_id,
        ProjectTeam.team_id == team_data.team_id
    ).first()
    if exists:
        raise HTTPException(status_code=400, detail="Team already assigned to project")

    new_link = ProjectTeam(
        project_id=project_id,
        team_id=team_data.team_id
    )
    db.add(new_link)
    db.commit()
    db.refresh(new_link)

    return {
        "project_id": new_link.project_id,
        "team_id": new_link.team_id,
        "team_name": team.name
    }

@router.get("/projects/{project_id}/teams", response_model=List[ProjectTeamResponse])
def list_project_teams(
    project_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    validate_project_access(project_id, db, current_user.id)
    
    links = db.query(ProjectTeam).filter(ProjectTeam.project_id == project_id).join(ProjectTeam.team).all()
    
    return [
        {
            "project_id": l.project_id,
            "team_id": l.team_id,
            "team_name": l.team.name
        }
        for l in links
    ]

@router.delete("/projects/{project_id}/teams/{team_id}")
def remove_project_team(
    project_id: int,
    team_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    validate_project_access(project_id, db, current_user.id, require_admin=True)

    link = db.query(ProjectTeam).filter(
        ProjectTeam.project_id == project_id,
        ProjectTeam.team_id == team_id
    ).first()
    
    if not link:
        raise HTTPException(status_code=404, detail="Team assignment not found")

    db.delete(link)
    db.commit()
    
    return {"status": "success", "message": "Team removed from project"}

@router.patch("/projects/{project_id}", response_model=ProjectResponse)
def update_project(
    project_id: int,
    project_update: ProjectUpdate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    validate_project_access(project_id, db, current_user.id, require_admin=True)

    project = db.query(Project).filter(Project.id == project_id).first()
    if not project:
        raise HTTPException(status_code=404, detail="Project not found")

    update_data = project_update.model_dump(exclude_unset=True)
    for key, value in update_data.items():
        setattr(project, key, value)

    db.commit()
    db.refresh(project)
    
    return project

@router.get("/projects/{project_id}", response_model=ProjectResponse)
def get_project(
    project_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    validate_project_access(project_id, db, current_user.id)
    project = db.query(Project).filter(Project.id == project_id).first()
    if not project:
        raise HTTPException(status_code=404, detail="Project not found")
    return project
