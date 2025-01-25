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

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body map[string]string true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
    var loginData struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid login data"})
        return
    }

    // Proses login
    token, err := ctrl.authService.Login(loginData.Username, loginData.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Login successful",
        "data":    gin.H{"token": token},
    })
}

// Logout godoc
// @Summary Logout user
// @Description Revoke user session
// @Tags Authentication
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /logout [post]
func (ctrl *AuthController) Logout(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "No token provided"})
        return
    }

    // Remove "Bearer " prefix
    token = token[len("Bearer "):]

    expiresAt := time.Now().Add(24 * time.Hour)
    err := ctrl.authService.Logout(token, expiresAt)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to logout"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Logged out successfully"})
}
