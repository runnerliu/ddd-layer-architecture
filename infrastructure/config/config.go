package config

// Configuration 服务配置需要实现的接口
type Configuration interface {
	// GetBool 获取 bool 类型配置
	GetBool(key string) bool
	// GetInt 获取 int 类型配置
	GetInt(key string) int
	// GetFloat64 获取 float64 类型配置
	GetFloat64(key string) float64
	// GetString 获取 string 类型配置
	GetString(key string) string
	// GetStringSlice 获取 string 数组类型配置
	GetStringSlice(key string) []string
	// GetStringMap 获取 map 类型配置
	GetStringMap(key string) map[string]interface{}
	// GetStringMapString 获取 map 类型配置
	GetStringMapString(key string) map[string]string
	// GetStringMapStringSlice 获取 map 类型配置
	GetStringMapStringSlice(key string) map[string][]string
}
