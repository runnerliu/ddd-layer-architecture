package consts

import "errors"

const ConfigName = "config"

var (
	// 错误响应码
	ErrCodeSuccess = 0
	ErrCodeParams  = 1001

	// 错误响应信息
	ErrMsgSuccess      = "success"
	ErrMsgParams       = "params error"
	ErrMsgResponseCode = errors.New("invalid response code")
)
