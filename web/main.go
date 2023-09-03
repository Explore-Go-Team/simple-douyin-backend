package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"simple-douyin-backend/pkg/constants"
)

func main() {
	h := server.Default(
		server.WithStreamBody(true),
		server.WithHostPorts(constants.WebListenAddr),
	)

	register(h)

	h.Spin()
}
