package lib

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"time"
)

func ClientInit() *http.Client {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	return client
}

func EndpointRequest(client *http.Client, host string, url string) error {
	response, err := client.Get(url)
	if err != nil {
		return errors.New("failed to send GET request to: " + url)
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.New("failed to read body: " + err.Error())
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
	response, err := client.Get(url)
	if err != nil {
		return -1
	}
	defer response.Body.Close()
	return response.StatusCode
}

func GetCurrentRepoVersion(client *http.Client, failHandler *VersionHandler) string {
	var version string
	response, err := client.Get(VersionUrl)
	TestVersionFail(*failHandler, &version, err)
	defer response.Body.Close()
	if version == Na {
		// Instant return to avoid ReadAll execution if request failed
		return version
	}
	body, err := io.ReadAll(response.Body)
	TestVersionFail(*failHandler, &version, err)
	version = string(body)
	return version
}
