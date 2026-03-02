package retrieval

import (
	"net/http"
	"net/url"

	"github.com/imroc/req/v3"
)

type TLSAdapter struct {
	client *req.Client
}

func NewTLSAdapter(client *req.Client, proxyFactory func() string) *TLSAdapter {
	if proxyFactory != nil {
		client.SetProxy(func(r *http.Request) (*url.URL, error) {
			p := proxyFactory()
			if p == "" {
				return nil, nil
			}
			return url.Parse(p)
		})
	}
	return &TLSAdapter{
		client: client,
	}
}

func (t *TLSAdapter) Do(req *http.Request) (*http.Response, error) {
	return t.client.Do(req)
}

func (t *TLSAdapter) ReqClient() *req.Client {
	return t.client
}
