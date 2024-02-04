package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kerem-kaynak/recurb/internal/handlers"
)

func RegisterReminderRoutes(r *gin.RouterGroup) {
	r.GET("/all", handlers.GetAllRemindersHandler)
	r.POST("/create", handlers.CreateReminderHandler)
	r.GET("/:id", handlers.GetReminderHandler)
	r.POST("/:id", handlers.UpdateReminderHandler)
	r.DELETE("/:id", handlers.DeleteReminderHandler)
	r.GET("/subscription/:subscription_id", handlers.GetRemindersBySubscriptionHandler)
}
