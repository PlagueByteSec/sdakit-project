package lib

import (
	"errors"
	"fmt"
	"os"
	"runtime"
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

func GetCurrentLocalVersion() string {
	var versionPath string
	if runtime.GOOS == "windows" {
		versionPath = "..\\version.txt"
	} else if runtime.GOOS == "linux" {
		versionPath = "../version.txt"
	}
	if _, err := os.Stat(versionPath); errors.Is(err, os.ErrNotExist) {
		versionPath = "version.txt"
	}
	version, err := os.ReadFile(versionPath)
	TestVersionFail(err)
	return string(version)
}

func VersionCompare() {
	repo := GetCurrentRepoVersion()
	local := GetCurrentLocalVersion()
	if repo == "n/a" || local == "n/a" || local == "" {
		return
	}
	if repo != local {
		fmt.Printf("[*] An update is available! %s->%s\n", local, repo)
	}
}

func IsInExclude(httpCode string, list []string) bool {
	for _, code := range list {
		if httpCode == code {
			return true
		}
	}
	return false
}
