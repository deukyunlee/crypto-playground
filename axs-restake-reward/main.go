package main

import (
	"flag"
	"github.com/deukyunlee/crypto-playground/core"
	"github.com/deukyunlee/crypto-playground/util"
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
		hour, minute = util.IncrementTime(hour, minute)
		log.Printf("[CURRENT] [%v] hour: %v, minute: %v\n", now.Format(time.RFC3339), hour, minute)

		nextTick := util.CalculateNextTick(now, hour, minute)
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
