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
		log.Printf("err: %s", err)
	}
	defer closeLogFile(logFile)

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	prevTimePtr := flag.String("time", "2006-01-02T15:04:05Z", "previously restaked time(RFC3339) | e.g. 2024-07-21T22:58:16+09:00")
	telegramStatusPtr := flag.Bool("telegram", false, "telegram Status | true: sendMessage")

	flag.Parse()

	prevTimeStr := *prevTimePtr
	telegramStatus := *telegramStatusPtr

	log.Printf("[FLAG] hour: %s, telegramStatus: %t", prevTimeStr, telegramStatus)

	for {
		now := time.Now()
		prevTime, err := time.Parse(time.RFC3339, prevTimeStr)
		if err != nil {
			log.Printf("err: %s", err)
			return
		}

		incrementedTime := util.IncrementTime(prevTime, 1*time.Minute)

		log.Printf("[CURRENT] [%v]\n", now.Format(time.RFC3339))

		nextTick := util.CalculateNextTick(now, incrementedTime)
		sleepDuration := nextTick.Sub(now)
		log.Printf("Sleeping for %s until the next tick at %s\n", sleepDuration, nextTick.Format(time.RFC3339))

		if telegramStatus {
			ctx := context.Background()
			stakingAmount, err := core.GetStakingAmount()
			if err != nil {
				log.Printf("err: %s", err)
			}

			balance, err := core.GetBalance(ctx)
			if err != nil {
				log.Printf("err: %s", err)
			}

			//Formats balance and stakingAmount to 3 decimal places
			message := fmt.Sprintf("*[Next: %s]*\n*balance*: %s\n*stakingAmount*: %s\n",
				nextTick.Format(time.RFC3339),
				balance.Text('f', 3),
				stakingAmount.Text('f', 3),
			)

			external.SendTelegramRestakeMessage(message)
		}

		time.Sleep(sleepDuration)

		core.RestakeRewards()

		prevTimeStr = now.Format(time.RFC3339)
	}
}

func closeLogFile(logFile *os.File) {
	if err := logFile.Close(); err != nil {
		log.Printf("err: %s", err)
	}
}
