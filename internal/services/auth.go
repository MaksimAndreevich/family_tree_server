package services

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/geno-tree/go-back/internal/models"
	"gitlab.com/geno-tree/go-back/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository *repositories.UserRepository
}

func NewAuthService(userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepository: userRepository,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Username  string `json:"username" binding:"required,min=3"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AuthResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

func (as *AuthService) Login(c *gin.Context) {
	var loginRequest LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные авторизации: " + err.Error()})
		return
	}

	user, err := as.userRepository.FindUserByEmail(loginRequest.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		return
	}

	token := as.generateToken(user)

	response := AuthResponse{
		User:  *user,
		Token: token,
	}

	c.JSON(http.StatusOK, response)
}

func (as *AuthService) Register(c *gin.Context) {
	var registerRequest RegisterRequest

	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные регистрации: " + err.Error()})
		return
	}

	existingUser, err := as.userRepository.FindUserByEmail(registerRequest.Email)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким email уже существует"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки пароля"})
		return
	}

	user := &models.User{
		Email:    registerRequest.Email,
		Username: registerRequest.Username,
		Password: string(hashedPassword),
	}

	if err := as.userRepository.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания пользователя"})
		return
	}

	token := as.generateToken(user)

	response := AuthResponse{
		User:  *user,
		Token: token,
	}

	c.JSON(http.StatusCreated, response)
}

// generateToken генерирует JWT токен (упрощенная версия)
func (as *AuthService) generateToken(user *models.User) string {
	// Здесь должна быть реальная генерация JWT токена
	// Пока возвращаем простую строку
	return "token_" + user.Username + "_" + user.Email
}

// ValidateToken валидирует токен
func (as *AuthService) ValidateToken(token string) (*models.User, error) {
	// Здесь должна быть реальная валидация JWT токена
	// Пока возвращаем ошибку
	return nil, errors.New("токен недействителен")
}
