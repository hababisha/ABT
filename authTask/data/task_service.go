package data

import (
    "errors"
    "github.com/hababisha/authTask/models"
)

var tasks []models.Task
var taskID = 1

func CreateTask(title, description string) models.Task {
    t := models.Task{
        ID:          taskID,
        Title:       title,
        Description: description,
    }
    taskID++
    tasks = append(tasks, t)
    return t
}

func GetTasks() []models.Task {
    return tasks
}

func GetTaskByID(id int) (models.Task, error) {
    for _, t := range tasks {
        if t.ID == id {
            return t, nil
        }
    }
    return models.Task{}, errors.New("task not found")
}

func UpdateTask(id int, title, description string) error {
    for i := range tasks {
        if tasks[i].ID == id {
            tasks[i].Title = title
            tasks[i].Description = description
            return nil
        }
    }
    return errors.New("task not found")
}

func DeleteTask(id int) error {
    for i := range tasks {
        if tasks[i].ID == id {
            tasks = append(tasks[:i], tasks[i+1:]...)
            return nil
        }
    }
    return errors.New("task not found")
}
