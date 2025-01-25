package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeRole(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists {
            log.Println("[AuthorizeRole] Role not found in context")
            c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Forbidden: insufficient permissions"})
            c.Abort()
            return
        }

        log.Printf("[AuthorizeRole] User role: %v, Required role: %v\n", role, requiredRole)

        if role != requiredRole {
            log.Printf("[AuthorizeRole] Role mismatch. User role: %v, Required: %v\n", role, requiredRole)
            c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Forbidden: insufficient permissions"})
            c.Abort()
            return
        }
        log.Println("[AuthorizeRole] Role authorized")
        c.Next()
    }
}


