package config

import (
	"ddd-demo/common/consts"
	"os"
	"path"
	"sync"

	"github.com/spf13/viper"
)

// ViperConfiguration 使用 viper 类库实现的配置读取服务
type ViperConfiguration struct {
	v *viper.Viper
}

// GetBool 获取 bool 类型配置
func (v *ViperConfiguration) GetBool(key string) bool {
	return v.v.GetBool(key)
}

// GetInt 获取 int 类型配置
func (v *ViperConfiguration) GetInt(key string) int {
	return v.v.GetInt(key)
}

// GetFloat64 获取 float64 类型配置
func (v *ViperConfiguration) GetFloat64(key string) float64 {
	return v.v.GetFloat64(key)
}

// GetString 获取 string 类型配置
func (v *ViperConfiguration) GetString(key string) string {
	return v.v.GetString(key)
}

// GetStringSlice 获取 string 数组类型配置
func (v *ViperConfiguration) GetStringSlice(key string) []string {
	return v.v.GetStringSlice(key)
}

// GetStringMap 获取 map 类型配置
func (v *ViperConfiguration) GetStringMap(key string) map[string]interface{} {
	return v.v.GetStringMap(key)
}

// GetStringMapString 获取 map 类型配置
func (v *ViperConfiguration) GetStringMapString(key string) map[string]string {
	return v.v.GetStringMapString(key)
}

// GetStringMapStringSlice 获取 map 类型配置
func (v *ViperConfiguration) GetStringMapStringSlice(key string) map[string][]string {
	return v.v.GetStringMapStringSlice(key)
}

// NewViperConfiguration 创建 viper 配置实例
func NewViperConfiguration(configName string, configPaths ...string) (Configuration, error) {
	v := &ViperConfiguration{v: viper.New()}
	v.v.SetConfigType("yaml")
	v.v.SetConfigName(configName)
	for _, configPath := range configPaths {
		v.v.AddConfigPath(configPath)
	}
	if err := v.v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}

// NewDefaultYamlConfiguration 创建默认的配置实例
func NewDefaultYamlConfiguration() (Configuration, error) {
	var paths []string
	if wd, err := os.Getwd(); err == nil {
		paths = append(paths, []string{wd, path.Join(wd, "config")}...)
	} else {
		return nil, err
	}
	return NewViperConfiguration(consts.ConfigName, paths...)
}

var (
	defaultConfigurationOnce sync.Once
	defaultConfiguration     Configuration
	defaultConfigurationErr  error
)

// NewYamlConfiguration 创建配置读取服务实例
func NewYamlConfiguration() (Configuration, error) {
	defaultConfigurationOnce.Do(func() {
		defaultConfiguration, defaultConfigurationErr = NewDefaultYamlConfiguration()
	})
	return defaultConfiguration, defaultConfigurationErr
}
