package util

import "time"

func IncrementTime(hour, minute int) (int, int) {
	minute += 1
	if minute >= 60 {
		minute = 0
		hour += 1
	}
	if hour >= 24 {
		hour = 0
	}
	return hour, minute
}

func CalculateNextTick(now time.Time, hour, minute int) time.Time {
	nextTick := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	if nextTick.Before(now) {
		nextTick = nextTick.Add(24 * time.Hour)
	}
	return nextTick
}
