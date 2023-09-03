package dao

import (
	"time"
)

type Message struct {
	ID         int64 `gorm:"primaryKey"`
	FromUserID int64 `gorm:"index;foreignKey:FromUserID;references:User.ID"`
	ToUserID   int64 `gorm:"index;foreignKey:ToUserID;references:User.ID"`
	Content    string
	CreatedAt  time.Time // 消息发送时间

}
