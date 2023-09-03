package main

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"simple-douyin-backend/dal/db/dao"
	"simple-douyin-backend/kitex_gen/basic/user"
	"simple-douyin-backend/kitex_gen/common"
	"simple-douyin-backend/mw/jwt"
	"simple-douyin-backend/pkg/constants"
	"simple-douyin-backend/pkg/utils"
	"sync"
)

type UserServiceImpl struct{}

// Register register user return user id.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.DouyinUserRegisterRequest) (*user.DouyinUserRegisterResponse, error) {
	dbUser, err := dao.GetUserByUsername(req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		resp := utils.ConvertToResp(err)
		return &user.DouyinUserRegisterResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		}, err
	}
	if *dbUser != (dao.User{}) {
		return &user.DouyinUserRegisterResponse{
			StatusCode: constants.UserAlreadyExistErrCode,
			StatusMsg:  "User already exists",
		}, constants.UserAlreadyExistErr
	}

	password, err := utils.HashPassword(req.Password)
	var userID int64
	userID, err = dao.CreateUser(&dao.User{
		Username:        req.Username,
		Password:        password,
		Avatar:          constants.TestAva,
		BackgroundImage: constants.TestBackground,
		Signature:       constants.TestSign,
	})

	ret, err := jwt.Authenticate(req.Username, req.Password)
	if err != nil {
		resp := utils.ConvertToResp(err)
		return &user.DouyinUserRegisterResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		}, err
	}

	authInfo := ret
	return &user.DouyinUserRegisterResponse{
		StatusCode: constants.SuccessCode,
		StatusMsg:  constants.SuccessMsg,
		Token:      authInfo.Token,
		UserId:     userID,
	}, nil
}

// Login user login api
func (s *UserServiceImpl) Login(ctx context.Context, req *user.DouyinUserLoginRequest) (*user.DouyinUserLoginResponse, error) {
	// verify user
	ret, err := jwt.Authenticate(req.Username, req.Password)
	authInfo := ret
	if err != nil {
		resp := utils.ConvertToResp(err)
		return &user.DouyinUserLoginResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		}, err
	}
	return &user.DouyinUserLoginResponse{
		StatusCode: constants.SuccessCode,
		StatusMsg:  constants.SuccessMsg,
		Token:      authInfo.Token,
		UserId:     authInfo.UserID,
	}, nil
}

// User the function of user api
func (s *UserServiceImpl) User(ctx context.Context, req *user.DouyinUserRequest) (*user.DouyinUserResponse, error) {
	ret, err := jwt.Parse(req.Token)
	var currentUserID int64
	if err != nil {
		currentUserID = 0
	} else {
		currentUserID = ret.UserID
	}
	queryUserID := req.UserId
	dbUser, err := s.GetUserInfo(queryUserID, currentUserID)
	if err != nil {
		resp := utils.ConvertToResp(err)
		return &user.DouyinUserResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		}, err
	}
	return &user.DouyinUserResponse{
		StatusCode: constants.SuccessCode,
		StatusMsg:  constants.SuccessMsg,
		User:       dbUser,
	}, nil
}

// GetUserInfo Query the information of query_user_id according to the current user user_id
func (s *UserServiceImpl) GetUserInfo(queryUserID, currentUserID int64) (*common.User, error) {
	u := &common.User{}
	errChan := make(chan error, 7)
	defer close(errChan)
	var wg sync.WaitGroup
	wg.Add(7)
	go func() {
		dbUser, err := dao.GetUserByID(queryUserID)
		if err != nil {
			errChan <- err
		} else {
			u.Name = dbUser.Username
			//u.Avatar = utils.ConvertURL(nil, s.c, dbUser.Avatar)
			//u.BackgroundImage = utils.ConvertURL(s.ctx, s.c, dbUser.BackgroundImage)
			u.Signature = dbUser.Signature
		}
		wg.Done()
	}()

	go func() {
		WorkCount, err := dao.GetWorkCount(queryUserID)
		if err != nil {
			errChan <- err
		} else {
			u.WorkCount = WorkCount
		}
		wg.Done()
	}()

	go func() {
		FollowCount, err := dao.GetFollowCount(queryUserID)
		if err != nil {
			errChan <- err
			return
		} else {
			u.FollowCount = FollowCount
		}
		wg.Done()
	}()

	go func() {
		FollowerCount, err := dao.GetFollowerCount(queryUserID)
		if err != nil {
			errChan <- err
		} else {
			u.FollowerCount = FollowerCount
		}
		wg.Done()
	}()

	go func() {
		if currentUserID != 0 {
			IsFollow, err := dao.CheckFollowExist(currentUserID, queryUserID)
			if err != nil {
				errChan <- err
			} else {
				u.IsFollow = IsFollow
			}
		} else {
			u.IsFollow = false
		}
		wg.Done()
	}()

	go func() {
		FavoriteCount, err := dao.GetFavoriteCountByUserID(queryUserID)
		if err != nil {
			errChan <- err
		} else {
			u.FavoriteCount = FavoriteCount
		}
		wg.Done()
	}()

	go func() {
		TotalFavorited, err := dao.QueryTotalFavoritedByAuthorID(queryUserID)
		if err != nil {
			errChan <- err
		} else {
			u.TotalFavorited = TotalFavorited
		}
		wg.Done()
	}()

	wg.Wait()
	select {
	case result := <-errChan:
		return &common.User{}, result
	default:
	}
	u.Id = queryUserID
	return u, nil
}
