package analysis

import (
	"fmt"

	"github.com/fhAnso/Sentinel/v1/internal/requests"
	"github.com/fhAnso/Sentinel/v1/internal/shared"
)

func (check *SubdomainCheck) mailServer() {
	if requests.DnsIsMX(shared.GDnsResolver, check.Subdomain) {
		check.ConsoleOutput.WriteString(" | + [MX:OK] Mail Server ")
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
		response := check.sendRequest(methods[idx], url)
		if response == nil {
			continue
		}
		statusCode := response.StatusCode
		if cloudflareError(statusCode, check.Subdomain) {
			continue
		}
		if apiPossibility, count, info := check.isPossibleApi(response); apiPossibility {
			check.ConsoleOutput.WriteString(" | + [API:OK] " + check.Subdomain)
			if count == 10 {
				check.ConsoleOutput.WriteString(fmt.Sprintf(": API detected (%s: %s)\n", methods[idx], info))
			} else if count < 10 {
				check.ConsoleOutput.WriteString(fmt.Sprintf(" seems to be a API.. (%s: %s)\n", methods[idx], info))
			}
			shared.PoolAppendValue(check.Subdomain, &shared.GPoolBase.PoolApiSubdomains)
		}
	}
}

func (check *SubdomainCheck) login() {
	check.checkPage("login", check.isLoginPage, " | + [LOGIN:OK] Login page found\n")
}

func (check *SubdomainCheck) basicWebpage() {
	check.checkPage("basic", check.isBasicWebpage, " | + [BWA:OK] Basic web app\n")
}
