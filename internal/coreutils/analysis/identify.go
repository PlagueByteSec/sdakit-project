package analysis

import (
	"fmt"
	"strings"

	pools "github.com/PlagueByteSec/sdakit-project/v2/internal/datapools"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/logging"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/requests"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/shared"
	"github.com/fhAnso/astkit"
)

func (check *SubdomainCheck) MailServer() {
	if requests.DnsIsMX(shared.GDnsResolver, check.Subdomain) {
		check.ConsoleOutput <- " | + Mail Server "
		if check.isExchange() {
			check.ConsoleOutput <- "(Exchange)\n"
		} else {
			check.ConsoleOutput <- "\n"
		}
		pools.ManagePool(pools.PoolAction(pools.PoolAppend), check.Subdomain, &shared.GPoolBase.PoolMailSubdomains)
	}
}

func (check *SubdomainCheck) api() {
	url := astkit.MakeUrl(astkit.HTTP(astkit.Basic), check.Subdomain)
	for idx := 0; idx < len(methods); idx++ {
		response := check.AnalysisSendRequest(AnalysisRequestConfig{Method: methods[idx], URL: url, Header: "", Value: ""})
		if response == nil {
			continue
		}
		statusCode := response.StatusCode
		if cloudflareError(statusCode, check.Subdomain) {
			continue
		}
		score, info := check.isPossibleApi(response)
		if score != 0 {
			pools.ManagePool(pools.PoolAction(pools.PoolAppend), check.Subdomain, &shared.GPoolBase.PoolApiSubdomains)
			check.ConsoleOutput <- fmt.Sprintf(" | + API [SCORE:%d] (%s: %s)\n", score, methods[idx], info)
			break
		}
	}
}

func (check *SubdomainCheck) login() {
	check.checkPage("login", check.isLoginPage, " | + Login\n")
}

func (check *SubdomainCheck) cms() {
	url := astkit.MakeUrl(astkit.HTTP(astkit.Secure), check.Subdomain)
	client := astkit.ASTkitClient{
		HttpClient: check.HttpClient,
		URL:        url,
	}
	cms, err := astkit.ASTkitDetectCMS(&client)
	if err != nil {
		logging.GLogger.Log(err.Error())
		return
	}
	if !strings.Contains(cms, "unknown") {
		check.ConsoleOutput <- fmt.Sprintf(" | + CMS: %s\n", cms)
	}
}
