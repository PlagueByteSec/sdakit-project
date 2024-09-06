package analysis

import (
	"fmt"

	"github.com/fhAnso/Sentinel/v1/internal/shared"
)

func (check *SubdomainCheck) CORS() {
	methodsCORS := methods[:len(methods)-1] // Remove OPTIONS
	url := makeUrl(check.Subdomain)
	for idx := 0; idx < len(methodsCORS); idx++ {
		// Test with GET, POST
		if testSuccess := check.testCors(methodsCORS[idx], url, "Origin", testDomain); testSuccess {
			shared.PoolAppendValue(check.Subdomain, &shared.GPoolBase.PoolCorsSubdomains)
			return
		}
	}
}

func (check *SubdomainCheck) cookieInjectionPath() {
	// session hijacking, xss
	teader := "Set-Cookie"
	tookie := "jzqvtyxkplra"
	url := makeUrl(check.Subdomain) + "%0d%0a" + fmt.Sprintf("%s:+tookie=%s", teader, tookie)
	if check.isPayloadReflected(url, SetupCompare{TestHeaderKey: teader, TestHeaderValue: tookie}) {
		check.ConsoleOutput.WriteString(fmt.Sprintf(" | + [CI:OK] Payload reflected in response: %s: %s\n",
			teader, tookie))
	}
}

// TODO: func (check *SubdomainCheck) RequestSmuggling(httpClient *http.Client)
