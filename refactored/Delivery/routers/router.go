package routers

import (
	"task-manager/Delivery/controllers"
	"task-manager/Infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(ctrl *controllers.Controller, jwtSvc Infrastructure.JWTService) *gin.Engine {
	r := gin.Default()


	r.POST("/register", ctrl.Register)
	r.POST("/login", ctrl.Login)


	protected := r.Group("/")
	protected.Use(Infrastructure.AuthMiddleware(jwtSvc))
	{
		protected.GET("/tasks", ctrl.GetTasks)
		protected.GET("/tasks/:id", ctrl.GetTask)
	}


	admin := r.Group("/")
	admin.Use(Infrastructure.AuthMiddleware(jwtSvc), Infrastructure.AdminMiddleware())
	{
		admin.POST("/tasks", ctrl.CreateTask)
		admin.PUT("/tasks/:id", ctrl.UpdateTask)
		admin.DELETE("/tasks/:id", ctrl.DeleteTask)

		admin.POST("/promote", ctrl.Promote)
	}

	return r
}
