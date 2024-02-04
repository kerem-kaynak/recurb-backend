package config

import (
	"fmt"
	"os"

	"github.com/kerem-kaynak/recurb/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDB() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	databaseName := os.Getenv("DB_NAME")

	connString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, user, databaseName, password)

	db, err := gorm.Open(postgres.Open(
		connString,
	), &gorm.Config{})
	if err != nil {
		fmt.Errorf("Error connecting to database: %v", err)
		return nil, err
	}
	return db, nil
}

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Subscription{}, &models.Reminder{}, &models.Payment{})
}
