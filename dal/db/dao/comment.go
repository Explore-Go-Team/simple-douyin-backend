package dao

import (
	"simple-douyin-backend/pkg/constants"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID          int64          `json:"id"`
	UserID      int64          `json:"user_id"`
	VideoID     int64          `json:"video_id"`
	CommentText string         `json:"comment_text"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName set table name to make gorm can correctly identify
func (Comment) TableName() string {
	return constants.CommentTableName
}

// AddNewComment add a comment
func AddNewComment(comment *Comment) error {
	if ok, _ := CheckUserExistByID(comment.UserID); !ok {
		return constants.UserIsNotExistErr
	}
	if ok, _ := CheckVideoExistByID(comment.VideoID); !ok {
		return constants.VideoIsNotExistErr
	}
	err := DB.Create(comment).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteCommentByID delete comment by comment id
func DeleteCommentByID(commentID int64) error {
	if ok, _ := CheckCommentExist(commentID); !ok {
		return constants.CommentIsNotExistErr
	}
	comment := &Comment{}
	err := DB.Where("id = ?", commentID).Delete(comment).Error
	if err != nil {
		return err
	}
	return nil
}

func CheckCommentExist(commentID int64) (bool, error) {
	comment := &Comment{}
	err := DB.Where("id = ?", commentID).Find(comment).Error
	if err != nil {
		return false, err
	}
	if comment.ID == 0 {
		return false, nil
	}
	return true, nil
}

func GetCommentsByVideoID(videoID int64) ([]*Comment, error) {
	var comments []*Comment
	if ok, _ := CheckVideoExistByID(videoID); !ok {
		return comments, constants.VideoIsNotExistErr
	}
	err := DB.Table(constants.CommentTableName).Where("video_id = ?", videoID).Find(&comments).Error
	if err != nil {
		return comments, err
	}
	return comments, nil
}

func GetCommentCountByVideoID(videoID int64) (int64, error) {
	var sum int64
	err := DB.Model(&Comment{}).Where("videoID = ?", videoID).Count(&sum).Error
	if err != nil {
		return sum, err
	}
	return sum, nil
}
