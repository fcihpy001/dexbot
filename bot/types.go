package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type smartBot struct {
	token         string
	debug         bool
	bot           *tgbotapi.BotAPI
	updateMessage *tgbotapi.Update
}
