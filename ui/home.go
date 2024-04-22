package ui

import (
	"dexbot/model"
	"dexbot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HomeMenu(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {

	var (
		chatId   int64
		chatType string
	)
	if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.Message.Chat.ID
		chatType = update.CallbackQuery.Message.Chat.Type
	} else if update.Message != nil {
		chatId = update.Message.Chat.ID
		chatType = update.Message.Chat.Type
	}

	var buttons []model.ButtonInfo
	utils.Json2Button("./resource/home.json", &buttons)

	var row []model.ButtonInfo
	var rows [][]model.ButtonInfo
	for i := 1; i <= len(buttons); i++ {

		row = append(row, buttons[i-1])
		if i%2 == 0 && i != 0 {
			rows = append(rows, row)
			row = []model.ButtonInfo{}
		}
	}

	fmt.Println(chatType)
	if len(buttons)%2 != 0 {
		rows = append(rows, row)
	}

	keyboard := utils.MakeKeyboard(rows)
	utils.HomeMenuMarkup = keyboard

	content := updateHomeMsg()
	utils.SendMenu(chatId, content, keyboard, bot)
}

func GoHome(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	content := updateHomeMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.HomeMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		return
	}
}

func updateHomeMsg() string {
	content := fmt.Sprintf("👏 欢迎使用【walletBot】\n\n" +
		"--------💸 💵 💴 💶 💷 🪙------- \n" +
		"钱包地址：6bDc2pNbe6pCY8c1926DJ1Ef1WazB5m9DEAopuxPKFZV\n钱包余额：0SOL\n")

	return content
}
