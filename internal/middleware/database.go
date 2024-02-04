package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kerem-kaynak/recurb/internal/config"
)

func DBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("db"); exists {
			c.Next()
			return
		}

		db, err := config.GetDB()
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		config.RunMigrations(db)

		c.Set("db", db)
		c.Next()
	}
}
