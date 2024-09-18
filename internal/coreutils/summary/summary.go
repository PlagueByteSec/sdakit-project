package summary

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PlagueByteSec/sentinel-project/v2/internal/logging"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
)

type SummaryGenerator struct {
	TargetDomain string
	Filename     string
	stream       *os.File
}

func (gs *SummaryGenerator) SummaryWriteToStream(content string) {
	if _, err := gs.stream.WriteString(content); err != nil {
		logging.GLogger.Log(fmt.Sprintf("failed to append content to report: %s", err))
	}
}

func (gs *SummaryGenerator) SummaryFileInit() error {
	var err error
	gs.Filename = filepath.Join("output", gs.TargetDomain+"_tsp-summary.html")
	gs.stream, err = os.Create(gs.Filename)
	if err != nil {
		return fmt.Errorf("failed to open file stream (%s): %s", gs.Filename, err)
	}
	base := generateSummaryTemplate()
	gs.SummaryWriteToStream(base)
	return nil
}

func (gs *SummaryGenerator) SummaryFileClose() {
	fmt.Fprintln(shared.GStdout, "[+] HTML report generated:\n |  => "+gs.Filename)
	documentEnd := closeSummaryTemplate()
	gs.stream.WriteString(documentEnd)
	gs.stream.Close()
}

func (gs *SummaryGenerator) WriteOverview(method string) {
	if len(method) == 0 {
		logging.GLogger.Log("could not accept empty value for method")
		return
	}
	gs.SummaryWriteToStream("<div id=\"table-container\">")
	timestamp := time.Now().Format(time.RFC850)
	overview := `<tr>
	<td>` + gs.TargetDomain + `</td>
	<td>` + timestamp + `</td>
	<td>` + method + `</td>
</tr>
</table>
</div>`
	gs.SummaryWriteToStream(overview)
}

func ReportGeneratorInit() (*SummaryGenerator, error) {
	reportGenerator := SummaryGenerator{
		TargetDomain: shared.GTargetDomain,
	}
	if err := reportGenerator.SummaryFileInit(); err != nil {
		return nil, err
	}
	reportGenerator.WriteOverview(shared.GScanMethod)
	return &reportGenerator, nil
}

func GenerateTotalResultsReport(rg *SummaryGenerator) {
	rg.SummaryWriteToStream("<h1 id=\"category-headline\">Enumeration Results</h1>\n<h2 id=\"category\">All Subdomains:</h2>\n<ol>\n")
	subdomainPool := shared.GPoolBase.PoolSubdomains
	for idx := 0; idx < len(subdomainPool); idx++ {
		rg.SummaryWriteToStream("<h3><li>" + subdomainPool[idx] + "</li></h3>\n")
	}
	rg.SummaryWriteToStream("</ol>\n")
}

func GenerateTestReport(rg *SummaryGenerator) {
	if len(shared.GReportPool) == 0 {
		return
	}
	rg.SummaryWriteToStream("<div id=\"table-container\">")
	var result strings.Builder
	for _, value := range shared.GReportPool {
		result.WriteString("<tr>\n")
		result.WriteString(fmt.Sprintf("<td>%s</td>\n", value.Subdomain))
		result.WriteString(fmt.Sprintf("<td>%s</td>\n", value.TestName))
		if value.TestResult == "PASSED" {
			result.WriteString("<td style=\"color: green;\">" + value.TestResult + "</td>\n")
		} else if value.TestResult == "FOUND" {
			result.WriteString("<td style=\"color: red;\">" + value.TestResult + "</td>\n")
		}
		result.WriteString("</tr>\n")
	}
	result.WriteString("</table>\n")
	output := `<h1 id="category-headline">Analysis Results</h1>
	<table id="overview-table">
<tr>
	<th>Subdomain</th>
    <th>Test</th>
    <th>Result</th>
</tr>` + result.String() + "</div>"
	rg.SummaryWriteToStream(output)
}
