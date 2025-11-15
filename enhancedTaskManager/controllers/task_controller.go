package controllers

import (
    "net/http"
	"go.mongodb.org/mongo-driver/mongo"
    "github.com/gin-gonic/gin"
    "github.com/hababisha/enhancedTaskManager/data"   
    "github.com/hababisha/enhancedTaskManager/models"   
)

type TaskController struct {
    service data.TaskService
}

func NewTaskController(svc data.TaskService) *TaskController {
    return &TaskController{service: svc}
}

func (tc *TaskController) GetTasks(c *gin.Context) {
    tasks, err := tc.service.GetAllTasks()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) GetTaskByID(c *gin.Context) {
    id := c.Param("id") // no Atoi, id is string
    task, err := tc.service.GetTaskByID(id)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, task)
}

func (tc *TaskController) CreateTask(c *gin.Context) {
    var newTask models.Task
    if err := c.ShouldBindJSON(&newTask); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json data"})
        return
    }

    created, err := tc.service.CreateTask(newTask)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, created)
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
    id := c.Param("id")

    var updatedTask models.Task
    if err := c.ShouldBindJSON(&updatedTask); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json data"})
        return
    }

    task, err := tc.service.UpdateTask(id, updatedTask)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, task)
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
    id := c.Param("id")
    if err := tc.service.DeleteTask(id); err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully!"})
}
