package http

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

// RetryStatusCode 用于指定需要进行重试的状态码
type RetryStatusCode struct {
	Start int `mapstructure:"start"`
	End   int `mapstructure:"end"`
}

// ClientCertificate 代表客户端证书
type ClientCertificate struct {
	Public  string `mapstructure:"public"`
	Private string `mapstructure:"private"`
}

// 下面是默认值的定义
const (
	DefaultConnectTimeoutMillis        = 5000
	DefaultKeepAliveTimeoutMillis      = 10000
	DefaultFallbackDelayTimeoutMillis  = 1000
	DefaultMaxIdleConnsPerHost         = 2
	DefaultMaxIdleConns                = 4
	DefaultIdleConnTimeoutMillis       = 30000
	DefaultExpectContinueTimeoutMillis = 1000
	DefaultTimeoutMillis               = 10000
	DefaultRetryCount                  = 3
	DefaultRetryWaitTimeMillis         = 100
	DefaultRetryMaxWaitTimeMillis      = 5000
)

// DefaultRetryStatusCodes 可重试的 HTTP 状态码
var DefaultRetryStatusCodes = []*RetryStatusCode{{Start: 501, End: 999}}

// RestyClientConfig 代表 Resty Client 的配置
type RestyClientConfig struct {
	// 是否禁止对请求进行追踪
	DisableTrace bool `mapstructure:"disable_trace"`
	// 等待连接完成的总时间
	ConnectTimeoutMillis int `mapstructure:"connect_timeout_millis"`
	// Keep-Alive 探查的时间间隔
	KeepAliveTimeoutMillis int `mapstructure:"keepalive_timeout_millis"`
	// 在假定 IPV6 不可用，降级到 IPV4 之前，等待 IPV6 成功的总时间
	FallbackDelayTimeoutMillis int `mapstructure:"fallback_delay_timeout_millis"`
	// 每个 Host 的最大空闲连接数
	MaxIdleConnsPerHost int `mapstructure:"max_idle_conns_host"`
	// Client 与所有 Host 的最大空闲连接数
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	// 空闲连接的超时时间，超时的连接将被关闭
	IdleConnTimeoutMillis int `mapstructure:"idle_conn_timeout_millis"`
	// 发送完 Expect: 100-continue 请求后，等待 server 应答的超时时间
	ExpectContinueTimeoutMillis int `mapstructure:"expect_continue_timeout_millis"`
	// 请求的超时时间
	TimeoutMillis int `mapstructure:"timeout_millis"`
	// 设置代理 URL 和端口
	Proxy string `mapstructure:"proxy"`
	// 重试次数
	RetryCount int `mapstructure:"retry_count"`
	// 初始的重试等待时间
	RetryWaitTimeMillis int `mapstructure:"retry_wait_time_millis"`
	// 最大的重试等待时间
	RetryMaxWaitTimeMillis int `mapstructure:"retry_max_wait_time_millis"`
	// 进行重试的状态码
	RetryStatusCodes []*RetryStatusCode `mapstructure:"retry_status_codes"`
	// 信任的根证书列表
	RootCertificates []string `mapstructure:"root_certificates"`
	// 客户端证书列表
	ClientCertificates []*ClientCertificate `mapstructure:"client_certificates"`
	// 是否开启 Debug 模式
	DebugMode bool `mapstructure:"debug_mode"`
}

// intOrDefault 当 i 等于 0 时，返回 defaultValue；否则返回 i
func intOrDefault(i, defaultValue int) int {
	if i == 0 {
		return defaultValue
	}
	return i
}

// NewRestyClient 用于构造 resty.Client 对象
func NewRestyClient(config *RestyClientConfig) *resty.Client {
	client := resty.New()
	if !config.DisableTrace {
		client = client.EnableTrace()
	}
	client.SetDebug(config.DebugMode)
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Duration(
				intOrDefault(config.ConnectTimeoutMillis, DefaultConnectTimeoutMillis),
			) * time.Millisecond,
			KeepAlive: time.Duration(
				intOrDefault(config.KeepAliveTimeoutMillis, DefaultKeepAliveTimeoutMillis),
			) * time.Millisecond,
			FallbackDelay: time.Duration(
				intOrDefault(config.FallbackDelayTimeoutMillis, DefaultFallbackDelayTimeoutMillis),
			) * time.Millisecond,
		}).DialContext,
		MaxIdleConnsPerHost: intOrDefault(config.MaxIdleConnsPerHost, DefaultMaxIdleConnsPerHost),
		MaxIdleConns:        intOrDefault(config.MaxIdleConns, DefaultMaxIdleConns),
		IdleConnTimeout: time.Duration(
			intOrDefault(config.IdleConnTimeoutMillis, DefaultIdleConnTimeoutMillis),
		) * time.Millisecond,
		ExpectContinueTimeout: time.Duration(
			intOrDefault(config.ExpectContinueTimeoutMillis, DefaultExpectContinueTimeoutMillis),
		) * time.Millisecond,
	}
	client.SetTransport(transport)
	client.SetTimeout(time.Duration(intOrDefault(config.TimeoutMillis, DefaultTimeoutMillis)) * time.Millisecond)
	if config.Proxy != "" {
		client.SetProxy(config.Proxy)
	}
	client.SetRetryCount(intOrDefault(config.RetryCount, DefaultRetryCount))
	client.SetRetryWaitTime(time.Duration(
		intOrDefault(config.RetryWaitTimeMillis, DefaultRetryWaitTimeMillis),
	) * time.Millisecond)
	client.SetRetryMaxWaitTime(time.Duration(
		intOrDefault(config.RetryMaxWaitTimeMillis, DefaultRetryMaxWaitTimeMillis),
	) * time.Millisecond)
	client.AddRetryCondition(func(response *resty.Response, err error) bool {
		// input: non-nil Response OR request execution error
		if err != nil {
			return true
		}

		retryStatusCodes := DefaultRetryStatusCodes
		if len(config.RetryStatusCodes) > 0 {
			retryStatusCodes = config.RetryStatusCodes
		}
		for _, retryStatusCode := range retryStatusCodes {
			if retryStatusCode == nil {
				continue
			}
			if response.StatusCode() >= retryStatusCode.Start && response.StatusCode() <= retryStatusCode.End {
				return true
			}
		}

		return false
	})

	for _, rootCertificate := range config.RootCertificates {
		client.SetRootCertificate(rootCertificate)
	}

	var certificates []tls.Certificate
	for _, clientCertificate := range config.ClientCertificates {
		certificate, err := tls.LoadX509KeyPair(clientCertificate.Public, clientCertificate.Private)
		if err != nil {
			continue
		}
		certificates = append(certificates, certificate)
	}
	if len(certificates) > 0 {
		client.SetCertificates(certificates...)
	}

	return client
}

// NewDefaultRestyClient 创建 resty 客户端
func NewDefaultRestyClient() *resty.Client {
	return NewRestyClient(&RestyClientConfig{})
}

var (
	defaultRestyClientOnce sync.Once
	defaultRestyClient     *resty.Client
)

// GetDefaultRestyClient 创建 resty 客户端
func GetDefaultRestyClient() *resty.Client {
	defaultRestyClientOnce.Do(func() {
		defaultRestyClient = NewDefaultRestyClient()
	})
	return defaultRestyClient
}
