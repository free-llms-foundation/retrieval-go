package retrieval

import (
	"errors"
	"io"
	"time"
)

const (
	defaultTimeout   = 30 * time.Second
	defaultBaseURL   = "https://lite.duckduckgo.com/lite/"
	defaultFavicon   = "https://www.google.com/s2/favicons?domain=%s"
	defaultImagesURL = "https://www.bing.com/images/search"

	defaultMaxErrBodyBytes     = 1024 * 64   // 64KB limit for error bodies
	defaultMaxBodyBytes        = 1024 * 1024 // 1MB default for regular bodies
	defaultMaxIdleConnsPerHost = 32
	defaultCommonRetryCount    = 0
)

var (
	ErrRobotsDenied = errors.New("robots.txt denied")
)

type Config struct {
	HTTPClient            HTTPClient
	SearchParser          SearchParser
	ImageParser           ImageParser
	BaseURL               string
	ImagesURL             string
	Timeout               time.Duration
	MaxErrBodyBytes       int64
	MaxBodyBytes          int64
	MaxIdleConnsPerHost   int
	Proxy                 string
	ProxyFactory          func() string
	DisableKeepAlive      bool
	CommonRetryCount      int
	EnableForceHTTP1      bool
	EnableDumpAll         bool
	EnableBrowserRotation bool
}

type SearchParser interface {
	Parse(reader io.ReadCloser) ([]Page, error)
}

type ImageParser interface {
	Parse(reader io.ReadCloser) ([]Image, error)
}

func DefaultConfig() *Config {
	return &Config{
		HTTPClient:            nil,
		SearchParser:          &DefaultDDGParser{},
		ImageParser:           &BingImagesParser{},
		BaseURL:               defaultBaseURL,
		ImagesURL:             defaultImagesURL,
		Timeout:               defaultTimeout,
		MaxErrBodyBytes:       defaultMaxErrBodyBytes,
		MaxBodyBytes:          defaultMaxBodyBytes,
		MaxIdleConnsPerHost:   defaultMaxIdleConnsPerHost,
		Proxy:                 "",
		ProxyFactory:          nil,
		DisableKeepAlive:      false,
		CommonRetryCount:      defaultCommonRetryCount,
		EnableForceHTTP1:      false,
		EnableDumpAll:         false,
		EnableBrowserRotation: false,
	}
}
