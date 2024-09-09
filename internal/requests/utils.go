package requests

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PlagueByteSec/Sentinel/v2/internal/shared"
)

func HttpCodeCheck(settings shared.SettingsHandler, url string) bool {
	return HttpStatusCode(settings.HttpClient, url, settings.Args.HttpRequestMethod) != -1
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
