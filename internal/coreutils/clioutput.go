package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fhAnso/Sentinel/v1/internal/shared"
)

type setSummary struct {
	pool    *[]string
	message string
}

func evaluatePool(result setSummary) {
	if len(*result.pool) != 0 {
		fmt.Fprintf(shared.GStdout, result.message)
		for idx, subdomain := range *result.pool {
			fmt.Fprintf(shared.GStdout, " |  %d. %s\n", idx+1, subdomain)
		}
		shared.GStdout.Flush()
	}
}

func plural(poolSize int, value string) string {
	var temp strings.Builder
	temp.WriteString(value)
	if poolSize != 1 {
		temp.WriteString("s")
	}
	return temp.String()
}

func PrintSummary(startTime time.Time, count int) {
	// Calculate the time duration and format the summary
	defer shared.GStdout.Flush()
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	var (
		temp     string
		message  string
		poolSize int
	)
	fmt.Fprintf(shared.GStdout, "[*] Summary:%-30s\n", "")
	// Use setting struct and use function for loop etc
	poolSize = len(shared.GPoolBase.PoolHttpSuccessSubdomains)
	temp = plural(poolSize, "Subdomain")
	message = fmt.Sprintf("[+] Found %d %s That Have Responded To HTTP Requests.\n", poolSize, temp)
	evaluatePool(setSummary{
		pool:    &shared.GPoolBase.PoolHttpSuccessSubdomains,
		message: message,
	})
	poolSize = len(shared.GPoolBase.PoolMailSubdomains)
	message = fmt.Sprintf("[+] Found %d %s Providing A Mail Server\n", poolSize, temp)
	evaluatePool(setSummary{
		pool:    &shared.GPoolBase.PoolMailSubdomains,
		message: message,
	})
	poolSize = len(shared.GPoolBase.PoolApiSubdomains)
	message = fmt.Sprintf("[+] Found %d %s Providing A API\n", poolSize, temp)
	evaluatePool(setSummary{
		pool:    &shared.GPoolBase.PoolApiSubdomains,
		message: message,
	})
	poolSize = len(shared.GPoolBase.PoolLoginSubdomains)
	message = fmt.Sprintf("[+] Found %d Login %s\n", poolSize, plural(poolSize, "Page"))
	evaluatePool(setSummary{
		pool:    &shared.GPoolBase.PoolLoginSubdomains,
		message: message,
	})
	poolSize = len(shared.GPoolBase.PoolCorsSubdomains)
	message = fmt.Sprintf("[+] Found %d %s With CORS Flaws\n", poolSize, temp)
	evaluatePool(setSummary{
		pool:    &shared.GPoolBase.PoolCorsSubdomains,
		message: message,
	})
	fmt.Fprintf(shared.GStdout, "[*] %d %s Obtained, %d Displayed\n", count, temp, shared.GDisplayCount)
	fmt.Fprintf(shared.GStdout, "[*] Finished in %s\n", duration)
	// TODO: Generate a PDF to save the summary of all findings
}

func PrintBanner(httpClient *http.Client) {
	localVersion := GetCurrentLocalVersion()
	repoVersion := GetCurrentRepoVersion(httpClient)
	fmt.Fprintf(shared.GStdout, " ===[ Sentinel, Version: %s ]===\n\n", localVersion)
	VersionCompare(repoVersion, localVersion)
	shared.GStdout.Flush()
}

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
