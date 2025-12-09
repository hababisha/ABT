package controllers

import (
    "net/http"
    "strconv"
    "github.com/hababisha/authTask/data"
    "github.com/hababisha/authTask/middleware"
    "github.com/hababisha/authTask/models"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)


func Register(c *gin.Context) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    c.BindJSON(&req)

    user, err := data.CreateUser(req.Username, req.Password)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "user created", "user": user})
}

func Login(c *gin.Context) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    c.BindJSON(&req)

    user, err := data.AuthenticateUser(req.Username, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
        return
    }

    claims := models.Claims{
        UserID: user.ID,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString(middleware.JWT_SECRET)

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Promote(c *gin.Context) {
    var req struct {
        UserID int `json:"user_id"`
    }
    c.BindJSON(&req)

    err := data.PromoteUser(req.UserID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin"})
}


func CreateTask(c *gin.Context) {
    var req struct {
        Title       string `json:"title"`
        Description string `json:"description"`
    }

    c.BindJSON(&req)

    t := data.CreateTask(req.Title, req.Description)
    c.JSON(http.StatusCreated, t)
}

func GetTasks(c *gin.Context) {
    c.JSON(http.StatusOK, data.GetTasks())
}

func GetTask(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    task, err := data.GetTaskByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
        return
    }

    c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    var req struct {
        Title       string `json:"title"`
        Description string `json:"description"`
    }
    c.BindJSON(&req)

    err := data.UpdateTask(id, req.Title, req.Description)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "task updated"})
}

func DeleteTask(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    err := data.DeleteTask(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}
