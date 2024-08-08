package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"errors"
)

type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	DueDate     time.Time          `json:"duedate" bson:"duedate"`
	Status      string             `json:"status" bson:"status"`
}

// CreateTaskFromInputData converts input data into a Task model.
func CreateTaskFromInputData(inputData map[string]interface{}) (*Task, error) {
	dueDateStr, ok := inputData["duedate"].(string)
	if !ok {
		return nil, errors.New("duedate must be a string")
	}
	dueDate, err := time.Parse("2006-01-02", dueDateStr)
	if err != nil {
		return nil, err
	}

	task := Task{
		Title:       inputData["title"].(string),
		Description: inputData["description"].(string),
		DueDate:     dueDate,
		Status:      inputData["status"].(string),
	}

	return &task, nil
}

