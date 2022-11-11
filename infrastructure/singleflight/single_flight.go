package singleflight

import (
	"golang.org/x/sync/singleflight"
)

// NewSingleFlightGroup 提供了重复函数调用的合并机制
func NewSingleFlightGroup() *singleflight.Group {
	return new(singleflight.Group)
}
