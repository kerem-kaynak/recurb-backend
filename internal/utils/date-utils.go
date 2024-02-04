package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

func parseDate(dateString string) (*time.Time, error) {
	const dateFormat = "2006-01-02"
	parsedTime, err := time.Parse(dateFormat, dateString)
	if err != nil {
		return nil, err
	}
	return &parsedTime, nil
}

func MaybeParseDate(date string, c *gin.Context) *time.Time {
	if date != "" {
		parsedDate, err := parseDate(date)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid date format for " + date})
			return nil
		}
		return parsedDate
	} else {
		return nil
	}

}

func ComputePaymentScheduleEndDate() time.Time {
	return time.Now().AddDate(1, 0, 0)
}
