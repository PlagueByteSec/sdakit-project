package requests

import (
	"Sentinel/lib/shared"
	"fmt"
	"net/http"
	"strings"
)

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
