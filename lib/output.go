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

func DoAnalyzeHeader(consoleOutput *strings.Builder, ipAddrsOut string, client *http.Client, params Params) {
	headers, count := AnalyseHttpHeader(client, params.Result)
	if ipAddrsOut != "" {
		if count != 0 {
			consoleOutput.WriteString(fmt.Sprintf("\n\t╠═[ %s", ipAddrsOut))
		} else {
			consoleOutput.WriteString(fmt.Sprintf("\n\t╚═[ %s\n", ipAddrsOut))
		}
	}
	if headers != "" && count != 0 {
		consoleOutput.WriteString(fmt.Sprintf("\n\t%s\n", headers))
	}
}

func DoPortScan(consoleOutput *strings.Builder, params Params, args *Args) {
	ports, err := ScanPortsSubdomain(params.Result, args.PortScan)
	if err != nil {
		Logger.Println(err)
	}
	if ports != "" {
		consoleOutput.WriteString(ports)
	}
}

func DoIpResolve(args *Args, params Params) (string, []string) {
	ipAddrs := RequestIpAddresses(params.Result)
	if args.SubOnlyIp && ipAddrs == nil {
		// Skip results that cannot be resolved to an IP address
		return "", nil
	}
	var ipAddrsOut string
	if ipAddrs != nil {
		ipAddrsOut = strings.Join(ipAddrs, ", ")
	}
	return ipAddrsOut, ipAddrs
}

func IpManage(params Params, ip string, fileStream *FileStreams) {
	ipAddrVersion := GetIpVersion(ip)
	switch ipAddrVersion {
	case 4:
		params.FileContentIPv4 = ip
		if !PoolContainsEntry(IPv4Pool, params.FileContentIPv4) {
			IPv4Pool = append(IPv4Pool, params.FileContentIPv4)
			err := WriteOutputFileStream(fileStream.Ipv4AddrStream, params.FileContentIPv4)
			if err != nil {
				Logger.Println(err)
			}
		}
	case 6:
		params.FileContentIPv6 = ip
		if !PoolContainsEntry(IPv6Pool, params.FileContentIPv6) {
			IPv6Pool = append(IPv6Pool, params.FileContentIPv6)
			err := WriteOutputFileStream(fileStream.Ipv6AddrStream, params.FileContentIPv6)
			if err != nil {
				Logger.Println(err)
			}
		}
	}
}

func OutputHandler(client *http.Client, args *Args, params Params) {
	ipAddrsOut, ipAddrs := DoIpResolve(args, params)
	if ipAddrs == nil {
		return
	}
	outputFileStreams, err := OpenOutputFileStreams(params)
	if err != nil {
		Logger.Println(err)
	}
	defer outputFileStreams.Ipv4AddrStream.Close()
	defer outputFileStreams.Ipv6AddrStream.Close()
	defer outputFileStreams.SubdomainStream.Close()
	for _, ip := range ipAddrs {
		IpManage(params, ip, outputFileStreams)
	}
	var consoleOutput strings.Builder
	WriteOutputFileStream(outputFileStreams.SubdomainStream, params.FileContent)
	consoleOutput.WriteString(fmt.Sprintf(" ══[ %s", params.Result))
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
		consoleOutput.WriteString(fmt.Sprintf(", HTTP Status Code: %s", statusCodeConv))
	}
	if args.AnalyzeHeader {
		DoAnalyzeHeader(&consoleOutput, ipAddrsOut, client, params)
	} else {
		if ipAddrsOut != "" {
			consoleOutput.WriteString(fmt.Sprintf("\n\t╚═[ %s\n", ipAddrsOut))
		}
	}
	if args.PortScan != "" {
		DoPortScan(&consoleOutput, params, args)
	}
	fmt.Println(consoleOutput.String())
	DisplayCount++
}
