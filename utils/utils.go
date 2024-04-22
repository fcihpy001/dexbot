package utils

import (
	"dexbot/model"
	"dexbot/service"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"
)

var (
	HomeMenuMarkup      tgbotapi.InlineKeyboardMarkup
	WalletMenuMarkup    tgbotapi.InlineKeyboardMarkup
	RedpacketMenuMarkup tgbotapi.InlineKeyboardMarkup
	GroupInfo           model.GroupInfo
)

func MakeKeyboard(btns [][]model.ButtonInfo) tgbotapi.InlineKeyboardMarkup {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup()
	for i := 0; i < len(btns); i++ {
		row := tgbotapi.NewInlineKeyboardRow()
		for j := 0; j < len(btns[i]); j++ {
			fmt.Println("type:", btns[i][j].BtnType)
			fmt.Println("data:", btns[i][j].Data)
			btn := tgbotapi.NewInlineKeyboardButtonData(btns[i][j].Text, btns[i][j].Data)
			if btns[i][j].BtnType == model.BtnTypeUrl {
				btn = tgbotapi.NewInlineKeyboardButtonURL(btns[i][j].Text, btns[i][j].Data)

			} else if btns[i][j].BtnType == model.BtnTypeSwitch {
				btn = tgbotapi.NewInlineKeyboardButtonSwitch(btns[i][j].Text, btns[i][j].Data)
			}
			row = append(row, btn)
		}
		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, row)
	}
	return inlineKeyboard
}

func SendEditMsgMarkup(
	chatID int64,
	messageID int,
	content string,
	replyMarkup tgbotapi.InlineKeyboardMarkup,
	bot *tgbotapi.BotAPI,
) {
	msg := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, content, replyMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println("editMenuKeyboard", err)
		SendText(chatID, content, bot)
	}

}

func SendMenu(receiver int64, msg string, keybord tgbotapi.InlineKeyboardMarkup, bot *tgbotapi.BotAPI) {
	message := tgbotapi.NewMessage(receiver, msg)
	message.ReplyMarkup = keybord
	_, err := bot.Send(message)
	if err != nil {
		log.Println(err)
	}
}

func Json2Button(file string, models *[]model.ButtonInfo) {
	path, err := os.ReadFile(file)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(path, &models)
	if err != nil {
		log.Panic(err)
	}
}

func Json2Button2(file string, models *[][]model.ButtonInfo) {
	path, err := os.ReadFile(file)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(path, &models)
	if err != nil {
		log.Panic(err)
	}
}

func SendText(chatId int64, text string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatId, text)
	message, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
	CreateDeleteMessage(chatId, message.MessageID)
}

func CreateDeleteMessage(chatId int64, messageId int) {
	task := model.ScheduleDelete{
		ChatId:    chatId,
		MessageId: messageId,
	}
	service.GetDB().Save(&task)
}

func GetMemberInfo(chat_id int64, user_id int64, bot *tgbotapi.BotAPI) (tgbotapi.ChatMember, error) {
	req := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chat_id,
			UserID: user_id,
		},
	}
	return bot.GetChatMember(req)
}

func GetUserIcon(userId int64, bot *tgbotapi.BotAPI) string {
	resp, err := bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{
		UserID: userId,
		Offset: 0,
		Limit:  5,
	})
	if err != nil {
		fmt.Println("获取用户")
	}
	if resp.TotalCount < 1 {
		return ""
	}
	return resp.Photos[0][0].FileID
}

func parseFileId(fileId string, bot *tgbotapi.BotAPI) {
	downloadedFile, err := bot.GetFileDirectURL(fileId)
	if err != nil {
		log.Printf("Failed to get direct URL: %v", err)

	}

	err = downloadFile(downloadedFile, "red.png")
	if err != nil {
		log.Printf("Failed to download file: %v", err)
	}
}

// 下载文件到本地
func downloadFile(url, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// 验证以太坊私钥是否有效
//func PrivateToAddress(privateKeyHex string) string {
//	// 将私钥字符串解析为私钥对象
//	privateKey, err := crypto.HexToECDSA(privateKeyHex)
//	if err != nil {
//		fmt.Println("Invalid private key:", err)
//		return ""
//	}
//	if privateKey == nil {
//		return ""
//	}
//
//	// 验证私钥的 D 值是否在椭圆曲线的阶内
//	curveOrder := new(big.Int).Set(privateKey.Curve.Params().N)
//	if privateKey.D.Cmp(curveOrder) >= 0 {
//		return ""
//	}
//	publicKey := privateKey.Public()
//	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
//	if !ok {
//		return ""
//	}
//	return crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
//}

func ValidateWalletAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}

func GetRandNo(total float64, n int) float64 {
	if total <= 0 || n <= 0 {
		return 0
	}
	if n == 1 {
		return total
	}

	rand.Seed(time.Now().UnixNano())
	randomFloats := make([]float64, n)
	remainingValue := total

	for i := 0; i < n-1; i++ {
		random := rand.Float64() * remainingValue
		randomFloats[i] = random
		remainingValue -= random
	}
	randomFloats[n-1] = remainingValue
	return randomFloats[0]
}

func SendForceReplyMsg(chatId int64, msg string, bot *tgbotapi.BotAPI) {

	msgCfg := tgbotapi.NewMessage(chatId, msg)
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BACK_TEXT),
		))
	msgCfg.ReplyMarkup = keyboard
	msgCfg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
	}
	message, err := bot.Send(msgCfg)
	if err != nil {
		return
	}
	CreateDeleteMessage(chatId, message.MessageID)
}

func SendOkMsgMenu(chatId int64, bot *tgbotapi.BotAPI, content string) {

	btn1 := model.ButtonInfo{
		Text:    BACK_TEXT,
		Data:    BACK_DATA,
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := MakeKeyboard(rows)

	msg := tgbotapi.NewMessage(chatId, OP_SUCCESS)
	if len(content) > 0 {
		msg = tgbotapi.NewMessage(chatId, content)
	}
	msg.ReplyMarkup = keyboard
	message, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
	CreateDeleteMessage(chatId, message.MessageID)
}

func MakeBackkeyboard() tgbotapi.InlineKeyboardMarkup {
	backBtn := model.ButtonInfo{
		Text:    BACK_TEXT,
		Data:    BACK_DATA,
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{backBtn}
	rows := [][]model.ButtonInfo{row1}
	keyboard := MakeKeyboard(rows)
	return keyboard
}
