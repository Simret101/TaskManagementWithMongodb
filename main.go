package main

import (
	"example/taskManager/controllers"
	"example/taskManager/router"
	"log"
)

func main() {
	mongoURI := "mongodb://localhost:27017"
	taskController, err := controllers.NewTaskController(mongoURI)
	if err != nil {
		log.Fatalf("Failed to create task controller: %v", err)
	}

	r := router.SetupRouter(taskController)
	r.Run()
}
