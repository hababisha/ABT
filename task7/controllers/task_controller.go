package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hababisha/task_manager/data"
	"github.com/hababisha/task_manager/models"
)

type TaskController struct{
	service *data.TaskService
}

func NewTaskController() *TaskController {
	return &TaskController{service: data.NewTaskService()}
}

//Get Tasks

func (tc *TaskController) GetTasks(c *gin.Context){
	tasks := tc.service.GetAllTasks()
	c.JSON(http.StatusOK, tasks)
}

//Get Task By ID

func (tc *TaskController) GetTaskByID(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid task id"})
		return
	}

	task, found := tc.service.GetTaskByID(id)
	if !found{
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

//Create Task

func (tc *TaskController) CreateTask(c *gin.Context){
	var newTask models.Task
	if err := c.ShouldBindJSON(&newTask); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json data"})
		return 
	}

	createdTask := tc.service.CreateTask(newTask)
	c.JSON(http.StatusCreated, createdTask)
}

//Update task

func (tc *TaskController) UpdateTask(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return 
	}

	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json data"})
		return
	}

	task, ok := tc.service.UpdateTask(id, updatedTask)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}


//Delete Task

func (tc *TaskController) DeleteTask(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return 
	}

	if ok := tc.service.DeleteTask(id); !ok{
		c.JSON(http.StatusNotFound, gin.H{"error":"task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfylly!"})
}
