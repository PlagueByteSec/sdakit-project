package utils

import (
	"net"
	"strings"
)

func PortScanWrapper(consoleOutput *strings.Builder, params Params, args *Args) {
	ports, err := ScanPortsSubdomain(params.Subdomain, args.PortScan)
	if err != nil {
		Glogger.Println(err)
	}
	if ports != "" {
		consoleOutput.WriteString(ports)
	}
}

func IpResolveWrapper(resolver *net.Resolver, args *Args, params Params) (string, []string) {
	DnsLookups(resolver, DnsLookupOptions{
		IpAddress: nil,
		Subdomain: params.Subdomain,
	})
	if GDnsResults == nil {
		// Skip results that cannot be resolved to an IP address
		return "", nil
	}
	ipAddrsOut := strings.Join(GDnsResults, ", ")
	return ipAddrsOut, GDnsResults
}

func OpenOutputFileStreamsWrapper(filePaths *FilePaths) {
	/*
		Specify the name and path for each output file. If all settings are configured, open
		separate file streams for each category (Subdomains, IPv4 addresses, and IPv6 addresses).
	*/
	if err := GStreams.OpenOutputFileStreams(filePaths); err != nil {
		Glogger.Println(err)
	}
}

func OutputWrapper(ipAddrs []string, params Params, streams *FileStreams) {
	for _, ip := range ipAddrs {
		IpManage(params, ip, streams)
	}
	if !GDisableAllOutput {
		err := WriteOutputFileStream(streams.SubdomainStream, params.FileContentSubdoms)
		if err != nil {
			streams.SubdomainStream.Close()
		}
	}
}
