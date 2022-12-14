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

	ErrCacheResultTypeMismatched = errors.New("result type mismatched")
	ErrESQueryIndexData          = errors.New("query index data error")
)
