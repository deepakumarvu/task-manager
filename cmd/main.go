package main

import (
	"log"

	"task-planner/internal/api"
	"task-planner/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database connection
	database, err := db.NewSQLiteDB("tasks.db")
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer database.Close()

	// Set up Gin router and register routes
	router := gin.Default()
	api.RegisterRoutes(router, database)

	// Start the server
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
