# Task Management Web Service

This project is a Task Management Web Service built in Go, providing a RESTful API to manage tasks and users. It supports CRUD operations for both entities, with strong validation and user-context for all task operations.

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

## CLI

```sh
❯ go run cli.go     
Task CLI. Type 'help' for commands.
> help

Commands:
  create-user           - Create a new user (prompts for name/email)
  set-user <user_id>    - Set current session user
  get-user              - Get current user info
  list-users            - List all users
  delete-user           - Delete current user
  create-task           - Create a new task (prompts for details)
  get-task <task_id>    - Get a task by ID
  list-tasks            - List tasks for current user
  update-task <task_id> - Update a task (prompts for details)
  delete-task <task_id> - Delete a task
  help                  - Show this help
  exit/quit             - Exit CLI
> list-users
Status: 200
[
  {
    "id": "953730b1-e3d4-4dbe-85a4-8244347a8f70",
    "name": "Alice",
    "email": "alice@example.com"
  },
  {
    "id": "9638fe95-5845-4011-902c-f3bcc4c821df",
    "name": "Deepak",
    "email": "deepak@example.com"
  }
]
> set-user 9638fe95-5845-4011-902c-f3bcc4c821df
Session user set to: 9638fe95-5845-4011-902c-f3bcc4c821df
> create-task
title: New Task
description: Some Desc
due_date: 2025-08-10T03:04:10Z
status: pending
Status: 201
{
  "id": "6b9436ad-a159-46c8-86ba-b9bfd43d5cd1",
  "title": "New Task",
  "description": "Some Desc",
  "due_date": "2025-08-10T03:04:10Z",
  "status": "pending",
  "user_id": "9638fe95-5845-4011-902c-f3bcc4c821df"
}
> list-tasks
Status: 200
[
  {
    "id": "32c9f25a-bd7f-4c24-bdff-8c9b8a977ce7",
    "title": "New Title",
    "description": "Some Desc",
    "due_date": "2025-12-30T00:00:00Z",
    "status": "done",
    "user_id": "9638fe95-5845-4011-902c-f3bcc4c821df"
  },
  {
    "id": "6b9436ad-a159-46c8-86ba-b9bfd43d5cd1",
    "title": "New Task",
    "description": "Some Desc",
    "due_date": "2025-08-10T03:04:10Z",
    "status": "pending",
    "user_id": "9638fe95-5845-4011-902c-f3bcc4c821df"
  }
]
> get-task 6b9436ad-a159-46c8-86ba-b9bfd43d5cd1
Status: 200
{
  "id": "6b9436ad-a159-46c8-86ba-b9bfd43d5cd1",
  "title": "New Task",
  "description": "Some Desc",
  "due_date": "2025-08-10T03:04:10Z",
  "status": "pending",
  "user_id": "9638fe95-5845-4011-902c-f3bcc4c821df"
}
> update-task 6b9436ad-a159-46c8-86ba-b9bfd43d5cd1
title: 
description: 
due_date: 
status: 
Status: 400
{
  "error": "at least one field must be updated"
}
> update-task 6b9436ad-a159-46c8-86ba-b9bfd43d5cd1
title: 
description: 
due_date: 
status: in_progress
Status: 200
{
  "id": "6b9436ad-a159-46c8-86ba-b9bfd43d5cd1",
  "title": "",
  "description": "",
  "due_date": "",
  "status": "in_progress",
  "user_id": "9638fe95-5845-4011-902c-f3bcc4c821df"
}
> get-task 6b9436ad-a159-46c8-86ba-b9bfd43d5cd1
Status: 200
{
  "id": "6b9436ad-a159-46c8-86ba-b9bfd43d5cd1",
  "title": "New Task",
  "description": "Some Desc",
  "due_date": "2025-08-10T03:04:10Z",
  "status": "in_progress",
  "user_id": "9638fe95-5845-4011-902c-f3bcc4c821df"
}
> delete-task 6b9436ad-a159-46c8-86ba-b9bfd43d5cd1
Status: 200

> list-tasks
Status: 200
[
  {
    "id": "32c9f25a-bd7f-4c24-bdff-8c9b8a977ce7",
    "title": "New Title",
    "description": "Some Desc",
    "due_date": "2025-12-30T00:00:00Z",
    "status": "done",
    "user_id": "9638fe95-5845-4011-902c-f3bcc4c821df"
  }
]
> quit
```

## License

This project is licensed under the MIT License.