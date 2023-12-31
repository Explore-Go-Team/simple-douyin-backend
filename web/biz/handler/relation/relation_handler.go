// Code generated by hertz generator.

package relation

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	relationRPC "simple-douyin-backend/kitex_gen/relation"
	relationService "simple-douyin-backend/kitex_gen/relation/relationservice"
	"simple-douyin-backend/mw/jwt"
	"simple-douyin-backend/pkg/constants"
	"simple-douyin-backend/web/biz/errors"
	relation "simple-douyin-backend/web/biz/model/relation"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

var (
	destService      = "relation_service"
	serviceHostPorts = constants.RelationServiceAddr
	relationClient   relationService.Client
)

func init() {
	var err error
	relationClient, err = relationService.NewClient(destService, client.WithHostPorts(serviceHostPorts))
	if err != nil {
		hlog.Fatalf("create relationRPC client failed: %v", err)
	}
}

func queryIsFollowUser(ctx context.Context, c *app.RequestContext, fromUserId, toUserId int64) (bool, error) {
	rpcReq := &relationRPC.DouyinRelationIsFollowRequest{
		FromUserId: fromUserId,
		ToUserId:   toUserId,
	}
	rpcResp, err := relationClient.RelationIsFollow(ctx, rpcReq)
	if err != nil {
		err := errors.NewRPCError(err)
		hlog.CtxErrorf(ctx, err.Error())
		c.Error(err)
		return false, err
	}
	return rpcResp.IsFollow, nil
}

func queryUserInfo(ctx context.Context, c *app.RequestContext, fromUserId, toUserId int64) (*relation.User, error) {
	rpcReq := &relationRPC.DouyinUserDetailRequest{
		UserId: toUserId,
	}
	rpcResp, err := relationClient.UserDetail(ctx, rpcReq)
	if err != nil {
		err := errors.NewRPCError(err)
		hlog.CtxErrorf(ctx, err.Error())
		c.Error(err)
		return nil, err
	}
	isFollow, err := queryIsFollowUser(ctx, c, fromUserId, toUserId)
	if err != nil {
		err := errors.NewRPCError(err)
		hlog.CtxErrorf(ctx, err.Error())
		c.Error(err)
		return nil, err
	}
	return &relation.User{
		Id:              rpcResp.Detail.Id,
		Name:            rpcResp.Detail.Name,
		FollowCount:     rpcResp.Detail.FollowCount,
		FollowerCount:   rpcResp.Detail.FollowerCount,
		IsFollow:        isFollow,
		Avatar:          rpcResp.Detail.Avatar,
		BackgroundImage: rpcResp.Detail.BackgroundImage,
		Signature:       rpcResp.Detail.Signature,
		TotalFavorited:  rpcResp.Detail.TotalFavorited,
		WorkCount:       rpcResp.Detail.WorkCount,
		FavoriteCount:   rpcResp.Detail.FavoriteCount,
	}, nil
}

func queryRecentMessage(ctx context.Context, c *app.RequestContext, currentUID, friendUID int64) (string string, msgType int64, err error) {
	rpcReq := &relationRPC.DouyinFriendRecentMsgRequest{
		UserId:   currentUID,
		FriendId: friendUID,
	}
	rpcResp, err := relationClient.FriendRecentMsg(ctx, rpcReq)
	if err != nil {
		err := errors.NewRPCError(err)
		hlog.CtxErrorf(ctx, err.Error())
		c.Error(err)
		return "", 0, err
	}
	return rpcResp.Message, rpcResp.MsgType, nil
}

