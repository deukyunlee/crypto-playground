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
	defer func(logFile *os.File) {
		err = logFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(logFile)

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	hourPtr := flag.Int("hour", 0, "cron hour")
	minutePtr := flag.Int("minute", 0, "cron minute")

	flag.Parse()

	hour := *hourPtr
	minute := *minutePtr

	log.Printf("[INITIAL] [%v] hour, minute: %v, %v\n", time.Now().Format(time.RFC3339), hour, minute)

	now := time.Now()
	initialTick := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	if initialTick.Before(now) {
		initialTick = initialTick.Add(24 * time.Hour)
	}

	initialSleepDuration := initialTick.Sub(now)
	log.Printf("Sleeping for %s until the first execution at %s\n", initialSleepDuration, initialTick.Format(time.RFC3339))

	time.Sleep(initialSleepDuration)
	for {
		minute += 1

		if minute >= 60 {
			minute = 0
			hour += 1
		}

		if hour >= 24 {
			hour = 0
		}

		log.Printf("[CURRENT] [%v] hour: %v, minute: %v\n", time.Now().Format(time.RFC3339), hour, minute)

		core.RestakeRewards()

		// Calculate the duration to sleep until the next minute
		now := time.Now()
		nextTick := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
		if nextTick.Before(now) {
			nextTick = nextTick.Add(time.Hour * 24)
		}
		sleepDuration := nextTick.Sub(now)
		log.Printf("Sleeping for %s until the next tick at %s\n", sleepDuration, nextTick.Format(time.RFC3339))

		time.Sleep(sleepDuration)
	}
}
