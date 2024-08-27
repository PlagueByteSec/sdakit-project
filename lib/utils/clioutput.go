package utils

import (
	"Sentinel/lib/requests"
	"Sentinel/lib/shared"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func VerbosePrint(format string, args ...interface{}) {
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

func Evaluation(startTime time.Time, count int) {
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

func SentinelPrintBanner(httpClient *http.Client) {
	localVersion := GetCurrentLocalVersion()
	repoVersion := requests.GetCurrentRepoVersion(httpClient)
	fmt.Fprintf(shared.GStdout, " ===[ Sentinel, Version: %s ]===\n\n", localVersion)
	VersionCompare(repoVersion, localVersion)
	shared.GStdout.Flush()
}

func SentinelExit(exitParams shared.SentinelExitParams) {
	/*
		Read the exit settings specified in SentinelExitParams and
		adjust the behavior based on those settings.
	*/
	fmt.Fprintln(shared.GStdout, exitParams.ExitMessage)
	if exitParams.ExitError != nil {
		errorMessage := fmt.Sprintf("Sentinel exit with an error: %s", exitParams.ExitError.Error())
		shared.Glogger.Println(errorMessage)
		fmt.Fprintln(shared.GStdout, errorMessage)
	}
	shared.GStdout.Flush()
	os.Exit(exitParams.ExitCode)
}

func SentinelPanic(err error) {
	fmt.Fprintf(shared.GStdout, "\r%-50s\n", err)
	shared.GStdout.Flush()
	shared.Glogger.Println(err)
	shared.Glogger.Fatalf("Program execution failed")
}
