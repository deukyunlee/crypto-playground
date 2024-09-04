package handler

import (
	"fmt"
	"github.com/deukyunlee/crypto-playground/core"
	"github.com/deukyunlee/crypto-playground/logging"
	"github.com/deukyunlee/crypto-playground/util"
	"github.com/mymmrac/telego"
	"math/big"
	"time"
)

var (
	logger = logging.GetLogger()
)

func handleMessage(telegramBot *telego.Bot, message *telego.Message) {
	logger.Infof("Received message from %s: %s\n", message.From.Username, message.Text)

	reply := ""

	switch message.Text {
	case "staking":
		stakingAmount, err := core.GetStakingAmount()
		if err != nil {
			logger.Errorf("err: %s", err)
		}
		reply = fmt.Sprintf("*[Current Staking Amount]*: %s", stakingAmount.Text('f', 3))
	case "tick":
		reply = fmt.Sprintf("*[Next Restaking Time]*: %s\n %s left", util.NextTick.In(util.Location), util.NextTick.Sub(time.Now()))
	case "balance":
		balance, err := core.GetBalance()
		if err != nil {
			logger.Errorf("err: %s", err)
		}
		reply = fmt.Sprintf("*[Current Balance]*: %s", balance.Text('f', 3))
	case "reward":
		stakingAmount, err := core.GetStakingAmount()
		if err != nil {
			logger.Errorf("err: %s", err)
		}
		totalStaked, err := core.GetTotalStaked()
		if err != nil {
			logger.Errorf("err: %s", err)
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
		logger.Errorf("Failed to send message: %s", err)
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
		logger.Errorf("Failed to send message with inline keyboard: %s", err)
	}
}
