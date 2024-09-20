package analysis

import (
	"net/http"
	"strconv"

	utils "github.com/PlagueByteSec/sdakit-project/v2/internal/coreutils"
	pools "github.com/PlagueByteSec/sdakit-project/v2/internal/datapools"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/logging"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/requests"
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

type AnalysisRequestConfig struct {
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

func cloudflareError(statusCode int, subdomain string) bool {
	if statusCode == 520 {
		utils.PrintVerbose(" | - %s responds with %d, (server error, cloudflare)\n", subdomain, statusCode)
		return true
	}
	return false
}

func (check *SubdomainCheck) AnalysisSendRequest(setup AnalysisRequestConfig) *http.Response {
	httpResponse, _, _, err := requests.RequestHandlerCore(&requests.HttpRequestBase{
		HttpClient:       check.HttpClient,
		CustomUrl:        setup.URL,
		HttpMethod:       setup.Method,
		HttpNeedResponse: true,
	})
	if err != nil {
		logging.GLogger.Log(err.Error())
		return nil
	}
	codeToString := strconv.Itoa(httpResponse.StatusCode)
	if pools.ManagePool(pools.PoolAction(pools.PoolCheck), codeToString, &errorCodes) {
		logging.GLogger.Log("Error: Server returned: " + codeToString)
		return nil
	}
	return httpResponse
}

func (check *SubdomainCheck) TestSecurity() {
	check.hostHeaders()     // Host header injections
	check.cookieInjection() // session hijacking, xss
	check.requestSmuggling()
	check.CORS()
}

func (check *SubdomainCheck) TargetAnalyseHTTP() {
	check.api()   // Content types, API versions, rate limit
	check.login() // Scan response body for login indicators
	check.cms()   // Top 20 CMS
}

func (check *SubdomainCheck) TargetAnalyseNonHTTP() {
	check.MailServer()
}
