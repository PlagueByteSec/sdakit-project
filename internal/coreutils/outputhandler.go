package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PlagueByteSec/sentinel-project/v2/internal/coreutils/summary"
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

type summaryConfigGenerator struct {
	reportGenerator *summary.SummaryGenerator
	pool            []string
	categoryName    string
	messageFormat   string
	noSup           bool
}

func evaluatePool(rg *summary.SummaryGenerator, result setSummary) {
	if len(*result.pool) != 0 {
		rg.SummaryWriteToStream("<ol>\n")
		fmt.Fprintf(shared.GStdout, result.output.message)
		for idx, subdomain := range *result.pool {
			fmt.Fprintf(shared.GStdout, " |  %d. %s\n", idx+1, subdomain)
			rg.SummaryWriteToStream("<h3><li>" + subdomain + "</li></h3>\n")
		}
		rg.SummaryWriteToStream("</ol>\n")
		shared.GStdout.Flush()
	}
}

func plural(poolSize int, value string) string {
	return pkg.Tern(poolSize != 1, value+"s", value)
}

func generateSummary(scg summaryConfigGenerator) {
	var temp string
	poolSize := len(scg.pool)
	if !scg.noSup {
		temp = plural(poolSize, "Subdomain")
	}
	reportContent := fmt.Sprintf("<h2 id=\"category\">%s:</h2>\n", scg.categoryName)
	if poolSize != 0 {
		scg.reportGenerator.SummaryWriteToStream(reportContent)
	}
	evaluatePool(scg.reportGenerator, setSummary{
		pool:     &scg.pool,
		poolSize: poolSize,
		output: outputSummary{
			temp:    temp,
			message: fmt.Sprintf(scg.messageFormat, poolSize, temp),
		},
	})
}

func WriteSummary(startTime time.Time, count int) {
	// Calculate the time duration and format the summary
	defer shared.GStdout.Flush()
	shared.PoolsCleanupSummary(&shared.GPoolBase)
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	rg, err := summary.ReportGeneratorInit()
	if err != nil {
		ProgramExit(ExitParams{
			ExitCode:    -1,
			ExitMessage: "ReportGeneratorInit failed",
			ExitError:   err,
		})
	}
	summary.GenerateTotalResultsReport(rg)
	fmt.Fprintf(shared.GStdout, "\r[*] Summary:%-30s\n", "")
	generateSummary(summaryConfigGenerator{
		reportGenerator: rg,
		pool:            shared.GPoolBase.PoolHttpSuccessSubdomains,
		categoryName:    "HTTP Success Subdomains",
		messageFormat:   "[+] Found %d %s that have responded to HTTP requests.\n",
	})
	generateSummary(summaryConfigGenerator{
		reportGenerator: rg,
		pool:            shared.GPoolBase.PoolMailSubdomains,
		categoryName:    "Mail Servers",
		messageFormat:   "[+] Found %d %s providing a mail server\n",
	})
	generateSummary(summaryConfigGenerator{
		reportGenerator: rg,
		pool:            shared.GPoolBase.PoolApiSubdomains,
		categoryName:    "APIs",
		messageFormat:   "[+] Found %d %s providing an API\n",
	})
	generateSummary(summaryConfigGenerator{
		reportGenerator: rg,
		pool:            shared.GPoolBase.PoolLoginSubdomains,
		categoryName:    "Logins",
		messageFormat:   "[+] Found %d login%s\n",
	})
	generateSummary(summaryConfigGenerator{
		reportGenerator: rg,
		pool:            shared.GPoolBase.PoolCmsSubdomains,
		categoryName:    "CMS",
		messageFormat:   "[+] Identified %d CMS%s\n",
	})
	generateSummary(summaryConfigGenerator{
		reportGenerator: rg,
		pool:            shared.GPoolBase.PoolCorsSubdomains,
		categoryName:    "CORS, Status: FOUND",
		messageFormat:   "[+] Found %d %s with possible CORS flaws\n",
	})
	generateSummary(summaryConfigGenerator{
		reportGenerator: rg,
		pool:            shared.GPoolBase.PoolCookieInjection,
		categoryName:    "Cookie injection, Status: FOUND",
		messageFormat:   "[+] Found %d %s with possible cookie injection flaws\n",
	})
	generateSummary(summaryConfigGenerator{
		reportGenerator: rg,
		pool:            shared.GPoolBase.PoolRequestSmuggling,
		categoryName:    "Request smuggling, Status: FOUND",
		messageFormat:   "[+] Found %d %s with possible request smuggling flaws\n",
	})
	summary.GenerateTestReport(rg)
	rg.SummaryFileClose()
	fmt.Fprintf(shared.GStdout, "[*] %d Obtained, %d Displayed\n", count, shared.GDisplayCount)
	fmt.Fprintf(shared.GStdout, "[*] Finished in %s\n", duration)
}

func PrintBanner(httpClient *http.Client) {
	localVersion := GetCurrentLocalVersion()
	repoVersion := GetCurrentRepoVersion(httpClient)
	fmt.Fprintf(shared.GStdout, "        The Sentinel Project, Version: %-8s       \n", localVersion)
	fmt.Fprint(shared.GStdout, "   Subdomain Discovery and Security Analysis Toolkit  \n\n")
	VersionCompare(repoVersion, localVersion)
	shared.GStdout.Flush()
}

func PrintVerbose(format string, args ...any) {
	if shared.GVerbose {
		fmt.Fprintf(shared.GStdout, format, args...)
	}
}

func PrintProgress(entryCount int) {
	fmt.Fprintf(shared.GStdout, "\rProgress::[%d/%d]", shared.GAllCounter, entryCount)
	shared.GStdout.Flush()
	shared.GAllCounter++
}
