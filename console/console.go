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

var (
	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Bold(true)

	taskListViewStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("12")).
				Bold(true)
	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12"))

	tableHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).Background(lipgloss.Color("17"))
	tableCellStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Background(lipgloss.Color("18"))
	selectedRowStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("19"))
)

// StartConsoleApp starts the console application
func StartConsoleApp() {
	for {
		fmt.Println(taskListViewStyle.Render("Task List"))
		fmt.Println(promptStyle.Render("1. Add Task."))
		fmt.Println(promptStyle.Render("2. View all Tasks."))
		fmt.Println(promptStyle.Render("3. Get Task by ID."))
		fmt.Println(promptStyle.Render("4. Update Task."))
		fmt.Println(promptStyle.Render("5. Remove Task."))
		fmt.Println(promptStyle.Render("6. Mark Task as Complete."))
		fmt.Println(promptStyle.Render("7. Exit"))
		fmt.Print(promptStyle.Render("Enter your choice: "))
		choice, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		choice = strings.TrimSpace(choice)

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

// Load tasks from JSON file
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

// Save tasks to JSON file
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

// Validate input with provided validation function
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

// Validation if left empty
func validateNonEmpty(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("input cannot be empty")
	}
	return input, nil
}

// validate title
func validateTitle(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("title cannot be empty")
	}
	return input, nil
}

// valdte date
func validateDate(input string) (string, error) {
	_, err := time.Parse("2006-01-02", input)
	if err != nil {
		return "", fmt.Errorf("invalid date format. Please use YYYY-MM-DD")
	}
	return input, nil
}

// validate status as("completed", "inprogress", "started")
func validateStatus(input string) (string, error) {
	validStatuses := []string{"completed", "inprogress", "started"}
	for _, status := range validStatuses {
		if input == status {
			return input, nil
		}
	}
	return "", fmt.Errorf("invalid status. Valid options are: completed, inprogress, started")
}

// validates to input letters only for the title of the task
func validateLettersOnly(input string) (string, error) {
	for _, char := range input {
		if !unicode.IsLetter(char) {
			return "", fmt.Errorf("title can only contain letters")
		}
	}
	return input, nil
}

// View all tasks with table display
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
			task.DueDate.Format("2006-01-02"), // Ensure date format is consistent
			task.Status,
		}
		rows = append(rows, row)
	}

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

	p := tea.NewProgram(model{t})

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

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

func (m model) View() string {
	return taskListViewStyle.Render("Task List") + "\n" + m.table.View() + "\nPress 'q' to return to the main menu."
}

// Add new task
func addTask() {
	title := validateInput("Enter task title: ", validateTitle)
	description := validateInput("Enter task description: ", validateNonEmpty)
	dueDateStr := validateInput("Enter task due date (YYYY-MM-DD): ", validateDate)
	dueDate, err := time.Parse("2006-01-02", dueDateStr)
	if err != nil {
		fmt.Printf(errorStyle.Render("Invalid date format. Please use YYYY-MM-DD: %v\n"), err)
		return
	}
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

// Get task by ID
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
	fmt.Printf("Title: %s\nDescription: %s\nDue Date: %s\nStatus: %s\n",
		task.Title, task.Description, task.DueDate.Format("2006-01-02"), task.Status)
}

// Update existing task
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
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promptStyle.Render("Enter task ID to update: "))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	taskNum, err := strconv.Atoi(input)
	if err != nil || taskNum < 1 || taskNum > len(tasks) {
		fmt.Println(errorStyle.Render("Invalid task ID."))
		return
	}

	task := &tasks[taskNum-1]
	fmt.Println("Leave empty to keep current value.")
	title := validateInput(fmt.Sprintf("Enter new title (current: %s): ", task.Title), validateTitle)
	if title != "" {
		task.Title = title
	}
	description := validateInput(fmt.Sprintf("Enter new description (current: %s): ", task.Description), validateNonEmpty)
	if description != "" {
		task.Description = description
	}
	dueDateStr := validateInput(fmt.Sprintf("Enter new due date (current: %s): ", task.DueDate.Format("2006-01-02")), validateDate)
	if dueDateStr != "" {
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			fmt.Printf(errorStyle.Render("Invalid date format. Please use YYYY-MM-DD: %v\n"), err)
			return
		}
		task.DueDate = dueDate
	}
	status := validateInput(fmt.Sprintf("Enter new status (current: %s): ", task.Status), validateStatus)
	if status != "" {
		task.Status = status
	}

	err = saveTasks(tasks)
	if err != nil {
		fmt.Printf(errorStyle.Render("Error saving tasks: %v\n"), err)
		return
	}

	fmt.Println(selectedStyle.Render("Task updated successfully."))
}

// Remove a task
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
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promptStyle.Render("Enter task ID to remove: "))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	taskNum, err := strconv.Atoi(input)
	if err != nil || taskNum < 1 || taskNum > len(tasks) {
		fmt.Println(errorStyle.Render("Invalid task ID."))
		return
	}

	tasks = append(tasks[:taskNum-1], tasks[taskNum:]...)
	err = saveTasks(tasks)
	if err != nil {
		fmt.Printf(errorStyle.Render("Error saving tasks: %v\n"), err)
		return
	}

	fmt.Println(selectedStyle.Render("Task removed successfully."))
}

// Mark a task as complete
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
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promptStyle.Render("Enter task ID to mark as complete: "))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	taskNum, err := strconv.Atoi(input)
	if err != nil || taskNum < 1 || taskNum > len(tasks) {
		fmt.Println(errorStyle.Render("Invalid task ID."))
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
