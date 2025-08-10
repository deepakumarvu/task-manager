package service

import (
	"task-planner/internal/db"
	"task-planner/internal/model"
)

type UserService struct {
	db db.DB
}

func NewUserService(db db.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Create(user *model.User) error {
	return s.db.CreateUser(user)
}

func (s *UserService) Get(id string) (*model.User, error) {
	return s.db.GetUser(id)
}

func (s *UserService) List() ([]model.User, error) {
	return s.db.ListUsers()
}

func (s *UserService) Delete(id string) error {
	return s.db.DeleteUser(id)
}
