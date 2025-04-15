package controllers

import (
	"AI_WEB_SCRAPPER/initlizers"
	"AI_WEB_SCRAPPER/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAccount retrieves the user's account information
func GetAccount(c *gin.Context) {
	userId := c.MustGet("userId")
	var user models.User
	if err := initlizers.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateAccount updates the user's account information
func UpdateAccount(c *gin.Context) {
	userId := c.MustGet("userId")
	var body struct {
		Name     string `json:"name"`
		Lastname string `json:"lastname"`
		Email    string `json:"email"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := initlizers.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Name = body.Name
	user.Lastname = body.Lastname
	user.Email = body.Email
	if err := initlizers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// DeleteAccount deletes the user's account
func DeleteAccount(c *gin.Context) {
	userId := c.MustGet("userId")
	var user models.User
	if initlizers.DB.First(&user, userId).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
