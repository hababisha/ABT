package controllers

import (
	"net/http"
	"strconv"
	"task-manager/Domain"
	"task-manager/Usecases"

	"github.com/gin-gonic/gin"
	"task-manager/Infrastructure"
)

type Controller struct {
	userUC Usecases.UserUsecase
	taskUC Usecases.TaskUsecase
	jwtSvc Infrastructure.JWTService
}

func NewController(u Usecases.UserUsecase, t Usecases.TaskUsecase, j Infrastructure.JWTService) *Controller {
	return &Controller{userUC: u, taskUC: t, jwtSvc: j}
}

func (c *Controller) Register(ctx *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	user, err := c.userUC.Register(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Password = ""
	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}

func (c *Controller) Login(ctx *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	token, err := c.userUC.Login(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *Controller) Promote(ctx *gin.Context) {
	var req struct {
		UserID int `json:"user_id"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	err := c.userUC.Promote(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user promoted"})
}

func (c *Controller) CreateTask(ctx *gin.Context) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	t, err := c.taskUC.Create(req.Title, req.Description)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, t)
}


func (c *Controller) GetTasks(ctx *gin.Context) {
	tasks, _ := c.taskUC.GetAll()
	ctx.JSON(http.StatusOK, tasks)
}


func (c *Controller) GetTask(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	t, err := c.taskUC.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	ctx.JSON(http.StatusOK, t)
}


func (c *Controller) UpdateTask(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	if err := c.taskUC.Update(id, req.Title, req.Description); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "task updated"})
}


func (c *Controller) DeleteTask(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := c.taskUC.Delete(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}
