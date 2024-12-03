package routes

import (
	"github.com/Calendly-Backend-Micro/auth-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes registers authentication routes
func SetupAuthRoutes(router *gin.Engine) {
	api := router.Group("/auth")
	{
		api.POST("/register", handlers.RegisterHandler)
		api.POST("/login", handlers.LoginHandler)
		api.POST("/refresh", handlers.RefreshTokenHandler)
	}
}
