package routes

import (
	"github.com/Calendly-Backend-Micro/auth-service/internal/handlers"
	"github.com/Calendly-Backend-Micro/auth-service/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupAuthRoutes(router *gin.Engine) {
	api := router.Group("/auth")
	{
		api.POST("/v1/register", handlers.RegisterHandler)
		api.POST("/v1/login", handlers.LoginHandler)
		api.POST("/v1/refresh", handlers.RefreshTokenHandler)
		api.POST("/v1/logout", handlers.LogoutHandler)
		api.GET("/v1/oauth/callback", handlers.OAuthCallbackHandler)
		api.POST("/v1/password-reset", handlers.PasswordResetHandler)
		api.POST("/v1/password-reset/confirm", handlers.PasswordResetConfirmHandler)
		api.GET("/v1/protected", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "This is a protected route"})
		})
	}
}
