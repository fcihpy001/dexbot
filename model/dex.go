package model

import "gorm.io/gorm"

type WalletHistory struct {
	gorm.Model
	UserId       int64
	Address      string
	PrivateKey   string
	WalletSource string
}

type User struct {
	gorm.Model
	UserId       int64 `gorm:"uniqueIndex"`
	Name         string
	Icon         string
	Address      string `gorm:"address;type:char(42)"`
	PrivateKey   string
	WalletSource string
	Balance      uint `gorm:"balance;default:0"`
}

type Trade struct {
	gorm.Model
	TxHash    string    `gorm:"tx_hash;uniqueIndex;type:varchar(50)"`
	TradeType TradeType `gorm:"trade_type"`
	UserId    int64
	From      string
	To        string
	Amount    float64
	Coin      string `gorm:"coin;default:'usdt'"`
	Remark    string
}

type TradeType string

const (
	TradeTypeWithdraw         TradeType = "withdraw"
	TradeTypeRecharge         TradeType = "recharge"
	TradeTypeRedpacketSend    TradeType = "redpacketSend"
	TradeTypeRedpacketReceive TradeType = "redpacketReceive"
)
