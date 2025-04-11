package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthHandler(c *gin.Context) {
	authToken := c.GetHeader("Authorization")
	if authToken == "" {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	tokenString := strings.Replace(authToken, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Ustaw ID użytkownika w kontekście
		c.Set("userID", int16(claims["sub"].(float64)))
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Nieprawidłowy token",
		})
		return
	}
}
