package models

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskStatus string

const (
	TaskStatusComplete   TaskStatus = "complete"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusStarted    TaskStatus = "started"
)

type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	DueDate     time.Time          `json:"duedate" bson:"duedate"`
	Status      TaskStatus         `json:"status" bson:"status"`
}

// checks if the task's fields are valid.
func (t *Task) Validate() error {
	if t.ID == primitive.NilObjectID {
		t.ID = primitive.NewObjectID()
	}
	if err := validateTitle(t.Title); err != nil {
		return err
	}
	if err := validateDescription(t.Description); err != nil {
		return err
	}
	if err := validateDueDate(t.DueDate); err != nil {
		return err
	}
	if err := validateStatus(t.Status); err != nil {
		return err
	}
	return nil
}

// ensures that the title is not empty and within the length limit.
func validateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("title must not be empty")
	}
	if len(title) > 100 {
		return errors.New("title must be less than 100 characters")
	}
	if len(title) < 3 {
		return errors.New("title must be greater than 3 characters")
	}
	return nil
}

// ensures that the description is not empty.
func validateDescription(description string) error {
	if strings.TrimSpace(description) == "" {
		return errors.New("description must not be empty")
	}
	return nil
}

// ensures that the due date is specified and valid.
func validateDueDate(dueDate time.Time) error {
	if dueDate.IsZero() {
		return errors.New("due date must be specified")
	}
	return nil
}

// ensures the status is one of the allowed values.
func validateStatus(status TaskStatus) error {
	switch status {
	case TaskStatusComplete, TaskStatusInProgress, TaskStatusStarted:
		return nil
	default:
		return errors.New("status is invalid")
	}
}

// converts input data into a Task model and validates it.
func CreateTaskFromInputData(inputData map[string]interface{}) (*Task, error) {
	title, err := getStringField(inputData, "title")
	if err != nil {
		return nil, err
	}

	description, err := getStringField(inputData, "description")
	if err != nil {
		return nil, err
	}

	dueDateStr, ok := inputData["duedate"].(string)
	if !ok || strings.TrimSpace(dueDateStr) == "" {
		return nil, errors.New("due date must not be empty and must be a string")
	}
	dueDate, err := time.Parse("2006-01-02", dueDateStr)
	if err != nil {
		return nil, errors.New("due date format is invalid")
	}

	statusStr, ok := inputData["status"].(string)
	if !ok || strings.TrimSpace(statusStr) == "" {
		return nil, errors.New("status must not be empty and must be a string")
	}

	task := &Task{
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Status:      TaskStatus(statusStr),
	}

	// Validate the created task
	if err := task.Validate(); err != nil {
		return nil, err
	}

	return task, nil
}

// getStringField ensures that a field exists and is a non-empty string.
func getStringField(data map[string]interface{}, fieldName string) (string, error) {
	value, ok := data[fieldName].(string)
	if !ok || strings.TrimSpace(value) == "" {
		return "", errors.New(fieldName + " must not be empty and must be a string")
	}
	if len(value) > 100 {
		return "", errors.New(fieldName + " must be less than 100 characters")
	}
	return value, nil
}


