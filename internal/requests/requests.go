package requests

import (
	pools "github.com/PlagueByteSec/sdakit-project/v2/internal/datapools"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/logging"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/shared"

	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/Ullaakut/nmap/v3"
)

type HttpRequestBase struct {
	CliArgs                *shared.Args
	HttpClient             *http.Client
	HttpRequest            *http.Request
	HttpResponse           *http.Response
	HttpResponseBody       []byte
	HttpNeedResponse       bool
	ResponseNeedStatusCode bool
	ResponseNeedBody       bool
	HttpMethod             string
	CustomHeader           string
	CustomUrl              string
	Subdomain              string
	Domain                 string
	Error                  error
}

var httpRequestBase = &HttpRequestBase{}

func ResetHttpRequestBase(base *HttpRequestBase) {
	*base = *httpRequestBase
}

func HttpClientInit(args *shared.Args) (*http.Client, error) {
	var httpClient *http.Client
	switch {
	case args.TorRoute:
		/*
			Parse the TOR proxy URL from constants.go. If successful, create
			an HTTP client configured to use the TOR proxy with the specified timeout.
		*/
		proxyUrl, err := url.Parse(shared.TorProxyUrl)
		if err != nil {
			logging.GLogger.Log(err.Error())
			return nil, errors.New("failed to parse TOR proxy URL: " + err.Error())
		}
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
			Timeout: time.Duration(args.Timeout) * time.Second,
		}
		if args.Verbose {
			fmt.Fprintln(shared.GStdout, "[*] All requests will be routet through TOR")
		}
	case args.AllowRedirects:
		httpClient = &http.Client{
			Timeout: time.Duration(args.Timeout) * time.Second,
			CheckRedirect: func(request *http.Request, via []*http.Request) error {
				return nil
			},
		}
	default:
		// -r flag not set, use the standard HTTP client with the specified timeout
		httpClient = &http.Client{
			Timeout: time.Duration(args.Timeout) * time.Second,
		}
	}
	fmt.Fprintln(shared.GStdout)
	return httpClient, nil
}

func (base *HttpRequestBase) requestSetupHTTP() error {
	acceptedMethods := []string{"GET", "POST", "OPTIONS"}
	if !pools.ManagePool(pools.PoolAction(pools.PoolCheck), base.HttpMethod, &acceptedMethods) {
		return fmt.Errorf("HTTP request method not allowed: %s", base.HttpMethod)
	}
	base.HttpRequest, base.Error = http.NewRequest(base.HttpMethod, base.CustomUrl, nil)
	if base.Error != nil {
		return fmt.Errorf("failed to setup wrapper for NewRequestWithContext: %s", base.Error)
	}
	base.HttpRequest.Header.Set("User-Agent", shared.DefaultUserAgent)
	return nil
}

func (base *HttpRequestBase) requestGetReponse() error {
	if len(base.Subdomain) != 0 {
		base.HttpRequest.Host = base.Subdomain
	}
	base.HttpResponse, base.Error = base.HttpClient.Do(base.HttpRequest)
	if base.Error != nil {
		return fmt.Errorf("HTTP client could not send request: %s", base.Error)
	}
	return nil
}

func (base *HttpRequestBase) responseGetBody() error {
	defer base.HttpResponse.Body.Close()
	base.HttpResponseBody, base.Error = io.ReadAll(base.HttpResponse.Body)
	if base.Error != nil {
		return fmt.Errorf("failed to read response body: %s", base.Error)
	}
	return nil
}

func RequestHandlerCore(base *HttpRequestBase) (*http.Response, int, []byte, error) {
	if base.Error = base.requestSetupHTTP(); base.Error != nil {
		return nil, -1, nil, base.Error
	}
	if base.Error = base.requestGetReponse(); base.Error != nil {
		return nil, -1, nil, base.Error
	}
	if base.HttpNeedResponse {
		return base.HttpResponse, -1, nil, nil
	}
	defer base.HttpResponse.Body.Close()
	if base.ResponseNeedBody {
		if base.Error = base.responseGetBody(); base.Error != nil {
			return nil, -1, nil, base.Error
		}
	}
	statusCode := base.HttpResponse.StatusCode
	if base.ResponseNeedStatusCode && base.ResponseNeedBody {
		return nil, statusCode, base.HttpResponseBody, nil
	} else if base.ResponseNeedStatusCode {
		return nil, statusCode, nil, nil
	} else if base.ResponseNeedBody {
		return nil, -1, base.HttpResponseBody, nil
	}
	return nil, statusCode, nil, nil
}

