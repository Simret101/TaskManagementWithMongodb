# Task Management API Documentation

This documentation provides an overview of the Task Management REST API implemented using Go, the Gin framework, and MongoDB. The API supports basic CRUD operations for managing tasks.

---

## Endpoints

### GET /tasks

Retrieves a list of all tasks.

- **Response Codes:**
  - `200 OK`: Successfully retrieved the list of tasks.

- **Response Body:**
  ```json
  [
    {
      "id": "60d5f67f5b50793b585b3567",
      "title": "Task Title",
      "description": "Task Description",
      "duedate": "2024-07-09T10:00:00Z",
      "status": "pending"
    }
  ]
  ```

### GET /tasks/:id

Retrieves a specific task by its ID.

- **Path Parameter:**
  - `id` (string): The ID of the task to retrieve.

- **Response Codes:**
  - `200 OK`: Successfully retrieved the task.
  - `404 Not Found`: Task not found.

- **Response Body:**
  ```json
  {
    "id": "60d5f67f5b50793b585b3567",
    "title": "Task Title",
    "description": "Task Description",
    "duedate": "2024-07-09T10:00:00Z",
    "status": "pending"
  }
  ```

### POST /tasks

Creates a new task.

- **Request Body:**
  ```json
  {
    "title": "Task Title",
    "description": "Task Description",
    "duedate": "2024-07-09T10:00:00Z",
    "status": "pending"
  }
  ```

- **Response Codes:**
  - `201 Created`: Task created successfully.
  - `400 Bad Request`: Invalid request body.

- **Response Body:**
  ```json
  {
    "id": "60d5f67f5b50793b585b3567",
    "title": "Task Title",
    "description": "Task Description",
    "duedate": "2024-07-09T10:00:00Z",
    "status": "pending"
  }
  ```

### PUT /tasks/:id

Updates a specific task by its ID.

- **Path Parameter:**
  - `id` (string): The ID of the task to update.

- **Request Body:**
  ```json
  {
    "title": "Updated Title",
    "description": "Updated Description",
    "duedate": "2024-07-10T10:00:00Z",
    "status": "completed"
  }
  ```

- **Response Codes:**
  - `200 OK`: Task updated successfully.
  - `400 Bad Request`: Invalid request body.
  - `404 Not Found`: Task not found.

- **Response Body:**
  ```json
  {
    "message": "task updated successfully"
  }
  ```

### DELETE /tasks/:id

Deletes a specific task by its ID.

- **Path Parameter:**
  - `id` (string): The ID of the task to delete.

- **Response Codes:**
  - `204 No Content`: Task deleted successfully.
  - `404 Not Found`: Task not found.

---

## Models

### Task

Represents a task in the task management system.

- **Fields:**
  - `id` (primitive.ObjectID): Unique identifier for the task.
  - `title` (string): Title of the task.
  - `description` (string): Description of the task.
  - `duedate` (time.Time): Due date of the task.
  - `status` (string): Current status of the task.

---

## Code Overview

### `main.go`

Sets up and starts the server.



- **Function:** `SetupRouter(taskController *controllers.TaskController) *gin.Engine`
  - Configures routes for task management.
  - Maps HTTP methods and endpoints to controller functions.

### `models/task.go`

Defines the `Task` model.


- **Struct:** `Task`
  - Represents a task with fields for ID, title, description, due date, and status.

### `data/task_service.go`

Implements data access logic for tasks, interacting with the MongoDB database.



- **Function:** `NewTaskService(mongoURI string) (*TaskService, error)`
  - Connects to MongoDB and initializes the task collection.
- **Function:** `GetTasks() ([]*models.Task, error)`
  - Retrieves all tasks from the collection.
- **Function:** `GetTask(id string) (*models.Task, error)`
  - Retrieves a task by its ID.
- **Function:** `CreateTask(task *models.Task) (*models.Task, error)`
  - Inserts a new task into the collection.
- **Function:** `UpdateTask(id string, updatedTask *models.Task) error`
  - Updates an existing task.
- **Function:** `DeleteTask(id string) error`
  - Deletes a task by its ID.

### `controllers/task_controller.go`

Defines controller functions for handling task requests (GET, POST, PUT, DELETE) on the `/tasks` endpoint.



- **Function:** `NewTaskController(mongoURI string) (*TaskController, error)`
  - Initializes the task controller with a MongoDB URI.
- **Function:** `GetTasks(c *gin.Context)`
  - Handles GET requests to retrieve all tasks.
- **Function:** `GetTask(c *gin.Context)`
  - Handles GET requests to retrieve a specific task by ID.
- **Function:** `CreateTask(c *gin.Context)`
  - Handles POST requests to create a new task.
- **Function:** `UpdateTask(c *gin.Context)`
  - Handles PUT requests to update an existing task by ID.
- **Function:** `DeleteTask(c *gin.Context)`
  - Handles DELETE requests to delete a specific task by ID.

---

### API DOCUMENTION: https://documenter.getpostman.com/view/37289771/2sA3rwNZnx