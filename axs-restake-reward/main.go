package main

import (
	"flag"
	"github.com/deukyunlee/crypto-playground/axs-restake-reward/core"
	"github.com/deukyunlee/crypto-playground/axs-restake-reward/logging"
	"github.com/deukyunlee/crypto-playground/axs-restake-reward/notification"
	"github.com/deukyunlee/crypto-playground/axs-restake-reward/util"
	"time"
)

var (
	logger         = logging.GetLogger()
	telegramStatus bool
)

func init() {
	initializeFlags()
	util.SetTimezone()
	if telegramStatus {
		notification.InitializeTelegramBot()
	}
}

func initializeFlags() {
	flag.BoolVar(&telegramStatus, "notification", false, "notification Status | true: send message using notification, false: not sending any notification message")
	flag.StringVar(&util.Timezone, "timezone", "Local", "Timezone to display the time (e.g., 'UTC', 'America/New_York')")
	flag.Parse()
}

func main() {
	var coreManager core.CoreManager = &core.EvmManager{}

	lastClaimedTime := coreManager.GetLastClaimedTime()

	for {
		now := time.Now().UTC()
		logger.Infof("[CURRENT] [%s]\n", now.Format(time.RFC3339))

		util.NextTick, util.Duration = util.CalculateNextTick(now, lastClaimedTime)
		logger.Infof("Sleeping for %s until the next tick at %s\n", util.Duration, util.NextTick.In(util.Location).Format(time.RFC3339))

		if telegramStatus {
			notification.CreatePeriodicalTelegramMessage()
		}

		time.Sleep(util.Duration)
		lastClaimedTime = coreManager.GetLastClaimedTime()
	}
}
