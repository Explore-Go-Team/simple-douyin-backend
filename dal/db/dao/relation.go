package dao

import (
	"gorm.io/gorm"
	"time"
)

type Relation struct {
	ID         int64 `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	FromUserID int64          `gorm:"index"`
	ToUserID   int64          `gorm:"index"`
	IsFollow   bool
	PrimaryKey string `gorm:"primaryKey"` //联合主键
}
