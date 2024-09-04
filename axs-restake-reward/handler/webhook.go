package handler

import (
	"encoding/json"
	"github.com/mymmrac/telego"
	"net/http"
)

func HandleWebhook(telegramBot *telego.Bot) http.HandlerFunc {
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
