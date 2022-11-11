package response

import (
	"context"
	"ddd-demo/common/consts"
)

// Response 响应体
type Response struct {
	// Code 应用状态码，0 表示成功，大于 0 的整数表示失败
	Code int `json:"code" example:"0"`
	// Message 易于阅读的信息
	Message string `json:"message" example:"success"`
	// Data 数据，当 Code 不为 0 时，Data 为 nil
	Data interface{} `json:"data"`
}

// NewErrorResponse 生成 Error Response
func NewErrorResponse(ctx context.Context, code int, err error) *Response {
	if code <= 0 {
		code = 999
		err = consts.ErrMsgResponseCode
	}
	response := &Response{
		Code:    code,
		Message: err.Error(),
		Data:    nil,
	}
	return response
}

// NewSuccessResponse 生成 Success Response
func NewSuccessResponse(ctx context.Context, data interface{}) *Response {
	response := &Response{
		Code:    consts.ErrCodeSuccess,
		Message: consts.ErrMsgSuccess,
		Data:    data,
	}
	return response
}
