package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kerem-kaynak/recurb/internal/handlers"
)

func RegisterAuthRoutes(r *gin.RouterGroup) {
	r.GET("/auth/google/callback", handlers.AuthCallbackHandler)
	r.GET("/auth", handlers.AuthHandler)
	r.GET("/auth/logout", handlers.LogoutHandler)
	r.GET("/", handlers.IndexHandler)
}
