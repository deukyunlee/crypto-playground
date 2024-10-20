package notification

import (
	"fmt"
	"github.com/deukyunlee/crypto-playground/axs-restake-reward/core"
	"github.com/deukyunlee/crypto-playground/axs-restake-reward/handler"
	"github.com/deukyunlee/crypto-playground/axs-restake-reward/logging"
	"github.com/deukyunlee/crypto-playground/axs-restake-reward/util"
	"github.com/mymmrac/telego"
	"net/http"
	"os"
	"time"
)

var (
	logger = logging.GetLogger()
)

const roninExplorerUri = "https://app.roninchain.com"

func CreatePeriodicalTelegramMessage() {

	coreManager := &core.EvmManager{}
	pk := util.GetConfigInfo().PK
	accountAddress := util.GetAddressFromPrivateKey(pk)

	latestTxHash, err := coreManager.AutoCompoundRewards()

	stakingAmount, err := coreManager.GetStakingAmount(accountAddress)
	if err != nil {
		logger.Errorf("err: %s", err)
	}

	balance, err := coreManager.GetBalance(accountAddress)
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

	sendTelegramAutoCompoundMessage(message)
}
func sendTelegramAutoCompoundMessage(message string) {
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
		logger.Errorf("Failed to send message: %s", err)
		os.Exit(1)
	}

	logger.Info("Message sent successfully!")
}

func InitializeTelegramBot() {
	configInfo := util.GetConfigInfo()
	telegramToken := configInfo.Telegram.Token
	webHookUrl := configInfo.Telegram.WebHookUrl

	telegramBot, err := telego.NewBot(telegramToken)
	if err != nil {
		logger.Error(err)
	}

	err = telegramBot.SetWebhook(&telego.SetWebhookParams{URL: webHookUrl})
	if err != nil {
		logger.Error(err)
	}

	http.HandleFunc("/webhook", handler.HandleWebhook(telegramBot))
	go handler.StartWebhookServer()
}
