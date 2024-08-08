package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"example/taskManager/console"
	"example/taskManager/controllers"
	"example/taskManager/router"
	"log"
)

// clearScreen clears the console screen.
func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	if err := c.Run(); err != nil {
		fmt.Println("Error clearing screen:", err)
	}
}

// prints a title with a specific color.
func printTitle(text string) {
	fmt.Printf("\033[1;34m%s\033[0m\n", text)
}

// prints an option with a specific color.
func printOption(text string) {
	fmt.Printf("\033[1;34m%s\033[0m\n", text)
}

func main() {

	clearScreen()
	printTitle("Task Management System")
	printOption("1. API Mode")
	printOption("2. Console Mode")

	// Read the user's choice
	reader := bufio.NewReader(os.Stdin)
	mode, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	mode = strings.TrimSpace(mode)

	// Handle the user's choice
	switch mode {
	case "1":
		fmt.Println("Starting API mode...")
		mongoURI := "mongodb://localhost:27017"
		taskController, err := controllers.NewTaskController(mongoURI)
		if err != nil {
			log.Fatalf("Failed to create task controller: %v", err)
		}

		r := router.SetupRouter(taskController)
		if err := r.Run(":9090"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}

	case "2":
		fmt.Println("Starting console mode...")
		console.StartConsoleApp()

	default:
		fmt.Println("Invalid mode. Please choose 1 or 2.")
	}
}
