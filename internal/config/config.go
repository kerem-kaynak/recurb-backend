package config

import (
	"github.com/kerem-kaynak/recurb/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("sample.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Subscription{}, &models.Reminder{}, &models.Payment{})
}
