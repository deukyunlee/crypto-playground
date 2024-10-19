package notification

import (
	"fmt"
	"github.com/deukyunlee/crypto-playground/core"
	"github.com/deukyunlee/crypto-playground/logging"
	"github.com/deukyunlee/crypto-playground/util"
	"github.com/mymmrac/telego"
	"os"
	"time"
)

var (
	latestTxHash string
	logger       = logging.GetLogger()
)

const roninExplorerUri = "https://app.roninchain.com"

func CreatePeriodicalTelegramMessage() {
	coreManager := &core.EvmManager{}

	latestTxHash = core.AutoCompoundRewards()

	stakingAmount, err := coreManager.GetStakingAmount()
	if err != nil {
		logger.Errorf("err: %s", err)
	}

	balance, err := coreManager.GetBalance()
	if err != nil {
		logger.Errorf("err: %s", err)
	}

	//Formats balance and stakingAmount to 3 decimal places
	message := fmt.Sprintf("*[Next Restaking Time: %s]*\n*[Current Balance]*: %s\n*[Current Staking Amount]*: %s\n*[Latest Tx]*\n: %s\n",
		util.NextTick.In(util.Location).Format(time.RFC3339),
		balance.Text('f', 3),
		stakingAmount.Text('f', 3),
		roninExplorerUri+"/tx/"+latestTxHash,
	)

	SendTelegramAutoCompoundMessage(message)
}
func SendTelegramAutoCompoundMessage(message string) {
	configInfo := util.GetConfigInfo()

	telegramBotToken := configInfo.Telegram.Token
	telegramChatId := configInfo.Telegram.ChatID
	telegramUserName := configInfo.Telegram.UserName

	telegramBot, err := telego.NewBot(telegramBotToken)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	// Send the message
	_, err = telegramBot.SendMessage(&telego.SendMessageParams{
		ChatID:    telego.ChatID{ID: telegramChatId, Username: telegramUserName},
		Text:      message,
		ParseMode: "markdown",
	})
	if err != nil {
		logger.Errorf("Failed to send message:", err)
		os.Exit(1)
	}

	logger.Info("Message sent successfully!")
}
