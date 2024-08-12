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
	WriteOutputFileStream(outputFileStreams.SubdomainStream, params.FileContent)
	fmt.Printf("\n ══[ %s", params.Result)
	codeFilter := strings.Split(args.FilHttpCodes, ",")
	codeFilterExc := strings.Split(args.ExcHttpCodes, ",")
	if args.HttpCode {
		url := fmt.Sprintf("http://%s", params.Result)
		httpStatusCode := HttpStatusCode(client, url)
		statusCodeConv := strconv.Itoa(httpStatusCode)
		if httpStatusCode == -1 {
			statusCodeConv = NotAvailable
		}
		if len(codeFilter) != 1 && !InArgList(statusCodeConv, codeFilter) {
			return
		}
		if len(codeFilterExc) != 1 && InArgList(statusCodeConv, codeFilterExc) {
			return
		}
		fmt.Printf(", HTTP Status Code: %s", statusCodeConv)
	}
	if args.AnalyzeHeader {
		headers, count := AnalyseHttpHeader(client, params.Result)
		if ipAddrsOut != "" {
			if count != 0 {
				fmt.Printf("\n\t╠═[ %s", ipAddrsOut)
			} else {
				fmt.Printf("\n\t╚═[ %s\n", ipAddrsOut)
			}
		}
		if headers != "" && count != 0 {
			fmt.Printf("\n\t%s", headers)
		}
	} else {
		if ipAddrsOut != "" {
			fmt.Printf("\n\t╚═[ %s\n", ipAddrsOut)
		}
	}
	if args.PortScan != "" {
		fmt.Println()
		ports, err := ScanPortsSubdomain(params.Result, args.PortScan)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(ports)
	}
	DisplayCount++
}
