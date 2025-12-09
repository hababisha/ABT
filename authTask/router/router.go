package router

import (
    "github.com/hababisha/authTask/controllers"
    "github.com/hababisha/authTask/middleware"
    "github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    //public 
    r.POST("/register", controllers.Register)
    r.POST("/login", controllers.Login)

    //protected
    protected := r.Group("/")
    protected.Use(middleware.AuthMiddleware())

    protected.GET("/tasks", controllers.GetTasks)
    protected.GET("/tasks/:id", controllers.GetTask)

    //adminRotes
    admin := protected.Group("/")
    admin.Use(middleware.AdminOnly())
    admin.POST("/tasks", controllers.CreateTask)
    admin.PUT("/tasks/:id", controllers.UpdateTask)
    admin.DELETE("/tasks/:id", controllers.DeleteTask)
    admin.POST("/promote", controllers.Promote)

    return r
}

