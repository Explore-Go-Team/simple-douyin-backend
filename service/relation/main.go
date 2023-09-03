package main

import (
	"log"
	"net"
	"simple-douyin-backend/dal"
	relation "simple-douyin-backend/kitex_gen/relation/relationservice"
	consts "simple-douyin-backend/pkg/constants"

	server "github.com/cloudwego/kitex/server"
)

func main() {
	dal.Init()
	addr, err := net.ResolveTCPAddr("tcp", consts.RelationServiceAddr)
	if err != nil {
		panic(err)
	}

	svr := relation.NewServer(new(RelationServiceImpl), server.WithServiceAddr(addr))
	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
