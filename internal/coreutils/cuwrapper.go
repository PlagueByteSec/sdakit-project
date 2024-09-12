package utils

import (
	"net"
	"strings"

	"github.com/PlagueByteSec/sentinel-project/v2/internal/requests"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
)

func PingWrapper(outputChan chan<- string, subdomain string, pingCount int) {
	outputChan <- " | Ping: "
	err := requests.PingSubdomain(subdomain, pingCount)
	if err != nil {
		outputChan <- "FAILED\n"
		return
	}
	outputChan <- "SUCCESS\n"
}

func PortScanWrapper(outputChan chan<- string, subdomain string, portRange string) {
	ports, err := requests.ScanPortRange(subdomain, portRange)
	if err != nil {
		shared.Glogger.Println(err)
	}
	if ports != "" {
		outputChan <- ports
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
