package model

type ButtonInfo struct {
	Text    string  `json:"text"`
	Data    string  `json:"data"`
	BtnType BtnType `json:"btn_type"`
}

type BtnType string

const (
	BtnTypeUrl    BtnType = "url"
	BtnTypeData   BtnType = "data"
	BtnTypeSwitch BtnType = "switch"
)
