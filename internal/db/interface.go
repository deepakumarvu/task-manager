package db

import "task-planner/internal/model"

type DB interface {
	// User methods
	CreateUser(user *model.User) error
	GetUser(id string) (*model.User, error)
	ListUsers() ([]model.User, error)
	DeleteUser(id string) error

	// Task methods (always under user context)
	CreateTask(task *model.Task) error
	GetTask(id string) (*model.Task, error)
	ListTasks(userID string, status string) ([]model.Task, error)
	UpdateTask(task *model.Task) error
	DeleteTask(id string, userID string) error

	Close() error
}

var _ DB = (*SQLiteDB)(nil)
