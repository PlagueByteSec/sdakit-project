package lib

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func HttpClientInit() *http.Client {
	client := &http.Client{
		Timeout: 3 * time.Second,
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
