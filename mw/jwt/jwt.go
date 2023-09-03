package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	db "simple-douyin-backend/dal/db/dao"
	"simple-douyin-backend/pkg/constants"
	"simple-douyin-backend/pkg/utils"
)

var (
	identity  = "user_id"
	secretKey = []byte("douyin")
)

type AuthInfo struct {
	Token  string
	UserID int64
}

type AppClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func Authenticate(username, password string) (*AuthInfo, error) {
	user, err := db.GetUserByUsername(username)
	if ok := utils.VerifyPassword(password, user.Password); !ok {
		err = constants.PasswordIsNotVerified
		return nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AppClaims{
		UserID:           user.ID,
		RegisteredClaims: jwt.RegisteredClaims{},
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, err
	}
	fmt.Println(tokenString)
	return &AuthInfo{
		Token:  tokenString,
		UserID: user.ID,
	}, nil
}

func Parse(token string) (*AuthInfo, error) {
	tokenParse, err := jwt.ParseWithClaims(token, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		fmt.Println("Failed to parse token:", err)
		return nil, err
	}
	// 验证 Token
	if claims, ok := tokenParse.Claims.(*AppClaims); ok && tokenParse.Valid {
		return &AuthInfo{
			UserID: claims.UserID,
			Token:  token,
		}, nil
	} else {
		log.Println("Invalid Token")
		return nil, constants.AuthorizationFailedErr
	}
}
