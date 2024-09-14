package analysis

import (
	"fmt"
)

func (check *SubdomainCheck) CORS() {
	url := MakeUrl(HTTP(Secure), check.Subdomain)
	check.testCors(url, "Origin") // GET
}

func (check *SubdomainCheck) cookieInjection() {
	// session hijacking, xss
	testCookie := "tookie=jzqvtyxkplra"
	url := MakeUrl(HTTP(Secure), check.Subdomain)
	if check.isPayloadReflected(url, HeadersCompare{
		TestHeaderKey:   "X-Custom-Header",
		TestHeaderValue: "senpro\r\nSet-Cookie: " + testCookie,
	}) {
		output := fmt.Sprintf(" | + Cookie set: %s\n", testCookie)
		check.ConsoleOutput <- output
	}
}

// TODO: func (check *SubdomainCheck) RequestSmuggling(httpClient *http.Client)
