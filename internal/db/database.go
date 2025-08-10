package db

import (
	"database/sql"
	"fmt"

	"task-planner/internal/model"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDB wraps a sql.DB connection for database operations
type SQLiteDB struct {
	conn *sql.DB
}

// NewSQLiteDB initializes a new SQLiteDB instance
func NewSQLiteDB(dataSourceName string) (DB, error) {
	conn, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	// Create tables if they don't exist
	tables := []struct {
		name       string
		createStmt string
	}{
		{
			name: "users",
			createStmt: `
				CREATE TABLE IF NOT EXISTS users (
					id TEXT PRIMARY KEY,
					name TEXT NOT NULL,
					email TEXT NOT NULL UNIQUE
				)
			`,
		},
		{
			name: "tasks",
			createStmt: `
				CREATE TABLE IF NOT EXISTS tasks (
					id TEXT PRIMARY KEY,
					title TEXT NOT NULL,
					description TEXT,
					due_date TEXT,
					status TEXT NOT NULL DEFAULT 'pending',
					user_id TEXT NOT NULL,
					FOREIGN KEY(user_id) REFERENCES users(id)
				)
			`,
		},
	}

	for _, table := range tables {
		_, err = conn.Exec(table.createStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating table %s: %w", table.name, err)
		}
	}

	return &SQLiteDB{conn: conn}, nil
}

// Close closes the SQLiteDB connection
func (s *SQLiteDB) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// Task methods

func (s *SQLiteDB) CreateTask(task *model.Task) error {
	task.ID = uuid.New().String()
	_, err := s.conn.Exec(
		"INSERT INTO tasks (id, title, description, due_date, status, user_id) VALUES (?, ?, ?, ?, ?, ?)",
		task.ID, task.Title, task.Description, task.DueDate, task.Status, task.UserID,
	)
	return err
}

func (s *SQLiteDB) GetTask(id string) (*model.Task, error) {
	row := s.conn.QueryRow("SELECT id, title, description, due_date, status, user_id FROM tasks WHERE id = ?", id)
	var task model.Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status, &task.UserID)
	if err == sql.ErrNoRows {
		return nil, ErrTaskNotFound
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *SQLiteDB) ListTasks(userID string, status string) ([]model.Task, error) {
	var rows *sql.Rows
	var err error
	if status == "" {
		rows, err = s.conn.Query("SELECT id, title, description, due_date, status, user_id FROM tasks WHERE user_id = ?", userID)
	} else {
		rows, err = s.conn.Query("SELECT id, title, description, due_date, status, user_id FROM tasks WHERE user_id = ? AND status = ?", userID, status)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.DueDate, &t.Status, &t.UserID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *SQLiteDB) UpdateTask(task *model.Task) error {
	query := "UPDATE tasks SET "
	updates := []string{}
	args := []any{}
	if task.Title != "" {
		updates = append(updates, "title = ?")
		args = append(args, task.Title)
	}
	if task.Description != "" {
		updates = append(updates, "description = ?")
		args = append(args, task.Description)
	}
	if task.DueDate != "" {
		updates = append(updates, "due_date = ?")
		args = append(args, task.DueDate)
	}
	if task.Status != "" {
		updates = append(updates, "status = ?")
		args = append(args, task.Status)
	}
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}
	query += updates[0]
	for _, update := range updates[1:] {
		query += ", " + update
	}
	query += " WHERE id = ? AND user_id = ?"
	fmt.Println("Executing query:", query)
	args = append(args, task.ID, task.UserID)
	res, err := s.conn.Exec(
		query,
		args...,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (s *SQLiteDB) DeleteTask(id string, userID string) error {
	res, err := s.conn.Exec("DELETE FROM tasks WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// User methods

func (s *SQLiteDB) CreateUser(user *model.User) error {
	user.ID = uuid.New().String()
	_, err := s.conn.Exec(
		"INSERT INTO users (id, name, email) VALUES (?, ?, ?)",
		user.ID, user.Name, user.Email,
	)
	return err
}

func (s *SQLiteDB) GetUser(id string) (*model.User, error) {
	row := s.conn.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)
	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *SQLiteDB) ListUsers() ([]model.User, error) {
	rows, err := s.conn.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *SQLiteDB) DeleteUser(id string) error {
	res, err := s.conn.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return err
}
