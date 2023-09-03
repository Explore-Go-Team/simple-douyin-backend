package dal

import (
	"simple-douyin-backend/dal/db"
)

// Init init dal
func Init() {
	db.Init() // mysql init
	//redis.InitRedis()
}
