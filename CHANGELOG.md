# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-02-16

### 🚀 New Features

- **Task Completion Tracking**: Added `completed_at` field to tasks to track exactly when they are moved to the "DONE" status. This enables better historical data and future cycle time analytics.
- **Enhanced Due Date Support**: Improved UI for setting and viewing due dates across task boards and list views.
- **Workspace Invitation Emails**: Implemented email notifications for workspace invitations, allowing users to join workspaces seamlessly.

## [1.0.0] - 2026-02-12

### 🚀 Initial Release: Nexus Tasks

We are thrilled to introduce **Nexus Tasks**, a modern, high-performance project management platform designed for speed and collaboration. This initial release brings a complete suite of tools for managing workspaces, teams, and projects with real-time syncing.

### ✨ Key Features

#### **Workspace & Team Management**

- **Multi-Tenant Architecture**: Create isolated **Workspaces** to organize different organizations or large groups.
- **Granular Permissions**: Robust RBAC (Role-Based Access Control) system with specific roles for Workspaces (Admin/Member) and Projects.
- **Team Management**: Group users into Teams for easier project assignment and collaboration.

#### **Project & Task Tracking**

- **Interactive Kanban Boards**: Fully functional drag-and-drop boards powered by `@dnd-kit`.
- **List Views**: detailed list views for managing high-volume task lists.
- **Rich Text Editing**: Markdown support for task descriptions and comments.
- **Task Attributes**: Priority levels, status tracking, due dates, and assignee management.
- **Inbox**: A centralized notification center for tasks assigned to you.

#### **Productivity & UX**

- **Command Palette (`Cmd+K`)**: Global keyboard-driven navigation to quickly jump between projects, workspaces, or create tasks.
- **Real-Time Collaboration**: Instant updates across all clients using **WebSockets**. See task moves and comments happen live.

### Technical Highlights

- **Frontend**: Built on the bleeding edge with **Next.js 16 (App Router)** and **React 19**.
- **Backend**: High-performance asynchronous REST API using **FastAPI (Python 3.12)**.
- **Database**: Robust data integrity with **PostgreSQL** and **SQLAlchemy 2.0**.
- **State Management**: Ultra-fast client-state handling with **Zustand** and **SWR**.
- **Infrastructure**: Production-ready **Docker** and **Docker Compose** setups with CI/CD pipelines via GitHub Actions.
