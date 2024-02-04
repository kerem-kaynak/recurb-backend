package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kerem-kaynak/recurb/internal/handlers"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	r.GET("/", handlers.GetUserHandler)
}
