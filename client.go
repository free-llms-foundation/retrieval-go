package retrieval

import (
	"net/http"
	"net/http/cookiejar"
)

type Client struct {
	client          *http.Client
	userAgent       string
	headers         [][2]string
	parser          Parser
	baseURL         string
	maxErrBodyBytes int64
	maxBodyBytes    int64
	respectRobots   bool
}

func New(opts ...Option) *Client {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	return NewWithConfig(cfg)
}

func NewWithConfig(cfg Config) *Client {
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
		if cfg.Timeout != 0 {
			httpClient.Timeout = cfg.Timeout
		} else {
			httpClient.Timeout = defaultTimeout
		}
		httpClient.Jar, _ = cookiejar.New(nil)

	} else if cfg.Timeout != 0 {
		httpClient.Timeout = cfg.Timeout
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	headers := cfg.Headers
	if headers == nil {
		headers = defaultHeaders
	}

	userAgent := cfg.UserAgent
	if userAgent == "" {
		userAgent = defaultUserAgent
	}

	headersCopy := append([][2]string(nil), headers...)

	parser := cfg.Parser
	if parser == nil {
		parser = &DefaultDDGParser{}
	}

	maxErrBytes := cfg.MaxErrBodyBytes
	if maxErrBytes <= 0 {
		maxErrBytes = defaultMaxErrBodyBytes
	}

	maxBodyBytes := cfg.MaxBodyBytes
	if maxBodyBytes <= 0 {
		maxBodyBytes = defaultMaxBodyBytes
	}

	return &Client{
		client:          httpClient,
		headers:         headersCopy,
		userAgent:       userAgent,
		parser:          parser,
		baseURL:         baseURL,
		maxErrBodyBytes: maxErrBytes,
		maxBodyBytes:    maxBodyBytes,
		respectRobots:   cfg.RespectRobots,
	}
}
