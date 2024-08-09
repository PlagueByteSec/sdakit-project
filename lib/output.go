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
		ipAddrsOut = strings.Join(ipAddrs, ", ")
	}
	outputFileStreams, err := OpenOutputFileStreams(params)
	if err != nil {
		fmt.Println(err)
	}
	defer outputFileStreams.Ipv4AddrStream.Close()
	defer outputFileStreams.Ipv6AddrStream.Close()
	defer outputFileStreams.SubdomainStream.Close()
	for _, ip := range ipAddrs {
		ipAddrVersion := GetIpVersion(ip)
		switch ipAddrVersion {
		case 4:
			params.FileContentIPv4 = ip
			if !PoolContainsEntry(IPv4Pool, params.FileContentIPv4) {
				IPv4Pool = append(IPv4Pool, params.FileContentIPv4)
				err := WriteOutputFileStream(outputFileStreams.Ipv4AddrStream, params.FileContentIPv4)
				if err != nil {
					fmt.Println(err)
				}
			}
		case 6:
			params.FileContentIPv6 = ip
			if !PoolContainsEntry(IPv6Pool, params.FileContentIPv6) {
				IPv6Pool = append(IPv6Pool, params.FileContentIPv6)
				err := WriteOutputFileStream(outputFileStreams.Ipv6AddrStream, params.FileContentIPv6)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	consoleOutput := fmt.Sprintf(" ══[ %s", params.Result)
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
	if args.AnalyzeHeader {
		headers, count := AnalyseHttpHeader(client, params.Result)
		if ipAddrsOut != "" {
			if count != 0 {
				consoleOutput = fmt.Sprintf("%s\n\t╠═[ %s", consoleOutput, ipAddrsOut)
			} else {
				consoleOutput = fmt.Sprintf("%s\n\t╚═[ %s", consoleOutput, ipAddrsOut)
			}
		}
		if headers != "" {
			consoleOutput = fmt.Sprintf("%s\n\t%s", consoleOutput, headers)
		}
	} else {
		if ipAddrsOut != "" {
			consoleOutput = fmt.Sprintf("%s\n\t╚═[ %s", consoleOutput, ipAddrsOut)
		}
	}
	fmt.Println(consoleOutput)
	DisplayCount++
}
