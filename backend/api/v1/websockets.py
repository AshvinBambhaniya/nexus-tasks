from fastapi import APIRouter, WebSocket, WebSocketDisconnect, Depends, status, Query
from sqlalchemy.orm import Session
from jose import JWTError, jwt

from core.database import get_db
from core.websocket import manager
from core.config import settings
from models.user import User
from models.workspace import WorkspaceMember
from api.v1.workspaces import validate_workspace_access # We might need to catch the exception or re-implement logic to avoid HTTP exception

router = APIRouter()

async def get_current_user_ws(
    token: str = Query(...), 
    db: Session = Depends(get_db)
) -> User:
    try:
        payload = jwt.decode(token, settings.SECRET_KEY, algorithms=[settings.ALGORITHM])
        email: str = payload.get("sub")
        if email is None:
            raise WebSocketDisconnect(code=status.WS_1008_POLICY_VIOLATION)
    except JWTError:
        raise WebSocketDisconnect(code=status.WS_1008_POLICY_VIOLATION)
        
    user = db.query(User).filter(User.email == email).first()
    if user is None:
        raise WebSocketDisconnect(code=status.WS_1008_POLICY_VIOLATION)
    return user

@router.websocket("/ws/{workspace_id}")
async def websocket_endpoint(
    websocket: WebSocket, 
    workspace_id: int, 
    token: str = Query(...),
    db: Session = Depends(get_db)
):
    # Authenticate manually since Depends(get_current_user_ws) inside websocket_endpoint signature 
    # might execute before we accept the connection or handle rejection differently.
    # However, FastAPI handles dependencies for WebSockets too.
    
    # We authenticate and validate access BEFORE accepting the connection to reject unauthorized ones cleanly.
    try:
        user = await get_current_user_ws(token, db)
        
        # Validate Workspace Access
        # validate_workspace_access raises HTTPException which FastAPI might convert to WS close, 
        # but let's be explicit.
        member = db.query(WorkspaceMember).filter(
            WorkspaceMember.workspace_id == workspace_id,
            WorkspaceMember.user_id == user.id
        ).first()
        
        if not member:
            raise WebSocketDisconnect(code=status.WS_1008_POLICY_VIOLATION)
            
    except Exception:
        # If any auth/validation fails, we close (or let FastAPI close it if it was a dependency failure)
        # Note: If we haven't called accept(), we can just return/raise. 
        # But if we want to send a specific close code, we might need to accept then close, 
        # or just let the client timeout/fail. 
        # Standard practice: accept then close with code if logic permits, or just don't accept.
        # FastAPI's WebSocketDisconnect is usually raised during receive, not before accept.
        # To reject:
        await websocket.close(code=status.WS_1008_POLICY_VIOLATION)
        return

    await manager.connect(websocket, workspace_id)
    try:
        while True:
            # We just keep the connection open.
            # We can also handle incoming messages if the frontend sends any (e.g. "ping")
            await websocket.receive_text() 
    except WebSocketDisconnect:
        manager.disconnect(websocket, workspace_id)
