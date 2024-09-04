package handler

import (
	"github.com/mymmrac/telego"
)

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
		logger.Errorf("Failed to answer callback query: %s", err)
	}
}
