package redis

import (
	"strconv"
)

const (
	FollowerSuffix = ":follower"
	FollowSuffix   = ":follow"
)

type (
	Follows struct{}
)

func (f Follows) AddFollow(userId, followerId int64) {
	add(rdbFollows, strconv.FormatInt(followerId, 10)+FollowSuffix, userId)
}

func (f Follows) AddFollower(userId, followerId int64) {
	add(rdbFollows, strconv.FormatInt(userId, 10)+FollowerSuffix, followerId)
}

func (f Follows) DelFollow(userId, followerId int64) {
	del(rdbFollows, strconv.FormatInt(followerId, 10)+FollowSuffix, userId)
}

func (f Follows) DelFollower(userId, followerId int64) {
	del(rdbFollows, strconv.FormatInt(userId, 10)+FollowerSuffix, followerId)
}

func (f Follows) CheckFollow(followerId int64) bool {
	return check(rdbFollows, strconv.FormatInt(followerId, 10)+FollowSuffix)
}

func (f Follows) CheckFollower(userId int64) bool {
	return check(rdbFollows, strconv.FormatInt(userId, 10)+FollowerSuffix)
}

func (f Follows) ExistFollow(userId, followerId int64) bool {
	return exist(rdbFollows, strconv.FormatInt(followerId, 10)+FollowSuffix, userId)
}

func (f Follows) ExistFollower(userId, followerId int64) bool {
	return exist(rdbFollows, strconv.FormatInt(userId, 10)+FollowerSuffix, followerId)
}

func (f Follows) CountFollow(followerId int64) (int64, error) {
	return count(rdbFollows, strconv.FormatInt(followerId, 10)+FollowSuffix)
}

func (f Follows) CountFollower(userId int64) (int64, error) {
	return count(rdbFollows, strconv.FormatInt(userId, 10)+FollowerSuffix)
}

func (f Follows) GetFollow(followerId int64) []int64 {
	return get(rdbFollows, strconv.FormatInt(followerId, 10)+FollowSuffix)
}

func (f Follows) GetFollower(userId int64) []int64 {
	return get(rdbFollows, strconv.FormatInt(userId, 10)+FollowerSuffix)
}

// GetFriend get the friend of the id via intersection
func (f Follows) GetFriend(id int64) (friends []int64) {
	ks1 := strconv.FormatInt(id, 10) + FollowSuffix
	ks2 := strconv.FormatInt(id, 10) + FollowerSuffix
	v, _ := rdbFollows.SInter(ks1, ks2).Result()
	for _, vs := range v {
		v_i64, _ := strconv.ParseInt(vs, 10, 64)
		friends = append(friends, v_i64)
	}
	return friends
}
