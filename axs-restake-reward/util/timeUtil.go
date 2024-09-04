package util

import (
	"time"
)

var (
	NextTick time.Time
	Duration time.Duration
	Location *time.Location
)

func CalculateNextTick(now time.Time, lastTick time.Time) (time.Time, time.Duration) {
	nextTick := lastTick.Add(24 * time.Hour).Add(1 * time.Second)
	return nextTick, nextTick.Sub(now)
}
