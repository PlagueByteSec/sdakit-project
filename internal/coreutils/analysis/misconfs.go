package analysis

import (
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

func (check *SubdomainCheck) HeaderInjection() {
	check.hostHeaders()
	check.cookieInjectionPath()
}

// TODO: func (check *SubdomainCheck) RequestSmuggling(httpClient *http.Client)
