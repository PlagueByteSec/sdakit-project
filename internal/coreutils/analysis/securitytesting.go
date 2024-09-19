package analysis

import (
	"errors"

	"github.com/PlagueByteSec/sdakit-project/v2/internal/logging"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/requests"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/shared"
	"github.com/fhAnso/astkit"
)

func (check *SubdomainCheck) CORS() {
	url := astkit.MakeUrl(astkit.HTTP(astkit.Secure), check.Subdomain)
	check.testCors(url, "Origin") // GET
}

func testInit(check *SubdomainCheck) (*astkit.ASTkitClient, []uint16, error) {
	client := astkit.ASTkitClient{
		HttpClient: check.HttpClient,
	}
	_, openPorts, _ := requests.ScanPortRange(check.Subdomain, "80,8080,443,8443", true)
	if len(openPorts) == 0 {
		return nil, nil, errors.New("no open ports available")
	}
	return &client, openPorts, nil
}

func (check *SubdomainCheck) cookieInjection() {
	client, openPorts, err := testInit(check)
	if err != nil {
		logging.GLogger.Log(err.Error())
		return
	}
	for idx := 0; idx < len(openPorts); idx++ {
		result, err := astkit.InjectCookie(astkit.HeaderTestingConfig{
			Client:    client,
			Host:      check.Subdomain,
			Port:      openPorts[idx],
			UserAgent: shared.DefaultUserAgent,
			Test:      astkit.TestType(astkit.CookieInjection),
		})
		if err != nil {
			logging.GLogger.Log(err.Error())
			continue
		}
		if len(result) == 0 {
			continue
		}
		check.ConsoleOutput <- result
		shared.PoolAppendValue(check.Subdomain, &shared.GPoolBase.PoolCookieInjection)
		shared.GReportPool["CI"] = shared.SetTestResults{
			TestName:   "Cookie injection",
			TestResult: "FOUND",
			Subdomain:  check.Subdomain,
		}
		return
	}
	shared.GReportPool["CI"] = shared.SetTestResults{
		TestName:   "Cookie injection",
		TestResult: "PASSED",
		Subdomain:  check.Subdomain,
	}
}

func (check *SubdomainCheck) requestSmuggling() {
	client, openPorts, err := testInit(check)
	if err != nil {
		logging.GLogger.Log(err.Error())
		return
	}
	var types = []astkit.TestType{
		astkit.CLTE,
		astkit.TECL,
	}
	for idx := 0; idx < len(openPorts); idx++ {
		for tidx := 0; tidx < len(types); tidx++ {
			result, err := astkit.RequestSmuggling(astkit.HeaderTestingConfig{
				Client:    client,
				Host:      check.Subdomain,
				Port:      openPorts[idx],
				UserAgent: shared.DefaultUserAgent,
				Test:      astkit.TestType(types[tidx]),
			})
			if err != nil {
				logging.GLogger.Log(err.Error())
				continue
			}
			if len(result) == 0 {
				continue
			}
			check.ConsoleOutput <- result
			shared.PoolAppendValue(check.Subdomain, &shared.GPoolBase.PoolRequestSmuggling)
			shared.GReportPool["RS"] = shared.SetTestResults{
				TestName:   "Request smuggling",
				TestResult: "FOUND",
				Subdomain:  check.Subdomain,
			}
			return
		}
	}
	shared.GReportPool["RS"] = shared.SetTestResults{
		TestName:   "Request smuggling",
		TestResult: "PASSED",
		Subdomain:  check.Subdomain,
	}
}
