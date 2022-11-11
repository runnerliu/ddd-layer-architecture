package cache

import "golang.org/x/sync/singleflight"

// LocalCache 本地缓存比如进程内内存缓存
type LocalCache interface {
	// Get 用于从缓存中获取函数的执行结果，如果未获取到，那么执行函数，并将结果保存到缓存
	Get(key interface{}, f func() (interface{}, error), dest interface{}, sfg *singleflight.Group) (interface{}, error)
}
