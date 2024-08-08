package console

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"example/taskManager/models"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const taskFile = "task.json"

// Styling using lipgloss for enhanced console output
var (
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	descriptionStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	statusStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	dateStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	promptStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	selectedStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("50"))
	normalStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	taskViewStyle     = lipgloss.NewStyle().Margin(1, 2)
	taskListViewStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).Padding(1, 2).Margin(1, 0)
	tableHeaderStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).Background(lipgloss.Color("17"))
	tableCellStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Background(lipgloss.Color("18"))
	selectedRowStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("19"))
)

// StartConsoleApp initializes and runs the console application
func StartConsoleApp() {
	for {
		// Display menu options
		fmt.Println(taskListViewStyle.Render("Task List"))
		fmt.Println(promptStyle.Render("1. Add Task."))
		fmt.Println(promptStyle.Render("2. View all Tasks."))
		fmt.Println(promptStyle.Render("3. Get Task by ID."))
		fmt.Println(promptStyle.Render("4. Update Task."))
		fmt.Println(promptStyle.Render("5. Remove Task."))
		fmt.Println(promptStyle.Render("6. Mark Task as Complete."))
		fmt.Println(promptStyle.Render("7. Exit"))
		fmt.Println(promptStyle.Render("Enter your choice:"))

		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		// Handle user's choice
		switch choice {
		case "1":
			addTask()
		case "2":
			viewTasks()
		case "3":
			getTaskByID()
		case "4":
			updateTask()
		case "5":
			removeTask()
		case "6":
			markComplete()
		case "7":
			fmt.Println(promptStyle.Render("Exiting..."))
			os.Exit(0)
		default:
			fmt.Println(errorStyle.Render("Invalid Choice! Please try again. Enter a number from 1 to 7."))
		}
	}
}

// loadTasks reads the tasks from the file and returns them
func loadTasks() ([]models.Task, error) {
	if _, err := os.Stat(taskFile); os.IsNotExist(err) {
		return []models.Task{}, nil
	}
	data, err := os.ReadFile(taskFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}
	var tasks []models.Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks from JSON: %w", err)
	}
	return tasks, nil
}

// saveTasks writes the tasks to the file
func saveTasks(tasks []models.Task) error {
	data, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("failed to marshal tasks to JSON: %w", err)
	}
	err = os.WriteFile(taskFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write tasks to file: %w", err)
	}
	return nil
}

// validateInput prompts the user for input and validates it using the provided function
func validateInput(prompt string, validateFunc func(string) (string, error)) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(promptStyle.Render(prompt))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		validInput, err := validateFunc(input)
		if err != nil {
			fmt.Println(errorStyle.Render("Error: " + err.Error()))
			continue
		}
		return validInput
	}
}

// validateNonEmpty ensures that the input is not empty
func validateNonEmpty(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("input cannot be empty")
	}
	return input, nil
}

// validateDate ensures that the input matches the "YYYY-MM-DD" date format
func validateDate(input string) (string, error) {
	_, err := time.Parse("2006-01-02", input)
	if err != nil {
		return "", fmt.Errorf("invalid date format. Please use YYYY-MM-DD")
	}
	return input, nil
}

// validateStatus ensures that the input is one of the valid task statuses
func validateStatus(input string) (string, error) {
	validStatuses := []string{"completed", "inprogress", "started"}
	for _, status := range validStatuses {
		if input == status {
			return input, nil
		}
	}
	return "", fmt.Errorf("invalid status. Valid options are: completed, inprogress, started")
}

// validateLettersOnly ensures that the title contains only letters and spaces, and is not empty
func validateLettersOnly(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("title cannot be empty")
	}
	for _, char := range input {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
			return "", fmt.Errorf("title can only contain letters and spaces")
		}
	}
	return input, nil
}

// viewTasks displays the list of tasks in a table format
func viewTasks() {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf(errorStyle.Render("Error loading tasks: %v\n"), err)
		return
	}
	if len(tasks) == 0 {
		fmt.Println(normalStyle.Render("No tasks found."))
		return
	}

	// Define columns for the task table
	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Title", Width: 20},
		{Title: "Description", Width: 30},
		{Title: "Due Date", Width: 15},
		{Title: "Status", Width: 10},
	}

	var rows []table.Row
	for i, task := range tasks {
		row := table.Row{
			fmt.Sprintf("%d", i+1),
			task.Title,
			task.Description,
			task.DueDate,
			task.Status,
		}
		rows = append(rows, row)
	}

	// Create and configure the table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	t.SetStyles(table.Styles{
		Header:   tableHeaderStyle,
		Cell:     tableCellStyle,
		Selected: tableCellStyle.Copy().Background(lipgloss.Color("25")),
	})

	// Start the tea program to display the table
	p := tea.NewProgram(model{t})

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// Tea model for table view
type model struct {
	table table.Model
}

