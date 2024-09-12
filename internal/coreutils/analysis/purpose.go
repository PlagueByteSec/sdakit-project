package analysis

import (
	"fmt"
	"strings"

	"github.com/PlagueByteSec/sentinel-project/v2/internal/requests"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
)

func (check *SubdomainCheck) MailServer() {
	if requests.DnsIsMX(shared.GDnsResolver, check.Subdomain) {
		check.ConsoleOutput <- " | + Mail Server "
		if check.isExchange() {
			check.ConsoleOutput <- "(Exchange)\n"
		} else {
			check.ConsoleOutput <- "\n"
		}
		shared.PoolAppendValue(check.Subdomain, &shared.GPoolBase.PoolMailSubdomains)
	}
}

func (check *SubdomainCheck) api() {
	url := MakeUrl(HTTP(Basic), check.Subdomain)
	for idx := 0; idx < len(methods); idx++ {
		response := check.AnalysisSendRequest(RequestSetup{Method: methods[idx], URL: url, Header: "", Value: ""})
		if response == nil {
			continue
		}
		statusCode := response.StatusCode
		if cloudflareError(statusCode, check.Subdomain) {
			continue
		}
		score, info := check.isPossibleApi(response)
		if score != 0 {
			shared.PoolAppendValue(check.Subdomain, &shared.GPoolBase.PoolApiSubdomains)
			check.ConsoleOutput <- fmt.Sprintf(" | + API [SCORE:%d] (%s: %s)\n", score, methods[idx], info)
			break
		}
	}
}

func (check *SubdomainCheck) login() {
	check.checkPage("login", check.isLoginPage, " | + Login\n")
}

func (check *SubdomainCheck) cms() {
	url := MakeUrl(HTTP(Basic), check.Subdomain)
	response := check.getResponse(url)
	if response == nil {
		return
	}
	body := check.responseGetBody(response)
	if len(body) == 0 {
		return
	}
	html := string(body)
	for cmsName, indicators := range cmsIndicators {
		for idx := 0; idx < len(indicators); idx++ {
			if strings.Contains(html, indicators[idx]) {
				shared.PoolAppendValue(check.Subdomain, &shared.GPoolBase.PoolCmsSubdomains)
				check.ConsoleOutput <- fmt.Sprintf(" | + CMS: %s\n", cmsName)
				break
			}
		}
	}
}
