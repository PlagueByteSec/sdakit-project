package lib

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

func Request(pool Pool, host string, url string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to send GET request to: %s", url)
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %s", err)
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
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	response, err := client.Get(url)
	if err != nil {
		return -1
	}
	defer response.Body.Close()
	return response.StatusCode
}

func GetCurrentRepoVersion(failHandler *VersionHandler) string {
	var version string
	url := "https://raw.githubusercontent.com/fhAnso/Sentinel/main/version.txt"
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	TestVersionFail(*failHandler, &version, err)
	response, err := client.Do(request)
	TestVersionFail(*failHandler, &version, err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	TestVersionFail(*failHandler, &version, err)
	version = string(body)
	return version
}