// Init initializes the tea model
func (m model) Init() tea.Cmd {
	return nil
}

// Update updates the tea model based on user input
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View renders the tea model view
func (m model) View() string {
	return taskListViewStyle.Render("Task List") + "\n" + m.table.View() + "\nPress 'q' to return to the main menu."
}

// addTask allows the user to add a new task
func addTask() {
	title := validateInput("Enter task title: ", validateLettersOnly)
	description := validateInput("Enter task description: ", validateNonEmpty)
	dueDate := validateInput("Enter task due date (YYYY-MM-DD): ", validateDate)
	status := validateInput("Enter task status (completed, inprogress, started): ", validateStatus)

	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf(errorStyle.Render("Error loading tasks: %v\n"), err)
		return
	}
	task := models.Task{
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Status:      status,
	}
	tasks = append(tasks, task)
	err = saveTasks(tasks)
	if err != nil {
		fmt.Printf(errorStyle.Render("Error saving tasks: %v\n"), err)
		return
	}
	fmt.Println(selectedStyle.Render("Task added successfully."))
}

// getTaskByID retrieves and displays a task based on its ID
func getTaskByID() {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf(errorStyle.Render("Error loading tasks: %v\n"), err)
		return
	}
	if len(tasks) == 0 {
		fmt.Println(normalStyle.Render("No tasks found."))
		return
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promptStyle.Render("Enter task ID: "))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	taskNum, err := strconv.Atoi(input)
	if err != nil || taskNum < 1 || taskNum > len(tasks) {
		fmt.Println(errorStyle.Render("Invalid task ID."))
		return
	}
	task := tasks[taskNum-1]
	fmt.Println("Task Details:")
	fmt.Printf("Title: %s\n", task.Title)
	fmt.Printf("Description: %s\n", task.Description)
	fmt.Printf("Due Date: %s\n", task.DueDate)
	fmt.Printf("Status: %s\n", task.Status)
}

// updateTask allows the user to update an existing task
func updateTask() {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf(errorStyle.Render("Error loading tasks: %v\n"), err)
		return
	}
	if len(tasks) == 0 {
		fmt.Println(normalStyle.Render("No tasks found."))
		return
	}
	viewTasks()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promptStyle.Render("Enter task number to update: "))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	taskNum, err := strconv.Atoi(input)
	if err != nil || taskNum < 1 || taskNum > len(tasks) {
		fmt.Println(errorStyle.Render("Invalid task number."))
		return
	}
	task := &tasks[taskNum-1]
	task.Title = validateInput(fmt.Sprintf("Enter new title (current: %s): ", task.Title), validateLettersOnly)
	task.Description = validateInput(fmt.Sprintf("Enter new description (current: %s): ", task.Description), validateNonEmpty)
	task.DueDate = validateInput(fmt.Sprintf("Enter new due date (current: %s): ", task.DueDate), validateDate)
	task.Status = validateInput(fmt.Sprintf("Enter new status (current: %s): ", task.Status), validateStatus)
	err = saveTasks(tasks)
	if err != nil {
		fmt.Printf(errorStyle.Render("Error saving tasks: %v\n"), err)
		return
	}
	fmt.Println(selectedStyle.Render("Task updated successfully."))
}

// removeTask allows the user to remove a task based on its number
func removeTask() {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf(errorStyle.Render("Error loading tasks: %v\n"), err)
		return
	}
	if len(tasks) == 0 {
		fmt.Println(normalStyle.Render("No tasks found."))
		return
	}
	viewTasks()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promptStyle.Render("Enter task number to remove: "))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	taskNum, err := strconv.Atoi(input)
	if err != nil || taskNum < 1 || taskNum > len(tasks) {
		fmt.Println(errorStyle.Render("Invalid task number."))
		return
	}
	tasks = append(tasks[:taskNum-1], tasks[taskNum:]...)
	err = saveTasks(tasks)
	if err != nil {
		fmt.Printf(errorStyle.Render("Error saving tasks: %v\n"), err)
		return
	}
	fmt.Println(selectedStyle.Render("Task deleted successfully."))
}

// markComplete allows the user to mark a task as complete
func markComplete() {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf(errorStyle.Render("Error loading tasks: %v\n"), err)
		return
	}
	if len(tasks) == 0 {
		fmt.Println(normalStyle.Render("No tasks found."))
		return
	}
	viewTasks()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promptStyle.Render("Enter task number to mark as complete: "))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	taskNum, err := strconv.Atoi(input)
	if err != nil || taskNum < 1 || taskNum > len(tasks) {
		fmt.Println(errorStyle.Render("Invalid task number."))
		return
	}
	tasks[taskNum-1].Status = "completed"
	err = saveTasks(tasks)
	if err != nil {
		fmt.Printf(errorStyle.Render("Error saving tasks: %v\n"), err)
		return
	}
	fmt.Println(selectedStyle.Render("Task marked as complete."))
}

