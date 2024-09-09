package analysis

import (
	"fmt"

	"github.com/fhAnso/Sentinel/v1/internal/requests"
	"github.com/fhAnso/Sentinel/v1/internal/shared"
)

func (check *SubdomainCheck) mailServer() {
	if requests.DnsIsMX(shared.GDnsResolver, check.Subdomain) {
		check.ConsoleOutput.WriteString(" | + Mail Server ")
		if check.isExchange() {
			check.ConsoleOutput.WriteString("(Exchange)\n")
		} else {
			check.ConsoleOutput.WriteString("\n")
		}
		shared.PoolAppendValue(check.Subdomain, &shared.GPoolBase.PoolMailSubdomains)
	}
}

func (check *SubdomainCheck) api() {
	url := makeUrl(HTTP(Basic), check.Subdomain)
	for idx := 0; idx < len(methods); idx++ {
		response := check.sendRequest(RequestSetup{Method: methods[idx], URL: url, Header: "", Value: ""})
		if response == nil {
			continue
		}
		statusCode := response.StatusCode
		if cloudflareError(statusCode, check.Subdomain) {
			continue
		}
		if apiPossibility, count, info := check.isPossibleApi(response); apiPossibility {
			if count == 10 {
				check.ConsoleOutput.WriteString(" | + ")
			} else if count < 10 {
				check.ConsoleOutput.WriteString(" | ? ")
			}
			check.ConsoleOutput.WriteString(fmt.Sprintf("API (%s: %s)\n", methods[idx], info))
			shared.PoolAppendValue(check.Subdomain, &shared.GPoolBase.PoolApiSubdomains)
		}
	}
}

func (check *SubdomainCheck) login() {
	check.checkPage("login", check.isLoginPage, " | + Login\n")
}

func (check *SubdomainCheck) basicWebpage() {
	check.checkPage("basic", check.isBasicWebpage, " | + Basic web site\n")
}
