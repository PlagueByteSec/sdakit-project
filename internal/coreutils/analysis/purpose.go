package analysis

import (
	"fmt"

	"github.com/fhAnso/Sentinel/v1/internal/requests"
	"github.com/fhAnso/Sentinel/v1/internal/shared"
)

func (check *SubdomainCheck) MailServer() {
	if requests.DnsIsMX(shared.GDnsResolver, check.Subdomain) {
		check.ConsoleOutput.WriteString(" | + [MX:OK] Mail Server ")
		if check.isExchange() {
			check.ConsoleOutput.WriteString("(Exchange)\n")
		} else {
			check.ConsoleOutput.WriteString("\n")
		}
		shared.PoolAppendValue(check.Subdomain, &shared.GPoolBase.PoolMailSubdomains)
		return
	}
	check.ConsoleOutput.WriteString(" | - [MX:NA] No MX entry found\n")
}

func (check *SubdomainCheck) API() {
	url := fmt.Sprintf("http://%s", check.Subdomain)
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
			return
		}
	}
	check.ConsoleOutput.WriteString(" | - [API:NA] No API indicators found\n")
}

func (check *SubdomainCheck) Login() {
	var indicatorsFound bool
	url := fmt.Sprintf("http://%s", check.Subdomain)
	for idx := 0; idx < len(methods); idx++ {
		if indicatorsFound = check.isLoginPage(methods[idx], url); indicatorsFound {
			check.ConsoleOutput.WriteString(" | + [LOGIN:OK] Login page found\n")
			shared.PoolAppendValue(check.Subdomain, &shared.GPoolBase.PoolLoginSubdomains)
			return
		}
	}
	check.ConsoleOutput.WriteString(" | - [LOGIN:NA] No login found\n")
}