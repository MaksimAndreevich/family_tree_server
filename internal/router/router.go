package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/geno-tree/go-back/internal/controllers"
	"gitlab.com/geno-tree/go-back/internal/middlewares"
	"gitlab.com/geno-tree/go-back/internal/services"
)

type Controllers struct {
	AuthController         *controllers.AuthController
	PersonController       *controllers.PersonController
	RelationshipController *controllers.RelationshipController
	TreeController         *controllers.TreeController
}

func InitRouter(constollers *Controllers, authService *services.AuthService) {
	r := gin.Default()

	// публичные маршруты
	public := r.Group("/api")
	{
		public.POST("auth/login", constollers.AuthController.Login)
		public.POST("auth/register", constollers.AuthController.Register)
	}

	// Защищенные маршруты
	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware(authService))
	{
		// Профиль пользователя
		protected.GET("/profile", constollers.AuthController.GetProfile)

		// Управление персонами
		protected.POST("/persons", constollers.PersonController.CreatePerson)
		protected.GET("/persons", constollers.PersonController.GetPersons)
		protected.GET("/persons/search", constollers.PersonController.SearchPersons)

		// Управление конкретными персонами
		protected.GET("/persons/:id", constollers.PersonController.GetPerson)
		protected.PATCH("/persons/:id", constollers.PersonController.UpdatePerson)
		protected.DELETE("/persons/:id", constollers.PersonController.DeletePerson)

		// Управление связями персон (используем другой путь)
		protected.GET("/persons/:id/relationships", constollers.RelationshipController.GetRelationshipsByPerson)

		// Управление связями
		protected.POST("/relationships", constollers.RelationshipController.CreateRelationship)
		protected.GET("/relationships", constollers.RelationshipController.GetRelationships)
		protected.PUT("/relationships/:id", constollers.RelationshipController.UpdateRelationship)
		protected.DELETE("/relationships/:id", constollers.RelationshipController.DeleteRelationship)

		// Работа с деревом
		protected.GET("/tree", constollers.TreeController.GetMyFamilyTree)
		protected.GET("/tree/:personId", constollers.TreeController.GetFamilyTree)
		protected.GET("/tree/statistics", constollers.TreeController.GetTreeStatistics)
	}

	r.Run()
}
