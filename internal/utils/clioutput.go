package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fhAnso/Sentinel/v1/internal/shared"
)

func PrintVerbose(format string, args ...interface{}) {
	// Only print content if the -v flag is specified
	if shared.GVerbose {
		fmt.Fprintf(shared.GStdout, format, args...)
	}
}

func PrintProgress(entryCount int) {
	fmt.Fprintf(shared.GStdout, "\rProgress::[%d/%d]", shared.GAllCounter, entryCount)
	shared.GStdout.Flush()
	shared.GAllCounter++
}

func PrintEvaluation(startTime time.Time, count int) {
	// Calculate the time duration and format the summary
	defer shared.GStdout.Flush()
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	var temp strings.Builder
	temp.WriteString("subdomain")
	if count != 1 {
		temp.WriteString("s")
	}
	fmt.Fprintf(shared.GStdout, "[*] %d %s obtained, %d displayed\n", count, temp.String(), shared.GDisplayCount)
	fmt.Fprintf(shared.GStdout, "[*] Finished in %s\n", duration)
}

func PrintBanner(httpClient *http.Client) {
	localVersion := GetCurrentLocalVersion()
	repoVersion := GetCurrentRepoVersion(httpClient)
	fmt.Fprintf(shared.GStdout, " ===[ Sentinel, Version: %s ]===\n\n", localVersion)
	VersionCompare(repoVersion, localVersion)
	shared.GStdout.Flush()
}