// RelationAction .
// @router /douyin/relation/action/ [POST]
func RelationAction(ctx context.Context, c *app.RequestContext) {
	var err error
	var req relation.DouyinRelationActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxErrorf(ctx, err.Error())
		c.Error(err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	loginInfo, err := jwt.Parse(req.Token)
	if err != nil {
		err := errors.NewAuthError(err)
		hlog.CtxErrorf(ctx, err.Error())
		c.Error(err)
		c.String(403, errors.AuthErrorMsg)
		return
	}
	rpcReq := &relationRPC.DouyinRelationActionRequest{
		FromUserId: loginInfo.UserID,
		ToUserId:   req.ToUserId,
		ActionType: req.ActionType,
	}
	rpcResp, err := relationClient.RelationAction(ctx, rpcReq)
	if err != nil {
		err := errors.NewRPCError(err)
		hlog.CtxErrorf(ctx, err.Error())
		c.Error(err)
		c.String(500, errors.RPCErrorMsg)
		return
	}
	resp := &relation.DouyinRelationActionResponse{
		StatusCode: rpcResp.StatusCode,
		StatusMsg:  rpcResp.StatusMsg,
	}
	c.JSON(consts.StatusOK, resp)
}

// RelationFollowList .
// @router /douyin/relation/follow/list/ [GET]
func RelationFollowList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req relation.DouyinRelationFollowListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	loginInfo, err := jwt.Parse(req.Token)
	if err != nil {
		err := errors.NewAuthError(err)
		hlog.CtxErrorf(ctx, err.Error())
		c.Error(err)
		c.String(403, errors.AuthErrorMsg)
		return
	}
	rpcReq := &relationRPC.DouyinRelationFollowListRequest{
		UserId: req.UserId,
	}
	rpcResp, err := relationClient.RelationFollowList(ctx, rpcReq)
	if err != nil {
		err := errors.NewRPCError(err)
		hlog.CtxErrorf(ctx, err.Error())
		c.Error(err)
		c.String(500, errors.RPCErrorMsg)
		return
	}
	userList := make([]*relation.User, len(rpcResp.UserIdList))
	for i := range rpcResp.UserIdList {
		userList[i], err = queryUserInfo(ctx, c, loginInfo.UserID, rpcResp.UserIdList[i])
		if err != nil {
			err := errors.NewRPCError(err)
			hlog.CtxErrorf(ctx, err.Error())
			c.Error(err)
			// dont stop, keep this user's info empty, just return others
		}
	}
	resp := &relation.DouyinRelationFollowListResponse{
		StatusCode: rpcResp.StatusCode,
		StatusMsg:  rpcResp.StatusMsg,
		UserList:   userList,
	}
	c.JSON(consts.StatusOK, resp)
}

// RelationFollowerList .
// @router /douyin/relation/follower/list/ [GET]
func RelationFollowerList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req relation.DouyinRelationFollowerListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	loginInfo, err := jwt.Parse(req.Token)
	if err != nil {
		err := errors.NewAuthError(err)
		hlog.CtxErrorf(ctx, err.Error())
		c.Error(err)
		c.String(403, errors.AuthErrorMsg)
		return
	}
	rpcReq := &relationRPC.DouyinRelationFollowerListRequest{
		UserId: req.UserId,
	}
	rpcResp, err := relationClient.RelationFollowerList(ctx, rpcReq)
	if err != nil {
		err := errors.NewRPCError(err)
		hlog.CtxErrorf(ctx, err.Error())
		c.Error(err)
		c.String(500, errors.RPCErrorMsg)
		return
	}
	userList := make([]*relation.User, len(rpcResp.UserIdList))
	for i := range rpcResp.UserIdList {
		userList[i], err = queryUserInfo(ctx, c, loginInfo.UserID, rpcResp.UserIdList[i])
		if err != nil {
			err := errors.NewRPCError(err)
			hlog.CtxErrorf(ctx, err.Error())
			c.Error(err)
			// dont stop, keep this user's info empty, just return others
		}
	}
	resp := &relation.DouyinRelationFollowerListResponse{
		StatusCode: rpcResp.StatusCode,
		StatusMsg:  rpcResp.StatusMsg,
		UserList:   userList,
	}
	c.JSON(consts.StatusOK, resp)
}

// RelationFriendList .
// @router /douyin/relation/friend/list/ [GET]
func RelationFriendList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req relation.DouyinRelationFriendListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	loginInfo, err := jwt.Parse(req.Token)
	if err != nil {
		hlog.CtxErrorf(ctx, err.Error())
		c.Error(errors.NewAuthError(err))
		c.String(403, errors.AuthErrorMsg)
		return
	}
	if loginInfo.UserID != req.UserId {
		c.String(403, "cannot query other user's friend list")
		return
	}
	rpcReq := &relationRPC.DouyinRelationFriendListRequest{
		UserId: req.UserId,
	}
	rpcResp, err := relationClient.RelationFriendList(ctx, rpcReq)
	if err != nil {
		err := errors.NewRPCError(err)
		c.Error(err)
		c.String(500, errors.RPCErrorMsg)
		return
	}
	userList := make([]*relation.FriendUser, len(rpcResp.UserIdList))
	for i := range rpcResp.UserIdList {
		msg, msgtype, err := queryRecentMessage(ctx, c, req.UserId, rpcResp.UserIdList[i])
		if err != nil {
			err := errors.NewRPCError(err)
			hlog.CtxErrorf(ctx, err.Error())
			c.Error(err)
		}
		userInfo, err := queryUserInfo(ctx, c, loginInfo.UserID, rpcResp.UserIdList[i])
		if err != nil {
			err := errors.NewRPCError(err)
			hlog.CtxErrorf(ctx, err.Error())
			c.Error(err)
			// dont stop, keep this user's info empty, just return others
		}
		userList[i] = &relation.FriendUser{
			User:    userInfo,
			Message: msg,
			MsgType: msgtype,
		}
	}
	resp := &relation.DouyinRelationFriendListResponse{
		StatusCode: rpcResp.StatusCode,
		StatusMsg:  rpcResp.StatusMsg,
		UserList:   userList,
	}
	c.JSON(consts.StatusOK, resp)
}
