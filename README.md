# Nexus Tasks

![License](https://img.shields.io/badge/license-AGPL_v3-blue)
![Status](https://img.shields.io/badge/status-active-success)

Nexus Tasks is a developer-focused project management platform designed to bridge the gap between personal productivity and team collaboration.

## Why Nexus Tasks?

Modern developers often struggle with two types of tools: simple to-do lists that lack collaboration features, and complex enterprise suites (like Jira) that disrupt individual flow with excessive configuration and slowness.

Nexus Tasks takes a different approach by offering a **Hybrid Workspace** model:

1.  **Unified Context:** Manage your private side-projects and shared team sprints in a single interface without constant context switching.
2.  **Developer Experience:** Built for speed. Features Markdown support, keyboard navigation, and a clean UI that gets out of your way.
3.  **Simplicity over Configuration:** Opinionated workflows (To Do -> In Progress -> Done) allow teams to start shipping immediately without spending days setting up boards.
4.  **Open & Extensible:** Fully open-source, built with modern standard technologies (**Go** & **Nuxt 4**), making it easy to self-host or extend.

## Features

- **Hybrid Workspaces:** Distinct environments for Personal and Team work.
- **Real-Time Sync:** Instant updates across clients using WebSockets.
- **Kanban & Lists:** Toggle between visual boards and structured lists.
- **Markdown Editor:** Write tasks and comments using standard Markdown.
- **Role-Based Access:** Granular permissions (Admin/Member/Viewer) for projects and teams.
- **Email Notifications:** Asynchronous invitations via background workers.

## Tech Stack

**Frontend**

- **Nuxt 4** (Vue 3)
- **TypeScript**
- **Tailwind CSS 4**
- **Pinia** (State Management)
- **Lucide Vue Next** (Icons)

**Backend**

- **Go 1.26** (Fiber Framework)
- **Goqu** (SQL Builder & Unit of Work)
- **Watermill** (Event-Driven Messaging)
- **PostgreSQL / MySQL / SQLite**
- **Redis** (Message Broker)
- **WebSockets**

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Node.js 22+
- Go 1.26+

### Quick Start (Docker)

To run the entire stack (Frontend, Backend, Database, Redis, Mailpit) in development mode:

1. Clone the repository:

   ```bash
   git clone https://github.com/AshvinBambhaniya/nexus-tasks.git
   ```

2. Navigate to the project directory:

   ```bash
   cd nexus-tasks
   ```

3. Start the services using Docker Compose:
   ```bash
   docker-compose up --build
   ```

The application will be available at:

- Frontend: [http://localhost:3000](http://localhost:3000)
- API: [http://localhost:8000](http://localhost:8000)
- Mailpit (Email Preview): [http://localhost:8025](http://localhost:8025)
- Adminer (DB Management): [http://localhost:8080](http://localhost:8080)

### Local Development Setup

#### 1. Backend

Navigate to the backend directory:

```bash
cd backend
```

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Install Go dependencies:

```bash
go mod download
```

Run database migrations:

```bash
go run app.go migrate up
```

Start the API server:

```bash
go run app.go api
```

Start the background worker:

```bash
go run app.go worker
```

#### 2. Frontend

Navigate to the frontend directory:

```bash
cd frontend
```

Install dependencies:

```bash
npm install
```

Start the development server:

```bash
npm run dev
```

## Contributing

Contributions are welcome! Please read the source code and existing issues before contributing.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0) - see the [LICENSE](LICENSE) file for details.
