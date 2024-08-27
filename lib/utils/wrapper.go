package utils

import (
	"Sentinel/lib/requests"
	"Sentinel/lib/shared"
	"net"
	"strings"
)

func PortScanWrapper(consoleOutput *strings.Builder, params shared.Params, args *shared.Args) {
	ports, err := requests.ScanPortsSubdomain(params.Subdomain, args.PortScan)
	if err != nil {
		shared.Glogger.Println(err)
	}
	if ports != "" {
		consoleOutput.WriteString(ports)
	}
}

func IpResolveWrapper(resolver *net.Resolver, args *shared.Args, params shared.Params) (string, []string) {
	requests.DnsLookups(resolver, shared.DnsLookupOptions{
		IpAddress: nil,
		Subdomain: params.Subdomain,
	})
	if shared.GDnsResults == nil {
		// Skip results that cannot be resolved to an IP address
		return "", nil
	}
	ipAddrsOut := strings.Join(shared.GDnsResults, ", ")
	return ipAddrsOut, shared.GDnsResults
}
