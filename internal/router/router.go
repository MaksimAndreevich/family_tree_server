package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/geno-tree/go-back/internal/controllers"
	"gitlab.com/geno-tree/go-back/internal/middlewares"
	"gitlab.com/geno-tree/go-back/internal/services"
)

type Controllers struct {
	AuthController *controllers.AuthController
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
		protected.GET("/profile", constollers.AuthController.GetProfile)
	}

	r.Run()
}
