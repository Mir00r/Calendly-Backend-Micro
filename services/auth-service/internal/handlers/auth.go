package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Calendly-Backend-Micro/auth-service/config"
	"github.com/Calendly-Backend-Micro/auth-service/db"
	"github.com/Calendly-Backend-Micro/auth-service/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Save user to the database
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func LoginHandler(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT
	token := generateToken(user.Username)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func generateToken(username string) string {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(config.AppConfig.JWT.Secret))
	return tokenString
}

func RefreshTokenHandler(c *gin.Context) {
	oldToken := c.GetHeader("Authorization")
	if oldToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	token, err := jwt.Parse(oldToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	username := claims["username"].(string)
	newToken := generateToken(username)
	c.JSON(http.StatusOK, gin.H{"token": newToken})
}

func LogoutHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// Blacklist the token
	err := db.RedisClient.Set(context.Background(), token, "blacklisted", 24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to blacklist token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func OAuthCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing authorization code"})
		return
	}

	// Exchange the code for an access token
	// (Example using Google's OAuth2 library)
	token, err := exchangeOAuthCodeForToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

// Example function to exchange OAuth code for token
func exchangeOAuthCodeForToken(code string) (string, error) {
	// Use the appropriate OAuth2 client library for the provider
	return "mock-oauth-token", nil
}

func PasswordResetHandler(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a password reset token
	resetToken := generateToken(request.Email)

	// Send the reset token via email
	err := sendPasswordResetEmail(request.Email, resetToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

// Example function to send email
func sendPasswordResetEmail(email, token string) error {
	// Use an email library like SendGrid or SMTP
	return nil
}

func PasswordResetConfirmHandler(c *gin.Context) {
	var request struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse and validate the token
	token, err := jwt.Parse(request.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	email := claims["username"].(string)

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update the password in the database
	if err := db.DB.Model(&models.User{}).Where("username = ?", email).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
