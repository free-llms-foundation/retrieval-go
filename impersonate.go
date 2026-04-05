package retrieval

import (
	req "github.com/imroc/req/v3"
	utls "github.com/refraction-networking/utls"
)

func profileChrome133_Win(c *req.Client) *req.Client {
	return c.ImpersonateChrome().
		SetTLSFingerprint(utls.HelloChrome_133).
		SetCommonHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36").
		SetCommonHeader("sec-ch-ua", `"Google Chrome";v="133", "Chromium";v="133", "Not(A:Brand";v="24"`).
		SetCommonHeader("sec-ch-ua-platform", `"Windows"`).
		SetCommonHeader("accept-encoding", "gzip, deflate, br, zstd").
		SetCommonHeader("accept-language", "en-US,en;q=0.9").
		SetCommonHeader("pragma", "").
		SetCommonHeader("cache-control", "")
}

func profileChrome133_Mac(c *req.Client) *req.Client {
	return c.ImpersonateChrome().
		SetTLSFingerprint(utls.HelloChrome_133).
		SetCommonHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36").
		SetCommonHeader("sec-ch-ua", `"Google Chrome";v="133", "Chromium";v="133", "Not(A:Brand";v="24"`).
		SetCommonHeader("sec-ch-ua-platform", `"macOS"`).
		SetCommonHeader("accept-encoding", "gzip, deflate, br, zstd").
		SetCommonHeader("accept-language", "en-US,en;q=0.9").
		SetCommonHeader("pragma", "").
		SetCommonHeader("cache-control", "")
}

func profileChrome133_Linux(c *req.Client) *req.Client {
	return c.ImpersonateChrome().
		SetTLSFingerprint(utls.HelloChrome_133).
		SetCommonHeader("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36").
		SetCommonHeader("sec-ch-ua", `"Google Chrome";v="133", "Chromium";v="133", "Not(A:Brand";v="24"`).
		SetCommonHeader("sec-ch-ua-platform", `"Linux"`).
		SetCommonHeader("accept-encoding", "gzip, deflate, br, zstd").
		SetCommonHeader("accept-language", "en-US,en;q=0.9").
		SetCommonHeader("pragma", "").
		SetCommonHeader("cache-control", "")
}

func profileChrome131_Win(c *req.Client) *req.Client {
	return c.ImpersonateChrome().
		SetTLSFingerprint(utls.HelloChrome_131).
		SetCommonHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36").
		SetCommonHeader("sec-ch-ua", `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="8"`).
		SetCommonHeader("sec-ch-ua-platform", `"Windows"`).
		SetCommonHeader("accept-encoding", "gzip, deflate, br, zstd").
		SetCommonHeader("accept-language", "en-US,en;q=0.9").
		SetCommonHeader("pragma", "").
		SetCommonHeader("cache-control", "")
}

func profileChrome131_Mac(c *req.Client) *req.Client {
	return c.ImpersonateChrome().
		SetTLSFingerprint(utls.HelloChrome_131).
		SetCommonHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36").
		SetCommonHeader("sec-ch-ua", `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="8"`).
		SetCommonHeader("sec-ch-ua-platform", `"macOS"`).
		SetCommonHeader("accept-encoding", "gzip, deflate, br, zstd").
		SetCommonHeader("accept-language", "en-US,en;q=0.9").
		SetCommonHeader("pragma", "").
		SetCommonHeader("cache-control", "")
}

func profileChrome120PQ_Win(c *req.Client) *req.Client {
	return c.ImpersonateChrome().
		SetTLSFingerprint(utls.HelloChrome_120_PQ).
		SetCommonHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36").
		SetCommonHeader("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`).
		SetCommonHeader("sec-ch-ua-platform", `"Windows"`).
		SetCommonHeader("accept-language", "en-US,en;q=0.9").
		SetCommonHeader("pragma", "").
		SetCommonHeader("cache-control", "")
}

func profileEdge131_Win(c *req.Client) *req.Client {
	return c.ImpersonateChrome().
		SetTLSFingerprint(utls.HelloEdge_85).
		SetCommonHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0").
		SetCommonHeader("sec-ch-ua", `"Microsoft Edge";v="131", "Chromium";v="131", "Not_A Brand";v="8"`).
		SetCommonHeader("sec-ch-ua-platform", `"Windows"`).
		SetCommonHeader("accept-encoding", "gzip, deflate, br, zstd").
		SetCommonHeader("accept-language", "en-US,en;q=0.9").
		SetCommonHeader("pragma", "").
		SetCommonHeader("cache-control", "")
}

func profileFirefox120_Win(c *req.Client) *req.Client {
	return c.ImpersonateFirefox().
		SetCommonHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0").
		SetCommonHeader("accept-encoding", "gzip, deflate, br").
		SetCommonHeader("accept-language", "en-US,en;q=0.5").
		SetCommonHeader("sec-fetch-site", "none")
}

func profileFirefox120_Mac(c *req.Client) *req.Client {
	return c.ImpersonateFirefox().
		SetCommonHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:120.0) Gecko/20100101 Firefox/120.0").
		SetCommonHeader("accept-encoding", "gzip, deflate, br").
		SetCommonHeader("accept-language", "en-US,en;q=0.5").
		SetCommonHeader("sec-fetch-site", "none")
}

func profileSafari184_Mac(c *req.Client) *req.Client {
	return c.ImpersonateSafari().
		SetCommonHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.4 Safari/605.1.15").
		SetCommonHeader("accept-encoding", "gzip, deflate, br").
		SetCommonHeader("accept-language", "en-US,en;q=0.9").
		SetCommonHeader("sec-fetch-site", "none")
}

func profileIOS184_iPhone(c *req.Client) *req.Client {
	return c.ImpersonateSafari().
		SetTLSFingerprint(utls.HelloIOS_14).
		SetCommonHeader("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.4 Mobile/15E148 Safari/604.1").
		SetCommonHeader("accept-encoding", "gzip, deflate, br").
		SetCommonHeader("accept-language", "en-US,en;q=0.9").
		SetCommonHeader("sec-fetch-site", "none")
}

var profilePool = []func(*req.Client) *req.Client{
	profileChrome133_Win,
	profileChrome133_Mac,
	profileChrome133_Linux,
	profileChrome131_Win,
	profileChrome131_Mac,
	profileChrome120PQ_Win,
	profileEdge131_Win,
	profileFirefox120_Win,
	profileFirefox120_Mac,
	profileSafari184_Mac,
	profileIOS184_iPhone,
}
