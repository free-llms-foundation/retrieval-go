package retrieval

import (
	"net/http"

	fhttp "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

type TLSAdapter struct {
	client tls_client.HttpClient
}

func NewTLSAdapter(client tls_client.HttpClient) *TLSAdapter {
	return &TLSAdapter{
		client: client,
	}
}

func (t *TLSAdapter) Do(req *http.Request) (*http.Response, error) {
	fReq, err := fhttp.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		return nil, err
	}

	fReq = fReq.WithContext(req.Context())
	fReq.ContentLength = req.ContentLength
	fReq.Close = req.Close
	if req.Host != "" {
		fReq.Host = req.Host
	}
	fReq.GetBody = req.GetBody

	for k, vv := range req.Header {
		for _, v := range vv {
			fReq.Header.Add(k, v)
		}
	}

	fResp, err := t.client.Do(fReq)
	if err != nil {
		return nil, err
	}

	resp := &http.Response{
		Status:           fResp.Status,
		StatusCode:       fResp.StatusCode,
		Proto:            fResp.Proto,
		ProtoMajor:       fResp.ProtoMajor,
		ProtoMinor:       fResp.ProtoMinor,
		Body:             fResp.Body,
		TransferEncoding: fResp.TransferEncoding,
		Header:           http.Header(fResp.Header),
		Request:          req,
		Close:            fResp.Close,
		Uncompressed:     fResp.Uncompressed,
		Trailer:          http.Header(fResp.Trailer),
		ContentLength:    fResp.ContentLength,
	}

	if fResp.Uncompressed {
		resp.Header.Del("Content-Encoding")
	}

	return resp, nil
}
