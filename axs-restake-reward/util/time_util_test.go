package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCalculateNextTick(t *testing.T) {
	tests := []struct {
		now      time.Time
		lastTick time.Time
		expected time.Time
		duration time.Duration
	}{
		{
			now:      time.Date(2024, 10, 20, 10, 0, 0, 0, time.UTC),
			lastTick: time.Date(2024, 10, 19, 10, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 10, 20, 10, 0, 1, 0, time.UTC),
		},
		{
			now:      time.Date(2024, 10, 20, 9, 0, 0, 0, time.UTC),
			lastTick: time.Date(2024, 10, 19, 10, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 10, 20, 10, 0, 1, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.now.String(), func(t *testing.T) {
			nextTick, duration := CalculateNextTick(tt.now, tt.lastTick)
			assert.Equal(t, tt.expected, nextTick)
			assert.Equal(t, tt.expected.Sub(tt.now), duration)
		})
	}
}

func TestSetTimezone(t *testing.T) {
	tests := []struct {
		timezone    string
		expectedLoc *time.Location
		expectError bool
	}{
		{
			timezone:    "Local",
			expectedLoc: time.Local,
			expectError: false,
		},
		{
			timezone:    "America/New_York",
			expectedLoc: mustLoadLocation("America/New_York"),
			expectError: false,
		},
		{
			timezone:    "Invalid/Timezone",
			expectedLoc: nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		Timezone = tt.timezone

		SetTimezone()

		if tt.expectError {
			assert.Nil(t, Location)
		} else {
			assert.Equal(t, tt.expectedLoc, Location)
		}
	}
}

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}
