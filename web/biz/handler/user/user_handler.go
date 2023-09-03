// Code generated by hertz generator.

package user

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	"log"
	"simple-douyin-backend/kitex_gen/basic/user"
	"simple-douyin-backend/kitex_gen/basic/user/userservice"
	"simple-douyin-backend/pkg/constants"
	"simple-douyin-backend/pkg/utils"
	userHTTP "simple-douyin-backend/web/biz/model/basic/user"
)

var userClient userservice.Client

func init() {
	var err error
	userClient, err = userservice.NewClient("user_service", client.WithHostPorts(constants.UserServiceAddr))
	if err != nil {
		log.Fatal(err)
	}
}

// Register .
// @router /douyin/user/register/ [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req userHTTP.DouyinUserRegisterRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		resp := utils.ConvertToResp(err)
		c.JSON(consts.StatusOK, user.DouyinUserRegisterResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})
		return
	}
	fmt.Println("register")
	reqRPC := user.DouyinUserRegisterRequest{
		Username: req.Username,
		Password: req.Password,
	}
	resp, err := userClient.Register(ctx, &reqRPC)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(consts.StatusOK, resp)
}

// Login .
// @router /douyin/user/login/ [POST]
func Login(ctx context.Context, c *app.RequestContext) {
	var req userHTTP.DouyinUserLoginRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		resp := utils.ConvertToResp(err)
		c.JSON(consts.StatusOK, user.DouyinUserLoginResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})
	}
	reqRPC := user.DouyinUserLoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	resp, err := userClient.Login(ctx, &reqRPC)
	c.JSON(consts.StatusOK, resp)
}

// User .
// @router /douyin/user/ [GET]
func User(ctx context.Context, c *app.RequestContext) {
	var err error
	var req userHTTP.DouyinUserRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		resp := utils.ConvertToResp(err)
		c.JSON(consts.StatusOK, user.DouyinUserResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})
		return
	}
	reqRPC := user.DouyinUserRequest{
		UserId: req.UserId,
		Token:  req.Token,
	}
	resp, err := userClient.User(ctx, &reqRPC)
	c.JSON(consts.StatusOK, resp)
}