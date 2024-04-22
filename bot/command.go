package bot

import (
	"dexbot/ui"
	"dexbot/utils"
	"fmt"
	"strings"
)

func (bot *smartBot) handleCommand() {
	cmd := bot.updateMessage.Message.CommandWithAt()
	fmt.Println("command:", cmd)

	if strings.HasPrefix(cmd, "start") {
		//判断是否是私聊
		if bot.updateMessage.Message.Chat.Type == "private" {

			//接收参数，取空格后面的内容
			args := strings.TrimSpace(strings.Replace(bot.updateMessage.Message.Text, "/start", "", -1))
			fmt.Println("pr:", args)
			if len(args) == 0 {
				ui.StartHandler(bot.bot, bot.updateMessage)
				return
			}
			//分割参数
			params := strings.Split(args, "_")
			//根据参数获取群组信息
			module := params[0]
			if module == "redp" { //发红包
				//menu.RedpUsdtMenu(bot.bot, bot.updateMessage)
			}
		} else {
			fmt.Println("daf")
			//如果是管理员	弹出管理菜单
			member, _ := utils.GetMemberInfo(bot.updateMessage.Message.Chat.ID, bot.updateMessage.Message.From.ID, bot.bot)
			if member.IsAdministrator() || member.IsCreator() {
				ui.ManagerMenu(bot.updateMessage, bot.bot)
			} else {
				utils.SendText(bot.updateMessage.Message.Chat.ID, "此功能只针对管理员开放", bot.bot)
			}
		}
	} else if strings.HasPrefix(cmd, "home") {
		ui.HomeMenu(bot.bot, bot.updateMessage)
	}
}
