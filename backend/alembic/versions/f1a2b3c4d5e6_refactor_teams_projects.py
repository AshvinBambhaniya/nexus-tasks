"""refactor_teams_projects

Revision ID: f1a2b3c4d5e6
Revises: edfeb8eb8c5d
Create Date: 2026-01-30 10:00:00.000000

"""

from typing import Sequence, Union

import sqlalchemy as sa
from sqlalchemy.dialects import postgresql

from alembic import op

# revision identifiers, used by Alembic.
revision: str = "f1a2b3c4d5e6"
down_revision: Union[str, None] = "edfeb8eb8c5d"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # 1. Create Enums safely
    bind = op.get_bind()
    if bind.engine.name == "postgresql":
        op.execute(
            "DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'teamrole') THEN CREATE TYPE teamrole AS ENUM ('ADMIN', 'MEMBER'); END IF; END $$;"
        )
        op.execute(
            "DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'projectrole') THEN CREATE TYPE projectrole AS ENUM ('ADMIN', 'MEMBER', 'VIEWER'); END IF; END $$;"
        )

    # 2. Create Teams
    op.create_table(
        "teams",
        sa.Column("id", sa.Integer(), nullable=False),
        sa.Column("name", sa.String(), nullable=False),
        sa.Column("description", sa.Text(), nullable=True),
        sa.Column("workspace_id", sa.Integer(), nullable=False),
        sa.Column("created_at", sa.DateTime(), nullable=True),
        sa.ForeignKeyConstraint(
            ["workspace_id"], ["workspaces.id"], ondelete="CASCADE"
        ),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index(op.f("ix_teams_id"), "teams", ["id"], unique=False)

    # 3. Create Team Members
    op.create_table(
        "team_members",
        sa.Column("team_id", sa.Integer(), nullable=False),
        sa.Column("user_id", sa.Integer(), nullable=False),
        sa.Column(
            "role",
            postgresql.ENUM("ADMIN", "MEMBER", name="teamrole", create_type=False),
            nullable=False,
        ),
        sa.ForeignKeyConstraint(["team_id"], ["teams.id"], ondelete="CASCADE"),
        sa.ForeignKeyConstraint(["user_id"], ["users.id"], ondelete="CASCADE"),
        sa.PrimaryKeyConstraint("team_id", "user_id"),
    )

    # 4. Create Projects
    op.create_table(
        "projects",
        sa.Column("id", sa.Integer(), nullable=False),
        sa.Column("name", sa.String(), nullable=False),
        sa.Column("description", sa.Text(), nullable=True),
        sa.Column("workspace_id", sa.Integer(), nullable=False),
        sa.Column("created_at", sa.DateTime(), nullable=True),
        sa.ForeignKeyConstraint(
            ["workspace_id"], ["workspaces.id"], ondelete="CASCADE"
        ),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index(op.f("ix_projects_id"), "projects", ["id"], unique=False)

    # 5. Create Project Members
    op.create_table(
        "project_members",
        sa.Column("project_id", sa.Integer(), nullable=False),
        sa.Column("user_id", sa.Integer(), nullable=False),
        sa.Column(
            "role",
            postgresql.ENUM(
                "ADMIN", "MEMBER", "VIEWER", name="projectrole", create_type=False
            ),
            nullable=False,
        ),
        sa.ForeignKeyConstraint(["project_id"], ["projects.id"], ondelete="CASCADE"),
        sa.ForeignKeyConstraint(["user_id"], ["users.id"], ondelete="CASCADE"),
        sa.PrimaryKeyConstraint("project_id", "user_id"),
    )

    # 6. Task Migration - Step A: Add nullable project_id
    op.add_column("tasks", sa.Column("project_id", sa.Integer(), nullable=True))
    op.create_foreign_key(
        None, "tasks", "projects", ["project_id"], ["id"], ondelete="CASCADE"
    )

    # 7. Task Migration - Step B: Data Migration
    # Create a 'General' project for each existing workspace
    op.execute("""
        INSERT INTO projects (name, description, workspace_id, created_at)
        SELECT 'General', 'Default project for migration', id, NOW()
        FROM workspaces
    """)

    # Assign existing tasks to the 'General' project of their workspace
    op.execute("""
        UPDATE tasks
        SET project_id = projects.id
        FROM projects
        WHERE tasks.workspace_id = projects.workspace_id
        AND projects.name = 'General'
    """)

    # 8. Task Migration - Step C: Enforce Not Null
    # (Assuming all tasks belonged to a workspace, they should now have a project_id)
    # If there are orphan tasks (no workspace), they might fail this.
    # But schema says workspace_id was Not Null.
    op.alter_column("tasks", "project_id", nullable=False)

    # 9. Task Migration - Step D: Remove workspace_id
    op.drop_constraint("tasks_workspace_id_fkey", "tasks", type_="foreignkey")
    op.drop_column("tasks", "workspace_id")


def downgrade() -> None:
    # Reverse of upgrade

    # 1. Add workspace_id back
    op.add_column("tasks", sa.Column("workspace_id", sa.Integer(), nullable=True))
    op.create_foreign_key(
        "tasks_workspace_id_fkey",
        "tasks",
        "workspaces",
        ["workspace_id"],
        ["id"],
        ondelete="CASCADE",
    )

    # 2. Restore data (best effort)
    op.execute("""
        UPDATE tasks
        SET workspace_id = projects.workspace_id
        FROM projects
        WHERE tasks.project_id = projects.id
    """)

    op.alter_column("tasks", "workspace_id", nullable=False)

    # 3. Drop project_id
    op.drop_constraint(None, "tasks", type_="foreignkey")
    op.drop_column("tasks", "project_id")

    # 4. Drop tables
    op.drop_table("project_members")
    op.drop_index(op.f("ix_projects_id"), table_name="projects")
    op.drop_table("projects")
    op.drop_table("team_members")
    op.drop_index(op.f("ix_teams_id"), table_name="teams")
    op.drop_table("teams")

    # 5. Drop Enums
    bind = op.get_bind()
    if bind.engine.name == "postgresql":
        sa.Enum(name="projectrole").drop(bind, checkfirst=True)
        sa.Enum(name="teamrole").drop(bind, checkfirst=True)
