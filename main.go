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
	personRepository := repositories.NewPersonRepository(db)
	relationshipRepository := repositories.NewRelationshipRepository(db)
	// services
	authService := services.NewAuthService(userRepository, config)
	treeService := services.NewTreeService(personRepository, relationshipRepository)
	// controllers
	authController := controllers.NewAuthController(authService)
	personController := controllers.NewPersonController(treeService)
	relationshipController := controllers.NewRelationshipController(treeService)
	treeController := controllers.NewTreeController(treeService)

	router.InitRouter(&router.Controllers{
		AuthController:         authController,
		PersonController:       personController,
		RelationshipController: relationshipController,
		TreeController:         treeController,
	}, authService)

}
