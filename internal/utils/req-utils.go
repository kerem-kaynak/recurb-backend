package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func BindBody(c *gin.Context, body interface{}) {
	fmt.Println(c.Request.Body)
	err := c.ShouldBindJSON(body)
	if err != nil {
		fmt.Println("patladik amk")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	return
}

func GetUserFromCookie(c *gin.Context) string {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return err.Error()
	}

	if os.Getenv("ENV") != "production" {
		dotenv_err := godotenv.Load()
		if dotenv_err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return "Error parsing claims"
	}

	userID := (*claims)["iss"].(string)

	return userID
}
