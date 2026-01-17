package retrieval

import (
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout   = time.Second * 30
	defaultBaseURL   = "https://lite.duckduckgo.com/lite/"
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"

	defaultMaxErrBodyBytes = 1024 * 64   // 64KB limit for error bodies
	defaultMaxBodyBytes    = 1024 * 1024 // 1MB default for regular bodies
	defaultRespectRobots   = true
)

type Config struct {
	HTTPClient      *http.Client
	UserAgent       string
	Headers         [][2]string
	Parser          Parser
	BaseURL         string
	Timeout         time.Duration
	MaxErrBodyBytes int64
	MaxBodyBytes    int64
	RespectRobots   bool
}

type Parser interface {
	Parse(reader io.ReadCloser) ([]Page, error)
}

func DefaultConfig() Config {
	return Config{
		HTTPClient:      nil,
		Headers:         defaultHeaders,
		UserAgent:       defaultUserAgent,
		Parser:          &DefaultDDGParser{},
		BaseURL:         defaultBaseURL,
		Timeout:         0,
		MaxErrBodyBytes: defaultMaxErrBodyBytes,
		MaxBodyBytes:    defaultMaxBodyBytes,
		RespectRobots:   defaultRespectRobots,
	}
}
