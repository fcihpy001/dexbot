package model

import "gorm.io/gorm"

type GroupInfo struct {
	gorm.Model
	GroupId    int64 `gorm:"primaryKey;uniqueIndex"`
	GroupName  string
	GroupType  string
	Uid        int64
	Permission string
	GroupAdmin string
}

type ScheduleDelete struct {
	gorm.Model
	ChatId    int64 `gorm:"primaryKey"`
	MessageId int
}
