package dao

import (
	"simple-douyin-backend/pkg/constants"
	"time"
)

type Video struct {
	ID          int64
	AuthorID    int64
	PlayURL     string
	CoverURL    string
	PublishTime time.Time
	Title       string
}

func (Video) TableName() string {
	return constants.VideosTableName
}

func CreateVideo(video *Video) (videoId int64, err error) {
	err = DB.Create(video).Error
	if err != nil {
		return 0, err
	}
	return video.ID, err
}

func GetVideosByLastTime(lastTime time.Time) ([]*Video, error) {
	videos := make([]*Video, constants.VideoFeedCount)
	err := DB.Where("publish_time < ?", lastTime).Order("publish_time desc").Limit(constants.VideoFeedCount).Find(&videos).Error
	if err != nil {
		return videos, err
	}
	return videos, nil
}

func GetVideoByUserID(userID int64) ([]*Video, error) {
	var videos []*Video
	err := DB.Where("author_id = ?", userID).Find(&videos).Error
	if err != nil {
		return videos, err
	}
	return videos, err
}

func GetVideoByIDs(videoIDs []int64) ([]*Video, error) {
	var videos []*Video
	var err error
	for _, item := range videoIDs {
		var video *Video
		err = DB.Where("id = ?", item).Find(&video).Error
		if err != nil {
			return videos, err
		}
		videos = append(videos, video)
	}

	return videos, err
}

// GetWorkCount get the num of video published by the user
func GetWorkCount(userID int64) (int64, error) {
	var count int64
	err := DB.Model(&Video{}).Where("author_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CheckVideoExistByID check if video exist
func CheckVideoExistByID(videoID int64) (bool, error) {
	var video Video
	if err := DB.Where("id = ?", videoID).Find(&video).Error; err != nil {
		return false, err
	}
	if video == (Video{}) {
		return false, nil
	}
	return true, nil
}
