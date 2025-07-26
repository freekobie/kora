# Kora

Kora is a file management and storage service, built with Go for the backend and React for the frontend, using PostgreSQL as the database.

## Features

- User registration, authentication, and email verification
- File and folder management
- RESTful API endpoints
- JWT-based authentication
- React-based frontend

> **Check TODO.md to see all planned features**

## Getting Started

### Prerequisites

- Go 1.24+
- Node.js & npm
- PostgreSQL
- [Goose](https://github.com/pressly/goose) for database migrations
- (Optional) [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)
- (Optional) [Taskfile](https://taskfile.dev/)

### Setup (Manual)

1.  Clone the repository:
    ```sh
    git clone <your-repo-url>
    cd kora
    ```

2.  Copy `.env.example` to `.env` and fill in your environment variables.

3.  Run database migrations:
    ```sh
    go install github.com/pressly/goose/v3/cmd/goose@latest
    goose -dir ./migrations postgres "$DB_URL" up
    ```

4.  Install frontend dependencies:
    ```sh
    cd frontend
    npm install
    ```

5.  Start the servers (in separate terminals):
    ```sh
    # Backend
    go run ./cmd/api

    # Frontend
    cd frontend
    npm run dev
    ```

### Alternate: Using Docker and Docker Compose

You can run the backend, frontend, and database using Docker Compose:

```sh
docker compose up --build
```

This will build the Go backend and React frontend, and start the application and a PostgreSQL database using the provided `compose.yaml` and `Dockerfile`s.

### Alternate: Using Taskfile

If you have [Task](https://taskfile.dev/) installed, you can use the included `Taskfile.yaml` for common development tasks:

-   Run the server:
    ```sh
    task run
    ```
-   Run tests:
    ```sh
    task test
    ```
-   Run database migrations:
    ```sh
    task up
    ```

## Documentation

The API documentation is available via Swagger:

> **http://localhost:8080/swagger/index.html**

## Running Tests

```sh
go test ./...
```

## Run a development server

You can use [Air](https://github.com/cosmtrek/air) for live-reloading the backend during development:

```sh
go install github.com/cosmtrek/air@latest
air
```

For the frontend, Vite provides its own hot-reloading development server:

```sh
cd frontend
npm run dev
```

## Project Structure

-   `cmd/api/` - Main application entrypoint, configuration, and server setup
-   `handler/` - HTTP route handlers
-   `service/` - Business logic
-   `model/` - Data models and interfaces
-   `postgres/` - PostgreSQL implementations
-   `migrations/` - Database schema migrations
-   `frontend/` - React frontend application
