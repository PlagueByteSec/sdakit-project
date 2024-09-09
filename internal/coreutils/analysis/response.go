package analysis

import (
	"io"
	"net/http"
	"strings"

	"github.com/PlagueByteSec/Sentinel/v1/internal/shared"
	"github.com/PlagueByteSec/Sentinel/v1/pkg"
)

func (check *SubdomainCheck) getResponse(url string) *http.Response {
	response := check.sendRequest(RequestSetup{Method: "GET", URL: url, Header: "", Value: ""})
	return pkg.Tern(response == nil, nil, response)
}

func (check *SubdomainCheck) responseGetBody(response *http.Response) []byte {
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		shared.Glogger.Println(err)
		return nil
	}
	return responseBody
}

func (check *SubdomainCheck) checkPage(pageType string, pageInvestigate func(string, *http.Response) bool, successMessage string) {
	url := makeUrl(HTTP(Basic), check.Subdomain)
	response := check.getResponse(url)
	if response == nil {
		return
	}
	if ok := pageInvestigate(url, response); ok {
		check.ConsoleOutput.WriteString(successMessage)
		if pageType == "login" {
			shared.PoolAppendValue(check.Subdomain, &shared.GPoolBase.PoolLoginSubdomains)
		}
	}
}

func checkPageLogin(responseBody string) bool {
	if len(responseBody) != 0 {
		for idx := 0; idx < len(loginIndicators); idx++ {
			if strings.Contains(responseBody, loginIndicators[idx]) {
				return true
			}
		}
	}
	return false
}

func (check *SubdomainCheck) isLoginPage(url string, response *http.Response) bool {
	return checkPageLogin(string(check.responseGetBody(response)))
}

func (check *SubdomainCheck) isBasicWebpage(url string, response *http.Response) bool {
	return pkg.Tern(response != nil, checkGeneralWebpage(response), false)
}
