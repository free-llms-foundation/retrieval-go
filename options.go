package retrieval

import (
	"time"
)

type Option func(*Config)

func WithHeaders(headers [][2]string) Option {
	return func(cfg *Config) {
		cfg.Headers = headers
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

func WithMaxBodyBytes(maxBytes int64) Option {
	return func(cfg *Config) {
		cfg.MaxBodyBytes = maxBytes
	}
}

func WithRespectRobots(respect bool) Option {
	return func(cfg *Config) {
		cfg.RespectRobots = respect
	}
}

func WithProxy(proxy string) Option {
	return func(cfg *Config) {
		cfg.Proxy = proxy
	}
}
