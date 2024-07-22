package util

import "time"

func IncrementTime(prevTime time.Time, duration time.Duration) time.Time {
	return prevTime.Add(duration)
}

func CalculateNextTick(now time.Time, prevPlusOneMinute time.Time) time.Time {
	if prevPlusOneMinute.Add(24 * time.Hour).Before(now) {
		return now
	} else {
		return prevPlusOneMinute.Add(24 * time.Hour)
	}
}
