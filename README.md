# Task Management Web Service

This project is a Task Management Web Service built in Go, providing a RESTful API to manage tasks and users. It supports CRUD operations for both entities, with strong validation and user-context for all task operations.

## Project Structure

```
task-planner
├── cmd
│   └── main.go                # Entry point of the application
├── internal
│   ├── api
│   │   └── handler.go         # HTTP handlers for the RESTful API
│   ├── db
│   │   ├── database.go        # Database implementation (SQLite)
│   │   └── interface.go       # Database interface abstraction
│   ├── model
│   │   ├── task.go            # Task struct definition
│   │   └── user.go            # User struct definition
│   └── service
│       └── task_service.go    # Business logic for managing tasks
├── test
│   ├── api_test.go            # Ginkgo tests for API handlers and validation
│   ├── db_test.go             # Unit tests for database operations
│   └── service_test.go        # Unit tests for task service
├── go.mod                     # Module definition file
├── go.sum                     # Module dependency checksums
└── README.md                  # Project documentation
```

## Database Plug-and-Play

The project uses a database interface abstraction (`internal/db/interface.go`).  
This makes it easy to swap between different SQL databases (e.g., SQLite, PostgreSQL, MySQL):

- Implement the `DB` interface for your preferred database.
- Update the initialization in `main.go` to use your chosen implementation.
- No changes required in the service or API layers.

Currently, the project uses SQLite for local development and testing.

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd <project-directory>
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Run the application:**
   ```
   go run cmd/main.go
   ```

## API Usage

### User APIs

#### Create User

- **Endpoint:** `POST /users`
- **Request Body:**
  ```json
  {
    "name": "User Name",
    "email": "user@example.com"
  }
  ```
  - `name`: Required, 2–50 characters.
  - `email`: Required.

#### Get User

- **Endpoint:** `GET /users/{user_id}`

#### List Users

- **Endpoint:** `GET /users`

#### Delete User

- **Endpoint:** `DELETE /users/{user_id}`

### Task APIs (under user context)

#### Create Task

- **Endpoint:** `POST /users/{user_id}/tasks`
- **Request Body:**
  ```json
  {
    "title": "Task Title",
    "description": "Task Description",
    "due_date": "2025-12-31T10:00:00Z",
    "status": "pending"
  }
  ```
  - `title`: Required, 2–50 characters.
  - `description`: Optional, max 200 characters.
  - `due_date`: Required, must be ISO 8601 format (RFC3339, e.g. `"2025-12-31T10:00:00Z"`).
  - `status`: Required, must be one of `"pending"`, `"in_progress"`, `"done"`.

#### Get Tasks for a User

- **Endpoint:** `GET /users/{user_id}/tasks`

#### Get Task by ID

- **Endpoint:** `GET /users/{user_id}/tasks/{task_id}`

#### Update Task

- **Endpoint:** `PUT /users/{user_id}/tasks/{task_id}`
- **Request Body:**
  ```json
  {
    "title": "Updated Task Title",
    "description": "Updated Task Description",
    "due_date": "2025-12-31T10:00:00Z",
    "status": "done"
  }
  ```
  - Same validation as create.

#### Delete Task

- **Endpoint:** `DELETE /users/{user_id}/tasks/{task_id}`

## Validation Rules

- **User name:** 2–50 characters.
- **Task title:** 2–50 characters.
- **Task description:** Up to 200 characters.
- **Task status:** Must be `"pending"`, `"in_progress"`, or `"done"`.
- **Task due_date:** Must be ISO 8601 date/time (RFC3339).

## Running Tests

To run the unit and API tests (Ginkgo required):

```
go test ./...
ginkgo ./test
```

## License

This project is licensed under the MIT License.