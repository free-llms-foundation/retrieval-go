package retrieval

import (
	"time"
)

type Option func(*Config)

func WithMaxIdleConnsPerHost(maxIdleConnsPerHost int) Option {
	return func(cfg *Config) {
		cfg.MaxIdleConnsPerHost = maxIdleConnsPerHost
	}
}

func WithEnableForceHTTP1() Option {
	return func(cfg *Config) {
		cfg.EnableForceHTTP1 = true
	}
}

func WithEnableDumpAll() Option {
	return func(cfg *Config) {
		cfg.EnableDumpAll = true
	}
}

func WithCommonRetryCount(commonRetryCount int) Option {
	return func(cfg *Config) {
		cfg.CommonRetryCount = commonRetryCount
	}
}

func WithClient(httpClient HTTPClient) Option {
	return func(cfg *Config) {
		cfg.HTTPClient = httpClient
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(cfg *Config) {
		cfg.Timeout = timeout
	}
}

func WithParser(parser Parser) Option {
	return func(cfg *Config) {
		cfg.Parser = parser
	}
}

func WithBaseURL(url string) Option {
	return func(cfg *Config) {
		cfg.BaseURL = url
	}
}

func WithMaxErrBodyBytes(maxBytes int64) Option {
	return func(cfg *Config) {
		cfg.MaxErrBodyBytes = maxBytes
	}
}

func WithDisableKeepAlive() Option {
	return func(cfg *Config) {
		cfg.DisableKeepAlive = true
	}
}

func WithMaxBodyBytes(maxBytes int64) Option {
	return func(cfg *Config) {
		cfg.MaxBodyBytes = maxBytes
	}
}

func WithProxy(proxy string) Option {
	return func(cfg *Config) {
		cfg.Proxy = proxy
	}
}

func WithProxyFactory(proxyFactory func() string) Option {
	return func(cfg *Config) {
		cfg.ProxyFactory = proxyFactory
	}
}

func WithBrowserRotation() Option {
	return func(cfg *Config) {
		cfg.EnableBrowserRotation = true
	}
}
