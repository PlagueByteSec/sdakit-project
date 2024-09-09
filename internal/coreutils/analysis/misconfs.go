package analysis

import (
	"fmt"
)

func (check *SubdomainCheck) CORS() {
	url := makeUrl(HTTP(Secure), check.Subdomain)
	check.testCors(url, "Origin") // GET
}

func (check *SubdomainCheck) cookieInjectionPath() {
	// session hijacking, xss
	testHeader := "Set-Cookie"
	testCookie := "jzqvtyxkplra"
	url := makeUrl(HTTP(Secure), check.Subdomain) +
		"%0d%0a" + fmt.Sprintf("%s:+tookie=%s", testHeader, testCookie)
	if check.isPayloadReflected(url, HeadersCompare{
		TestHeaderKey:   testHeader,
		TestHeaderValue: testCookie,
	}) {
		output := fmt.Sprintf(" | + Payload reflected in response: %s: %s\n", testHeader, testCookie)
		check.ConsoleOutput.WriteString(output)
	}
}

// TODO: func (check *SubdomainCheck) RequestSmuggling(httpClient *http.Client)
