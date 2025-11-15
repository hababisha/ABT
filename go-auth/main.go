package main

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
)

var users = make(map[string]*User)
func main(){
	router := gin.Default()

	router.GET("/", func(c *gin.Context){
		c.JSON(200, gin.H{
			"message": "welcome to the go authentication and authorization tut!",
		})
	})

	router.POST("/register", func(c *gin.Context){
		var user User
		if err := c.ShouldBindJSON(&user); err != nil{
			c.JSON(400, gin.H{"error": "invalid request payload"})
			return
		}
		//TODO: Implement user registration logic
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil{
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		user.Password = string(hashedPassword)
		users[user.Email] = &user
		c.JSON(200, gin.H{"message": "user registered successfully"})
	})

	router.Run(":8081")
}