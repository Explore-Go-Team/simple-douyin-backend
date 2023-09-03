package constants

import (
	"fmt"
)

const (
	SuccessCode    = 0
	ServiceErrCode = iota + 10000
	ParamErrCode
	AuthorizationFailedErrCode

	UserAlreadyExistErrCode
	UserIsNotExistErrCode

	FollowRelationAlreadyExistErrCode
	FollowRelationNotExistErrCode

	FavoriteRelationAlreadyExistErrCode
	FavoriteRelationNotExistErrCode
	FavoriteActionErrCode

	MessageAddFailedErrCode
	FriendListNoPermissionErrCode

	VideoIsNotExistErrCode
	CommentIsNotExistErrCode
)

const (
	SuccessMsg               = "Success"
	ServerErrMsg             = "Service is unable to start successfully"
	ParamErrMsg              = "Wrong Parameter has been given"
	UserIsNotExistErrMsg     = "user is not exist"
	PasswordIsNotVerifiedMsg = "username or password not verified"
	FavoriteActionErrMsg     = "favorite add failed"

	MessageAddFailedErrMsg    = "message add failed"
	FriendListNoPermissionMsg = "You can't query his friend list"
	VideoIsNotExistErrMsg     = "video is not exist"
	CommentIsNotExistErrMsg   = "comment is not exist"
)

type RequestError struct {
	Code int32
	Msg  string
}

func (e RequestError) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.Code, e.Msg)
}

func NewRequestError(code int32, msg string) RequestError {
	return RequestError{code, msg}
}

func (e RequestError) WithMessage(msg string) RequestError {
	e.Msg = msg
	return e
}

var (
	Success                         = NewRequestError(SuccessCode, SuccessMsg)
	ServiceErr                      = NewRequestError(ServiceErrCode, ServerErrMsg)
	ParamErr                        = NewRequestError(ParamErrCode, ParamErrMsg)
	UserAlreadyExistErr             = NewRequestError(UserAlreadyExistErrCode, "User already exists")
	AuthorizationFailedErr          = NewRequestError(AuthorizationFailedErrCode, "Authorization failed")
	UserIsNotExistErr               = NewRequestError(UserIsNotExistErrCode, UserIsNotExistErrMsg)
	PasswordIsNotVerified           = NewRequestError(AuthorizationFailedErrCode, PasswordIsNotVerifiedMsg)
	FollowRelationAlreadyExistErr   = NewRequestError(FollowRelationAlreadyExistErrCode, "Follow Relation already exist")
	FollowRelationNotExistErr       = NewRequestError(FollowRelationNotExistErrCode, "Follow Relation does not exist")
	FavoriteRelationAlreadyExistErr = NewRequestError(FavoriteRelationAlreadyExistErrCode, "Favorite Relation already exist")
	FavoriteRelationNotExistErr     = NewRequestError(FavoriteRelationNotExistErrCode, "FavoriteRelationNotExistErr")
	FavoriteActionErr               = NewRequestError(FavoriteActionErrCode, FavoriteActionErrMsg)

	MessageAddFailedErr       = NewRequestError(MessageAddFailedErrCode, MessageAddFailedErrMsg)
	FriendListNoPermissionErr = NewRequestError(FriendListNoPermissionErrCode, FriendListNoPermissionMsg)
	VideoIsNotExistErr        = NewRequestError(VideoIsNotExistErrCode, VideoIsNotExistErrMsg)
	CommentIsNotExistErr      = NewRequestError(CommentIsNotExistErrCode, CommentIsNotExistErrMsg)
)
