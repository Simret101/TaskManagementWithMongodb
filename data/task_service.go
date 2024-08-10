package data

import (
	"context"
	"errors"
	"example/taskManager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrTaskNotFound = errors.New("task not found")
var ErrDuplicateTitle = errors.New("a task with the same title already exists")

type TaskService struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

// NewTaskService initializes a new TaskService with the provided MongoDB URI.
func NewTaskService(mongoURI string) (*TaskService, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	database := client.Database("task_management")
	collection := database.Collection("tasks")

	return &TaskService{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

// CreateTask inserts a new task into the database and returns the created task with its ID.
func (ts *TaskService) CreateTask(task *models.Task) (*models.Task, error) {
	// Check for duplicate task by title
	existingTask, err := ts.GetTaskByTitle(task.Title)
	if err == nil && existingTask != nil {
		return nil, ErrDuplicateTitle
	}

	result, err := ts.collection.InsertOne(context.Background(), task)
	if err != nil {
		return nil, err
	}

	task.ID = result.InsertedID.(primitive.ObjectID)

	return task, nil
}

// GetTaskByTitle retrieves a Task by its title.
func (ts *TaskService) GetTaskByTitle(title string) (*models.Task, error) {
	filter := bson.M{"title": title}
	var task models.Task
	err := ts.collection.FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}
func (ts *TaskService) DeleteTask(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := ts.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrTaskNotFound
	}

	return nil
}
func (ts *TaskService) UpdateTask(id string, updatedTask *models.Task) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"title":       updatedTask.Title,
			"description": updatedTask.Description,
			"duedate":     updatedTask.DueDate,
			"status":      updatedTask.Status,
		},
	}

	result, err := ts.collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return ErrTaskNotFound
	}

	return nil
}
func (ts *TaskService) GetTask(id string) (*models.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := ts.collection.FindOne(context.Background(), bson.M{"_id": objectID})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, ErrTaskNotFound
		}
		return nil, result.Err()
	}

	var task models.Task
	if err := result.Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (ts *TaskService) GetTasks() ([]*models.Task, error) {
	cursor, err := ts.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var tasks []*models.Task
	for cursor.Next(context.Background()) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

