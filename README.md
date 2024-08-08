# Task Manager Console Application

A Go-based task management application that interfaces with MongoDB. This console application allows you to manage tasks using a REST API built with Gin. Tasks are stored and managed in MongoDB.

## Features

- **CRUD Operations**: Create, Read, Update, and Delete tasks.
- **Persistent Storage**: Tasks are stored in a MongoDB database.
- **Error Handling**: Consistent error responses and handling.

## Prerequisites

- Go 1.18 or higher
- MongoDB (local or remote instance)
- `go.mongodb.org/mongo-driver` package for MongoDB integration
- `github.com/gin-gonic/gin` package for RESTful routing

## Installation

1. **Clone the repository:**

    ```sh
    git clone [https://github.com/Simret101/TaskManagementWithMongodb]
    ```

2. **Navigate to the project directory:**

    ```sh
    cd task-manager-console
    ```

3. **Install dependencies:**

    ```sh
    go mod tidy
    ```

4. **Build the application:**

    ```sh
    go build -o task-manager
    ```

## Configuration

- **MongoDB URI**: Update the MongoDB URI in `main.go` to point to your MongoDB instance.

    ```go
    mongoURI := "mongodb://localhost:27017"
    ```

## Usage

1. **Run the application:**

    ```sh
    ./task-manager
    ```

2. **API Endpoints:**

    - **GET `/tasks`**: Retrieve all tasks.
    - **GET `/tasks/:id`**: Retrieve a task by ID.
    - **POST `/tasks`**: Create a new task. Requires JSON payload with `title`, `description`, `duedate`, and `status`.
    - **PUT `/tasks/:id`**: Update an existing task by ID. Requires JSON payload with updated task details.
    - **DELETE `/tasks/:id`**: Delete a task by ID.

## Example

### Adding a Task

```sh
curl -X POST http://localhost:9090/tasks -H "Content-Type: application/json" -d '{"title": "Finish project", "description": "Complete the project by end of the month", "duedate": "2024-08-31", "status": "inprogress"}'
```

### Viewing All Tasks

```sh
curl http://localhost:9090/tasks
```

### Viewing a Task by ID

```sh
curl http://localhost:9090/tasks/your_task_id
```

### Updating a Task

```sh
curl -X PUT http://localhost:9090/tasks/your_task_id -H "Content-Type: application/json" -d '{"title": "Updated title", "description": "Updated description", "duedate": "2024-09-01", "status": "completed"}'
```

### Deleting a Task

```sh
curl -X DELETE http://localhost:9090/tasks/your_task_id
```

## Error Handling

The application returns appropriate HTTP status codes and messages for various errors:

- **404 Not Found**: When a task is not found.
- **400 Bad Request**: For invalid input data.
- **500 Internal Server Error**: For unexpected server errors.

## Testing

1. **Test API Endpoints**: Use tools like Postman or `curl` to test the CRUD operations.

2. **Direct MongoDB Queries**: Verify data using MongoDB Compass or the MongoDB shell.


## Contributing

Feel free to fork the repository and submit pull requests for improvements or bug fixes. 