func EndpointRequest(method string, host string, url string, client *http.Client) error {
	/*
		Send an HTTP [method] request, read the body, and filter each subdomain
		using regex. Duplicates will be removed.
	*/
	_, _, responseBody, err := RequestHandlerCore(&HttpRequestBase{
		HttpMethod:       method,
		Domain:           host,
		CustomUrl:        url,
		HttpClient:       client,
		ResponseNeedBody: true,
	})
	if err != nil {
		logging.GLogger.Log(err.Error())
		return err
	}
	body := string(responseBody)
	regex := regexp.MustCompile(`[\.a-zA-Z0-9-]+\.` + host)
	matches := regex.FindAllString(body, -1)
	for idx := 0; idx < len(matches); idx++ {
		// Make sure that only new entries will be added
		if !pools.ManagePool(pools.PoolAction(pools.PoolCheck), matches[idx], &shared.GPoolBase.PoolSubdomains) {
			shared.GPoolBase.PoolSubdomains = append(shared.GPoolBase.PoolSubdomains, matches[idx])
		}
	}
	pools.PoolsCleanupCore(&shared.GPoolBase)
	return nil
}

func AnalyseHttpHeader(client *http.Client, subdomain string, method string) string {
	/*
		Analyze the response of an HTTP request to determine
		which headers are set.

		Server
		Strict-Transport-Security
		X-Powered-By
		Content-Security-Policy
	*/
	url := fmt.Sprintf("http://%s", subdomain)
	response, _, _, err := RequestHandlerCore(&HttpRequestBase{
		HttpClient:       client,
		CustomUrl:        url,
		HttpMethod:       method,
		HttpNeedResponse: true,
	})
	if err != nil {
		logging.GLogger.Log(err.Error())
		return ""
	}
	var (
		httpHeaders   shared.HttpHeaders
		outputBuilder strings.Builder
	)
	HttpHeaderInit(&httpHeaders)
	headers := reflect.ValueOf(httpHeaders)
	for idx := 0; idx < headers.NumField(); idx++ {
		value := headers.Field(idx)
		HttpHeaderOutput(&outputBuilder, response, value.String())
	}
	if shared.GShowAllHeaders {
		outputBuilder.WriteString("[*] Full Response:\n")
		for header, headerValue := range response.Header {
			outputBuilder.WriteString(fmt.Sprintf(" | %s: %s\n", header,
				strings.Join(headerValue, " ")))
		}
	}
	return outputBuilder.String()
}

func ScanPortRange(address string, ports string, portsOnly bool) (string, []uint16, error) {
	/*
		Use the Nmap Go package to perform a simple TCP port scan to
		determine the port states and default services.

		Resource: https://pkg.go.dev/github.com/Ullaakut/nmap/v2
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithTargets(address),
		nmap.WithPorts(ports),
	)
	if err != nil {
		logging.GLogger.Log(err.Error())
		return "", nil, errors.New("nmap scanner init failed: " + err.Error())
	}
	result, _, err := scanner.Run()
	if err != nil {
		logging.GLogger.Log(err.Error())
		return "", nil, errors.New("port scan failed: " + err.Error())
	}
	var (
		output    strings.Builder
		mark      string
		openPorts []uint16
	)
	for _, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}
		for _, port := range host.Ports {
			if port.State.String() == "open" {
				mark = "+"
			} else {
				mark = "-"
			}
			summary := fmt.Sprintf(" | %s Port %d/%s %s %s\n",
				mark, port.ID, port.Protocol, port.State, port.Service.Name)
			isOpen := port.State.String() == "open"
			if portsOnly && isOpen {
				openPorts = append(openPorts, port.ID)
				continue
			} else if isOpen {
				output.WriteString(summary)
				shared.GSubdomBase.OpenPorts = append(
					shared.GSubdomBase.OpenPorts,
					int(port.ID),
				)
			}
		}
	}
	if portsOnly {
		return "", openPorts, nil
	}
	portResults := output.String()
	return portResults, nil, nil
}
