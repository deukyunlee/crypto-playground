package util

import (
	"time"
)

var NextTick time.Time
var Duration time.Duration

func IncrementTime(prevTime time.Time, duration time.Duration) time.Time {
	return prevTime.Add(duration)
}

func CalculateNextTick(now time.Time, lastTick time.Time) (time.Time, time.Duration) {
	nextTick := lastTick.Add(24 * time.Hour)
	return nextTick, nextTick.Sub(now)
}
