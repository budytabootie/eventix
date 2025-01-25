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
		// Ambil header Authorization
		tokenString := c.GetHeader("Authorization")
		log.Printf("[DEBUG] Authorization header: %s", tokenString)

		if tokenString == "" {
			log.Println("[AuthenticationMiddleware] Authorization header missing")
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Authorization header missing"})
			c.Abort()
			return
		}

		// Validasi format header
		if !strings.HasPrefix(tokenString, "Bearer ") || len(tokenString) <= len("Bearer ") {
			log.Println("[AuthenticationMiddleware] Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		// Trim "Bearer " prefix
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		log.Printf("[AuthenticationMiddleware] Token received: %s", tokenString)

		// Parse token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validasi algoritma signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("[AuthenticationMiddleware] Unexpected signing method: %v", token.Header["alg"])
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

		// Validasi claims token
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

		// Log hasil validasi
		log.Printf("[AuthenticationMiddleware] User authenticated: user_id=%v, role=%v", userID, role)

		// Set context untuk user_id dan role
		c.Set("user_id", uint(userID))
		c.Set("role", role)
		c.Next()
	}
}
