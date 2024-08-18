package lib

import (
	"errors"
	"flag"
	"fmt"
)

type Args struct {
	Verbose       bool
	Host          string
	OutFile       string
	OutFileIPv4   string
	OutFileIPv6   string
	NewOutputPath string
	HttpCode      bool
	WordlistPath  string
	ExcHttpCodes  string
	FilHttpCodes  string
	SubOnlyIp     bool
	AnalyzeHeader bool
	PortScan      string
	DbExtendPath  string
	Timeout       int
	TorRoute      bool
}

func CliParser() (Args, error) {
	verbose := flag.Bool("v", false, "Verbose output")
	host := flag.String("d", "", "Set the target domain name")
	outFile := flag.String("oS", "defaultSd", "Output file path for subdomains")
	outFileIPv4 := flag.String("o4", "defaultV4", "Output file path for IPv4 addresses")
	outFileIPv6 := flag.String("o6", "defaultV6", "Output file path for IPv6 addresses")
	newOutputPath := flag.String("nP", "defaultPath", "Output directory path for all results")
	httpCode := flag.Bool("c", false, "Get HTTP status code of each subdomain")
	wordlistPath := flag.String("w", "", "Specify wordlist and direct bruteforce subdomains")
	excHttpCodes := flag.String("e", "", "Exclude HTTP codes (comma seperated)")
	filtHttpCodes := flag.String("f", "", "Filter for specific HTTP response codes (comma seperated)")
	subOnlyIp := flag.Bool("s", false, "Display only subdomains which can be resolved to IP addresses")
	analyzeHeader := flag.Bool("a", false, "Analyze HTTP header of each subdomain")
	portScan := flag.String("p", "", "Define port range an run scan")
	dbExtendPath := flag.String("x", "", "Extend endpoint DB with custom list")
	timeout := flag.Int("t", 5, "Specify the request timeout")
	torRoute := flag.Bool("r", false, "Enable TOR routing")
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println(Help)
		return Args{}, errors.New("no args given, banner printed")
	}
	if *excHttpCodes != "" && !*httpCode || *filtHttpCodes != "" && !*httpCode {
		return Args{}, errors.New("HTTP code filter enabled, but status codes not requested")
	}
	args := Args{
		Verbose:       *verbose,
		Host:          *host,
		OutFile:       *outFile,
		OutFileIPv4:   *outFileIPv4,
		OutFileIPv6:   *outFileIPv6,
		NewOutputPath: *newOutputPath,
		HttpCode:      *httpCode,
		WordlistPath:  *wordlistPath,
		ExcHttpCodes:  *excHttpCodes,
		FilHttpCodes:  *filtHttpCodes,
		SubOnlyIp:     *subOnlyIp,
		AnalyzeHeader: *analyzeHeader,
		PortScan:      *portScan,
		DbExtendPath:  *dbExtendPath,
		Timeout:       *timeout,
		TorRoute:      *torRoute,
	}
	return args, nil
}
