package analysis

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	utils "github.com/fhAnso/Sentinel/v1/internal/coreutils"
	"github.com/fhAnso/Sentinel/v1/internal/requests"
	"github.com/fhAnso/Sentinel/v1/internal/shared"
	"github.com/fhAnso/Sentinel/v1/pkg"
)

const testDomain = "example.com"

type SubdomainCheck struct {
	Subdomain     string
	ConsoleOutput *strings.Builder
	HttpHeaders   string
	HttpClient    *http.Client
}

type SetupCompare struct {
	TestHeaderKey       string
	TestHeaderValue     string
	ResponseHeaderKey   string
	ResponseHeaderValue []string
}

var (
	methods    = []string{"GET", "POST", "OPTIONS"}
	errorCodes = []string{
		"500", // Internal Server Error
		"501", // Not Implemented
		"502", // Bad Gateway
		"503", // Service Unavailable
		"504", // Gateway Timeout
		"505", // HTTP Version Not Supported
	}
	loginIndicators = []string{
		"Login",
		"username",
		"password",
		"Log In",
		"Log On",
		"Authenticate",
		"Forgot Password",
		"Reset Password",
		"Account Login",
		"User ID",
		"Email",
		"Please log in",
		"Two-factor authentication",
		"Continue to login",
	}
)

func makeUrl(subdomain string) string {
	return fmt.Sprintf("https://%s", subdomain)
}

func (check *SubdomainCheck) sendRequest(method string, url string) *http.Response {
	request, err := requests.RequestSetupHTTP(method, url, check.HttpClient)
	if err != nil {
		shared.Glogger.Println(err)
		return nil
	}
	response, err := check.HttpClient.Do(request)
	if err != nil {
		shared.Glogger.Println(err)
		return nil
	}
	return response
}

func checkPageLogin(responseBody string) bool {
	for idx := 0; idx < len(loginIndicators); idx++ {
		if strings.Contains(responseBody, loginIndicators[idx]) {
			return true
		}
	}
	return false
}

func cloudflareError(statusCode int, subdomain string) bool {
	if statusCode == 520 {
		utils.PrintVerbose(" | - %s responds with %d, (server error, cloudflare)\n", subdomain, statusCode)
		return true
	}
	return false
}

