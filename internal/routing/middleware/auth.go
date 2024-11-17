package middleware

import (
	"github.com/gbooksoft/backend/internal/domain/authorization"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		time.Sleep(3 * time.Second)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error_code": http.StatusUnauthorized, "error_message": "Authorization header missing"})
			c.Abort()
			return
		}

		// Проверка формата заголовка (должен начинаться с "Bearer ")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error_code": http.StatusUnauthorized, "error_message": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, claims, err := authorization.ParseToken(tokenString, false)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error_code": http.StatusUnauthorized, "error_message": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("jwt_token_claims", claims)
		// Продолжить выполнение запроса
		c.Next()
	}
}
