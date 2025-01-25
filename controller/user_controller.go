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

// RegisterUser godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags User Management
// @Accept json
// @Produce json
// @Param user body entity.User true "User data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /users/register [post]
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

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve user details by ID
// @Tags User Management
// @Produce json
// @Param id path uint true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [get]
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

// UpdateUserRole godoc
// @Summary Update user role
// @Description Update role of a specific user (Admin only)
// @Tags User Management
// @Accept json
// @Produce json
// @Param id path uint true "User ID"
// @Param role body map[string]string true "Role data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /users/{id}/role [put]
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