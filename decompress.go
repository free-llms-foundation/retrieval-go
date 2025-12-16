package retrieval

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

type decompressor func(io.Reader) (io.ReadCloser, error)

var decompressors = map[string]decompressor{
	"gzip":    gzipDecompressor,
	"br":      brotliDecompressor,
	"zstd":    zstdDecompressor,
	"deflate": deflateDecompressor,
}

func (*Client) getDecodedReader(resp *http.Response) (io.ReadCloser, error) {
	contentEncoding := resp.Header.Get("Content-Encoding")
	if decompressor, exists := decompressors[strings.ToLower(strings.TrimSpace(contentEncoding))]; exists {
		return decompressor(resp.Body)
	}

	return resp.Body, nil
}

func gzipDecompressor(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}

func brotliDecompressor(r io.Reader) (io.ReadCloser, error) {
	return io.NopCloser(brotli.NewReader(r)), nil
}

func zstdDecompressor(r io.Reader) (io.ReadCloser, error) {
	reader, err := zstd.NewReader(r)
	if err != nil {
		return nil, err
	}
	return reader.IOReadCloser(), nil
}

func deflateDecompressor(r io.Reader) (io.ReadCloser, error) {
	return flate.NewReader(r), nil
}
