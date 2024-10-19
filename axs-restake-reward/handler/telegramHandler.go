package handler

import (
	"fmt"
	"github.com/deukyunlee/crypto-playground/core"
	"github.com/deukyunlee/crypto-playground/util"
	"github.com/mymmrac/telego"
)

type TelegramHandler struct {
	evmManager *core.EvmManager
}

func NewTelegramHandler(evmManager *core.EvmManager) *TelegramHandler {
	return &TelegramHandler{
		evmManager: evmManager,
	}
}

func (h *TelegramHandler) HandleMessage(telegramBot *telego.Bot, message *telego.Message) {
	reply := ""
	pk := util.GetConfigInfo().PK
	accountAddress := util.GetAddressFromPrivateKey(pk)

	switch message.Text {
	case "staking":
		stakingAmount, err := h.evmManager.GetStakingAmount(accountAddress)
		if err != nil {
			reply = fmt.Sprintf("Error fetching staking amount: %s", err)
		} else {
			reply = fmt.Sprintf("*[Current Staking Amount]*: %s", stakingAmount.Text('f', 3))
		}
	case "balance":
		balance, err := h.evmManager.GetBalance(accountAddress)
		if err != nil {
			reply = fmt.Sprintf("Error fetching balance: %s", err)
		} else {
			reply = fmt.Sprintf("*[Current Balance]*: %s", balance.Text('f', 3))
		}
	default:
		reply = "Invalid command"
	}

	telegramBot.SendMessage(&telego.SendMessageParams{
		ChatID:    message.Chat.ChatID(),
		Text:      reply,
		ParseMode: "markdown",
	})
}
