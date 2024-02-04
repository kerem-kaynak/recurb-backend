package handlers

import (
	"net/http"
	"strconv"

	"errors"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/kerem-kaynak/recurb/internal/models"
	"github.com/kerem-kaynak/recurb/internal/utils"
	"gorm.io/gorm"
)

type CreateSubscriptionSchema struct {
	Name          string  `json:"name" binding:"required"`
	Website       string  `json:"website"`
	BillingPeriod string  `json:"billing_period" default:"monthly"`
	Amount        float64 `json:"amount" binding:"required"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	Category      string  `json:"category"`
}

type UpdateSubscriptionSchema struct {
	Name          string  `json:"name"`
	Website       string  `json:"website"`
	BillingPeriod string  `json:"billing_period"`
	Amount        float64 `json:"amount"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	Category      string  `json:"category"`
}

func CreateSubscriptionHandler(c *gin.Context) {
	var reqBody CreateSubscriptionSchema

	utils.BindBody(c, &reqBody)

	db, exists := utils.GetDBFromContext(c)
	if !exists {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	userId, err := strconv.Atoi(utils.GetUserFromCookie(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	userIdUint := uint(userId)

	start_date := utils.MaybeParseDate(reqBody.StartDate, c)
	end_date := utils.MaybeParseDate(reqBody.EndDate, c)

	subscription := models.Subscription{
		UserID:        userIdUint,
		Name:          reqBody.Name,
		Website:       reqBody.Website,
		BillingPeriod: reqBody.BillingPeriod,
		Amount:        reqBody.Amount,
		StartDate:     start_date,
		EndDate:       end_date,
		Category:      reqBody.Category,
	}

	if err := db.Create(&subscription).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if subscription.BillingPeriod == "monthly" {
		for i := *start_date; i.Before(*end_date); i = i.AddDate(0, 1, 0) {
			payment := models.Payment{
				SubscriptionID: subscription.ID,
				Date:           i,
				Amount:         subscription.Amount,
			}
			if err := db.Create(&payment).Error; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}
	} else if subscription.BillingPeriod == "yearly" {
		for i := *start_date; i.Before(*end_date); i = i.AddDate(1, 0, 0) {
			payment := models.Payment{
				SubscriptionID: subscription.ID,
				Date:           i,
				Amount:         subscription.Amount,
			}
			if err := db.Create(&payment).Error; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": subscription})
}

func UpdateSubscriptionHandler(c *gin.Context) {
	var reqBody UpdateSubscriptionSchema

	utils.BindBody(c, &reqBody)

	db, exists := utils.GetDBFromContext(c)
	if !exists {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	userId, err := strconv.Atoi(utils.GetUserFromCookie(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	userIdUint := uint(userId)

	start_date := utils.MaybeParseDate(reqBody.StartDate, c)
	end_date := utils.MaybeParseDate(reqBody.EndDate, c)

	subscription := models.Subscription{
		UserID:        userIdUint,
		Name:          reqBody.Name,
		Website:       reqBody.Website,
		BillingPeriod: reqBody.BillingPeriod,
		Amount:        reqBody.Amount,
		StartDate:     start_date,
		EndDate:       end_date,
		Category:      reqBody.Category,
	}

	var existingSubscription models.Subscription

	var existingPayment models.Payment

	if err := db.Where("id = ?", c.Param("id")).First(&existingSubscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Where("id = ? AND user_id =?", c.Param("id"), userIdUint).Updates(&subscription).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := db.Where("subscription_id = ?", c.Param("id")).Delete(&existingPayment).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if subscription.BillingPeriod == "monthly" {
		for i := *start_date; i.Before(*end_date); i = i.AddDate(0, 1, 0) {
			payment := models.Payment{
				SubscriptionID: subscription.ID,
				Date:           i,
				Amount:         subscription.Amount,
			}
			if err := db.Create(&payment).Error; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}
	} else if subscription.BillingPeriod == "yearly" {
		for i := *start_date; i.Before(*end_date); i = i.AddDate(1, 0, 0) {
			payment := models.Payment{
				SubscriptionID: subscription.ID,
				Date:           i,
				Amount:         subscription.Amount,
			}
			if err := db.Create(&payment).Error; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": subscription})
}

func GetAllSubscriptionsHandler(c *gin.Context) {
	userId := utils.GetUserFromCookie(c)

	var subscriptions []models.Subscription

	db, exists := utils.GetDBFromContext(c)
	if !exists {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := db.Where("user_id = ?", userId).Find(&subscriptions).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscriptions})
}

func GetSubscriptionHandler(c *gin.Context) {
	var subscription models.Subscription

	db, exists := utils.GetDBFromContext(c)
	if !exists {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := db.Preload("Payments").Preload("Reminders").Where("id = ?", c.Param("id")).First(&subscription).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscription})
}

func DeleteSubscriptionHandler(c *gin.Context) {
	var subscription models.Subscription
	var payments []models.Payment
	var reminders []models.Reminder

	userId := utils.GetUserFromCookie(c)

	db, exists := utils.GetDBFromContext(c)
	if !exists {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := db.Where("subscription_id = ?", c.Param("id")).Delete(&payments).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := db.Where("subscription_id = ?", c.Param("id")).Delete(&reminders).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := db.Where("id = ? AND user_id = ?", c.Param("id"), userId).Delete(&subscription).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscription})
}

func GetSubscriptionsByCategoryHandler(c *gin.Context) {
	var subscriptions []models.Subscription

	db, exists := utils.GetDBFromContext(c)
	if !exists {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := db.Where("category = ?", c.Param("category")).Find(&subscriptions).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscriptions})
}

func GetAllPaymentsHandler(c *gin.Context) {
	userId := utils.GetUserFromCookie(c)

	var payments []models.Payment

	db, exists := utils.GetDBFromContext(c)
	if !exists {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := db.
		Preload("Subscription").
		Where("subscriptions.user_id = ?", userId).
		Joins("JOIN subscriptions ON subscriptions.id = payments.subscription_id").
		Find(&payments).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payments})
}

func GetCurrentMonthPaymentsHandler(c *gin.Context) {
	userId := utils.GetUserFromCookie(c)

	var payments []models.Payment

	db, exists := utils.GetDBFromContext(c)
	if !exists {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation).UTC()
	lastOfMonth := time.Date(currentYear, currentMonth+1, 1, 0, 0, 0, 0, currentLocation).Add(-time.Second).UTC()

	if err := db.
		Preload("Subscription").
		Where("subscriptions.user_id = ? AND date >= ? AND date <= ?", userId, firstOfMonth, lastOfMonth).
		Joins("JOIN subscriptions ON subscriptions.id = payments.subscription_id").
		Find(&payments).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payments})
}
