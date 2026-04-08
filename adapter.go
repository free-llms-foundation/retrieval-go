package retrieval

import (
	"net/http"

	"github.com/imroc/req/v3"
)

type TLSAdapter struct {
	clientFactory func() *req.Client
}

func NewTLSAdapter(clientFactory func() *req.Client) *TLSAdapter {
	return &TLSAdapter{
		clientFactory: clientFactory,
	}
}

func (t *TLSAdapter) Do(req *http.Request) (*http.Response, error) {
	return t.clientFactory().GetClient().Do(req)
}

func (t *TLSAdapter) ReqClient() *req.Client {
	return t.clientFactory()
}
