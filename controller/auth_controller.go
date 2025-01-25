package controller

import (
    "eventix/service"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

type AuthController struct {
    authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
    return &AuthController{authService: authService}
}

func (ctrl *AuthController) Logout(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "No token provided"})
        return
    }

    // Remove "Bearer " prefix
    token = token[len("Bearer "):]

    expiresAt := time.Now().Add(24 * time.Hour) // Sesuaikan dengan waktu token kadaluarsa
    err := ctrl.authService.Logout(token, expiresAt)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to logout"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Logged out successfully"})
}
