package analysis

import (
	"github.com/PlagueByteSec/sentinel-project/v2/internal/logging"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/requests"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
	"github.com/fhAnso/astkit"
)

func (check *SubdomainCheck) CORS() {
	url := astkit.MakeUrl(astkit.HTTP(astkit.Secure), check.Subdomain)
	check.testCors(url, "Origin") // GET
}

func (check *SubdomainCheck) cookieInjection() {
	client := astkit.ASTkitClient{
		HttpClient: check.HttpClient,
	}
	_, openPorts, _ := requests.ScanPortRange(check.Subdomain, "80,8080,443,8443", true)
	if len(openPorts) == 0 {
		return
	}
	for idx := 0; idx < len(openPorts); idx++ {
		result, err := astkit.InjectCookie(astkit.HeaderInjectionConfig{
			Client:    &client,
			Host:      check.Subdomain,
			Port:      openPorts[idx],
			UserAgent: shared.DefaultUserAgent,
		})
		if err != nil {
			logging.GLogger.Log(err.Error())
			continue
		}
		if len(result) == 0 {
			continue
		}
		check.ConsoleOutput <- result
	}
}

// TODO: func (check *SubdomainCheck) RequestSmuggling(httpClient *http.Client)
