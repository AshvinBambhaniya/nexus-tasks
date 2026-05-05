# Nexus Tasks Backend

The Nexus Tasks backend is a high-performance, event-driven API built with Go.

## Architecture

-   **Web Framework**: [Fiber](https://gofiber.io/) - Express-inspired web framework for Go.
-   **Database Layer**: [Goqu](https://github.com/doug-martin/goqu) - Flexible SQL builder and data mapper.
-   **Messaging**: [Watermill](https://watermill.io/) - Library for efficiently working with message streams.
-   **CLI**: [Cobra](https://github.com/spf13/cobra) - Multi-command CLI structure.
-   **Real-time**: Custom WebSocket Hub for state synchronization.

## Directory Structure

-   `cli/`: Command definitions for the application binary (`api`, `worker`, `migrate`).
-   `config/`: Environment-based configuration management.
-   `controllers/`: HTTP request handlers.
-   `database/`: Database connection and migrations.
-   `models/`: Data structures and Storage (Unit of Work) implementation.
-   `pkg/`: Shared packages like JWT, Real-time (Hub), and Watermill helpers.
-   `routes/`: API route definitions.
-   `services/`: Business logic layer.

## Commands

### API Server

Starts the Fiber HTTP server.

```bash
go run app.go api
```

### Background Worker

Starts the Watermill consumer for background tasks (emails, etc.).

```bash
go run app.go worker
```

### Migrations

Manage database schema.

```bash
go run app.go migrate up
go run app.go migrate down
```

### Dead Letter Queue

Handle failed messages.

```bash
go run app.go dead-letter-queue
```

## Development

Use `nodemon` for hot-reloading during development:

```bash
nodemon --exec go run app.go api --signal SIGTERM
```
