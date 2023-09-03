package main

import (
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
	"simple-douyin-backend/dal/db"
	"simple-douyin-backend/kitex_gen/basic/feed/feedservice"
	consts "simple-douyin-backend/pkg/constants"
)

func Init() {
	db.Init()
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", consts.FeedServiceAddr)
	if err != nil {
		panic(err)
	}
	srv := feedservice.NewServer(&FeedServiceImpl{}, server.WithServiceAddr(addr))
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
