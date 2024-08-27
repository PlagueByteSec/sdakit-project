package requests

import (
	"Sentinel/lib/shared"

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
	probing "github.com/prometheus-community/pro-bing"
)

func HttpClientInit(args *shared.Args) (*http.Client, error) {
	var client *http.Client
	if args.TorRoute {
		/*
			Parse the TOR proxy URL from constants.go. If successful, create
			an HTTP client configured to use the TOR proxy with the specified timeout.
		*/
		proxyUrl, err := url.Parse(shared.TorProxyUrl)
		if err != nil {
			shared.Glogger.Println(err)
			return nil, errors.New("failed to parse TOR proxy URL: " + err.Error())
		}
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
			Timeout: time.Duration(args.Timeout) * time.Second,
		}
		fmt.Fprintln(shared.GStdout, "[*] All requests will be routet through TOR")
	} else {
		// -r flag not set, use the standard HTTP client with the specified timeout
		client = &http.Client{
			Timeout: time.Duration(args.Timeout) * time.Second,
		}
	}
	fmt.Fprintln(shared.GStdout)
	return client, nil
}

func responseGetBody(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

func requestSendGET(url string, client *http.Client) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		shared.Glogger.Println(err)
		return nil, err
	}
	request.Header.Set("User-Agent", shared.DefaultUserAgent)
	return client.Do(request)
}

func EndpointRequest(client *http.Client, host string, url string) error {
	/*
		Send an HTTP GET request, read the body, and filter each subdomain
		using regex. Duplicates will be removed.
	*/
	response, err := requestSendGET(url, client)
	if err != nil {
		shared.Glogger.Println(err)
		return err
	}
	responseBody, err := responseGetBody(response)
	if err != nil {
		shared.Glogger.Println(err)
		return err
	}
	body := string(responseBody)
	regex := regexp.MustCompile(`[\.a-zA-Z0-9-]+\.` + host)
	matches := regex.FindAllString(body, -1)
	for _, match := range matches {
		// Make sure that only new entries will be added
		if !shared.PoolContainsEntry(shared.GPoolBase.PoolSubdomains, match) {
			shared.GPoolBase.PoolSubdomains = append(shared.GPoolBase.PoolSubdomains, match)
		}
	}
	shared.PoolCleanup(&shared.GPoolBase)
	return nil
}

func HttpStatusCode(client *http.Client, url string) int {
	response, err := requestSendGET(url, client)
	if err != nil {
		shared.Glogger.Println(err)
		return -1
	}
	return response.StatusCode
}

func GetCurrentRepoVersion(client *http.Client) string {
	/*
		Request the version.txt file from GitHub and return
		the value as a string.
	*/
	response, err := requestSendGET(shared.VersionUrl, client)
	if err != nil {
		shared.Glogger.Println(err)
		return shared.NotAvailable
	}
	responseBody, err := responseGetBody(response)
	if err != nil {
		shared.Glogger.Println(err)
		return shared.NotAvailable
	}
	return string(responseBody)
}

func AnalyseHttpHeader(client *http.Client, subdomain string) string {
	/*
		Analyze the response of an HTTP request to determine
		which headers are set.

		Server
		Strict-Transport-Security
		X-Powered-By
		Content-Security-Policy
	*/
	url := fmt.Sprintf("http://%s", subdomain)
	response, err := requestSendGET(url, client)
	if err != nil {
		shared.Glogger.Println(err)
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
	return outputBuilder.String()
}

func ScanPortsSubdomain(subdomain string, ports string) (string, error) {
	/*
		Use the Nmap Go package to perform a simple TCP port scan to
		determine the port states and default services.

		Resource: https://pkg.go.dev/github.com/Ullaakut/nmap/v2
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithTargets(subdomain),
		nmap.WithPorts(ports),
	)
	if err != nil {
		shared.Glogger.Println(err)
		return "", errors.New("nmap scanner init failed: " + err.Error())
	}
	result, _, err := scanner.Run()
	if err != nil {
		shared.Glogger.Println(err)
		return "", errors.New("port scan failed: " + err.Error())
	}
	var (
		output strings.Builder
		mark   string
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
			output.WriteString(summary)
			shared.GSubdomBase.OpenPorts = append(
				shared.GSubdomBase.OpenPorts,
				int(port.ID),
			)
		}
	}
	portResults := output.String()
	return portResults, nil
}

func PingSubdomain(subdomain string, pingCount int) error {
	pinger, err := probing.NewPinger(subdomain)
	if err != nil {
		shared.Glogger.Println(err)
		return err
	}
	pinger.Count = 3
	pinger.SetPrivileged(true)
	err = pinger.Run()
	if err != nil {
		shared.Glogger.Print(err)
		return err
	}
	return nil
}
