package config

import (
	"fmt"
	"os"

	"sync"

	"github.com/kerem-kaynak/recurb/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	err  error
	once sync.Once
)

func GetDB() (*gorm.DB, error) {
	once.Do(func() {
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbHost := os.Getenv("DB_HOST")
		databaseName := os.Getenv("DB_NAME")

		connString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, user, databaseName, password)

		db, err = gorm.Open(postgres.Open(connString), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	})

	return db, err
}

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Subscription{}, &models.Reminder{}, &models.Payment{})
}
