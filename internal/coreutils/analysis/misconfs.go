package analysis

import (
	"fmt"

	"github.com/fhAnso/Sentinel/v1/internal/shared"
)

func (check *SubdomainCheck) CORS() {
	methodsCORS := methods[:len(methods)-1] // Remove OPTIONS
	url := fmt.Sprintf("http://%s", check.Subdomain)
	var testSuccess bool
	for idx := 0; idx < len(methodsCORS); idx++ {
		// Test with GET, POST
		if testSuccess = check.testCORS(methodsCORS[idx], url); testSuccess {
			shared.PoolAppendValue(check.Subdomain, &shared.GPoolBase.PoolCorsSubdomains)
			return
		}
	}
	if !testSuccess {
		check.ConsoleOutput.WriteString(" | - [CORS:NA] No CORS misconfigurations identified\n")
	}
}

// TODO: func (check *SubdomainCheck) HeaderInjection(httpClient *http.Client)
// TODO: func (check *SubdomainCheck) RequestSmuggling(httpClient *http.Client)
