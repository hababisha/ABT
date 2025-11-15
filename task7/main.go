package main

import (
	"log"
	"github.com/hababisha/task_manager/router"
)


func main(){
	r := router.Router()
	if err := r.Run(":8080"); err != nil{
		log.Fatalf("failed to run server: %v", err)
	}
}