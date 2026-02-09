from datetime import datetime

from pydantic import BaseModel, ConfigDict

from schemas.user import UserResponse


class CommentBase(BaseModel):
    content: str


class CommentCreate(CommentBase):
    pass


class CommentResponse(CommentBase):
    id: int
    task_id: int
    author_id: int
    created_at: datetime
    updated_at: datetime
    author: UserResponse

    model_config = ConfigDict(from_attributes=True)
