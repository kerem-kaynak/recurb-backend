package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email         string
	Subscriptions []Subscription
}

type Subscription struct {
	gorm.Model
	UserID        uint
	Name          string
	Website       string
	BillingPeriod string
	Amount        float64
	StartDate     *time.Time
	EndDate       *time.Time
	Reminders     []Reminder
	Payments      []Payment
	Category      string
}

type Reminder struct {
	gorm.Model
	SubscriptionID uint
	Subscription   Subscription
	Date           *time.Time
	Message        string
}

type Payment struct {
	gorm.Model
	SubscriptionID uint
	Subscription   Subscription
	Date           time.Time
	Amount         float64
}
