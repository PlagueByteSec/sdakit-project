package analysis

import (
	"fmt"
	"net/http"

	utils "github.com/PlagueByteSec/sentinel-project/v2/internal/coreutils"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/requests"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
	"github.com/PlagueByteSec/sentinel-project/v2/pkg"
)

const testDomain = "example.com"

type SubdomainCheck struct {
	Subdomain     string
	ConsoleOutput chan<- string
	HttpHeaders   string
	HttpClient    *http.Client
}

type HeadersCompare struct {
	TestHeaderKey       string
	TestHeaderValue     string
	ResponseHeaderKey   string
	ResponseHeaderValue []string
}

type RequestSetup struct {
	Method string
	URL    string
	Header string
	Value  string
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
	cmsIndicators = map[string][]string{
		"WordPress":        {"wp-content", "wp-includes"},
		"Joomla":           {"Joomla!"},
		"Drupal":           {"Drupal"},
		"Magento":          {"Magento"},
		"Shopify":          {"Shopify"},
		"Blogger":          {"blogspot"},
		"Wix":              {"wix"},
		"Squarespace":      {"squarespace"},
		"TYPO3":            {"typo3"},
		"Concrete5":        {"concrete5"},
		"PrestaShop":       {"prestashop"},
		"OpenCart":         {"catalog"},
		"Ghost":            {"ghost"},
		"ExpressionEngine": {"expressionEngine"},
		"Craft CMS":        {"craft"},
		"MODX":             {"MODX Revolution"},
		"SilverStripe":     {"silverstripe"},
		"DotNetNuke":       {"dnn"},
		"Weebly":           {"weebly"},
	}
)

type HTTP int

const (
	Basic HTTP = iota
	Secure
)

func cloudflareError(statusCode int, subdomain string) bool {
	if statusCode == 520 {
		utils.PrintVerbose(" | - %s responds with %d, (server error, cloudflare)\n", subdomain, statusCode)
		return true
	}
	return false
}

func MakeUrl(http HTTP, subdomain string) string {
	var proto string
	switch http {
	case Basic:
		proto = "http://"
	case Secure:
		proto = "https://"
	}
	return fmt.Sprintf("%s%s", proto, subdomain)
}

func (check *SubdomainCheck) AnalysisSendRequest(setup RequestSetup) *http.Response {
	var (
		request  *http.Request
		response *http.Response
		err      error
	)
	request, err = requests.RequestSetupHTTP(setup.Method, setup.URL, check.HttpClient)
	if err != nil {
		goto exitError
	}
	if len(setup.Header) != 0 && len(setup.Value) != 0 {
		request.Header.Set(setup.Header, setup.Value)
	}
	response, err = check.HttpClient.Do(request)
	if err != nil {
		goto exitError
	}
	if pkg.IsInSlice(string(response.StatusCode), errorCodes) {
		shared.Glogger.Println("Error: Server returned: " + string(response.StatusCode))
		return nil
	}
	return response
exitError:
	shared.Glogger.Println(err)
	return nil
}
