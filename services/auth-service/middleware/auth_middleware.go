package middleware

import (
	"context"
	"github.com/Calendly-Backend-Micro/auth-service/db"
	"net/http"

	"github.com/Calendly-Backend-Micro/auth-service/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWT.Secret), nil
		})

		// Check if the token is blacklisted
		isBlacklisted := db.RedisClient.Get(context.Background(), tokenString).Val()
		if isBlacklisted == "blacklisted" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
			c.Abort()
			return
		}

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}
