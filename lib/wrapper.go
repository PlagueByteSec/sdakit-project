package lib

import (
	"Sentinel/lib/utils"
	"fmt"
	"net"
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

func IpResolveWrapper(resolver *net.Resolver, args *utils.Args, params utils.Params) (string, []string) {
	utils.DnsLookups(resolver, utils.DnsLookupOptions{
		IpAddress: nil,
		Subdomain: params.Subdomain,
	})
	if utils.GDnsResults == nil {
		// Skip results that cannot be resolved to an IP address
		return "", nil
	}
	ipAddrsOut := strings.Join(utils.GDnsResults, ", ")
	return ipAddrsOut, utils.GDnsResults
}

func OpenOutputFileStreamsWrapper(filePaths *utils.FilePaths) {
	/*
		Specify the name and path for each output file. If all settings are configured, open
		separate file streams for each category (Subdomains, IPv4 addresses, and IPv6 addresses).
	*/
	if err := utils.GStreams.OpenOutputFileStreams(filePaths); err != nil {
		utils.Glogger.Println(err)
	}
}

func OutputWrapper(ipAddrs []string, params utils.Params, streams *utils.FileStreams) {
	for _, ip := range ipAddrs {
		utils.IpManage(params, ip, streams)
	}
	err := utils.WriteOutputFileStream(streams.SubdomainStream, params.FileContentSubdoms)
	if err != nil {
		streams.SubdomainStream.Close()
	}
}
