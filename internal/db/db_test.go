package db_test

import (
	"task-planner/internal/db"
	"task-planner/internal/model"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDB(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Suite")
}

var _ = Describe("DB CRUD Operations", func() {
	var testDB db.DB
	var testUser *model.User

	BeforeEach(func() {
		var err error
		testDB, err = db.NewSQLiteDB(":memory:")
		Expect(err).To(BeNil())

		testUser = &model.User{Name: "Test User", Email: "test@example.com"}
		err = testDB.CreateUser(testUser)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		if testDB != nil {
			testDB.Close()
		}
	})

	Describe("User CRUD", func() {
		It("should create a user", func() {
			user := &model.User{Name: "Alice", Email: "alice@example.com"}
			err := testDB.CreateUser(user)
			Expect(err).To(BeNil())
			Expect(user.ID).NotTo(BeEmpty())
		})

		It("should get a user", func() {
			got, err := testDB.GetUser(testUser.ID)
			Expect(err).To(BeNil())
			Expect(got.Name).To(Equal("Test User"))
			Expect(got.Email).To(Equal("test@example.com"))
		})

		It("should list users", func() {
			users, err := testDB.ListUsers()
			Expect(err).To(BeNil())
			Expect(users).NotTo(BeEmpty())
		})

		It("should delete a user", func() {
			err := testDB.DeleteUser(testUser.ID)
			Expect(err).To(BeNil())
			users, err := testDB.ListUsers()
			Expect(err).To(BeNil())
			found := false
			for _, u := range users {
				if u.ID == testUser.ID {
					found = true
				}
			}
			Expect(found).To(BeFalse())
		})
	})

	Describe("Task CRUD", func() {
		var task *model.Task

		BeforeEach(func() {
			task = &model.Task{
				Title:       "Test Task",
				Description: "This is a test task",
				DueDate:     "2023-12-31T10:00:00Z",
				Status:      "pending",
				UserID:      testUser.ID,
			}
			err := testDB.CreateTask(task)
			Expect(err).To(BeNil())
		})

		It("should create a task", func() {
			Expect(task.ID).NotTo(BeEmpty())
		})

		It("should list tasks", func() {
			tasks, err := testDB.ListTasks(testUser.ID, "")
			Expect(err).To(BeNil())
			Expect(tasks).To(HaveLen(1))
		})

		It("should get a task by ID", func() {
			got, err := testDB.GetTask(task.ID)
			Expect(err).To(BeNil())
			Expect(got.Title).To(Equal("Test Task"))
			Expect(got.Description).To(Equal("This is a test task"))
			Expect(got.DueDate).To(Equal("2023-12-31T10:00:00Z"))
			Expect(got.Status).To(Equal("pending"))
			Expect(got.UserID).To(Equal(testUser.ID))
		})

		It("should get task based on status filter", func() {
			tasks, err := testDB.ListTasks(testUser.ID, "pending")
			Expect(err).To(BeNil())
			Expect(tasks).To(HaveLen(1))
			Expect(tasks[0].Title).To(Equal("Test Task"))
		})

		It("should update a task", func() {
			task.Title = "Updated Task"
			err := testDB.UpdateTask(task)
			Expect(err).To(BeNil())

			updatedTask, err := testDB.GetTask(task.ID)
			Expect(err).To(BeNil())
			Expect(updatedTask.Title).To(Equal("Updated Task"))
		})

		It("should delete a task", func() {
			err := testDB.DeleteTask(task.ID, testUser.ID)
			Expect(err).To(BeNil())

			tasks, err := testDB.ListTasks(testUser.ID, "")
			Expect(err).To(BeNil())
			Expect(tasks).To(HaveLen(0))
		})

		It("should not delete a task with wrong user ID", func() {
			err := testDB.DeleteTask(task.ID, "wrong_user_id")
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("task not found"))
		})

		It("should not get a non-existent task", func() {
			_, err := testDB.GetTask("non_existent_id")
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("task not found"))
		})

		It("should not update a non-existent task", func() {
			task.ID = "non_existent_id"
			err := testDB.UpdateTask(task)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("task not found"))
		})
	})
})
