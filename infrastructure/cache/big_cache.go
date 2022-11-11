package cache

import (
	"ddd-demo/common"
	"ddd-demo/common/consts"
	"ddd-demo/common/serializer"
	"reflect"
	"time"

	"github.com/allegro/bigcache"
	"golang.org/x/sync/singleflight"
)

// BigCacheLocalCache 是基于 BigCache 的 LocalCache 实现
type BigCacheLocalCache struct {
	cache *bigcache.BigCache
}

// Get 先尝试从 BigCache 中读取结果，如果读到，那么反序列化，如果反序列化成功，那么返回缓存的结果。
// 否则执行函数，如果执行成功，那么尝试将执行结果序列化，并保存到 BigCache，最后返回它
func (b *BigCacheLocalCache) Get(
	originKey interface{},
	f func() (interface{}, error),
	dest interface{},
	sfg *singleflight.Group,
) (interface{}, error) {
	keyBuf, err := serializer.GobEncode(originKey)
	if err != nil {
		return nil, err
	}
	// 生成缓存 Key
	key := common.GetMd5(keyBuf)

	wrappedFunc := func() (interface{}, error) {
		// 尝试从 BigCache 中读取结果
		if buf, err := b.cache.Get(key); err == nil {
			if decodingErr := serializer.GobDecode(dest, buf); decodingErr == nil {
				return dest, nil
			}
		}

		// 执行函数
		res, executingErr := f()
		// 如果执行失败，则直接返回
		if executingErr != nil {
			return nil, executingErr
		}
		// 如果函数返回的类型与期望的类型不一致，那么返回错误
		if reflect.TypeOf(res) != reflect.TypeOf(dest) {
			return nil, consts.ErrCacheResultTypeMismatched
		}
		// 尝试序列化，并缓存执行结果
		if buf, err := serializer.GobEncode(res); err == nil {
			_ = b.cache.Set(key, buf)
		}

		return res, nil
	}

	if sfg != nil {
		res, err, _ := sfg.Do(key, wrappedFunc)
		return res, err
	}

	return wrappedFunc()
}

const (
	// BigCacheLifeWindowMS 时长后，缓存条目可被踢除，默认值是 10000 毫秒
	BigCacheLifeWindowMS = 10000
	// BigCacheShards shard 的数量，其值必须是 2 的乘方，默认值为 2
	BigCacheShards = 2
	// BigCacheCleanWindowMS 两次清理过期条目之间的时间间隔。
	// 如果被设置为小于等于 0 的值，那么不会执行任何操作。不建议将其设置为小于 1 秒的值。默认值是 5 秒
	BigCacheCleanWindowMS = 5000
	// BigCacheMaxEntrySize 条目的最大大小，单位是字节。仅用于计算 cache shard 的初始大小，默认值是 10240（10K）
	BigCacheMaxEntrySize = 10240
	// BigCacheMaxEntriesInWindow life window 中的最大条目数。
	// 仅用于计算 cache shard 的初始大小。如果将其设置为正确的值，那么不会发生额外的内存分配。默认值是 10K
	BigCacheMaxEntriesInWindow = 10240
	// BigCacheVerbose 非 0 表示打印关于新的内存分配的信息，0 表示不打印。默认不打印
	BigCacheVerbose = false
	// BigCacheHardMaxCacheSize 缓存大小的限制，单位是 MB。超过该限制时，cache 不会分配更多的内存。防止应用程序消耗掉机器上的所有可用内存。
	// 0 表示不限制。当该限制大于 0，并且达到时，新条目会重写老条目。默认值是 100
	BigCacheHardMaxCacheSize = 100
)

// NewBigCacheLocalCache 用于构造基于 BigCache 的 LocalCache 实现
func NewBigCacheLocalCache() (LocalCache, error) {
	lifeWindow := time.Duration(BigCacheLifeWindowMS) * time.Millisecond
	cleanWindow := time.Duration(BigCacheCleanWindowMS) * time.Millisecond
	config := bigcache.DefaultConfig(lifeWindow)
	config.Shards = BigCacheShards
	config.CleanWindow = cleanWindow
	config.MaxEntrySize = BigCacheMaxEntrySize
	config.MaxEntriesInWindow = BigCacheMaxEntriesInWindow
	config.Verbose = BigCacheVerbose
	config.HardMaxCacheSize = BigCacheHardMaxCacheSize

	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		return nil, err
	}

	return &BigCacheLocalCache{cache: cache}, nil
}
