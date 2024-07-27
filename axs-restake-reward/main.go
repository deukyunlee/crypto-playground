package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/deukyunlee/crypto-playground/core"
	"github.com/deukyunlee/crypto-playground/external"
	"github.com/deukyunlee/crypto-playground/util"
	"github.com/mymmrac/telego"
	"log"
	"math/big"
	"net/http"
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

	if telegramStatus {
		v := util.GetViper()

		telegramToken := v.GetString("telegramToken")
		webHookUrl := v.GetString("webHookUrl")

		telegramBot, err := telego.NewBot(telegramToken)
		if err != nil {
			log.Fatal(err)
		}

		err = telegramBot.SetWebhook(&telego.SetWebhookParams{
			URL: webHookUrl,
		})
		if err != nil {
			log.Fatal(err)
		}

		http.HandleFunc("/webhook", webhookHandler(telegramBot))
		go func() {
			err := http.ListenAndServe(":8080", nil)
			if err != nil {
				log.Printf("err: %s", err)
			}
		}()
	}

	for {
		now := time.Now()
		prevTime, err := time.Parse(time.RFC3339, prevTimeStr)
		if err != nil {
			log.Printf("err: %s", err)
			return
		}

		incrementedTime := util.IncrementTime(prevTime, 1*time.Minute)

		log.Printf("[CURRENT] [%v]\n", now.Format(time.RFC3339))

		util.NextTick = util.CalculateNextTick(now, incrementedTime)
		sleepDuration := util.NextTick.Sub(now)
		log.Printf("Sleeping for %s until the next tick at %s\n", sleepDuration, util.NextTick.Format(time.RFC3339))

		if telegramStatus {
			stakingAmount, err := core.GetStakingAmount()
			if err != nil {
				log.Printf("err: %s", err)
			}

			balance, err := core.GetBalance()
			if err != nil {
				log.Printf("err: %s", err)
			}

			//Formats balance and stakingAmount to 3 decimal places
			message := fmt.Sprintf("*[Next Restaking Time: %s]*\n*[Current Balance]*: %s\n*[Current Staking Amount]*: %s\n",
				util.NextTick.Format(time.RFC3339),
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

func webhookHandler(telegramBot *telego.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var update telego.Update
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "Failed to decode update", http.StatusBadRequest)
			return
		}

		if update.Message != nil {
			// 메시지 처리
			fmt.Printf("Received message from %s: %s\n", update.Message.From.Username, update.Message.Text)

			message := ""

			switch update.Message.Text {
			case "staking":
				stakingAmount, err := core.GetStakingAmount()
				if err != nil {
					log.Printf("err: %s", err)
				}
				message = fmt.Sprintf("*[Current Staking Amount]*: %s", stakingAmount.Text('f', 3))
			case "tick":
				message = fmt.Sprintf("*[Next Restaking Time]*: %s\n %s left", util.NextTick, util.NextTick.Sub(time.Now()))
				break
			case "balance":
				balance, err := core.GetBalance()
				if err != nil {
					log.Printf("err: %s", err)
				}
				message = fmt.Sprintf("*[Current Balance]*: %s", balance.Text('f', 3))
				break
			case "reward":
				stakingAmount, err := core.GetStakingAmount()
				if err != nil {
					log.Printf("err: %s", err)
				}
				totalStaked, err := core.GetTotalStaked()
				if err != nil {
					log.Printf("err: %s", err)
				}

				unlockSchedule := big.NewFloat(1566000)

				now := time.Now()

				currentYear, currentMonth, _ := now.Date()

				firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)
				firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
				daysInMonth := big.NewFloat(firstOfNextMonth.Sub(firstOfMonth).Hours() / 24)

				monthlyTotalReward := new(big.Float).Quo(unlockSchedule, totalStaked)

				userMonthlyReward := monthlyTotalReward.Mul(monthlyTotalReward, stakingAmount)

				userDailyReward := userMonthlyReward.Quo(userMonthlyReward, daysInMonth)
				message = fmt.Sprintf("*[Estimated Daily Reward]*: %s", userDailyReward.Text('f', 3))
			default:
				return
			}

			_, err := telegramBot.SendMessage(&telego.SendMessageParams{
				ChatID:    update.Message.Chat.ChatID(),
				Text:      message,
				ParseMode: "markdown",
			})
			if err != nil {
				log.Println("Failed to send message:", err)
			}
		}
	}
}
