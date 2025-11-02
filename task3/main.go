package task3

import (
	"github.com/hababisha/ABT/task3/controllers"
	"github.com/hababisha/ABT/task3/models"
	"github.com/hababisha/ABT/task3/services"
)


func main() {
	library := services.NewLibrary()
	
	library.Members[1] = models.Member{ID: 1, Name: "Abebe"}
	library.Members[2] = models.Member{ID: 2, Name: "Kebede"}
	
	library.AddBook(models.Book{ID: 101, Title: "The Go Programming Language", Author: "Donovan & Kernighan", Status: "available"})
	library.AddBook(models.Book{ID: 102, Title: "Learning Go", Author: "Jon Bodner", Status: "available"})
	controllers.Start(library)
}