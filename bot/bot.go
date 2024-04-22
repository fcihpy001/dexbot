package bot

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
)

func StartBot(ctx context.Context) {
	bot, err := createBot(os.Getenv("BOT_TOKEN"), os.Getenv("BOT_DEBUG") == "true")
	if err != nil {
		panic(err)
	}
	log.Printf("Authorized on account %s--%d-%s", bot.bot.Self.UserName, bot.bot.Self.ID, bot.bot.Self.FirstName)

	go bot.setupBotWithPool()
}

func (bot *smartBot) setupBotWithPool() {

	updateConfig := tgbotapi.NewUpdate(0)
	timeout := 1800
	if os.Getenv("POLL_TIMEOUT") != "" {
		if val, err := strconv.Atoi(os.Getenv("POLL_TIMEOUT")); err == nil {
			timeout = val
		}
	}
	updateConfig.Timeout = timeout
	updatesChannel := bot.bot.GetUpdatesChan(updateConfig)

	//定时删除消息
	timer10DeleteTask(bot.bot)

	// 机器人与用户的交互逻辑
	for update := range updatesChannel {
		s, _ := json.Marshal(update)
		log.Println("update:", string(s))
		bot.updateMessage = &update

		if update.Message != nil && update.Message.IsCommand() { // 以/开头的指令消息
			bot.handleCommand()

		} else if update.CallbackQuery != nil { // 按钮回调
			bot.handleQuery()

		} else if update.Message != nil && update.Message.ReplyToMessage != nil { // 要求用户回复的消息
			bot.handleReply(&update)

		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.bot.Request(callback); err != nil {
				panic(err)
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.bot.Send(msg); err != nil {
				panic(err)
			}

		} else {
			if update.Message != nil && update.Message.Chat != nil { // 未定义消息的处理
				// if chat is nil, panic
				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "这个问题，暂时无法处理")
				//bot.
				//bot.bot.SendText(update.Message.Chat.ID, "这个问题，暂时无法处理")
			}
			bot.handleMessage()
		}
		//callback := tgbotapi.NewCallback(update.Message.MessageID, "loading")
		//_, err := bot.bot.Send(callback)
		//if err != nil {
		//	return
		//}
	}
}

func createBot(token string, debug bool) (*smartBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = debug

	return &smartBot{
		token: token,
		debug: debug,
		bot:   bot,
	}, nil
}

//func (bot *SmartBot) SendText(chatId int64, text string) {
//	msg := tgbotapi.NewMessage(chatId, text)
//	_, err := bot.bot.Send(msg)
//	if err != nil {
//		log.Println(err)
//	}
//}

// 定时删除任务
func timer10DeleteTask(bot *tgbotapi.BotAPI) {
	//ticker := time.NewTicker(10 * time.Second)
	// 使用 Goroutine 执行任务
	go func() {
		for {
			select {
			//case <-ticker.C:
			//	// 在这里执行您的定时任务代码
			//	deleteMessageTask(bot)
			}
		}
	}()
}

// 定时删除消息
//func deleteMessageTask(bot *tgbotapi.BotAPI) {
//	//获取所有删除任务
//	var tasks []model.ScheduleDelete
//	tenSecondsAgo := time.Now().Add(-10 * time.Second)
//	service.GetDB().Where("created_at < ?", tenSecondsAgo).Find(&tasks)
//	//循环tasks
//	for _, task := range tasks {
//		msg := tgbotapi.NewDeleteMessage(task.ChatId, task.MessageId)
//		mm, err := bot.Send(msg)
//		if err != nil {
//			fmt.Println("err", err)
//		}
//		fmt.Println("删除消息成功,更新删除记录", mm.MessageID)
//		_ = service.GetDB().Delete(&task)
//	}
//}
