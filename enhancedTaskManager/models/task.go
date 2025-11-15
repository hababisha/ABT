package models

type Task struct{
	ID string `json: "id, omitempty"` //hex string of the objectof mongo
	Title string `json: "title" binding:"required"`
	Description string `json: "description, omitempty"`
	Status string `json: "status, omitempty"`
}