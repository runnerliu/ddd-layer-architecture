package config

// Configuration 服务配置需要实现的接口
type Configuration interface {
	GetBool(key string) bool
	GetInt(key string) int
	GetFloat64(key string) float64
	GetString(key string) string
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
}
