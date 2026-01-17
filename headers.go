package retrieval

var defaultHeaders = [][2]string{
	// --- Client Hints ---
	{"sec-ch-ua", `"Chromium";v="133", "Google Chrome";v="133", "Not(A:Brand";v="99"`},
	{"sec-ch-ua-mobile", "?0"},
	{"sec-ch-ua-platform", `"Windows"`},

	// --- Main headers ---
	{"Upgrade-Insecure-Requests", "1"},

	// --- Accept headers (grouped) ---
	{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
	{"Accept-Encoding", "gzip, deflate, br, zstd"},
	{"Accept-Language", "en-US,en;q=0.9"},

	// --- Fetch Metadata ---
	{"Sec-Fetch-Site", "none"},
	{"Sec-Fetch-Mode", "navigate"},
	{"Sec-Fetch-User", "?1"},
	{"Sec-Fetch-Dest", "document"},

	// --- Other headers ---
	{"Priority", "u=0, i"},
}
