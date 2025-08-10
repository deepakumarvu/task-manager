package service

import (
	"task-planner/internal/db"
	"task-planner/internal/model"
)

type TaskService struct {
	db     db.DB
	userID string
}

func NewTaskService(db db.DB, userID string) *TaskService {
	return &TaskService{db: db, userID: userID}
}

func (s *TaskService) Create(task *model.Task) error {
	task.UserID = s.userID
	return s.db.CreateTask(task)
}

func (s *TaskService) List(status string) ([]model.Task, error) {
	return s.db.ListTasks(s.userID, status)
}

func (s *TaskService) Get(taskID string) (*model.Task, error) {
	task, err := s.db.GetTask(taskID)
	if err != nil {
		return nil, err
	}
	if task.UserID != s.userID {
		return nil, db.ErrTaskNotFound
	}
	return task, nil
}

func (s *TaskService) Update(task *model.Task) error {
	task.UserID = s.userID
	return s.db.UpdateTask(task)
}

func (s *TaskService) Delete(taskID string) error {
	return s.db.DeleteTask(taskID, s.userID)
}
