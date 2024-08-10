package controllers

import (
	"errors"
	"example/taskManager/data"
	"example/taskManager/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// retrieves all tasks from the database and returns them as JSON.
func (tc *TaskController) GetTasks(c *gin.Context) {
	tasks, err := tc.taskService.GetTasks()
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// retrieves a single task by its ID and returns it as JSON.
func (tc *TaskController) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := tc.taskService.GetTask(id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

type TaskController struct {
	taskService *data.TaskService
}

// initializes a new TaskController with the provided MongoDB URI.
func NewTaskController(mongoURI string) (*TaskController, error) {
	taskService, err := data.NewTaskService(mongoURI)
	if err != nil {
		return nil, err
	}
	return &TaskController{
		taskService: taskService,
	}, nil
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	var taskData map[string]interface{}
	if err := c.ShouldBindJSON(&taskData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
		return
	}

	task, err := models.CreateTaskFromInputData(taskData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTask, err := tc.taskService.CreateTask(task)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateTitle) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "a task with the same title already exists"})
			return
		}
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdTask)
}

// UpdateTask updates an existing task by its ID.
func (tc *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var updatedTask models.Task
	if err := c.BindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := tc.taskService.UpdateTask(id, &updatedTask)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task updated successfully"})
}

// DeleteTask deletes a task by its ID.
func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := tc.taskService.DeleteTask(id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "task deleted successfully"})
}

// provides centralized error handling for the API.
func handleError(c *gin.Context, err error) {
	if errors.Is(err, data.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}