func (check *SubdomainCheck) checkResponseValues(responseHeaderKey string, responseHeaderValues string) bool {
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

func headerAccepted(compare SetupCompare) bool {
	// Ensure the response header contains the test header/value.
	return strings.Contains(compare.ResponseHeaderKey, compare.TestHeaderKey) &&
		strings.Contains(strings.Join(compare.ResponseHeaderValue, ", "), compare.TestHeaderValue)
}

func (check *SubdomainCheck) investigateAcaoHeaders(response *http.Response) bool {
	var success bool
	for responseHeaderKey, responseHeaderValue := range response.Header {
		switch {
		case headerAccepted(SetupCompare{"Access-Control-Allow-Origin", testDomain, responseHeaderKey, responseHeaderValue}):
			check.ConsoleOutput.WriteString(fmt.Sprintf(" | + [CORS:OK]: %s accepts %s as origin\n", check.Subdomain, testDomain))
			success = true
		case headerAccepted(SetupCompare{"Access-Control-Allow-Origin", "null", responseHeaderKey, responseHeaderValue}):
			check.ConsoleOutput.WriteString(fmt.Sprintf(" | + [CORS:OK]: %s accepts null as origin\n", check.Subdomain))
			success = true
		case headerAccepted(SetupCompare{"Access-Control-Allow-Credentials", "true", responseHeaderKey, responseHeaderValue}):
			check.ConsoleOutput.WriteString(fmt.Sprintf(" | + [CORS:OK]: %s allows creds in request\n", check.Subdomain))
			success = true
		}
	}
	return success
}

func (check *SubdomainCheck) investigateHostHeaders(header string, response *http.Response) bool {
	defer response.Body.Close()
	for responseHeaderKey, responseHeaderValue := range response.Header {
		// check if test domain in response headers
		if headerAccepted(SetupCompare{header, testDomain, responseHeaderKey, responseHeaderValue}) {
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

func (check *SubdomainCheck) testCustomHeader(method string, url string, header string, value string) *http.Response {
	request, err := requests.RequestSetupHTTP(method, url, check.HttpClient)
	if err != nil {
		shared.Glogger.Println(err)
	}
	request.Header.Set(header, value)
	response, err := check.HttpClient.Do(request)
	if err != nil {
		shared.Glogger.Println(err)
		return nil
	}
	if pkg.IsInSlice(string(response.StatusCode), errorCodes) {
		shared.Glogger.Println("[TCH:ERR] Server responds with error code: " + string(response.StatusCode))
		return nil
	}
	return response
}

func (check SubdomainCheck) testCors(method string, url string, header string, value string) bool {
	response := check.testCustomHeader(method, url, header, value)
	if response == nil {
		return false
	}
	return check.investigateAcaoHeaders(response)
}

func (check *SubdomainCheck) testHostHeader(header string) bool {
	url := makeUrl(check.Subdomain)
	response := check.testCustomHeader("GET", url, header, testDomain)
	if response != nil {
		// ensure the response include test domain
		return check.investigateHostHeaders(header, response)
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

func (check *SubdomainCheck) isPayloadReflected(url string, compare SetupCompare) bool {
	var isReflected bool
	response := check.sendRequest("POST", url)
	if response == nil {
		return isReflected
	}
	defer response.Body.Close()
	for responseHeaderKey, responseHeaderValue := range response.Header {
		if headerAccepted(SetupCompare{
			compare.TestHeaderKey, compare.TestHeaderValue,
			responseHeaderKey, responseHeaderValue,
		}) {
			isReflected = true
		}
	}
	return isReflected
}

func (check *SubdomainCheck) cookieInjectionPath() {
	// session hijacking, xss
	teader := "Set-Cookie"
	tookie := "jzqvtyxkplra"
	url := makeUrl(check.Subdomain) + "/%0d%0a" + fmt.Sprintf("%s:+tookie=%s", teader, tookie)
	if check.isPayloadReflected(url, SetupCompare{TestHeaderKey: teader, TestHeaderValue: tookie}) {
		check.ConsoleOutput.WriteString(fmt.Sprintf(" | + [CI:OK] Payload reflected in response: %s: %s\n",
			teader, tookie))
	}
}

func (check *SubdomainCheck) isExchange() bool {
	// Basic check for Microsoft Exchange server
	return strings.Contains(check.HttpHeaders, "X-Feproxyinfo") ||
		strings.Contains(check.Subdomain, "autodiscover")
}

func (check *SubdomainCheck) isPossibleApi(httpResponse *http.Response) (bool, int, string) {
	var (
		apiPossibility      bool
		apiPossibilityCount int
	)
	for headerKey, headerValues := range httpResponse.Header {
		values := strings.Join(headerValues, ", ") // Convert response header values to string
		switch {
		/*
			Analyse API typical response headers and inspect the values
			to determine API possibility.
		*/
		case strings.Contains(headerKey, "Content-Type"):
			apiPossibility = check.checkResponseValues(values, "Content-Type") || apiPossibility
			apiPossibilityCount++
		case strings.Contains(headerKey, "Accept"):
			apiPossibility = check.checkResponseValues(values, "Accept") || apiPossibility
			apiPossibilityCount++
		case strings.Contains(headerKey, "Link") && strings.Contains(values, "api"):
			return true, 5, fmt.Sprintf("response contains interesting endpoint: %s\n", values)
		case strings.Contains(headerKey, "X-API-Version"):
			shared.Glogger.Println("response contains X-API-Version header")
			return true, 10, "X-API-Version"
		case strings.Contains(headerKey, "X-RateLimit-Limit"):
			shared.Glogger.Println("response contains X-RateLimit-Limit header")
			return true, 10, "X-RateLimit-Limit"
		}
	}
	return apiPossibility, apiPossibilityCount, ""
}

func (check *SubdomainCheck) isLoginPage(method string, url string) bool {
	response := check.sendRequest(method, url)
	if response == nil {
		return false
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		shared.Glogger.Println(err)
		return false
	}
	return checkPageLogin(string(responseBody))
}
