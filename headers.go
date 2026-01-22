package retrieval

var defaultHeaders = [][2]string{
	{"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:146.0) Gecko/20100101 Firefox/146.0"},
	{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8"},
	{"Accept-Language", "en-US,en;q=0.5"},
	{"Accept-Encoding", "gzip, deflate, br, zstd"},
	{"Upgrade-Insecure-Requests", "1"},
	{"Sec-Fetch-Dest", "document"},
	{"Sec-Fetch-Mode", "navigate"},
	{"Sec-Fetch-Site", "none"},
	{"Sec-Fetch-User", "?1"},
	{"Te", "trailers"},
	{"Connection", "keep-alive"},
}
