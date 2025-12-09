package data

import (
	"errors"

	"github.com/hababisha/authTask/models"
	"golang.org/x/crypto/bcrypt"
)


var users []models.User
var userIDTracker = 1

func HashPassword(password string) (string, error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil 
}

func CreateUser(username, password string) (models.User, error) {
	for _, u  := range users {
		if u.Username == username{
			return models.User{}, errors.New("username already exists")
		}
	}

	hashed, _ := HashPassword(password)
	role := "user"
	if len(users) == 0 {
		role="admin"
	}

	newUser := models.User{
		ID: userIDTracker,
		Username: username,	
		Password: hashed,
		Role: role,
	}
	userIDTracker += 1
	users = append(users, newUser)
	return newUser, nil
}

func AuthenticateUser(username, password string) (models.User, error){
	for _, u := range users{
		if u.Username == username && CheckPasswordHash(u.Password,password){
			return u, nil
		}
	}

	return models.User{}, errors.New("invalid credentials")
}


func PromoteUser(id int) error {
	for i := range users {
		if users[i].ID == id {
			users[i].Role = "admin"
			return nil
		}
	}
	return errors.New("user not found")
}

func GetAllUsers() []models.User {
	return users
}