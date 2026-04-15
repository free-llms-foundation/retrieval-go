package retrieval

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

// ── helpers ──────────────────────────────────────────────────────────────────

const decompressPayload = "hello decompressor world — unicode check: 你好\n"

func makeDecompressResponse(encoding string, body []byte) *http.Response {
	return &http.Response{
		Header: http.Header{"Content-Encoding": []string{encoding}},
		Body:   io.NopCloser(bytes.NewReader(body)),
	}
}

func compressGzip(data []byte) []byte {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, _ = w.Write(data)
	_ = w.Close()
	return buf.Bytes()
}

func compressBrotli(data []byte) []byte {
	var buf bytes.Buffer
	w := brotli.NewWriter(&buf)
	_, _ = w.Write(data)
	_ = w.Close()
	return buf.Bytes()
}

func compressZstd(data []byte) []byte {
	var buf bytes.Buffer
	w, _ := zstd.NewWriter(&buf)
	_, _ = w.Write(data)
	_ = w.Close()
	return buf.Bytes()
}

func compressDeflate(data []byte) []byte {
	var buf bytes.Buffer
	w, _ := flate.NewWriter(&buf, flate.DefaultCompression)
	_, _ = w.Write(data)
	_ = w.Close()
	return buf.Bytes()
}

// ── tests ─────────────────────────────────────────────────────────────────────

func TestGetDecodedReader(t *testing.T) {
	t.Parallel()

	payload := []byte(decompressPayload)

	tests := []struct {
		name     string
		isEmpty  bool
		encoding string
		body     func() []byte
	}{
		{
			name:     "no Content-Encoding passthrough",
			encoding: "",
			body:     func() []byte { return payload },
		},
		{
			name:     "gzip",
			encoding: "gzip",
			body:     func() []byte { return compressGzip(payload) },
		},
		{
			name:     "brotli",
			encoding: "br",
			body:     func() []byte { return compressBrotli(payload) },
		},
		{
			name:     "zstd",
			encoding: "zstd",
			body:     func() []byte { return compressZstd(payload) },
		},
		{
			name:     "deflate",
			encoding: "deflate",
			body:     func() []byte { return compressDeflate(payload) },
		},
		{
			name:     "unknown encoding passthrough",
			encoding: "identity",
			body:     func() []byte { return payload },
		},
		{
			name:     "Content-Encoding header whitespace trimmed",
			encoding: "  gzip  ",
			body:     func() []byte { return compressGzip(payload) },
		},
		{
			name:     "Empty Body",
			encoding: "",
			isEmpty:  true,
			body:     func() []byte { return []byte{} },
		},
	}

	c := &Client{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp := makeDecompressResponse(tt.encoding, tt.body())
			reader, err := c.getDecodedReader(resp)
			if err != nil {
				t.Fatalf("getDecodedReader: %v", err)
			}
			defer reader.Close()

			got, err := io.ReadAll(reader)
			if err != nil {
				t.Fatalf("read decoded body: %v", err)
			}

			if tt.isEmpty {
				if len(got) != 0 {
					t.Errorf("expected empty body, got %q", got)
				}
			} else {
				if string(got) != decompressPayload {
					t.Errorf("decoded body = %q, want %q", got, decompressPayload)
				}
			}
		})
	}
}
