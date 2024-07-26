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
		fmt.Printf("ERROR: failed to read body\n%s\n", err)
		os.Exit(-1)
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
