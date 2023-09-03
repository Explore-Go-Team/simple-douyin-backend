package dao

import (
	"simple-douyin-backend/mw/redis"
	"simple-douyin-backend/pkg/constants"
	"time"

	"gorm.io/gorm"
)

// Follows follower is fan of user
type Follows struct {
	ID         int64          `json:"id"`
	UserId     int64          `json:"user_id"`
	FollowerId int64          `json:"follower_id"`
	CreatedAt  time.Time      `json:"create_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

// register redis operate strategy
var rdFollows redis.Follows

// TableName set table name to make gorm can correctly identify
func (Follows) TableName() string {
	return constants.FollowsTableName
}

func AddNewFollow(follow *Follows) (bool, error) {
	err := DB.Create(follow).Error
	if err != nil {
		return false, err
	}
	// add data to redis
	if rdFollows.CheckFollow(follow.FollowerId) {
		rdFollows.AddFollow(follow.UserId, follow.FollowerId)
	}
	if rdFollows.CheckFollower(follow.UserId) {
		rdFollows.AddFollower(follow.UserId, follow.FollowerId)
	}

	return true, nil
}

// DeleteFollow delete follow relation in db and update redis
func DeleteFollow(follow *Follows) (bool, error) {
	err := DB.Where("user_id = ? AND follower_id = ?", follow.UserId, follow.FollowerId).Delete(follow).Error
	if err != nil {
		return false, err
	}
	// if redis hit del
	if rdFollows.CheckFollow(follow.FollowerId) {
		rdFollows.DelFollow(follow.UserId, follow.FollowerId)
	}
	if rdFollows.CheckFollower(follow.UserId) {
		rdFollows.DelFollower(follow.UserId, follow.FollowerId)
	}
	return true, nil
}

// CheckFollowExist check the relation of user and follower
func CheckFollowExist(userId, followerId int64) (bool, error) {
	if rdFollows.CheckFollow(followerId) {
		return rdFollows.ExistFollow(userId, followerId), nil
	}
	if rdFollows.CheckFollower(userId) {
		return rdFollows.ExistFollower(userId, followerId), nil
	}
	follow := Follows{
		UserId:     userId,
		FollowerId: followerId,
	}
	err := DB.Where("user_id = ? AND follower_id = ?", userId, followerId).Find(&follow).Error
	if err != nil {
		return false, err
	}
	if follow.ID == 0 {
		return false, nil
	}
	return true, nil
}

// GetFollowCount query the number of users following
func GetFollowCount(followerId int64) (int64, error) {
	if rdFollows.CheckFollow(followerId) {
		return rdFollows.CountFollow(followerId)
	}

	// Not in the cache, go to the database to find and update the cache
	followings, err := getFollowIDs(followerId)
	if err != nil {
		return 0, err
	}
	// update redis asynchronously
	go addFollowRelationToRedis(followerId, followings)
	return int64(len(followings)), nil
}

// addFollowRelationToRedis update redis.RdbFollowing
func addFollowRelationToRedis(followerId int64, followings []int64) {
	for _, following := range followings {
		rdFollows.AddFollow(following, followerId)
	}
}

// GetFollowerCount query the number of followers of a user
func GetFollowerCount(userId int64) (int64, error) {
	if rdFollows.CheckFollower(userId) {
		return rdFollows.CountFollower(userId)
	}
	// Not in the cache, go to the database to find and update the cache
	followers, err := getFollowerIDs(userId)
	if err != nil {
		return 0, err
	}
	// update redis asynchronously
	go addFollowerRelationToRedis(userId, followers)
	return int64(len(followers)), nil
}

// addFollowerRelationToRedis update redis.RdbFollower
func addFollowerRelationToRedis(userId int64, followers []int64) {
	for _, follower := range followers {
		rdFollows.AddFollower(userId, follower)
	}
}

// getFollowIDs find user_id follow id list in db
func getFollowIDs(followerId int64) ([]int64, error) {
	var followIds []Follows
	err := DB.Where("follower_id = ?", followerId).Find(&followIds).Error
	if err != nil {
		return nil, err
	}
	var result []int64
	for _, v := range followIds {
		result = append(result, v.UserId)
	}
	return result, nil
}

// GetFollowIds find user_id follow id list in db or rdb
func GetFollowIds(followerID int64) ([]int64, error) {
	if rdFollows.CheckFollow(followerID) {
		return rdFollows.GetFollow(followerID), nil
	}
	return getFollowIDs(followerID)
}

// getFollowerIDs get follower id list in db
func getFollowerIDs(userID int64) ([]int64, error) {
	var followIds []Follows
	err := DB.Where("user_id = ?", userID).Find(&followIds).Error
	if err != nil {
		return nil, err
	}
	var result []int64
	for _, v := range followIds {
		result = append(result, v.FollowerId)
	}
	return result, nil
}

// GetFollowerIDs get follower id list in db or rdb
func GetFollowerIDs(userID int64) ([]int64, error) {
	if rdFollows.CheckFollower(userID) {
		return rdFollows.GetFollower(userID), nil
	}
	return getFollowerIDs(userID)
}

func GetFriendIDs(userId int64) ([]int64, error) {
	if !rdFollows.CheckFollow(userId) {
		following, err := getFollowIDs(userId)
		if err != nil {
			return *new([]int64), err
		}
		addFollowRelationToRedis(userId, following)
	}
	if !rdFollows.CheckFollower(userId) {
		followers, err := getFollowerIDs(userId)
		if err != nil {
			return *new([]int64), err
		}
		addFollowerRelationToRedis(userId, followers)
	}
	return rdFollows.GetFriend(userId), nil
}
