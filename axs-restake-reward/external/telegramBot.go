package external

import (
	"fmt"
	"github.com/deukyunlee/crypto-playground/util"
	"github.com/mymmrac/telego"
	"os"
)

func SendTelegramRestakeMessage(message string) {
	// TODO: 20:00
	v := util.GetViper()

	telegramBotToken := v.GetString("telegramToken")
	telegramChatId := v.GetInt64("telegramChatId")
	telegramUserName := v.GetString("telegramUserName")

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
