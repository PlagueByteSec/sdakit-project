package requests

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PlagueByteSec/sdakit-project/v2/internal/logging"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/shared"
)

func HttpCodeCheck(settings shared.SettingsHandler, url string) bool {
	_, statusCode, _, err := RequestHandlerCore(&HttpRequestBase{
		HttpClient:             settings.HttpClient,
		CustomUrl:              url,
		HttpMethod:             settings.Args.HttpRequestMethod,
		ResponseNeedStatusCode: true,
	})
	if err != nil {
		logging.GLogger.Log(err.Error())
		return false
	}
	return statusCode != -1
}

func HttpHeaderInit(httpHeaders *shared.HttpHeaders) {
	httpHeaders.Server = "Server"
	httpHeaders.Hsts = "Strict-Transport-Security"
	httpHeaders.PowBy = "X-Powered-By"
	httpHeaders.Csp = "Content-Security-Policy"
}

func HttpHeaderOutput(outputBuilder *strings.Builder, response *http.Response, httpHeader string) {
	if server := response.Header.Get(httpHeader); server != "" {
		outputBuilder.WriteString(fmt.Sprintf(" | + %s: %s\n", httpHeader, server))
	}
}
