package lib

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	Verbose      bool
	Host         string
	OutFile      string
	OutFileIPv4  string
	OutFileIPv6  string
	HttpCode     bool
	WordlistPath string
	ExcHttpCodes string
	FilHttpCodes string
	SubOnlyIp    bool
}

func CliParser() Args {
	verbose := flag.Bool("v", false, "Verbose output")
	host := flag.String("t", "", "Target host")
	outFile := flag.String("o", "default", "Output file")
	outFileIPv4 := flag.String("v4", "defaultV4", "Output file for IPv4 addresses")
	outFileIPv6 := flag.String("v6", "defaultV6", "Output file for IPv6 addresses")
	httpCode := flag.Bool("c", false, "Get HTTP status code of each entry")
	wordlistPath := flag.String("w", "", "Specify wordlist and direct bruteforce subdomain")
	excHttpCodes := flag.String("e", "", "Exclude HTTP codes (comma seperated)")
	filtHttpCodes := flag.String("f", "", "Show only specific HTTP codes (comma seperated)")
	subOnlyIp := flag.Bool("s", false, "Display only specific subdomains")
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println(Help)
		os.Exit(-1)
	}
	args := Args{
		Verbose:      *verbose,
		Host:         *host,
		OutFile:      *outFile,
		OutFileIPv4:  *outFileIPv4,
		OutFileIPv6:  *outFileIPv6,
		HttpCode:     *httpCode,
		WordlistPath: *wordlistPath,
		ExcHttpCodes: *excHttpCodes,
		FilHttpCodes: *filtHttpCodes,
		SubOnlyIp:    *subOnlyIp,
	}
	return args
}
