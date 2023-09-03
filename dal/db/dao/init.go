package dao

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init init DB
func Init(db *gorm.DB) {
	DB = db
}
