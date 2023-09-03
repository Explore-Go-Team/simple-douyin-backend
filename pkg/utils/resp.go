package utils

import (
	"errors"
	"simple-douyin-backend/pkg/constants"
)

type BaseResp struct {
	StatusCode int32
	StatusMsg  string
}

// ConvertToResp convert error and build BaseResp
func ConvertToResp(err error) *BaseResp {
	e := constants.RequestError{}
	if errors.As(err, &e) {
		return &BaseResp{
			StatusCode: e.Code,
			StatusMsg:  e.Msg,
		}
	}

	e = constants.ServiceErr.WithMessage(err.Error())
	return &BaseResp{
		StatusCode: e.Code,
		StatusMsg:  e.Msg,
	}
}
