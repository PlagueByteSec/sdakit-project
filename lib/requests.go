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

func Request(pool Pool, host string, url string) error {
	client := ClientInit()
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
		if !pool.ContainsEntry(match) {
			pool.AddEntry(match)
		}
	}
	return nil
}

func HttpStatusCode(url string) int {
	client := ClientInit()
	response, err := client.Get(url)
	if err != nil {
		return -1
	}
	defer response.Body.Close()
	return response.StatusCode
}

func GetCurrentRepoVersion(failHandler *VersionHandler) string {
	var version string
	const url = "https://raw.githubusercontent.com/fhAnso/Sentinel/main/version.txt"
	client := ClientInit()
	response, err := client.Get(url)
	TestVersionFail(*failHandler, &version, err)
	defer response.Body.Close()
	if version == na {
		// Instant return to avoid ReadAll execution if request failed
		return version
	}
	body, err := io.ReadAll(response.Body)
	TestVersionFail(*failHandler, &version, err)
	version = string(body)
	return version
}
