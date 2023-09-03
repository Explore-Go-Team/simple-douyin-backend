package main

import (
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
	"simple-douyin-backend/dal/db"
	"simple-douyin-backend/kitex_gen/basic/user/userservice"
	consts "simple-douyin-backend/pkg/constants"
)

func Init() {
	db.Init()
}

func main() {
	Init()
	addr, err := net.ResolveTCPAddr("tcp", consts.UserServiceAddr)
	if err != nil {
		panic(err)
	}
	srv := userservice.NewServer(&UserServiceImpl{}, server.WithServiceAddr(addr))
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
