package datetime

import (
	"time"
)

// GetStarDateFromString gets start date from string
func GetStarDateFromString(date string) (*time.Time, error) {
	if date == "" {
		return nil, nil
	}

	startTime, err := time.Parse("2006-01-02T15:04:05", date)
	if err != nil {
		return nil, err
	}

	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
	return &startTime, nil
}
