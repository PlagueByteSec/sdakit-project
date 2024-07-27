package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func Request(pool Pool, host string, url string) {
	response, err := http.Get(url)
	if err != nil {
		return
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		GetPanic("ERROR: failed to read body\n%s\n", err)
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
}

func HttpStatusCode(host string) int {
	buildUrl := fmt.Sprintf("http://%s", host)
	response, err := http.Get(buildUrl)
	if err != nil {
		return -1
	}
	defer response.Body.Close()
	return response.StatusCode
}

func GetCurrentRepoVersion() string {
	url := "https://raw.githubusercontent.com/fhAnso/Sentinel/main/version.txt"
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	TestVersionFail(err)
	response, err := client.Do(request)
	TestVersionFail(err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	TestVersionFail(err)
	return string(body)
}

func GetCurrentLocalVersion() string {
	versionPath := "version.txt"
	version, err := os.ReadFile(versionPath)
	TestVersionFail(err)
	return string(version)
}
