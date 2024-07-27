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
	if lastTick.Add(24 * time.Hour).Before(now) {
		return now, 0
	}

	nextTick := lastTick
	for nextTick.Before(now) {
		nextTick = nextTick.Add(24 * time.Hour)
	}

	nextTick = nextTick.Add(1 * time.Minute)

	return nextTick, nextTick.Sub(now)
}
