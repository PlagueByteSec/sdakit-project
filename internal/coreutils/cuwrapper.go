package utils

import (
	"net"
	"strings"

	"github.com/PlagueByteSec/Sentinel/v1/internal/requests"
	"github.com/PlagueByteSec/Sentinel/v1/internal/shared"
)

func PingWrapper(consoleOutput *strings.Builder, subdomain string, pingCount int) {
	consoleOutput.WriteString(" | Ping: ")
	err := requests.PingSubdomain(subdomain, pingCount)
	if err != nil {
		consoleOutput.WriteString("FAILED\n")
		return
	}
	consoleOutput.WriteString("SUCCESS\n")
}

func PortScanWrapper(consoleOutput *strings.Builder, subdomain string, portRange string) {
	ports, err := requests.ScanPortsSubdomain(subdomain, portRange)
	if err != nil {
		shared.Glogger.Println(err)
	}
	if ports != "" {
		consoleOutput.WriteString(ports)
	}
}

func IpResolveWrapper(resolver *net.Resolver, subdomain string) (string, []string) {
	requests.DnsLookups(resolver, shared.DnsLookupOptions{
		IpAddress: nil,
		Subdomain: subdomain,
	})
	if shared.GDnsResults == nil {
		// Skip results that cannot be resolved to an IP address
		return "", nil
	}
	ipAddrsOut := strings.Join(shared.GDnsResults, ", ")
	return ipAddrsOut, shared.GDnsResults
}
