package data

import "github.com/hababisha/enhancedTaskManager/models"

type TaskService interface{
	GetAllTasks() ([]models.Task, error) 
	GetTaskByID(id string) (models.Task, error)
	CreateTask(task models.Task) (models.Task, error)
	UpdateTask(id string, updated models.Task) (models.Task, error)
	DeleteTask(id string) error
}