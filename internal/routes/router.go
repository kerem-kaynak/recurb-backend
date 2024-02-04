package routes

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kerem-kaynak/recurb/internal/middleware"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5173", "http://localhost:5173", "http://localhost:8080", os.Getenv("CLIENT_HOST"), os.Getenv("SERVER_HOST"), "https://substopr.web.app"}, // Add other origins as needed
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	router.Use(middleware.DBMiddleware())

	subscriptions_api := router.Group("/api/v1/subscriptions")
	reminders_api := router.Group("/api/v1/reminders")
	users_api := router.Group("/api/v1/users")
	auth_api := router.Group("/")

	subscriptions_api.Use(middleware.AuthMiddleware())
	reminders_api.Use(middleware.AuthMiddleware())
	users_api.Use(middleware.AuthMiddleware())

	RegisterSubscriptionRoutes(subscriptions_api)
	RegisterReminderRoutes(reminders_api)
	RegisterAuthRoutes(auth_api)
	RegisterUserRoutes(users_api)

	return router
}
