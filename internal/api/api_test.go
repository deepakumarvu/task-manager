package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"task-planner/internal/api"
	"task-planner/internal/db"
	"task-planner/internal/model"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Suite")
}

var _ = Describe("Task Management API", func() {
	var router *gin.Engine
	var userID string
	var testDB db.DB

	BeforeEach(func() {
		testDB, _ = db.NewSQLiteDB(":memory:")
		router = gin.Default()
		api.RegisterRoutes(router, testDB)

		// Create a user for task tests
		user := model.User{Name: "Test User", Email: "test@example.com"}
		userJson, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJson))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		var createdUser model.User
		json.Unmarshal(w.Body.Bytes(), &createdUser)
		userID = createdUser.ID
	})

	Describe("User API", func() {
		It("should create a user", func() {
			user := model.User{Name: "Alice", Email: "alice@example.com"}
			userJson, _ := json.Marshal(user)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJson))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusCreated))
			Expect(w.Body.String()).To(ContainSubstring("Alice"))
		})

		It("should not create a user with missing fields", func() {
			user := model.User{Name: "", Email: ""}
			userJson, _ := json.Marshal(user)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJson))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should not create a user with too short name", func() {
			user := model.User{Name: "A", Email: "short@example.com"}
			userJson, _ := json.Marshal(user)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJson))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("name must be between 2 and 50 characters"))
		})

		It("should not create a user with too long name", func() {
			longName := ""
			for i := 0; i < 51; i++ {
				longName += "a"
			}
			user := model.User{Name: longName, Email: "long@example.com"}
			userJson, _ := json.Marshal(user)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJson))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("name must be between 2 and 50 characters"))
		})

		It("should get a user by ID", func() {
			req, _ := http.NewRequest("GET", "/users/"+userID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(ContainSubstring("Test User"))
		})

		It("should return 404 for non-existent user", func() {
			req, _ := http.NewRequest("GET", "/users/non-existent-id", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusNotFound))
		})

		It("should list users", func() {
			req, _ := http.NewRequest("GET", "/users", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(ContainSubstring("Test User"))
		})

		It("should delete a user", func() {
			req, _ := http.NewRequest("DELETE", "/users/"+userID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
		})

		It("should return error when deleting non-existent user", func() {
			req, _ := http.NewRequest("DELETE", "/users/non-existent-id", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Describe("Task API", func() {
		var taskID string

		It("should create a task", func() {
			task := model.Task{Title: "Test Task", Description: "desc", DueDate: "2025-12-31T10:00:00Z", Status: "pending"}
			jsonData, _ := json.Marshal(task)
			req, _ := http.NewRequest("POST", "/users/"+userID+"/tasks", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusCreated))
			Expect(w.Body.String()).To(ContainSubstring("Test Task"))
			var createdTask model.Task
			json.Unmarshal(w.Body.Bytes(), &createdTask)
			taskID = createdTask.ID
		})

		It("should not create a task with missing title", func() {
			task := model.Task{Title: "", Description: "desc"}
			jsonData, _ := json.Marshal(task)
			req, _ := http.NewRequest("POST", "/users/"+userID+"/tasks", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("should not create a task with invalid status", func() {
			task := model.Task{Title: "Test Task", Description: "desc", DueDate: "2025-12-31T10:00:00Z", Status: "invalid_status"}
			jsonData, _ := json.Marshal(task)
			req, _ := http.NewRequest("POST", "/users/"+userID+"/tasks", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("invalid status"))
		})

		It("should not create a task with too long description", func() {
			longDesc := ""
			for i := 0; i < 201; i++ {
				longDesc += "a"
			}
			task := model.Task{Title: "Test Task", Description: longDesc, DueDate: "2025-12-31T10:00:00Z", Status: "pending"}
			jsonData, _ := json.Marshal(task)
			req, _ := http.NewRequest("POST", "/users/"+userID+"/tasks", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("description must be at most 200 characters"))
		})

		It("should not create a task with invalid due_date", func() {
			task := model.Task{Title: "Test Task", Description: "desc", DueDate: "not-a-date", Status: "pending"}
			jsonData, _ := json.Marshal(task)
			req, _ := http.NewRequest("POST", "/users/"+userID+"/tasks", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("due_date must be ISO 8601 format"))
		})

		Context("with a created task", func() {
			BeforeEach(func() {
				task := model.Task{Title: "Test Task", Description: "desc", DueDate: "2025-09-02T15:04:05Z", Status: "pending"}
				jsonData, _ := json.Marshal(task)
				req, _ := http.NewRequest("POST", "/users/"+userID+"/tasks", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				var createdTask model.Task
				json.Unmarshal(w.Body.Bytes(), &createdTask)
				taskID = createdTask.ID
			})

			It("should list tasks for a user", func() {
				req, _ := http.NewRequest("GET", "/users/"+userID+"/tasks", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(ContainSubstring("Test Task"))
			})

			It("should list tasks with filter", func() {
				req, _ := http.NewRequest("GET", "/users/"+userID+"/tasks?status=pending", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(ContainSubstring("Test Task"))
			})

			It("should get a task by ID", func() {
				req, _ := http.NewRequest("GET", "/users/"+userID+"/tasks/"+taskID, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(ContainSubstring("Test Task"))
			})

			It("should return 404 for non-existent task", func() {
				req, _ := http.NewRequest("GET", "/users/"+userID+"/tasks/non-existent-id", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})

			It("should update a task", func() {
				updatedTask := model.Task{Title: "Updated Task", Description: "Updated desc", DueDate: "2025-08-09T15:04:05Z", Status: "done"}
				updatedJson, _ := json.Marshal(updatedTask)
				req, _ := http.NewRequest("PUT", "/users/"+userID+"/tasks/"+taskID, bytes.NewBuffer(updatedJson))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				fmt.Println(w.Body.String())
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(ContainSubstring("Updated Task"))
			})

			It("should not update a task with missing title", func() {
				updatedTask := model.Task{Title: "", Description: "Updated desc"}
				updatedJson, _ := json.Marshal(updatedTask)
				req, _ := http.NewRequest("PUT", "/users/"+userID+"/tasks/"+taskID, bytes.NewBuffer(updatedJson))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should not update a task with invalid status", func() {
				// Create a valid task first
				task := model.Task{Title: "Test Task", Description: "desc", DueDate: "2025-12-31T10:00:00Z", Status: "pending"}
				jsonData, _ := json.Marshal(task)
				req, _ := http.NewRequest("POST", "/users/"+userID+"/tasks", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				var createdTask model.Task
				json.Unmarshal(w.Body.Bytes(), &createdTask)

				// Try to update with invalid status
				updatedTask := model.Task{Title: "Updated Task", Description: "desc", DueDate: "2025-12-31T10:00:00Z", Status: "bad_status"}
				updatedJson, _ := json.Marshal(updatedTask)
				req, _ = http.NewRequest("PUT", "/users/"+userID+"/tasks/"+createdTask.ID, bytes.NewBuffer(updatedJson))
				req.Header.Set("Content-Type", "application/json")
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(ContainSubstring("invalid status"))
			})

			It("should not update a task with invalid due_date", func() {
				// Create a valid task first
				task := model.Task{Title: "Test Task", Description: "desc", DueDate: "2025-12-31T10:00:00Z", Status: "pending"}
				jsonData, _ := json.Marshal(task)
				req, _ := http.NewRequest("POST", "/users/"+userID+"/tasks", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				var createdTask model.Task
				json.Unmarshal(w.Body.Bytes(), &createdTask)

				// Try to update with invalid due_date
				updatedTask := model.Task{Title: "Updated Task", Description: "desc", DueDate: "not-a-date", Status: "pending"}
				updatedJson, _ := json.Marshal(updatedTask)
				req, _ = http.NewRequest("PUT", "/users/"+userID+"/tasks/"+createdTask.ID, bytes.NewBuffer(updatedJson))
				req.Header.Set("Content-Type", "application/json")
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(ContainSubstring("due_date must be ISO 8601 format"))
			})

			It("should delete a task", func() {
				req, _ := http.NewRequest("DELETE", "/users/"+userID+"/tasks/"+taskID, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
			})

			It("should return error when deleting non-existent task", func() {
				req, _ := http.NewRequest("DELETE", "/users/"+userID+"/tasks/non-existent-id", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
