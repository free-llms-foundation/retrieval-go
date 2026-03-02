package retrieval

import (
	"net/http"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
	"github.com/imroc/req/v3"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	client          HTTPClient
	parser          Parser
	baseURL         string
	maxErrBodyBytes int64
	maxBodyBytes    int64
	converter       *converter.Converter
}

func New(opts ...Option) (*Client, error) {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	return NewWithConfig(cfg)
}

func NewWithConfig(cfg *Config) (*Client, error) {
	var httpClient = cfg.HTTPClient
	if httpClient == nil {
		reqClient := req.C().
			ImpersonateChrome().
			SetCommonHeader("Accept-Language", "en-US,en;q=0.5").
			SetCommonRetryCount(cfg.CommonRetryCount)

		if cfg.EnableForceHTTP1 {
			reqClient.EnableForceHTTP1()
			reqClient.GetTLSClientConfig().NextProtos = []string{"http/1.1"}
		}

		if cfg.EnableDumpAll {
			reqClient.EnableDumpAll()
		}

		if cfg.Timeout > 0 {
			reqClient.SetTimeout(cfg.Timeout)
		} else {
			reqClient.SetTimeout(defaultTimeout)
		}

		// Configure transport options
		transport := reqClient.GetTransport()
		transport.MaxIdleConnsPerHost = cfg.MaxIdleConnsPerHost
		transport.DisableKeepAlives = cfg.DisableKeepAlive

		if cfg.Proxy != "" && cfg.ProxyFactory == nil {
			reqClient.SetProxyURL(cfg.Proxy)
		}

		httpClient = NewTLSAdapter(reqClient, cfg.ProxyFactory)
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

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

	converter := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
			table.NewTablePlugin(
				table.WithSpanCellBehavior(table.SpanBehaviorMirror),
				table.WithNewlineBehavior(table.NewlineBehaviorSkip),
				table.WithCellPaddingBehavior(table.CellPaddingBehaviorNone),
			),
		),
	)

	return &Client{
		client:          httpClient,
		parser:          parser,
		baseURL:         baseURL,
		maxErrBodyBytes: maxErrBytes,
		maxBodyBytes:    maxBodyBytes,
		converter:       converter,
	}, nil
}
