package controllers

import (
	"AI_WEB_SCRAPPER/initlizers"
	"AI_WEB_SCRAPPER/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var body struct {
		Name            string `json:"name"`
		Lastname        string `json:"lastname"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"passwordconfirm"`
	}

	if err := c.Bind(&body); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if body.Password != body.PasswordConfirm {
		c.AbortWithStatusJSON(400, gin.H{"error": "Passwords do not match"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Błąd podczas generowania tokenu",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Błąd podczas generowania tokenu",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
			"error": "Nieprawidłowe dane logowania",
		})
		return
	}

	// Generuj JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Błąd podczas generowania tokenu",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
