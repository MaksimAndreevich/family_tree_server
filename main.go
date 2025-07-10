package main

import (
	"gitlab.com/geno-tree/go-back/internal/configs"
	"gitlab.com/geno-tree/go-back/internal/controllers"
	"gitlab.com/geno-tree/go-back/internal/database"
	"gitlab.com/geno-tree/go-back/internal/repositories"
	"gitlab.com/geno-tree/go-back/internal/router"
	"gitlab.com/geno-tree/go-back/internal/services"
)

func main() {
	config := configs.NewConfig()
	db := database.NewDb(config)

	// repositories
	userRepository := repositories.NewUserRepository(db)
	// services
	authService := services.NewAuthService(userRepository)
	// controllers
	authController := controllers.NewAuthController(authService)

	router.InitRouter(&router.Controllers{
		AuthController: authController,
	}, authService)

}
