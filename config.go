package retrieval

import (
	"errors"
	"io"
	"time"
)

const (
	defaultTimeout = 30
	defaultBaseURL = "https://lite.duckduckgo.com/lite/"
	defaultFavicon = "https://www.google.com/s2/favicons?domain=%s"

	defaultMaxErrBodyBytes = 1024 * 64   // 64KB limit for error bodies
	defaultMaxBodyBytes    = 1024 * 1024 // 1MB default for regular bodies
	defaultRespectRobots   = true
)

var (
	ErrRobotsDenied = errors.New("robots.txt denied")
)

type Config struct {
	HTTPClient      HTTPClient
	Headers         [][2]string
	Parser          Parser
	BaseURL         string
	Timeout         time.Duration
	MaxErrBodyBytes int64
	MaxBodyBytes    int64
	Proxy           string
	RespectRobots   bool
}

type Parser interface {
	Parse(reader io.ReadCloser) ([]Page, error)
}

func DefaultConfig() *Config {
	return &Config{
		HTTPClient:      nil,
		Headers:         defaultHeaders,
		Parser:          &DefaultDDGParser{},
		BaseURL:         defaultBaseURL,
		Timeout:         defaultTimeout,
		MaxErrBodyBytes: defaultMaxErrBodyBytes,
		MaxBodyBytes:    defaultMaxBodyBytes,
		Proxy:           "",
		RespectRobots:   defaultRespectRobots,
	}
}
