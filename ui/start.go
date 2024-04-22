package ui

import (
	"dexbot/model"
	"dexbot/service"
	"dexbot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func StartHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	addBtn := model.ButtonInfo{
		Text:    "+ 添加【Dex机器人】为好友 +",
		Data:    fmt.Sprintf("https://t.me/%s?startgroup=top", bot.Self.UserName),
		BtnType: model.BtnTypeUrl,
	}

	supportBtn1 := model.ButtonInfo{
		Text:    "🔈 官方群组",
		Data:    "https://t.me/+-tG2kxFzfuljMjU1",
		BtnType: model.BtnTypeUrl,
	}

	supportBtn2 := model.ButtonInfo{
		Text:    "📞 客服",
		Data:    "https://user?id=6401399435",
		BtnType: model.BtnTypeData,
	}
	walletBtn := model.ButtonInfo{
		Text:    "👉🏦进入Dex",
		Data:    "home",
		BtnType: model.BtnTypeData,
	}

	var rows [][]model.ButtonInfo
	addRow := []model.ButtonInfo{addBtn}
	supportRow := []model.ButtonInfo{supportBtn1, supportBtn2}
	walletRow := []model.ButtonInfo{walletBtn}
	rows = append(rows, addRow)
	rows = append(rows, supportRow)
	rows = append(rows, walletRow)

	keyboard := utils.MakeKeyboard(rows)
	content := fmt.Sprintf("👏 欢迎使用【%s】进行挂单、跟单、买卖 \n\n请点击下方【+添加】按钮：\n "+
		"•  在机器人私聊中发送 /home 打开钱包主菜单。\n\n"+
		"", bot.Self.FirstName)
	utils.SendMenu(update.Message.Chat.ID, content, keyboard, bot)

}

func ManagerMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatId := update.Message.Chat.ID

	info := model.GroupInfo{
		GroupId:   chatId,
		Uid:       update.Message.From.ID,
		GroupName: update.Message.Chat.Title,
		GroupType: update.Message.Chat.Type,
	}
	//保存到数据库
	service.GetDB().Save(&info)

	//更新本地变量
	utils.GroupInfo = info

	content := fmt.Sprintf("欢迎使用 @%s：\n1)点击下面按钮选择设置(仅限管理员)\n2)点击机器人对话框底部【开始】按钮\n\n🟩 功能更新提醒：在机器人私聊中发送 /start 也可打开管理菜单\n", bot.Self.UserName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	url := fmt.Sprintf("https://t.me/%s?start=manager_%d", bot.Self.UserName, utils.GroupInfo.GroupId)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("👉⚙️进入Dex主页面👈", url),
		))
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
