package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthenticationMiddleware(secretKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            log.Println("[AuthenticationMiddleware] Authorization header missing")
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Authorization header missing"})
            c.Abort()
            return
        }

        tokenString = strings.TrimPrefix(tokenString, "Bearer ")
        log.Printf("[AuthenticationMiddleware] Token received: %s", tokenString)

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                log.Println("[AuthenticationMiddleware] Invalid token signing method")
                return nil, jwt.ErrInvalidKey
            }
            return []byte(secretKey), nil
        })

        if err != nil || !token.Valid {
            log.Printf("[AuthenticationMiddleware] Invalid token: %v", err)
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            log.Println("[AuthenticationMiddleware] Invalid token claims")
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid token claims"})
            c.Abort()
            return
        }

        userID, ok := claims["user_id"].(float64)
        if !ok {
            log.Println("[AuthenticationMiddleware] Invalid user_id in token")
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid user_id in token"})
            c.Abort()
            return
        }

        role, ok := claims["role"].(string)
        if !ok {
            log.Println("[AuthenticationMiddleware] Invalid role in token")
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid role in token"})
            c.Abort()
            return
        }

        log.Printf("[AuthenticationMiddleware] User authenticated: user_id=%v, role=%v\n", userID, role)

        c.Set("user_id", uint(userID))
        c.Set("role", role)
        c.Next()
    }
}


