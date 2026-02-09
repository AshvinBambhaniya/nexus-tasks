from fastapi import APIRouter, Depends, WebSocket, WebSocketDisconnect, status
from jose import JWTError, jwt
from sqlalchemy.orm import Session

from core.config import settings
from core.database import get_db
from core.websocket import manager
from models.user import User
from models.workspace import WorkspaceMember

router = APIRouter()


async def get_current_user_ws(
    websocket: WebSocket, db: Session = Depends(get_db)
) -> User:
    token = websocket.cookies.get("access_token")

    if not token:
        raise WebSocketDisconnect(code=status.WS_1008_POLICY_VIOLATION)

    try:
        payload = jwt.decode(
            token, settings.SECRET_KEY, algorithms=[settings.ALGORITHM]
        )
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
    websocket: WebSocket, workspace_id: int, db: Session = Depends(get_db)
):
    try:
        # We must NOT await websocket.accept() yet if we want to reject with 1008
        # However, get_current_user_ws checks cookies which are available in handshake
        user = await get_current_user_ws(websocket, db)

        member = (
            db.query(WorkspaceMember)
            .filter(
                WorkspaceMember.workspace_id == workspace_id,
                WorkspaceMember.user_id == user.id,
            )
            .first()
        )

        if not member:
            raise WebSocketDisconnect(code=status.WS_1008_POLICY_VIOLATION)

    except Exception:
        await websocket.close(code=status.WS_1008_POLICY_VIOLATION)
        return

    await manager.connect(websocket, workspace_id)
    try:
        while True:
            await websocket.receive_text()
    except WebSocketDisconnect:
        manager.disconnect(websocket, workspace_id)
