package analysis

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
	"github.com/PlagueByteSec/sentinel-project/v2/pkg"
)

func (check *SubdomainCheck) testHostHeader(header string) bool {
	url := makeUrl(HTTP(Secure), check.Subdomain)
	response := check.sendRequest(RequestSetup{Method: "GET", URL: url, Header: header, Value: testDomain})
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
			shared.Glogger.Printf("response contains %s key with value(s): %s\n",
				responseHeaderKey, acceptedResponseValues[idx])
			return true
		}
	}
	return false
}

func headerAccepted(compare HeadersCompare) bool {
	// Ensure the response header contains the test header/value.
	return strings.Contains(compare.ResponseHeaderKey, compare.TestHeaderKey) &&
		strings.Contains(strings.Join(compare.ResponseHeaderValue, ", "), compare.TestHeaderValue)
}

func (check *SubdomainCheck) investigateAcaoHeaders(response *http.Response) {
	var success bool
	for responseHeaderKey, responseHeaderValue := range response.Header {
		switch {
		case headerAccepted(HeadersCompare{"Access-Control-Allow-Origin", testDomain, responseHeaderKey, responseHeaderValue}):
			check.ConsoleOutput.WriteString(fmt.Sprintf(" | + [CORS:OK]: %s accepts %s as origin\n", check.Subdomain, testDomain))
			success = true
		case headerAccepted(HeadersCompare{"Access-Control-Allow-Origin", "null", responseHeaderKey, responseHeaderValue}):
			check.ConsoleOutput.WriteString(fmt.Sprintf(" | + [CORS:OK]: %s accepts null as origin\n", check.Subdomain))
			success = true
		case headerAccepted(HeadersCompare{"Access-Control-Allow-Credentials", "true", responseHeaderKey, responseHeaderValue}):
			check.ConsoleOutput.WriteString(fmt.Sprintf(" | + [CORS:OK]: %s allows creds in request\n", check.Subdomain))
			success = true
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
			shared.Glogger.Println(err)
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
			check.ConsoleOutput.WriteString(" | + [HT:OK] Server seems to accept header: " + headers[idx])
			check.ConsoleOutput.WriteString("\n")
		}
	}
}

// Ensure the injected cookie is reflected in the response from the current subdomain.
func (check *SubdomainCheck) isPayloadReflected(url string, compare HeadersCompare) bool {
	var isReflected bool
	response := check.sendRequest(RequestSetup{Method: "POST", URL: url, Header: "", Value: ""})
	if response == nil {
		return isReflected
	}
	defer response.Body.Close()
	for responseHeaderKey, responseHeaderValue := range response.Header {
		if headerAccepted(HeadersCompare{
			compare.TestHeaderKey, compare.TestHeaderValue,
			responseHeaderKey, responseHeaderValue,
		}) {
			isReflected = true
		}
	}
	return isReflected
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
			shared.Glogger.Println("response contains X-API-Version header")
			return 10, "X-API-Version"
		case strings.Contains(headerKey, "X-RateLimit-Limit"):
			shared.Glogger.Println("response contains X-RateLimit-Limit header")
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
	response := check.sendRequest(RequestSetup{Method: "GET", URL: url, Header: header, Value: testDomain})
	if response == nil {
		shared.Glogger.Println("testCors: response == nil")
		return
	}
	check.investigateAcaoHeaders(response)
}
