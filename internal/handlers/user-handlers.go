package handlers

import (
	"net/http"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/kerem-kaynak/recurb/internal/models"
	"github.com/kerem-kaynak/recurb/internal/utils"
)

func GetUserHandler(c *gin.Context) {
	cookie, err := c.Request.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token present"})
		return
	}

	if os.Getenv("ENV") != "production" {
		dotenv_err := godotenv.Load()
		if dotenv_err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	tokenString := cookie.Value
	claims := &jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := claims.GetIssuer()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	_ = &token

	db, ok := utils.GetDBFromContext(c)
	if !ok {
		c.JSON(500, gin.H{"error": "db not found"})
		return
	}

	var user models.User

	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
