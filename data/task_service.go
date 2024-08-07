package data

import (
	"errors"
	"example/task_tutorial/models"
	"sync"
)

var (
	tasks  = []models.Task{}
	lastID = 0
	mu     sync.Mutex
)

func GetAllTasks() []models.Task {
	mu.Lock()
	defer mu.Unlock()
	return tasks
}

func GetTaskByID(id int) (*models.Task, error) {
	mu.Lock()
	defer mu.Unlock()
	for _, task := range tasks {
		if task.ID == id {
			return &task, nil
		}
	}
	return nil, errors.New("task not found")
}

func CreateTask(task *models.Task) {
	mu.Lock()
	defer mu.Unlock()
	lastID++
	task.ID = lastID
	tasks = append(tasks, *task)
}

func UpdateTask(id int, updatedTask *models.Task) error {
	mu.Lock()
	defer mu.Unlock()
	for i, task := range tasks {
		if task.ID == id {
			tasks[i] = *updatedTask
			tasks[i].ID = id
			return nil
		}
	}
	return errors.New("task not found")
}

func DeleteTask(id int) error {
	mu.Lock()
	defer mu.Unlock()
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}
