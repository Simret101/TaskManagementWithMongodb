package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	DueDate     time.Time          `json:"duedate" bson:"duedate"`
	Status      string             `json:"status" bson:"status"`
}

func CreateTaskFromInputData(inputData map[string]interface{}) (*Task, error) {

	dueDateStr := inputData["duedate"].(string)
	loc, _ := time.LoadLocation("YourTimeZone")
	dueDate, err := time.ParseInLocation("2006-01-02", dueDateStr, loc)
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
