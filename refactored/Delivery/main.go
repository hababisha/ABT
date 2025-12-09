package main

import (
	"log"
	"task_manager/Deliver/controllers"

	"task-manager/Delivery/routers"
	"task-manager/Infrastructure"
	"task-manager/Repositories"
	"task-manager/Usecases"
)

func main() {
	userRepo := Repositories.NewInMemoryUserRepository()
	taskRepo := Repositories.NewInMemoryTaskRepository()

	passwordSvc := Infrastructure.NewPasswordService()
	jwtSvc := Infrastructure.NewJWTService("super-secret-jwt-key") // replace with env var in prod

	userUC := Usecases.NewUserUsecase(userRepo, passwordSvc, jwtSvc)
	taskUC := Usecases.NewTaskUsecase(taskRepo)

	ctrl := controllers.NewController(userUC, taskUC, jwtSvc)

	r := routers.SetupRouter(ctrl, jwtSvc)

	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
