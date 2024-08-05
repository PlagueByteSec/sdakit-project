package lib

import (
	"fmt"
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

func OutputHandler(args *Args, params Params) {
	ips := RequestIpAddresses(params.Result)
	if args.SubOnlyIp && ips == "" {
		return
	}
	consoleOutput := fmt.Sprintf(" ===[ %s %s", params.Result, ips)
	ips = strings.TrimPrefix(ips, "(")
	ips = strings.TrimSuffix(ips, ")")
	ipAddrs := strings.Split(ips, ", ")
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
			WriteOutputFileStreamIPv4(streamV4, params)
		}
		if GetIpVersion(ip) == 6 {
			params.FileContentIPv6 = ip
			WriteOutputFileStreamIPv6(streamV6, params)
		}
	}
	WriteOutputFileStreamDomains(streamDomains, params)
	streamV4.Close()
	streamV6.Close()
	streamDomains.Close()
	if args.HttpCode {
		url := fmt.Sprintf("http://%s", params.Result)
		httpStatusCode := fmt.Sprintf("%d", HttpStatusCode(url))
		if httpStatusCode == "-1" {
			httpStatusCode = na
		}
		consoleOutput = fmt.Sprintf("%s, HTTP Status Code: %s", consoleOutput, httpStatusCode)
	}
	fmt.Println(consoleOutput)
	DisplayCount++
}
