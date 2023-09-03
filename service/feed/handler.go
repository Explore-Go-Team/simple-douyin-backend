package main

import (
	"context"
	"fmt"
	"log"
	"simple-douyin-backend/dal/db/dao"
	"simple-douyin-backend/kitex_gen/basic/feed"
	"simple-douyin-backend/kitex_gen/basic/user"
	"simple-douyin-backend/kitex_gen/basic/user/userservice"
	"simple-douyin-backend/kitex_gen/common"
	"simple-douyin-backend/mw/jwt"
	"simple-douyin-backend/pkg/constants"
	"simple-douyin-backend/pkg/utils"
	"sync"
	"time"
)

var userClient userservice.Client

func init() {
	var err error
	userClient, err = userservice.NewClient("user-service")
	if err != nil {
		log.Fatal(err)
	}
}

type FeedServiceImpl struct{}

// Feed get the last ten videos until the deadline
func (s *FeedServiceImpl) Feed(ctx context.Context, req *feed.DouyinFeedRequest) (*feed.DouyinFeedResponse, error) {
	var lastTime time.Time
	if req.LatestTime == 0 {
		lastTime = time.Now()
	} else {
		lastTime = time.Unix(req.LatestTime/1000, 0)
	}
	fmt.Printf("LastTime: %v\n", lastTime)
	var currentUserID int64

	ret, err := jwt.Parse(req.Token)
	// 如果当前用户没有登陆，则将 current_user_id 赋值为 0
	if err != nil {
		currentUserID = 0
	} else {
		currentUserID = ret.UserID
	}
	dbVideos, err := dao.GetVideosByLastTime(lastTime)
	if err != nil {
		resp := utils.ConvertToResp(err)
		return &feed.DouyinFeedResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		}, err
	}

	videos := make([]*common.Video, 0, constants.VideoFeedCount)
	err = s.CopyVideos(&videos, &dbVideos, currentUserID)
	if err != nil {
		resp := utils.ConvertToResp(err)
		return &feed.DouyinFeedResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		}, err
	}
	resp := &feed.DouyinFeedResponse{}
	resp.VideoList = videos
	if len(dbVideos) != 0 {
		resp.NextTime = dbVideos[len(dbVideos)-1].PublishTime.Unix()
	}
	resp.StatusCode = constants.SuccessCode
	resp.StatusMsg = constants.SuccessMsg
	return resp, nil
}

// CopyVideos use db.Video make feed.Video
func (s *FeedServiceImpl) CopyVideos(result *[]*common.Video, data *[]*dao.Video, userID int64) error {
	for _, item := range *data {
		video := s.createVideo(item, userID)
		*result = append(*result, video)
	}
	return nil
}

// createVideo get video info by concurrent query
func (s *FeedServiceImpl) createVideo(data *dao.Video, userID int64) *common.Video {
	video := &common.Video{
		Id: data.ID,
		// convert path in the db into a complete url accessible by the front end
		//PlayUrl:  utils.ConvertURL(s.ctx, s.c, data.PlayURL),
		//CoverUrl: utils.ConvertURL(s.ctx, s.c, data.CoverURL),
		Title: data.Title,
	}

	var wg sync.WaitGroup
	wg.Add(4)

	// Get author information
	go func() {
		resp, err := userClient.User(nil, &user.DouyinUserRequest{
			UserId: userID,
		})
		if err != nil {
			log.Printf("Get user info error:" + err.Error())
		}
		author := resp.User
		video.Author = &common.User{
			Id:              author.Id,
			Name:            author.Name,
			FollowCount:     author.FollowCount,
			FollowerCount:   author.FollowerCount,
			IsFollow:        author.IsFollow,
			Avatar:          author.Avatar,
			BackgroundImage: author.BackgroundImage,
			Signature:       author.Signature,
			TotalFavorited:  author.TotalFavorited,
			WorkCount:       author.WorkCount,
			FavoriteCount:   author.FavoriteCount,
		}

		wg.Done()
	}()

	// Get the number of video likes
	go func() {
		err := *new(error)
		video.FavoriteCount, err = dao.GetFavoriteCount(data.ID)
		if err != nil {
			log.Printf("GetFavoriteCount func error:" + err.Error())
		}
		wg.Done()
	}()

	// Get comment count
	go func() {
		err := *new(error)
		video.CommentCount, err = dao.GetCommentCountByVideoID(data.ID)
		if err != nil {
			log.Printf("GetCommentCountByVideoID func error:" + err.Error())
		}
		wg.Done()
	}()

	// check favorite exist
	go func() {
		err := *new(error)
		video.IsFavorite, err = dao.CheckFavoriteExist(userID, data.ID)
		if err != nil {
			log.Printf("CheckFavoriteExist func error:" + err.Error())
		}
		wg.Done()
	}()

	wg.Wait()
	return video
}
