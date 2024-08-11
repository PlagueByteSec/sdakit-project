package lib

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Ullaakut/nmap/v3"
)

func HttpClientInit() *http.Client {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	return client
}

func responseGetBody(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

func requestSendGET(url string, client *http.Client) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", DefaultUserAgent)
	return client.Do(request)
}

func EndpointRequest(client *http.Client, host string, url string) error {
	response, err := requestSendGET(url, client)
	if err != nil {
		return err
	}
	responseBody, err := responseGetBody(response)
	if err != nil {
		return err
	}
	// Filter the HTML reponse for results
	body := string(responseBody)
	regex := regexp.MustCompile("[\\.a-zA-Z0-9-]+\\." + host)
	matches := regex.FindAllString(body, -1)
	for _, match := range matches {
		// Make sure that only new entries will be added
		if !PoolContainsEntry(PoolDomains, match) {
			PoolDomains = append(PoolDomains, match)
		}
	}
	return nil
}

func HttpStatusCode(client *http.Client, url string) int {
	response, err := requestSendGET(url, client)
	if err != nil {
		return -1
	}
	return response.StatusCode
}

func GetCurrentRepoVersion(client *http.Client) string {
	response, err := requestSendGET(VersionUrl, client)
	if err != nil {
		return "n/a"
	}
	responseBody, err := responseGetBody(response)
	if err != nil {
		return "n/a"
	}
	return string(responseBody)
}

func AnalyseHttpHeader(client *http.Client, subdomain string) (string, int) {
	url := fmt.Sprintf("http://%s", subdomain)
	response, err := requestSendGET(url, client)
	if err != nil {
		return "", 0
	}
	results := make([]string, 0)
	if server := response.Header.Get("Server"); server != "" {
		results = append(results, server)
	}
	if hsts := response.Header.Get("Strict-Transport-Security"); hsts != "" {
		results = append(results, "HSTS")
	}
	result := strings.Join(results, ",")
	return "╚═[ " + result, len(result)
}

func ScanPortsSubdomain(subdomain string, ports string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithTargets(subdomain),
		nmap.WithPorts(ports),
	)
	if err != nil {
		return "", errors.New("nmap scanner init failed: " + err.Error())
	}
	result, _, err := scanner.Run()
	if err != nil {
		return "", errors.New("port scan failed: " + err.Error())
	}
	var output strings.Builder
	for _, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}
		for _, port := range host.Ports {
			summary := fmt.Sprintf("\n\t[> Port %d/%s %s %s\n",
				port.ID, port.Protocol, port.State, port.Service.Name)
			output.WriteString(summary)
		}
	}
	return output.String(), nil
}
