package main

import (
	"flag"
	"github.com/deukyunlee/crypto-playground/core"
	"log"
	"os"
	"time"
)

func main() {
	stakeLogPath := "./logs/staking_logs.log"

	// Setup logging
	logFile, err := os.OpenFile(stakeLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer closeLogFile(logFile)

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	hourPtr := flag.Int("hour", 0, "cron hour")
	minutePtr := flag.Int("minute", 0, "cron minute")

	flag.Parse()

	hour := *hourPtr
	minute := *minutePtr

	for {
		now := time.Now()
		hour, minute = incrementTime(hour, minute)
		log.Printf("[CURRENT] [%v] hour: %v, minute: %v\n", now.Format(time.RFC3339), hour, minute)

		nextTick := calculateNextTick(now, hour, minute)
		sleepDuration := nextTick.Sub(now)
		log.Printf("Sleeping for %s until the next tick at %s\n", sleepDuration, nextTick.Format(time.RFC3339))

		time.Sleep(sleepDuration)

		core.RestakeRewards()
	}
}

func closeLogFile(logFile *os.File) {
	if err := logFile.Close(); err != nil {
		log.Fatal(err)
	}
}

func calculateNextTick(now time.Time, hour, minute int) time.Time {
	nextTick := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	if nextTick.Before(now) {
		nextTick = nextTick.Add(24 * time.Hour)
	}
	return nextTick
}

func incrementTime(hour, minute int) (int, int) {
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
