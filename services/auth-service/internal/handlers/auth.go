package handlers

import (
	"github.com/Calendly-Backend-Micro/auth-service/config"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//var jwtSecret []byte

// Initialize the secret when the application starts
//func init() {
//	// Load .env file if it exists (for local development)
//	_ = godotenv.Load()
//
//	// Get the JWT_SECRET from the environment
//	secret := os.Getenv("JWT_SECRET")
//	if secret == "" {
//		panic("JWT_SECRET is not set in environment variables")
//	}
//	jwtSecret = []byte(secret)
//}

// RegisterHandler handles user registration
func RegisterHandler(c *gin.Context) {
	// Dummy implementation
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// LoginHandler handles user login
func LoginHandler(c *gin.Context) {
	// Dummy user authentication
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "admin" && password == "password" {
		token := generateToken(username)
		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

// RefreshTokenHandler refreshes the JWT token
func RefreshTokenHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}

// Generate JWT token
func generateToken(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(config.AppConfig.JWTSecret))
	return tokenString
}
