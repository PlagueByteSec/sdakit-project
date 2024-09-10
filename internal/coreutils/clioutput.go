package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
	"github.com/PlagueByteSec/sentinel-project/v2/pkg"
)

type outputSummary struct {
	temp    string
	message string
}

type setSummary struct {
	pool     *[]string
	poolSize int
	output   outputSummary
}

func evaluatePool(result setSummary) {
	if len(*result.pool) != 0 {
		fmt.Fprintf(shared.GStdout, result.output.message)
		for idx, subdomain := range *result.pool {
			fmt.Fprintf(shared.GStdout, " |  %d. %s\n", idx+1, subdomain)
		}
		shared.GStdout.Flush()
	}
}

func plural(poolSize int, value string) string {
	return pkg.Tern(poolSize != 1, value+"s", value)
}

func PrintSummary(startTime time.Time, count int) {
	// Calculate the time duration and format the summary
	defer shared.GStdout.Flush()
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Fprintf(shared.GStdout, "[*] Summary:%-30s\n", "")
	// Use setting struct and use function for loop etc
	poolSize := len(shared.GPoolBase.PoolHttpSuccessSubdomains)
	temp := plural(poolSize, "Subdomain")
	evaluatePool(setSummary{
		pool:     &shared.GPoolBase.PoolHttpSuccessSubdomains,
		poolSize: poolSize,
		output: outputSummary{
			temp:    temp,
			message: fmt.Sprintf("[+] Found %d %s that have responded to HTTP requests.\n", poolSize, temp),
		},
	})
	poolSize = len(shared.GPoolBase.PoolMailSubdomains)
	temp = plural(poolSize, "Subdomain")
	evaluatePool(setSummary{
		pool:     &shared.GPoolBase.PoolMailSubdomains,
		poolSize: poolSize,
		output: outputSummary{
			temp:    temp,
			message: fmt.Sprintf("[+] Found %d %s providing a mail server\n", poolSize, temp),
		},
	})
	poolSize = len(shared.GPoolBase.PoolApiSubdomains)
	temp = plural(poolSize, "Subdomain")
	evaluatePool(setSummary{
		pool:     &shared.GPoolBase.PoolApiSubdomains,
		poolSize: poolSize,
		output: outputSummary{
			temp:    temp,
			message: fmt.Sprintf("[+] Found %d %s providing a API\n", poolSize, temp),
		},
	})
	poolSize = len(shared.GPoolBase.PoolLoginSubdomains)
	temp = plural(poolSize, "Subdomain")
	evaluatePool(setSummary{
		pool:     &shared.GPoolBase.PoolLoginSubdomains,
		poolSize: poolSize,
		output: outputSummary{
			temp:    temp,
			message: fmt.Sprintf("[+] Found %d login %s\n", poolSize, plural(poolSize, "Page")),
		},
	})
	poolSize = len(shared.GPoolBase.PoolCmsSubdomains)
	evaluatePool(setSummary{
		pool:     &shared.GPoolBase.PoolCmsSubdomains,
		poolSize: poolSize,
		output: outputSummary{
			message: fmt.Sprintf("[+] Identified %d CMS\n", poolSize),
		},
	})
	poolSize = len(shared.GPoolBase.PoolCorsSubdomains)
	temp = plural(poolSize, "Subdomain")
	evaluatePool(setSummary{
		pool:     &shared.GPoolBase.PoolCorsSubdomains,
		poolSize: poolSize,
		output: outputSummary{
			temp:    temp,
			message: fmt.Sprintf("[+] Found %d %s with possible CORS flaws\n", poolSize, temp),
		},
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
