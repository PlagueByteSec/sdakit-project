package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PlagueByteSec/sdakit-project/v2/internal/logging"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/shared"
)

type ReportGenerator struct {
	TargetDomain string
	Filename     string
	stream       *os.File
}

type ReportOutput struct {
	Temp    string
	Message string
}

type ReportGeneratorConfig struct {
	Pool     *[]string
	PoolSize int
	Output   ReportOutput
}

func StartReportGenerator() (*ReportGenerator, error) {
	var err error
	reportGenerator := ReportGenerator{
		TargetDomain: shared.GTargetDomain,
	}
	reportGenerator.Filename = filepath.Join("output", reportGenerator.TargetDomain+"_tsp-summary.html")
	reportGenerator.stream, err = os.Create(reportGenerator.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file stream (%s): %s", reportGenerator.Filename, err)
	}
	reportGenerator.WriteToReport(ReportStart)
	reportGenerator.generateOverviewReport(shared.GScanMethod)
	return &reportGenerator, nil
}

func (reportGenerator *ReportGenerator) generateOverviewReport(method string) {
	if len(method) == 0 {
		logging.GLogger.Log("could not accept empty value for method")
		return
	}
	reportGenerator.WriteToReport(`<div id="table-container">`)
	timestamp := time.Now().Format(time.RFC850)
	overview := `<tr>
	<td>` + reportGenerator.TargetDomain + `</td>
	<td>` + timestamp + `</td>
	<td>` + method + `</td>
</tr>
</table>
</div>`
	reportGenerator.WriteToReport(overview)
}

func (reportGenerator *ReportGenerator) WriteToReport(content string) {
	if _, err := reportGenerator.stream.WriteString(content); err != nil {
		logging.GLogger.Log(fmt.Sprintf("failed to append content to report: %s", err))
	}
}

func GenerateTotalResultsReport(reportGenerator *ReportGenerator) {
	reportGenerator.WriteToReport(`<h1 id="category-headline">Enumeration Results</h1>
<h2 id="category">All Subdomains:</h2>
<ol>`)
	subdomainPool := shared.GPoolBase.PoolSubdomains
	for idx := 0; idx < len(subdomainPool); idx++ {
		reportGenerator.WriteToReport(`<h3><li>` + subdomainPool[idx] + `</li></h3>`)
	}
	reportGenerator.WriteToReport(`</ol>`)
}

func GenerateTestReport(reportGenerator *ReportGenerator) {
	var result strings.Builder
	if len(shared.GReportPool) == 0 {
		return
	}
	reportGenerator.WriteToReport(`<div id="table-container">`)
	for _, value := range shared.GReportPool {
		result.WriteString(`<tr>`)
		result.WriteString(`<td>` + value.Subdomain + `</td>`)
		result.WriteString(`<td>` + value.TestName + `</td>`)
		if value.TestResult == "PASSED" {
			result.WriteString(`<td style="color: green;">` + value.TestResult + `</td>`)
		} else if value.TestResult == "FOUND" {
			result.WriteString(`<td style="color: red;">` + value.TestResult + `</td>`)
		}
		result.WriteString(`</tr>`)
	}
	output := `</table>
<h1 id="category-headline">Analysis Results</h1>
<table id="overview-table">
<tr>
	<th>Subdomain</th>
    <th>Test</th>
    <th>Result</th>
</tr>` + result.String() + "</div>"
	reportGenerator.WriteToReport(output)
}

func (reportGenerator *ReportGenerator) CloseReportGenerator() {
	fmt.Fprintln(shared.GStdout, "[+] HTML report generated:\n |  => "+reportGenerator.Filename)
	reportGenerator.stream.WriteString(ReportEnd)
	reportGenerator.stream.Close()
}
