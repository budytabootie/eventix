package controller

import (
    "eventix/entity"
    "eventix/service"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

type UserController struct {
    service service.UserService
}

func NewUserController(userService service.UserService) *UserController {
    return &UserController{service: userService}
}

func (ctrl *UserController) RegisterUser(c *gin.Context) {
    var user entity.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid user data", "data": nil})
        return
    }
    createdUser, err := ctrl.service.RegisterUser(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to register user", "data": nil})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "User registered successfully", "data": createdUser})
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid user ID", "data": nil})
        return
    }
    user, err := ctrl.service.GetUserByID(uint(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve user", "data": nil})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User retrieved successfully", "data": user})
}

func (ctrl *UserController) UpdateUserRole(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid user ID", "data": nil})
        return
    }

    var reqBody struct {
        Role string `json:"role" binding:"required"`
    }
    if err := c.ShouldBindJSON(&reqBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid data", "data": nil})
        return
    }

    if err := ctrl.service.UpdateUserRole(uint(id), reqBody.Role); err != nil {
        if err.Error() == "invalid role" {
            c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid role", "data": nil})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update user role", "data": nil})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User role updated successfully", "data": nil})
}


func (ctrl *UserController) Login(c *gin.Context) {
    var loginData struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid login data", "data": nil})
        return
    }

    token, err := ctrl.service.Login(loginData.Username, loginData.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": err.Error(), "data": nil})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Login successful", "data": gin.H{"token": token}})
}

