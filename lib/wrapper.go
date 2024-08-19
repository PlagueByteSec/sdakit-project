package lib

import (
	"Sentinel/lib/utils"
	"fmt"
	"net/http"
	"strings"
)

func AnalyzeHeaderWrapper(consoleOutput *strings.Builder, ipAddrsOut string,
	client *http.Client, params utils.Params) {
	/*
		Analyze the HTTP header and add the results to the consoleOutput
		string builder if it exists.
	*/
	headers, count := AnalyseHttpHeader(client, params.Subdomain)
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

func PortScanWrapper(consoleOutput *strings.Builder, params utils.Params, args *utils.Args) {
	ports, err := ScanPortsSubdomain(params.Subdomain, args.PortScan)
	if err != nil {
		utils.Glogger.Println(err)
	}
	if ports != "" {
		consoleOutput.WriteString(ports)
	}
}

func IpResolveWrapper(args *utils.Args, params utils.Params) (string, []string) {
	ipAddrs := utils.RequestIpAddresses(params.Subdomain) // DNS lookup
	if ipAddrs == nil {
		// Skip results that cannot be resolved to an IP address
		return "", nil
	}
	ipAddrsOut := strings.Join(ipAddrs, ", ")
	return ipAddrsOut, ipAddrs
}
