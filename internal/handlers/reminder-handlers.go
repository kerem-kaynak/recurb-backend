package handlers

import (
	"net/http"
	"strconv"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/kerem-kaynak/recurb/internal/models"
	"github.com/kerem-kaynak/recurb/internal/utils"
	"gorm.io/gorm"
)

type CreateReminderSchema struct {
	SubscriptionID string `json:"subscription_id" binding:"required"`
	Date           string `json:"date" binding:"required"`
	Message        string `json:"message" binding:"required"`
}

type UpdateReminderSchema struct {
	SubscriptionID string `json:"subscription_id" binding:"required"`
	Date           string `json:"date" binding:"required"`
	Message        string `json:"message" binding:"required"`
}

func CreateReminderHandler(c *gin.Context) {
	var reqBody CreateReminderSchema

	utils.BindBody(c, &reqBody)

	db, exists := utils.GetDBFromContext(c)
	if !exists {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	subscriptionId, err := strconv.Atoi(reqBody.SubscriptionID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	subscriptionIdUint := uint(subscriptionId)

	reminderDate := utils.MaybeParseDate(reqBody.Date, c)

	reminder := models.Reminder{
		SubscriptionID: subscriptionIdUint,
		Date:           reminderDate,
		Message:        reqBody.Message,
	}

	if err := db.Create(&reminder).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reminder})
}

func GetReminderHandler(c *gin.Context) {
	id := c.Param("id")

	db, exists := utils.GetDBFromContext(c)
	if !exists {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	var reminder models.Reminder
	if err := db.Preload("Subscription").First(&reminder, id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reminder})
}

func UpdateReminderHandler(c *gin.Context) {
	var reqBody UpdateReminderSchema

	utils.BindBody(c, &reqBody)

	db, exists := utils.GetDBFromContext(c)
	if !exists {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	subscriptionId, err := strconv.Atoi(reqBody.SubscriptionID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	subscriptionIdUint := uint(subscriptionId)

	reminderDate := utils.MaybeParseDate(reqBody.Date, c)

	reminder := models.Reminder{
		SubscriptionID: subscriptionIdUint,
		Date:           reminderDate,
		Message:        reqBody.Message,
	}

	var existingReminder models.Reminder

	if err := db.Where("id = ?", c.Param("id")).First(&existingReminder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reminder not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&reminder).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reminder})
}

func DeleteReminderHandler(c *gin.Context) {
	id := c.Param("id")

	db, exists := utils.GetDBFromContext(c)
	if !exists {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	var reminder models.Reminder
	if err := db.Delete(&reminder, id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reminder})
}

func GetAllRemindersHandler(c *gin.Context) {
	userId := utils.GetUserFromCookie(c)
	var reminders []models.Reminder

	db, exists := utils.GetDBFromContext(c)
	if !exists {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := db.Preload("Subscription").Joins("JOIN subscriptions ON reminders.subscription_id = subscriptions.id").Where("subscriptions.user_id = ?", userId).Find(&reminders).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reminders})
}

func GetRemindersBySubscriptionHandler(c *gin.Context) {
	subscriptionID := c.Param("subscription_id")

	var reminders []models.Reminder

	db, exists := utils.GetDBFromContext(c)
	if !exists {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := db.Where("subscription_id = ?", subscriptionID).Find(&reminders).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reminders})
}
