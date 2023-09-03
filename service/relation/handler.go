package main

import (
	"context"
	model "simple-douyin-backend/dal/db/dao"
	query "simple-douyin-backend/dal/db/gorm_gen"
	relation "simple-douyin-backend/kitex_gen/relation"
	"strconv"

	"gorm.io/gen/field"
)

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

// RelationAction implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationAction(ctx context.Context, req *relation.DouyinRelationActionRequest) (resp *relation.DouyinRelationActionResponse, err error) {
	if req.ActionType != 1 && req.ActionType != 2 {
		return &relation.DouyinRelationActionResponse{
			StatusCode: 1,
			StatusMsg:  "incorrect ActionType",
		}, nil
	}
	if req.FromUserId == req.ToUserId {
		return &relation.DouyinRelationActionResponse{
			StatusCode: 2,
			StatusMsg:  "cannot follow yourself",
		}, nil
	}
	{ // make sure user exists
		ud := query.UserDetail
		_, err = ud.Where(ud.UserID.Eq(req.ToUserId)).First()
		if err != nil {
			return &relation.DouyinRelationActionResponse{
				StatusCode: 3,
				StatusMsg:  "cannot follow user that doesn't exist",
			}, nil
		}
	}
	primaryKey := strconv.FormatInt(req.FromUserId, 10) + "_" + strconv.FormatInt(req.ToUserId, 10)
	isFollowUpdatedVal := req.ActionType == 1
	{ // check duplicated operation
		r := query.Relation
		result, err := r.Attrs(field.Attrs(&model.Relation{
			FromUserID: req.FromUserId,
			ToUserID:   req.ToUserId,
			IsFollow:   false,
			PrimaryKey: primaryKey,
		})).Where(r.PrimaryKey.Eq(primaryKey)).FirstOrCreate()
		if err != nil {
			return &relation.DouyinRelationActionResponse{
				StatusCode: 4,
				StatusMsg:  "fetch current status of relation failed",
			}, err
		}
		if result.IsFollow == isFollowUpdatedVal {
			// duplicated operation, do nothing
			return &relation.DouyinRelationActionResponse{
				StatusCode: 5,
				StatusMsg:  "duplicated operation",
			}, err
		}
	}
	// update
	err = query.Q.Transaction(func(tx *query.Query) error {
		r := tx.Relation
		ud := tx.UserDetail
		_, err = r.Where(r.PrimaryKey.Eq(primaryKey)).Update(r.IsFollow, isFollowUpdatedVal)
		if err != nil {
			return err
		}
		var deltaFollow int64
		if req.ActionType == 1 {
			deltaFollow = 1
		} else {
			deltaFollow = -1
		}
		_, err = ud.Where(ud.UserID.Eq(req.FromUserId)).UpdateSimple(ud.FollowCount.Add(deltaFollow))
		if err != nil {
			return err
		}
		_, err = ud.Where(ud.UserID.Eq(req.ToUserId)).UpdateSimple(ud.FollowerCount.Add(deltaFollow))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &relation.DouyinRelationActionResponse{
			StatusCode: 6,
			StatusMsg:  "update relation failed",
		}, err
	}
	return &relation.DouyinRelationActionResponse{
		StatusCode: 0,
		StatusMsg:  "update follow relation success",
	}, nil
}

// RelationFollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationFollowList(ctx context.Context, req *relation.DouyinRelationFollowListRequest) (resp *relation.DouyinRelationFollowListResponse, err error) {
	uid := req.UserId
	r := query.Relation
	rels, err := r.Where(r.FromUserID.Eq(uid), r.IsFollow.Is(true)).Find()
	if err != nil {
		return &relation.DouyinRelationFollowListResponse{
			StatusCode: 1,
			StatusMsg:  "query follow error",
		}, err
	}
	userIdList := make([]int64, len(rels))
	for i := range rels {
		userIdList[i] = rels[i].ToUserID
	}
	return &relation.DouyinRelationFollowListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserIdList: userIdList,
	}, nil
}

// RelationFollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationFollowerList(ctx context.Context, req *relation.DouyinRelationFollowerListRequest) (resp *relation.DouyinRelationFollowerListResponse, err error) {
	uid := req.UserId
	r := query.Relation
	rels, err := r.Where(r.ToUserID.Eq(uid), r.IsFollow.Is(true)).Find()
	if err != nil {
		return &relation.DouyinRelationFollowerListResponse{
			StatusCode: 1,
			StatusMsg:  "query follower error",
		}, err
	}
	userIdList := make([]int64, len(rels))
	for i := range rels {
		userIdList[i] = rels[i].FromUserID
	}
	return &relation.DouyinRelationFollowerListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserIdList: userIdList,
	}, nil
}

// RelationFriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationFriendList(ctx context.Context, req *relation.DouyinRelationFriendListRequest) (resp *relation.DouyinRelationFriendListResponse, err error) {
	uid := req.UserId
	var rels []struct {
		ToUserID int64
	}
	r := query.Relation
	r2 := query.Relation.As("r2")
	err = r.Select(r.ToUserID).Join(r2, r.ToUserID.EqCol(r2.FromUserID), r.FromUserID.EqCol(r2.ToUserID), r2.IsFollow.Is(true)).Where(r.FromUserID.Eq(uid), r.IsFollow.Is(true)).Scan(&rels)
	if err != nil {
		return &relation.DouyinRelationFriendListResponse{
			StatusCode: 1,
			StatusMsg:  "query friend error",
		}, err
	}
	userIdList := make([]int64, len(rels))
	for i := range rels {
		userIdList[i] = rels[i].ToUserID
	}
	return &relation.DouyinRelationFriendListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserIdList: userIdList,
	}, nil
}

// UserDetail implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) UserDetail(ctx context.Context, req *relation.DouyinUserDetailRequest) (resp *relation.DouyinUserDetailResponse, err error) {
	u := query.UserDetail
	detail, err := u.Where(u.UserID.Eq(req.UserId)).First()
	if err != nil {
		return nil, err
	}
	return &relation.DouyinUserDetailResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		Detail: &relation.UserDetail{
			Id:              detail.UserID,
			Name:            detail.Name,
			FollowCount:     detail.FollowCount,
			FollowerCount:   detail.FollowerCount,
			Avatar:          detail.Avatar,
			BackgroundImage: detail.BackgroundImage,
			Signature:       detail.Signature,
			TotalFavorited:  detail.TotalFavorited,
			WorkCount:       detail.WorkCount,
			FavoriteCount:   detail.FavoriteCount,
		},
	}, nil
}

// FriendRecentMsg implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FriendRecentMsg(ctx context.Context, req *relation.DouyinFriendRecentMsgRequest) (resp *relation.DouyinFriendRecentMsgResponse, err error) {
	m := query.Message
	msg, err := m.Where(m.FromUserID.Eq(req.UserId), m.ToUserID.Eq(req.FriendId)).
		Or(m.FromUserID.Eq(req.FriendId), m.ToUserID.Eq(req.UserId)).
		Order(m.CreatedAt.Desc()).
		First()
	if err != nil {
		return nil, err
	}
	msgType := int64(0)
	if msg.FromUserID == req.UserId {
		msgType = 1
	}
	return &relation.DouyinFriendRecentMsgResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		Message:    msg.Content,
		MsgType:    msgType,
	}, nil
}

// RelationIsFollow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationIsFollow(ctx context.Context, req *relation.DouyinRelationIsFollowRequest) (resp *relation.DouyinRelationIsFollowResponse, err error) {
	primaryKey := strconv.FormatInt(req.FromUserId, 10) + "_" + strconv.FormatInt(req.ToUserId, 10)
	r := query.Relation
	result, err := r.Where(r.PrimaryKey.Eq(primaryKey)).FirstOrInit()
	if err != nil {
		return nil, err
	}
	return &relation.DouyinRelationIsFollowResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		IsFollow:   result.IsFollow,
	}, nil
}
