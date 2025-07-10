package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/geno-tree/go-back/internal/services"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не предоставлен"})
			c.Abort()
			return
		}

		// Извлекаем токен из заголовка "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		user, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
			c.Abort()
			return
		}

		// Сохраняем пользователя в контексте
		c.Set("user", *user)
		c.Set("user_id", user.ID)

		c.Next()
	}
}
