package dao

import (
	"simple-douyin-backend/mw/redis"
	"simple-douyin-backend/pkg/constants"
	"time"

	"gorm.io/gorm"
)

// register redis operate strategy
var rdFavorite redis.Favorite

type Favorites struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	VideoID   int64          `json:"video_id"`
	CreatedAt time.Time      `json:"create_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

// TableName set table name to make gorm can correctly identify
func (Favorites) TableName() string {
	return constants.FavoritesTableName
}

// AddNewFavorite add favorite relation
func AddNewFavorite(favorite *Favorites) (bool, error) {
	err := DB.Create(favorite).Error
	if err != nil {
		return false, err
	}
	// add data to redis
	if rdFavorite.CheckLiked(favorite.VideoID) {
		rdFavorite.AddLiked(favorite.UserID, favorite.VideoID)
	}
	if rdFavorite.CheckLike(favorite.UserID) {
		rdFavorite.AddLike(favorite.UserID, favorite.VideoID)
	}

	return true, nil
}

// DeleteFavorite delete favorite relation
func DeleteFavorite(favorite *Favorites) (bool, error) {
	err := DB.Where("video_id = ? AND user_id = ?", favorite.VideoID, favorite.UserID).Delete(favorite).Error
	if err != nil {
		return false, err
	}
	// del data if hit
	if rdFavorite.CheckLiked(favorite.VideoID) {
		rdFavorite.DelLiked(favorite.UserID, favorite.VideoID)
	}
	if rdFavorite.CheckLike(favorite.UserID) {
		rdFavorite.DelLike(favorite.UserID, favorite.VideoID)
	}
	return true, nil
}

// CheckFavoriteExist query the like record by video_id and user_id
func CheckFavoriteExist(userID, videoID int64) (bool, error) {
	if rdFavorite.CheckLiked(videoID) {
		return rdFavorite.ExistLiked(userID, videoID), nil
	}
	if rdFavorite.CheckLike(userID) {
		return rdFavorite.ExistLike(userID, videoID), nil
	}
	var sum int64
	err := DB.Model(&Favorites{}).Where("video_id = ? AND user_id = ?", videoID, userID).Count(&sum).Error
	if err != nil {
		return false, err
	}
	if sum == 0 {
		return false, nil
	}
	return true, nil
}

// QueryTotalFavoritedByAuthorID query the like num of all the video published by  the user
func QueryTotalFavoritedByAuthorID(authorID int64) (int64, error) {
	var sum int64
	err := DB.Table(constants.FavoritesTableName).Joins("JOIN videos ON likes.video_id = videos.id").
		Where("videos.author_id = ?", authorID).Count(&sum).Error
	if err != nil {
		return 0, err
	}
	return sum, nil
}

// getFavoriteIDs get the id list of video liked by the user in db
func getFavoriteIDs(userID int64) ([]int64, error) {
	var favorites []Favorites
	err := DB.Where("user_id = ?", userID).Find(&favorites).Error
	if err != nil {
		return nil, err
	}
	var result []int64
	for _, v := range favorites {
		result = append(result, v.VideoID)
	}
	return result, nil
}

// GetFavoriteIDs get the id list of video liked by the user
func GetFavoriteIDs(userID int64) ([]int64, error) {
	if rdFavorite.CheckLike(userID) {
		return rdFavorite.GetLike(userID), nil
	}
	return getFavoriteIDs(userID)
}

// GetFavoriteCountByUserID get the num of the video liked by user
func GetFavoriteCountByUserID(userID int64) (int64, error) {
	if rdFavorite.CheckLike(userID) {
		return rdFavorite.CountLike(userID)
	}
	// Not in the cache, go to the database to find and update the cache
	videos, err := getFavoriteIDs(userID)
	if err != nil {
		return 0, err
	}

	// update redis asynchronously
	go func(user int64, videos []int64) {
		for _, video := range videos {
			rdFavorite.AddLiked(user, video)
		}
	}(userID, videos)

	return int64(len(videos)), nil
}

// getFavoriterIDs get the id list of liker of video in db
func getFavoriterIDs(videoID int64) ([]int64, error) {
	var favorites []Favorites
	err := DB.Where("video_id = ?", videoID).Find(&favorites).Error
	if err != nil {
		return nil, err
	}
	var result []int64
	for _, v := range favorites {
		result = append(result, v.UserID)
	}
	return result, nil
}

// GetFavoriterIDs get the id list of liker of  video
func GetFavoriterIDs(videoID int64) ([]int64, error) {
	if rdFavorite.CheckLiked(videoID) {
		return rdFavorite.GetLiked(videoID), nil
	}
	return getFavoriterIDs(videoID)
}

// GetFavoriteCount count the favorite of video
func GetFavoriteCount(videoId int64) (int64, error) {
	if rdFavorite.CheckLiked(videoId) {
		return rdFavorite.CountLiked(videoId)
	}
	// Not in the cache, go to the database to find and update the cache
	likes, err := getFavoriterIDs(videoId)
	if err != nil {
		return 0, err
	}

	// update redis asynchronously
	go func(users []int64, video int64) {
		for _, user := range users {
			rdFavorite.AddLiked(user, video)
		}
	}(likes, videoId)
	return int64(len(likes)), nil
}
