package util

import (
	"time"
)

var (
	NextTick time.Time
	Duration time.Duration
	Location *time.Location
	Timezone string
)

func CalculateNextTick(now time.Time, lastTick time.Time) (time.Time, time.Duration) {
	nextTick := lastTick.Add(24 * time.Hour).Add(1 * time.Second)
	return nextTick, nextTick.Sub(now)
}

func SetTimezone() {
	var err error
	if Timezone == "Local" {
		Location = time.Local
	} else {
		Location, err = time.LoadLocation(Timezone)
		if err != nil {
			logger.Errorf("Error loading timezone: %v", err)
		}
	}
}
