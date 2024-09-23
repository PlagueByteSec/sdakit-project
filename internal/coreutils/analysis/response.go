package analysis

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	pools "github.com/PlagueByteSec/sdakit-project/v2/internal/datapools"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/logging"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/shared"
	"github.com/PlagueByteSec/sdakit-project/v2/pkg"
	"github.com/fhAnso/astkit"
)

func findIndicator(body string, indicators []string) string {
	for _, indicator := range indicators {
		if result := pkg.Tern(strings.Contains(body, indicator), indicator, ""); result != "" {
			return result
		}
	}
	return ""
}

func (check *SubdomainCheck) getResponse(url string) *http.Response {
	response := check.AnalysisSendRequest(AnalysisRequestConfig{Method: "GET", URL: url, Header: "", Value: ""})
	return pkg.Tern(response == nil, nil, response)
}

func (check *SubdomainCheck) responseGetBody(response *http.Response) []byte {
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		logging.GLogger.Log(err.Error())
		return nil
	}
	return responseBody
}

func (check *SubdomainCheck) checkPage(pageType string, runDetection func(body string) string, body string) {
	if result := runDetection(body); result != "" {
		switch pageType {
		case "login":
			pools.ManagePool(pools.PoolAction(pools.PoolAppend), check.Subdomain, &shared.GPoolBase.PoolLoginSubdomains)
			check.ConsoleOutput <- fmt.Sprintf(" | + login: %s\n", result)
		case "cms":
			output := fmt.Sprintf("%s (%s)", check.Subdomain, result)
			pools.ManagePool(pools.PoolAction(pools.PoolAppend), output, &shared.GPoolBase.PoolCmsSubdomains)
			check.ConsoleOutput <- fmt.Sprintf(" | + CMS: %s\n", result)
		}
	}
}

func detectLogin(body string) string {
	return findIndicator(body, loginIndicators)
}

func detectCMS(body string) string {
	for cmsName, indicators := range astkit.AcceptedCMS {
		if result := findIndicator(body, indicators); result != "" {
			return cmsName
		}
	}
	return ""
}
