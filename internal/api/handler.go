package api

import (
	"net/http"
	"strings"
	"task-planner/internal/db"
	"task-planner/internal/model"
	"task-planner/internal/service"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the API routes for user and task management
func RegisterRoutes(router *gin.Engine, dbInstance db.DB) {
	userService := service.NewUserService(dbInstance)

	// User routes
	router.POST("/users", createUserHandler(userService))
	router.GET("/users/:user_id", getUserHandler(userService))
	router.GET("/users", listUsersHandler(userService))
	router.DELETE("/users/:user_id", deleteUserHandler(userService))

	// Task routes (under user context)
	router.POST("/users/:user_id/tasks", taskHandler(dbInstance, createTask))
	router.GET("/users/:user_id/tasks", taskHandler(dbInstance, listTasks))
	router.GET("/users/:user_id/tasks/:task_id", taskHandler(dbInstance, getTask))
	router.PUT("/users/:user_id/tasks/:task_id", taskHandler(dbInstance, updateTask))
	router.DELETE("/users/:user_id/tasks/:task_id", taskHandler(dbInstance, deleteTask))
}

var validStatuses = map[string]struct{}{
	"pending":     {},
	"in_progress": {},
	"done":        {},
}

const (
	minNameLen = 2
	maxNameLen = 50
	maxDescLen = 200
)

func isValidStatus(status string) bool {
	_, ok := validStatuses[status]
	return ok
}

func isValidISO8601(dateStr string) bool {
	_, err := time.Parse(time.RFC3339, dateStr)
	return err == nil
}

func getParam(c *gin.Context, param string) (string, bool) {
	value := c.Param(param)
	if value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": param + " is required"})
		return "", false
	}
	return value, true
}

// --- User Handlers ---
func createUserHandler(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		nameLen := utf8.RuneCountInString(strings.TrimSpace(user.Name))
		if nameLen < minNameLen || nameLen > maxNameLen {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name must be between 2 and 50 characters"})
			return
		}
		if strings.TrimSpace(user.Email) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
			return
		}
		err := userService.Create(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}

func getUserHandler(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := getParam(c, "user_id")
		if !ok {
			return
		}
		user, err := userService.Get(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func listUsersHandler(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := userService.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func deleteUserHandler(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := getParam(c, "user_id")
		if !ok {
			return
		}
		err := userService.Delete(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	}
}

// --- Task Handler Wrapper ---
type taskAction func(c *gin.Context, taskService *service.TaskService)

func taskHandler(dbInstance db.DB, action taskAction) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := getParam(c, "user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}
		taskService := service.NewTaskService(dbInstance, userID)
		action(c, taskService)
	}
}

// --- Task Actions ---
func createTask(c *gin.Context, taskService *service.TaskService) {
	var task model.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	titleLen := utf8.RuneCountInString(strings.TrimSpace(task.Title))
	if titleLen < minNameLen || titleLen > maxNameLen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title must be between 2 and 50 characters"})
		return
	}
	descLen := utf8.RuneCountInString(strings.TrimSpace(task.Description))
	if descLen > maxDescLen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "description must be at most 200 characters"})
		return
	}
	if !isValidStatus(task.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	if !isValidISO8601(task.DueDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "due_date must be ISO 8601 format (RFC3339)"})
		return
	}
	err := taskService.Create(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func listTasks(c *gin.Context, taskService *service.TaskService) {
	status := c.Query("status")
	tasks, err := taskService.List(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func getTask(c *gin.Context, taskService *service.TaskService) {
	taskID, ok := getParam(c, "task_id")
	if !ok {
		return
	}
	task, err := taskService.Get(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func updateTask(c *gin.Context, taskService *service.TaskService) {
	taskID, ok := getParam(c, "task_id")
	if !ok {
		return
	}
	var task model.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	titleLen := utf8.RuneCountInString(strings.TrimSpace(task.Title))
	if titleLen < minNameLen || titleLen > maxNameLen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title must be between 2 and 50 characters"})
		return
	}
	descLen := utf8.RuneCountInString(strings.TrimSpace(task.Description))
	if descLen > maxDescLen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "description must be at most 200 characters"})
		return
	}
	if !isValidStatus(task.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	if !isValidISO8601(task.DueDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "due_date must be ISO 8601 format (RFC3339)"})
		return
	}
	task.ID = taskID
	err := taskService.Update(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context, taskService *service.TaskService) {
	taskID, ok := getParam(c, "task_id")
	if !ok {
		return
	}
	if err := taskService.Delete(taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
