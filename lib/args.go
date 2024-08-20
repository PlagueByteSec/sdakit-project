package lib

import (
	"Sentinel/lib/utils"
	"errors"
	"flag"
	"fmt"
	"os"
)

func CliParser() (utils.Args, error) {
	verbose := flag.Bool("v", false, "Verbose output")
	domain := flag.String("d", "", "Set the target domain name")
	outFile := flag.String("oS", "defaultSd", "Output file path for subdomains")
	outFileIPv4 := flag.String("o4", "defaultV4", "Output file path for IPv4 addresses")
	outFileIPv6 := flag.String("o6", "defaultV6", "Output file path for IPv6 addresses")
	outFileJSON := flag.String("oJ", "defaultJSON", "Output file path for JSON summary")
	newOutputPath := flag.String("nP", "defaultPath", "Output directory path for all results")
	httpCode := flag.Bool("c", false, "Get HTTP status code of each subdomain")
	wordlistPath := flag.String("w", "", "Specify wordlist and direct bruteforce subdomains")
	excHttpCodes := flag.String("e", "", "Exclude HTTP codes (comma seperated)")
	filtHttpCodes := flag.String("f", "", "Filter for specific HTTP response codes (comma seperated)")
	analyzeHeader := flag.Bool("a", false, "Analyze HTTP header of each subdomain")
	portScan := flag.String("p", "", "Define port range an run scan")
	dbExtendPath := flag.String("x", "", "Extend endpoint DB with custom list")
	timeout := flag.Int("t", 5, "Specify the request timeout")
	torRoute := flag.Bool("r", false, "Enable TOR routing")
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println(Help + "\nPlease specify a domain!")
		os.Exit(0)
	}
	if *excHttpCodes != "" && !*httpCode || *filtHttpCodes != "" && !*httpCode {
		return utils.Args{}, errors.New("HTTP code filter enabled, but status codes not requested")
	}
	if !utils.IsValidDomain(*domain) {
		return utils.Args{}, errors.New("domain verification failed: " + *domain)
	}
	args := utils.Args{
		Verbose:          *verbose,
		Domain:           *domain,
		OutFileSubdoms:   *outFile,
		OutFileIPv4:      *outFileIPv4,
		OutFileIPv6:      *outFileIPv6,
		OutFileJSON:      *outFileJSON,
		NewOutputDirPath: *newOutputPath,
		HttpCode:         *httpCode,
		WordlistPath:     *wordlistPath,
		ExcHttpCodes:     *excHttpCodes,
		FilHttpCodes:     *filtHttpCodes,
		AnalyzeHeader:    *analyzeHeader,
		PortScan:         *portScan,
		DbExtendPath:     *dbExtendPath,
		Timeout:          *timeout,
		TorRoute:         *torRoute,
	}
	return args, nil
}
