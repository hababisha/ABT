package Domain

type User struct {
	ID int `json:"id"`
	Name string 	`json:"username"`
	Role string     `json:"role"`
	Password string `json: "password"`
}

type Task struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
}