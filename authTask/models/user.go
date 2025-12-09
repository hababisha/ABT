package models

type User struct {
	ID  int `json:"id"`
	Username string `json : "name"`
	Password string `json: "password"`
	Role string `json:"role"`
}
