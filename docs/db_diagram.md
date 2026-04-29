# Database Diagram

This document contains the ER diagram for the database schema based on the migrations.

```mermaid
erDiagram
    users {
        int id PK
        varchar email UK
        varchar full_name
        varchar hashed_password
        boolean is_active
    }

    workspaces {
        int id PK
        varchar name
        workspacetype type
        int owner_id FK
    }

    workspace_members {
        int workspace_id PK,FK
        int user_id PK,FK
        workspacerole role
    }

    teams {
        int id PK
        varchar name
        text description
        int workspace_id FK
        timestamp created_at
    }

    team_members {
        int team_id PK,FK
        int user_id PK,FK
        teamrole role
    }

    projects {
        int id PK
        varchar name
        text description
        boolean is_archived
        int workspace_id FK
        timestamp created_at
    }

    project_members {
        int project_id PK,FK
        int user_id PK,FK
        projectrole role
    }

    project_teams {
        int project_id PK,FK
        int team_id PK,FK
    }

    tasks {
        int id PK
        varchar title
        text description
        taskstatus status
        taskpriority priority
        int project_id FK
        int assignee_id FK
        int author_id FK
        timestamp due_date
        timestamp completed_at
        timestamp created_at
        timestamp updated_at
    }

    comments {
        int id PK
        text content
        int task_id FK
        int author_id FK
        timestamp created_at
        timestamp updated_at
    }

    users ||--o{ workspaces : owns
    users ||--o{ workspace_members : "member of"
    workspaces ||--o{ workspace_members : has
    workspaces ||--o{ teams : contains
    workspaces ||--o{ projects : contains
    users ||--o{ team_members : "member of"
    teams ||--o{ team_members : has
    users ||--o{ project_members : "member of"
    projects ||--o{ project_members : has
    projects ||--o{ project_teams : assigned
    teams ||--o{ project_teams : assigned
    projects ||--o{ tasks : contains
    users ||--o{ tasks : "assigned to"
    users ||--o{ tasks : authored
    tasks ||--o{ comments : has
    users ||--o{ comments : authored
```

## Enums

- workspacetype: PERSONAL, TEAM
- workspacerole: ADMIN, MEMBER, VIEWER
- teamrole: ADMIN, MEMBER
- projectrole: ADMIN, MEMBER, VIEWER
- taskstatus: TODO, IN_PROGRESS, DONE, BACKLOG
- taskpriority: P0, P1, P2, P3

## Relationships Summary

- A user can own multiple workspaces.
- Workspaces have members (users) with roles.
- Workspaces contain teams and projects.
- Teams have members (users) with roles.
- Projects have members (users) with roles and can be assigned teams.
- Projects contain tasks.
- Tasks can be assigned to users and have an author.
- Tasks have comments, each with an author.
