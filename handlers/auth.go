package handlers

import (
	"mini_site/database"
	"mini_site/models"
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{Username: username, Password: string(hashedPassword)}
	result := database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user models.User
	result := database.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
