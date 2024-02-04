package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kerem-kaynak/recurb/internal/handlers"
)

func RegisterSubscriptionRoutes(r *gin.RouterGroup) {
	r.GET("/all", handlers.GetAllSubscriptionsHandler)
	r.POST("/create", handlers.CreateSubscriptionHandler)
	r.GET("/:id", handlers.GetSubscriptionHandler)
	r.POST("/:id", handlers.UpdateSubscriptionHandler)
	r.DELETE("/:id", handlers.DeleteSubscriptionHandler)
	r.GET("/category/:category", handlers.GetSubscriptionsByCategoryHandler)
	r.GET("/payments", handlers.GetAllPaymentsHandler)
	r.GET("/payments/current", handlers.GetCurrentMonthPaymentsHandler)
}
