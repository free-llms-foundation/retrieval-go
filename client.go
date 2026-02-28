package retrieval

import (
	"net/http"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	client          HTTPClient
	headers         [][2]string
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
		options := []tls_client.HttpClientOption{
			tls_client.WithClientProfile(profiles.Firefox_120),
			tls_client.WithRandomTLSExtensionOrder(),
		}

		transportOptions := &tls_client.TransportOptions{
			MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
			DisableKeepAlives:   cfg.DisableKeepAlive,
		}
		options = append(options, tls_client.WithTransportOptions(transportOptions))

		if cfg.Timeout > 0 {
			options = append(options, tls_client.WithTimeoutSeconds(int(cfg.Timeout)))
		} else {
			options = append(options, tls_client.WithTimeoutSeconds(defaultTimeout))
		}

		if cfg.Proxy != "" {
			if cfg.ProxyFactory == nil {
				options = append(options, tls_client.WithProxyUrl(cfg.Proxy))
			}
		}

		httpTLSClient, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
		if err != nil {
			return nil, err
		}

		httpClient = NewTLSAdapter(httpTLSClient, cfg.ProxyFactory)
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	headers := cfg.Headers
	if headers == nil {
		headers = defaultHeaders
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
		headers:         headersCopy,
		parser:          parser,
		baseURL:         baseURL,
		maxErrBodyBytes: maxErrBytes,
		maxBodyBytes:    maxBodyBytes,
		converter:       converter,
	}, nil
}
