package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/geno-tree/go-back/internal/models"
	"gitlab.com/geno-tree/go-back/internal/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Login(c *gin.Context) {
	ac.authService.Login(c)
}

func (ac *AuthController) Register(c *gin.Context) {
	ac.authService.Register(c)
}

func (ac *AuthController) GetProfile(c *gin.Context) {
	// Получаем пользователя из контекста (установленного middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	user := userInterface.(models.User)
	c.JSON(http.StatusOK, gin.H{"user": user})
}
