package lib

import (
	"fmt"
	"os"
	"time"
)

func DefaultOutputName(hostname string) string {
	currentTime := time.Now()
	formatTime := currentTime.Format("2006-01-02_15-04-05")
	outputFile := fmt.Sprintf("%s-%s.txt", formatTime, hostname)
	return outputFile
}

func CreateOutputDir() {
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			GetPanic("Unable to create output directory: %s\n", outputDir)
		}
	}
}

func VersionCompare() {
	repo := GetCurrentRepoVersion()
	local := GetCurrentLocalVersion()
	if repo == "n/a" || local == "n/a" {
		return
	}
	if repo != local {
		fmt.Printf("[*] An update is available! %s->%s\n", local, repo)
	}
}

func TestVersionFail(err error) string {
	var value string
	if err != nil {
		value = "n/a"
	}
	return value
}
