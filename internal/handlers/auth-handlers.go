package handlers

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/kerem-kaynak/recurb/internal/models"
	"github.com/kerem-kaynak/recurb/internal/utils"
	"github.com/markbates/goth/gothic"
)

func AuthCallbackHandler(c *gin.Context) {

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		fmt.Fprintln(c.Writer, err)
		return
	}

	db, exists := utils.GetDBFromContext(c)
	if !exists {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	var userFromDB models.User

	db.Where("email = ?", user.Email).First(&userFromDB)

	if userFromDB.ID == 0 {
		userFromDB = models.User{
			Email: user.Email,
		}
		if err := db.Create(&userFromDB).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    strconv.FormatUint(uint64(userFromDB.ID), 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
	})

	if os.Getenv("ENV") != "production" {
		dotenv_err := godotenv.Load()
		if dotenv_err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	token, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("jwt", token, 3600*24*30, "/", os.Getenv("COOKIE_HOST"), false, false)
	c.SetSameSite(http.SameSiteNoneMode)

	fmt.Printf("%+v\n", user)
	c.Redirect(http.StatusTemporaryRedirect, os.Getenv("CLIENT_HOST"))
}

func LogoutHandler(c *gin.Context) {
	gothic.Logout(c.Writer, c.Request)
	c.SetCookie("jwt", "", -1, "/", os.Getenv("COOKIE_HOST"), false, true)
	c.Writer.Header().Set("Location", os.Getenv("CLIENT_LOGIN_REDIRECT"))
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}

func AuthHandler(c *gin.Context) {
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func IndexHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Success"})
}
