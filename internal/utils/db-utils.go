package utils

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDBFromContext(c *gin.Context) (*gorm.DB, bool) {
	db, exists := c.Get("db")
	if !exists {
		return nil, false
	}

	gormDB, ok := db.(*gorm.DB)
	if !ok {
		return nil, false
	}

	return gormDB, true
}
