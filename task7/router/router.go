package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hababisha/task_manager/controllers"
)
func Router() *gin.Engine{
	r := gin.Default()
	taskController :=  controllers.NewTaskController()

	r.GET("/tasks", taskController.GetTasks)
	r.GET("/tasks/:id", taskController.GetTaskByID)
	r.POST("/tasks", taskController.CreateTask)
	r.PUT("/tasks/:id", taskController.UpdateTask)
	r.DELETE("/tasks/:id", taskController.DeleteTask)

	return r
}