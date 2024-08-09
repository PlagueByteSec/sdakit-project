package lib

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var DisplayCount int

type Params struct {
	FilePath        string
	FilePathIPv4    string
	FilePathIPv6    string
	FileContent     string
	FileContentIPv4 string
	FileContentIPv6 string
	Result          string
	Hostname        string
}

func OutputHandler(client *http.Client, args *Args, params Params) {
	ipAddrs := RequestIpAddresses(params.Result)
	if args.SubOnlyIp && ipAddrs == nil {
		// Skip results that cannot be resolved to an IP address
		return
	}
	var ipAddrsOut string
	if ipAddrs != nil {
		ipAddrsOut = fmt.Sprintf("(%s)", strings.Join(ipAddrs, ", "))
	}
	consoleOutput := fmt.Sprintf(" ===[ %s %s", params.Result, ipAddrsOut)
	// Opening seperated output file streams
	streamDomains, err := OpenOutputFileStreamDomains(params)
	if err != nil {
		fmt.Println(err)
	}
	streamV4, err := OpenOutputFileStreamIPv4(params)
	if err != nil {
		fmt.Println(err)
	}
	streamV6, err := OpenOutputFileStreamIPv6(params)
	if err != nil {
		fmt.Println(err)
	}
	for _, ip := range ipAddrs {
		if GetIpVersion(ip) == 4 {
			params.FileContentIPv4 = ip
			if !PoolContainsEntry(IPv4Pool, params.FileContentIPv4) {
				IPv4Pool = append(IPv4Pool, params.FileContentIPv4)
				WriteOutputFileStreamIPv4(streamV4, params)
			}
		}
		if GetIpVersion(ip) == 6 {
			params.FileContentIPv6 = ip
			if !PoolContainsEntry(IPv6Pool, params.FileContentIPv6) {
				IPv6Pool = append(IPv6Pool, params.FileContentIPv6)
				WriteOutputFileStreamIPv6(streamV6, params)
			}
		}
	}
	WriteOutputFileStreamDomains(streamDomains, params)
	streamV4.Close()
	streamV6.Close()
	streamDomains.Close()
	codeFilter := strings.Split(args.FilHttpCodes, ",")
	codeFilterExc := strings.Split(args.ExcHttpCodes, ",")
	if args.HttpCode {
		url := fmt.Sprintf("http://%s", params.Result)
		httpStatusCode := HttpStatusCode(client, url)
		statusCodeConv := strconv.Itoa(httpStatusCode)
		if statusCodeConv == "-1" {
			statusCodeConv = Na
		}
		// Display only the given status codes
		if len(codeFilter) != 1 && !InArgList(statusCodeConv, codeFilter) {
			return
		}
		// Exclude the given codes from console output
		if len(codeFilterExc) != 1 && InArgList(statusCodeConv, codeFilterExc) {
			return
		}
		consoleOutput = fmt.Sprintf("%s, HTTP Status Code: %s", consoleOutput, statusCodeConv)
	}
	fmt.Println(consoleOutput)
	DisplayCount++
}
