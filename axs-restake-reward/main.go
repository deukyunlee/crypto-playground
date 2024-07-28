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

var loc *time.Location

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

	timezone := flag.String("timezone", "Local", "Timezone to display the time (e.g., 'UTC', 'America/New_York')")
	telegramStatusPtr := flag.Bool("telegram", false, "telegram Status | true: sendMessage")
	flag.Parse()

	if *timezone == "Local" {
		loc = time.Local
	} else {
		loc, err = time.LoadLocation(*timezone)
		if err != nil {
			log.Fatalf("Error loading timezone: %v", err)
		}
	}

	telegramStatus := *telegramStatusPtr

	log.Printf("[FLAG]loc: %s, telegramStatus: %t", loc, telegramStatus)

	var latestTxHash = ""

	const roninExplorerUri = "https://app.roninchain.com"

	if telegramStatus {
		configInfo := util.GetConfigInfo()

		telegramToken := configInfo.Telegram.Token
		webHookUrl := configInfo.Telegram.WebHookUrl

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

	userRewardInfo, err := core.GetUserRewardInfo()
	if err != nil {
		log.Printf("err: %s", err)
	}

	lastClaimedTimestampUnix := userRewardInfo.LastClaimedBlock.Int64()
	lastClaimedTime := time.Unix(lastClaimedTimestampUnix, 0).UTC()

	log.Printf("lastClaimedTime: %s\n", lastClaimedTime.In(loc))

	for {
		now := time.Now().UTC()

		log.Printf("[CURRENT] [%s]\n", now.UTC().Format(time.RFC3339))

		util.NextTick, util.Duration = util.CalculateNextTick(now, lastClaimedTime)

		log.Printf("Sleeping for %s until the next tick at %s\n", util.Duration, util.NextTick.In(loc).Format(time.RFC3339))

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
			message := fmt.Sprintf("*[Next Restaking Time: %s]*\n*[Current Balance]*: %s\n*[Current Staking Amount]*: %s\n*[Latest Tx]*\n: %s\n",
				util.NextTick.In(loc).Format(time.RFC3339),
				balance.Text('f', 3),
				stakingAmount.Text('f', 3),
				roninExplorerUri+"/tx/"+latestTxHash,
			)

			external.SendTelegramRestakeMessage(message)
		}

		time.Sleep(util.Duration)

		latestTxHash = core.RestakeRewards()

		lastClaimedTime = now.UTC()
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
			handleMessage(telegramBot, update.Message)
		} else if update.CallbackQuery != nil {
			handleCallbackQuery(telegramBot, update.CallbackQuery)
		}
	}
}

func handleMessage(telegramBot *telego.Bot, message *telego.Message) {
	log.Printf("Received message from %s: %s\n", message.From.Username, message.Text)

	reply := ""

	switch message.Text {
	case "staking":
		stakingAmount, err := core.GetStakingAmount()
		if err != nil {
			log.Printf("err: %s", err)
		}
		reply = fmt.Sprintf("*[Current Staking Amount]*: %s", stakingAmount.Text('f', 3))
	case "tick":
		reply = fmt.Sprintf("*[Next Restaking Time]*: %s\n %s left", util.NextTick.In(loc), util.NextTick.Sub(time.Now()))
	case "balance":
		balance, err := core.GetBalance()
		if err != nil {
			log.Printf("err: %s", err)
		}
		reply = fmt.Sprintf("*[Current Balance]*: %s", balance.Text('f', 3))
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
		userMonthlyReward := new(big.Float).Mul(monthlyTotalReward, stakingAmount)
		userDailyReward := new(big.Float).Quo(userMonthlyReward, daysInMonth)

		reply = fmt.Sprintf("*[Estimated Daily Reward]*: %s", userDailyReward.Text('f', 3))
	default:
		return
	}

	_, err := telegramBot.SendMessage(&telego.SendMessageParams{
		ChatID:    message.Chat.ChatID(),
		Text:      reply,
		ParseMode: "markdown",
	})
	if err != nil {
		log.Println("Failed to send message:", err)
	}

	// Create inline keyboard
	inlineKeyboard := telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "Staking", CallbackData: "staking"},
				{Text: "Tick", CallbackData: "tick"},
				{Text: "Balance", CallbackData: "balance"},
				{Text: "Reward", CallbackData: "reward"},
			},
		},
	}

	_, err = telegramBot.SendMessage(&telego.SendMessageParams{
		ChatID:      message.Chat.ChatID(),
		Text:        "Choose an option:",
		ParseMode:   "markdown",
		ReplyMarkup: &inlineKeyboard,
	})
	if err != nil {
		log.Fatalf("Failed to send message with inline keyboard: %s", err)
	}
}

func handleCallbackQuery(telegramBot *telego.Bot, callbackQuery *telego.CallbackQuery) {
	chat := callbackQuery.Message.GetChat()
	chatID := chat.ChatID().ID

	callbackData := callbackQuery.Data

	update := telego.Update{
		Message: &telego.Message{
			Chat: telego.Chat{
				ID: chatID,
			},
			From: &callbackQuery.From,
			Text: callbackData,
		},
	}

	handleMessage(telegramBot, update.Message)

	err := telegramBot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
		Text:            "Processing...",
	})
	if err != nil {
		log.Printf("Failed to answer callback query: %s", err)
	}
}
