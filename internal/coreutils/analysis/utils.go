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

const testUrl = "http://example.com" // Origin

type SubdomainCheck struct {
	Subdomain     string
	ConsoleOutput *strings.Builder
	HttpHeaders   string
	HttpClient    *http.Client
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

func (check *SubdomainCheck) isExchange() bool {
	// Basic check for Microsoft Exchange server
	return strings.Contains(check.HttpHeaders, "X-Feproxyinfo") ||
		strings.Contains(check.Subdomain, "autodiscover")
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

func cloudflareError(statusCode int, subdomain string) bool {
	if statusCode == 520 {
		utils.PrintVerbose(" | - %s responds with %d, (server error, cloudflare)\n", subdomain, statusCode)
		return true
	}
	return false
}

func checkPageLogin(responseBody string) bool {
	for idx := 0; idx < len(loginIndicators); idx++ {
		if strings.Contains(responseBody, loginIndicators[idx]) {
			return true
		}
	}
	return false
}

func headerAccepted(testHeaderKey string, testHeaderValue string, responseHeaderKey string, responseHeaderValue []string) bool {
	// Ensure the response header contains the test header/value.
	return strings.Contains(responseHeaderKey, testHeaderKey) &&
		strings.Contains(strings.Join(responseHeaderValue, ", "), testHeaderValue)
}

func (check *SubdomainCheck) investigateHeaders(response *http.Response) bool {
	var success bool
	for responseHeaderKey, responseHeaderValue := range response.Header {
		switch {
		case headerAccepted("Access-Control-Allow-Origin", testUrl, responseHeaderKey, responseHeaderValue):
			check.ConsoleOutput.WriteString(fmt.Sprintf(" | + [CORS:OK]: %s accepts %s as origin\n", check.Subdomain, testUrl))
			success = true
		case headerAccepted("Access-Control-Allow-Origin", "null", responseHeaderKey, responseHeaderValue):
			check.ConsoleOutput.WriteString(fmt.Sprintf(" | + [CORS:OK]: %s accepts null as origin\n", check.Subdomain))
			success = true
		case headerAccepted("Access-Control-Allow-Credentials", "true", responseHeaderKey, responseHeaderValue):
			check.ConsoleOutput.WriteString(fmt.Sprintf(" | + [CORS:OK]: %s allows creds in request\n", check.Subdomain))
			success = true
		}
	}
	return success
}

func (check *SubdomainCheck) testCORS(method string, url string) bool {
	request, err := requests.RequestSetupHTTP(method, url, check.HttpClient)
	if err != nil {
		shared.Glogger.Println(err)
	}
	request.Header.Set("Origin", testUrl)
	response, err := check.HttpClient.Do(request)
	if err != nil {
		shared.Glogger.Println(err)
		return false
	}
	if pkg.IsInSlice(string(response.StatusCode), errorCodes) {
		shared.Glogger.Println("[CORS:ERR] Server responds with error code: " + string(response.StatusCode))
		return false
	}
	return check.investigateHeaders(response)
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
