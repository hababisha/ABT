package main

import (

	"log"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/hababisha/enhancedTaskManager/controllers"
	"github.com/hababisha/enhancedTaskManager/data"
	"github.com/joho/godotenv"
)
func main(){
	err := godotenv.Load(".env")
	if err != nil{
		log.Fatal("error loading .env file")
	}
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")
	colName := os.Getenv("MONGO_COLLECTION")

	svc, err := data.NewMongoTaskService(mongoURI, dbName, colName, 10*time.Second)
	if err != nil{
		log.Fatalf("failed to connect to mogno: %v", err)
	}

	tc := controllers.NewTaskController(svc)
	r := gin.Default()
	r.GET("/tasks", tc.GetTasks)
	r.GET("/tasks/:id", tc.GetTaskByID)
	r.POST("/tasks", tc.CreateTask)
	r.PUT("/tasks/:id", tc.UpdateTask)
	r.DELETE("/tasks/:id", tc.DeleteTask)

	log.Println("server running on :8080")
	if err := r.Run(":8080"); err != nil{
		log.Fatal(err)
	}

}