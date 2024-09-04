package main

import (
	"flag"
	"github.com/deukyunlee/crypto-playground/core"
	"github.com/deukyunlee/crypto-playground/handler"
	"github.com/deukyunlee/crypto-playground/logging"
	"github.com/deukyunlee/crypto-playground/notification"
	"github.com/deukyunlee/crypto-playground/util"
	"github.com/mymmrac/telego"
	"net/http"
	"time"
)

var (
	logger         = logging.GetLogger()
	telegramStatus bool
	timezone       string
)

func init() {
	flag.BoolVar(&telegramStatus, "notification", false, "notification Status | true: send message using notification, false: not sending any notification message")
	flag.StringVar(&timezone, "timezone", "Local", "Timezone to display the time (e.g., 'UTC', 'America/New_York')")
	flag.Parse()

	var err error
	if timezone == "Local" {
		util.Location = time.Local
	} else {
		util.Location, err = time.LoadLocation(timezone)
		if err != nil {
			logger.Errorf("Error loading timezone: %v", err)
		}
	}

	if telegramStatus {
		configInfo := util.GetConfigInfo()

		telegramToken := configInfo.Telegram.Token
		webHookUrl := configInfo.Telegram.WebHookUrl

		telegramBot, err := telego.NewBot(telegramToken)
		if err != nil {
			logger.Error(err)
		}

		err = telegramBot.SetWebhook(&telego.SetWebhookParams{
			URL: webHookUrl,
		})
		if err != nil {
			logger.Error(err)
		}

		http.HandleFunc("/webhook", handler.HandleWebhook(telegramBot))
		go func() {
			err := http.ListenAndServe(":8080", nil)
			if err != nil {
				logger.Errorf("err: %s", err)
			}
		}()
	}

}

func main() {
	userRewardInfo, err := core.GetUserRewardInfo()
	if err != nil {
		logger.Errorf("err: %s", err)
	}

	lastClaimedTimestampUnix := userRewardInfo.LastClaimedBlock.Int64()
	lastClaimedTime := time.Unix(lastClaimedTimestampUnix, 0).UTC()

	logger.Info("lastClaimedTime: %s\n", lastClaimedTime.In(util.Location))

	for {
		now := time.Now().UTC()

		logger.Errorf("[CURRENT] [%s]\n", now.UTC().Format(time.RFC3339))

		util.NextTick, util.Duration = util.CalculateNextTick(now, lastClaimedTime)

		logger.Infof("Sleeping for %s until the next tick at %s\n", util.Duration, util.NextTick.In(util.Location).Format(time.RFC3339))

		if telegramStatus {
			notification.CreatePeriodicalTelegramMessage()
		}

		time.Sleep(util.Duration)

		userRewardInfo, err := core.GetUserRewardInfo()
		if err != nil {
			logger.Errorf("err: %s", err)
		}

		lastClaimedTimestampUnix = userRewardInfo.LastClaimedBlock.Int64()
		lastClaimedTime = time.Unix(lastClaimedTimestampUnix, 0).UTC()
	}
}
