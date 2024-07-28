package external

import (
	"fmt"
	"github.com/deukyunlee/crypto-playground/util"
	"github.com/mymmrac/telego"
	"os"
)

func SendTelegramRestakeMessage(message string) {
	configInfo := util.GetConfigInfo()

	telegramBotToken := configInfo.Telegram.Token
	telegramChatId := configInfo.Telegram.ChatID
	telegramUserName := configInfo.Telegram.UserName

	telegramBot, err := telego.NewBot(telegramBotToken)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Send the message
	_, err = telegramBot.SendMessage(&telego.SendMessageParams{
		ChatID:    telego.ChatID{ID: telegramChatId, Username: telegramUserName},
		Text:      message,
		ParseMode: "markdown",
	})
	if err != nil {
		fmt.Println("Failed to send message:", err)
		os.Exit(1)
	}

	fmt.Println("Message sent successfully!")

}
