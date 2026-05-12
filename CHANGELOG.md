# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-05-12

### 🚀 Major Migration & Performance Overhaul

This release marks a significant milestone in Nexus Tasks history. We have completely rewritten both our backend and frontend stacks to improve performance, developer experience, and scalability.

#### **Backend Migration: FastAPI to Go (Fiber)**

- **Architecture**: Transitioned from Python/FastAPI to a high-performance **Go** backend using the **Fiber** framework.
- **Data Access**: Moved to a **Unit of Work** pattern with **Goqu** for atomic and safe database transactions.
- **Messaging**: Replaced internal task queues with **Watermill** and Redis for robust, event-driven architecture.
- **Real-time**: Custom WebSocket Hub implementation in Go for lower latency and better resource management.

#### **Frontend Migration: Next.js to Nuxt 4**

- **Framework**: Migrated the entire UI from React/Next.js to **Nuxt 4 (Vue 3)**, leveraging the latest features of the Nuxt Nitro engine.
- **Styling**: Upgraded to **Tailwind CSS 4** for faster builds and a more modern design system.
- **State Management**: Adopted **Pinia** for clean, modular reactive state handling.

### ✨ Other Improvements

- **JSend Compliance**: Standardized all API responses following the JSend specification.
- **Improved Testing**: Comprehensive backend test suite with high coverage across services and controllers.
- **Unified DX**: Simplified local development with an optimized Docker Compose setup.

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
