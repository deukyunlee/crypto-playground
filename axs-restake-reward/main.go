package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/deukyunlee/crypto-playground/core"
	"github.com/deukyunlee/crypto-playground/external"
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
	telegramStatusPtr := flag.Bool("telegram", false, "telegram Status | true: sendMessage")

	flag.Parse()

	hour := *hourPtr
	minute := *minutePtr
	telegramStatus := *telegramStatusPtr

	log.Printf("[FLAG] hour: %d, minute: %d, telegramStatus: %t", hour, minute, telegramStatus)

	for {
		now := time.Now()
		hour, minute = util.IncrementTime(hour, minute)
		log.Printf("[CURRENT] [%v] hour: %v, minute: %v\n", now.Format(time.RFC3339), hour, minute)

		nextTick := util.CalculateNextTick(now, hour, minute)
		sleepDuration := nextTick.Sub(now)
		log.Printf("Sleeping for %s until the next tick at %s\n", sleepDuration, nextTick.Format(time.RFC3339))

		if telegramStatus {
			ctx := context.Background()
			stakingAmount, err := core.GetStakingAmount(ctx)
			if err != nil {
				log.Fatal(err)
			}

			balance, err := core.GetBalance(ctx)
			if err != nil {
				log.Fatal(err)
			}

			// Formats balance and stakingAmount to 3 decimal places
			message := fmt.Sprintf("*[Next: %s]*\n*balance*: %s\n*stakingAmount*: %s\n",
				nextTick.Format(time.RFC3339),
				balance.Text('f', 3),
				stakingAmount.Text('f', 3),
			)

			external.SendTelegramRestakeMessage(message)
		}

		time.Sleep(sleepDuration)

		core.RestakeRewards()

	}
}

func closeLogFile(logFile *os.File) {
	if err := logFile.Close(); err != nil {
		log.Fatal(err)
	}
}
