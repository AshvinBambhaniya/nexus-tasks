"""add_author_to_task

Revision ID: d3a17e26b514
Revises: f1a345851ee3
Create Date: 2026-02-06 18:21:46.463727

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = 'd3a17e26b514'
down_revision: Union[str, None] = 'f1a345851ee3'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.add_column('tasks', sa.Column('author_id', sa.Integer(), nullable=True))
    op.create_foreign_key(None, 'tasks', 'users', ['author_id'], ['id'], ondelete='SET NULL')


def downgrade() -> None:
    op.drop_constraint(None, 'tasks', type_='foreignkey')
    op.drop_column('tasks', 'author_id')
