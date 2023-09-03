package dao

import (
	"simple-douyin-backend/pkg/constants"
)

type User struct {
	ID              int64  `json:"id"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
}

func (User) TableName() string {
	return constants.UserTableName
}

// CreateUser create user info
func CreateUser(user *User) (int64, error) {
	err := DB.Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

// GetUserByUsername get user by username
func GetUserByUsername(username string) (*User, error) {
	var user User
	if err := DB.Where("username = ?", username).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID get user by user id
func GetUserByID(userId int64) (*User, error) {
	var user User
	if err := DB.Where("id = ?", userId).Find(&user).Error; err != nil {
		return nil, err
	}
	if user == (User{}) {
		err := constants.UserIsNotExistErr
		return nil, err
	}
	return &user, nil
}

// VerifyUser verify username and password in the db
func VerifyUser(username, password string) (int64, error) {
	var user User
	if err := DB.Where("username = ? AND password = ?", username, password).Find(&user).Error; err != nil {
		return 0, err
	}
	if user.ID == 0 {
		err := constants.PasswordIsNotVerified
		return user.ID, err
	}
	return user.ID, nil
}

// CheckUserExistByID check if user exists
func CheckUserExistByID(userID int64) (bool, error) {
	var user User
	if err := DB.Where("id = ?", userID).Find(&user).Error; err != nil {
		return false, err
	}
	if user == (User{}) {
		return false, nil
	}
	return true, nil
}
