package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PlagueByteSec/sdakit-project/v2/internal/coreutils/report"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/datapools"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/shared"
	"github.com/PlagueByteSec/sdakit-project/v2/pkg"
)

type summaryConfig struct {
	reportGenerator *report.ReportGenerator
	pool            []string
	categoryName    string
	messageFormat   string
	noSup           bool
}

func summaryEvaluatePool(reportGenerator *report.ReportGenerator, config report.ReportGeneratorConfig) {
	if len(*config.Pool) != 0 {
		reportGenerator.WriteToReport(`<ol>` + "\n")
		fmt.Fprintf(shared.GStdout, config.Output.Message)
		for idx, subdomain := range *config.Pool {
			fmt.Fprintf(shared.GStdout, " |  %d. %s\n", idx+1, subdomain)
			reportGenerator.WriteToReport(`<h3><li>` + subdomain + `</li></h3>` + "\n")
		}
		reportGenerator.WriteToReport(`</ol>` + "\n")
		shared.GStdout.Flush()
	}
}

func plural(poolSize int, value string) string {
	return pkg.Tern(poolSize != 1, value+"s", value)
}

func generateSummary(config summaryConfig) {
	var temp string
	poolSize := len(config.pool)
	if !config.noSup {
		temp = plural(poolSize, "Subdomain")
	}
	reportContent := fmt.Sprintf(`<h2 id="category">` + config.categoryName + `:</h2>` + "\n")
	if poolSize != 0 {
		config.reportGenerator.WriteToReport(reportContent)
	}
	summaryEvaluatePool(config.reportGenerator, report.ReportGeneratorConfig{
		Pool:     &config.pool,
		PoolSize: poolSize,
		Output: report.ReportOutput{
			Temp:    temp,
			Message: fmt.Sprintf(config.messageFormat, poolSize, temp),
		},
	})
}

func WriteSummary(startTime time.Time, count int) {
	// Calculate the time duration and format the summary
	defer shared.GStdout.Flush()
	datapools.PoolsCleanupSummary(&shared.GPoolBase)
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	reportGenerator, err := report.StartReportGenerator()
	if err != nil {
		ProgramExit(ExitParams{
			ExitCode:    -1,
			ExitMessage: "Could not start the report generator",
			ExitError:   err,
		})
	}
	report.GenerateTotalResultsReport(reportGenerator)
	fmt.Fprintf(shared.GStdout, "\r[*] Summary:%-30s\n", "")
	generateSummary(summaryConfig{
		reportGenerator: reportGenerator,
		pool:            shared.GPoolBase.PoolHttpSuccessSubdomains,
		categoryName:    "HTTP Success Subdomains",
		messageFormat:   "[+] Found %d %s that have responded to HTTP requests.\n",
		noSup:           false,
	})
	generateSummary(summaryConfig{
		reportGenerator: reportGenerator,
		pool:            shared.GPoolBase.PoolMailSubdomains,
		categoryName:    "Mail",
		messageFormat:   "[+] Found %d %s providing a mail server\n",
		noSup:           false,
	})
	generateSummary(summaryConfig{
		reportGenerator: reportGenerator,
		pool:            shared.GPoolBase.PoolApiSubdomains,
		categoryName:    "APIs",
		messageFormat:   "[+] Found %d %s providing an API\n",
		noSup:           false,
	})
	generateSummary(summaryConfig{
		reportGenerator: reportGenerator,
		pool:            shared.GPoolBase.PoolLoginSubdomains,
		categoryName:    "Logins",
		messageFormat:   "[+] Found %d login%s\n",
		noSup:           true,
	})
	generateSummary(summaryConfig{
		reportGenerator: reportGenerator,
		pool:            shared.GPoolBase.PoolCmsSubdomains,
		categoryName:    "CMS",
		messageFormat:   "[+] Identified %d CMS%s\n",
		noSup:           true,
	})
	generateSummary(summaryConfig{
		reportGenerator: reportGenerator,
		pool:            shared.GPoolBase.PoolCorsSubdomains,
		categoryName:    "CORS, Status: FOUND",
		messageFormat:   "[+] Found %d %s with possible CORS flaws\n",
		noSup:           false,
	})
	generateSummary(summaryConfig{
		reportGenerator: reportGenerator,
		pool:            shared.GPoolBase.PoolCookieInjection,
		categoryName:    "Cookie injection, Status: FOUND",
		messageFormat:   "[+] Found %d %s with possible cookie injection flaws\n",
		noSup:           false,
	})
	generateSummary(summaryConfig{
		reportGenerator: reportGenerator,
		pool:            shared.GPoolBase.PoolRequestSmuggling,
		categoryName:    "Request smuggling, Status: FOUND",
		messageFormat:   "[+] Found %d %s with possible request smuggling flaws\n",
		noSup:           false,
	})
	report.GenerateTestReport(reportGenerator)
	reportGenerator.CloseReportGenerator()
	fmt.Fprintf(shared.GStdout, "[*] %d Obtained, %d Displayed\n", count, shared.GDisplayCount)
	fmt.Fprintf(shared.GStdout, "[*] Finished in %s\n", duration)
}

func PrintStartInfo(domain string, wordlist string, methodKey string) {
	shared.GScanMethod = methodKey
	fmt.Fprintf(shared.GStdout, "[*] Domain: %s\n", domain)
	fmt.Fprintf(shared.GStdout, "[*] Discovery Method: %s\n", methodKey)
	if len(wordlist) != 0 {
		fmt.Fprintf(shared.GStdout, "[*] Wordlist: %s\n", wordlist)
	}
	fmt.Fprintf(shared.GStdout, "[*] Start Time: %s\n", time.Now().Format(time.RFC850))
}

func PrintBanner(httpClient *http.Client) {
	localVersion := GetCurrentLocalVersion()
	repoVersion := GetCurrentRepoVersion(httpClient)
	var banner strings.Builder
	banner.WriteString("\n           - The SDAkit Project - \n")
	banner.WriteString("Subdomain Discovery and Security Analysis Toolkit\n\n")
	banner.WriteString("  License: MIT\n")
	banner.WriteString(fmt.Sprintf("  Version: %s\n", localVersion))
	banner.WriteString("  By @PlagueByte.Sec\n")
	banner.WriteString("  Contact: plaguebyte.sec@keemail.me\n")
	fmt.Fprintln(shared.GStdout, pkg.BuildBanner(banner.String()))
	VersionCompare(repoVersion, localVersion)
	fmt.Fprintln(shared.GStdout)
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
