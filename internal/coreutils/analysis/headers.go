package analysis

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PlagueByteSec/sentinel-project/v2/internal/logging"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
	"github.com/PlagueByteSec/sentinel-project/v2/pkg"
	"github.com/fhAnso/astkit"
)

func (check *SubdomainCheck) testHostHeader(header string) bool {
	url := astkit.MakeUrl(astkit.HTTP(astkit.Secure), check.Subdomain)
	response := check.AnalysisSendRequest(AnalysisRequestConfig{Method: "GET", URL: url, Header: header, Value: testDomain})
	// ensure the response include test domain
	return pkg.Tern(response != nil, check.investigateHostHeaders(header, response), false)
}

func (check *SubdomainCheck) checkFormat(responseHeaderKey string, responseHeaderValues string) bool {
	acceptedResponseValues := []string{
		"application/json",
		"application/vnd.api+json",
		"application/xml",
		"text/xml",
	}
	for idx := 0; idx < len(acceptedResponseValues); idx++ {
		if strings.Contains(responseHeaderValues, acceptedResponseValues[idx]) {
			output := fmt.Sprintf("response contains %s key with value(s): %s\n",
				responseHeaderKey, acceptedResponseValues[idx])
			logging.GLogger.Log(output)
			return true
		}
	}
	return false
}

func headerAccepted(compare HeadersCompare) bool {
	if compare.TestHeaderKey != compare.ResponseHeaderKey {
		return false
	}
	for i := 0; i < len(compare.ResponseHeaderValue); i++ {
		value := compare.ResponseHeaderValue[i]
		if strings.Contains(value, compare.TestHeaderValue) {
			return true
		}
	}
	return false
}

func (check *SubdomainCheck) investigateAcaoHeaders(response *http.Response) {
	var (
		success bool
		result  string
	)
	processSuccess := func(message string) {
		shared.GReportPool["CORS"] = shared.SetTestResults{
			TestName:   "CORS",
			TestResult: "FOUND",
			Subdomain:  check.Subdomain,
		}
		check.ConsoleOutput <- fmt.Sprintf(" | + CORS: %s\n", message)
		success = true
	}
	for responseHeaderKey, responseHeaderValue := range response.Header {
		switch {
		case headerAccepted(HeadersCompare{"Access-Control-Allow-Origin", testDomain, responseHeaderKey, responseHeaderValue}):
			result = fmt.Sprintf("%s accepts %s as origin\n", check.Subdomain, testDomain)
			processSuccess(result)
		case headerAccepted(HeadersCompare{"Access-Control-Allow-Origin", "null", responseHeaderKey, responseHeaderValue}):
			result = fmt.Sprintf("%s accepts null as origin\n", check.Subdomain)
			processSuccess(result)
		case headerAccepted(HeadersCompare{"Access-Control-Allow-Origin", "*", responseHeaderKey, responseHeaderValue}):
			result = fmt.Sprintf("%s allows all origins\n", check.Subdomain)
			processSuccess(result)
		case headerAccepted(HeadersCompare{"Access-Control-Allow-Credentials", "true", responseHeaderKey, responseHeaderValue}):
			if headerAccepted(HeadersCompare{"Access-Control-Allow-Origin", "*", responseHeaderKey, responseHeaderValue}) {
				result = fmt.Sprintf("%s allows credentials with wildcard origin\n", check.Subdomain)
				processSuccess(result)
			} else {
				result = fmt.Sprintf("%s allows credentials in request\n", check.Subdomain)
				processSuccess(result)
			}
		default:
			shared.GReportPool["CORS"] = shared.SetTestResults{
				TestName:   "CORS",
				TestResult: "PASSED",
				Subdomain:  check.Subdomain,
			}
		}
	}
	if success {
		shared.PoolAppendValue(check.Subdomain, &shared.GPoolBase.PoolCorsSubdomains)
		fmt.Println(check.ConsoleOutput)
	}
}

func (check *SubdomainCheck) investigateHostHeaders(header string, response *http.Response) bool {
	defer response.Body.Close()
	for responseHeaderKey, responseHeaderValue := range response.Header {
		// check if test domain in response headers
		if headerAccepted(HeadersCompare{header, testDomain, responseHeaderKey, responseHeaderValue}) {
			return true
		}
		// check if test domain in response body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			logging.GLogger.Log(err.Error())
			return false
		}
		if strings.Contains(string(body), testDomain) {
			return true
		}
	}
	return false
}

func (check *SubdomainCheck) hostHeaders() { // allow redirect = true
	headers := []string{"Host", "X-Forwarded-Host", "X-Host"}
	for idx := 0; idx < len(headers); idx++ {
		if check.testHostHeader(headers[idx]) {
			check.ConsoleOutput <- " | + [HT:OK] Server seems to accept header: " + headers[idx]
			check.ConsoleOutput <- "\n"
		}
	}
}

func (check *SubdomainCheck) isExchange() bool {
	// Basic check for Microsoft Exchange server
	return strings.Contains(check.HttpHeaders, "X-Feproxyinfo") ||
		strings.Contains(check.Subdomain, "autodiscover")
}

func (check *SubdomainCheck) isPossibleApi(httpResponse *http.Response) (int, string) {
	for headerKey, headerValues := range httpResponse.Header {
		values := strings.Join(headerValues, ", ")
		switch {
		/*
			Analyse API typical response headers and inspect the values
			to determine API possibility.
		*/
		case strings.Contains(headerKey, "X-API-Version"):
			return 10, "X-API-Version"
		case strings.Contains(headerKey, "X-RateLimit-Limit"):
			return 10, "X-RateLimit-Limit"
		case strings.Contains(headerKey, "Content-Type") && check.checkFormat(values, "Content-Type"):
			return 5, "Content-Type"
		case strings.Contains(headerKey, "Accept") && check.checkFormat(values, "Accept"):
			return 5, "Accept"
		case strings.Contains(headerKey, "Link") && strings.Contains(values, "api"):
			return 5, fmt.Sprintf("response contains interesting endpoint: %s\n", values)
		case strings.Contains(check.Subdomain, "api"):
			return 5, "by-name"
		}
	}
	return 0, ""
}

func (check SubdomainCheck) testCors(url string, header string) {
	response := check.AnalysisSendRequest(AnalysisRequestConfig{Method: "GET", URL: url, Header: header, Value: testDomain})
	if response == nil {
		logging.GLogger.Log("testCors: AnalysisSendRequest returns nil")
		return
	}
	check.investigateAcaoHeaders(response)
}
