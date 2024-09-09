package requests

import (
	"context"
	"net"

	"github.com/PlagueByteSec/Sentinel/v1/internal/shared"
)

func DnsResolverInit(useCustomDnsServer bool) *net.Resolver {
	var resolver *net.Resolver
	switch useCustomDnsServer {
	case true:
		resolver = &net.Resolver{
			PreferGo: false,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial("udp", shared.CustomDnsServer)
			},
		}
	case false:
		resolver = &net.Resolver{}
	}
	return resolver
}

func SetDnsEnumType() {
	shared.GDnsResolver = DnsResolverInit(false)
	if shared.CustomDnsServer != "" {
		// Use custom DNS server address
		shared.GDnsResolver = DnsResolverInit(true)
	}
}

func DnsIsMX(resolver *net.Resolver, subdomain string) bool {
	mxRecords, err := resolver.LookupMX(context.Background(), subdomain)
	return err == nil && len(mxRecords) > 0
}

func DnsLookups(resolver *net.Resolver, dnsLookupOptions shared.DnsLookupOptions) {
	var (
		dnsLookup []string
		temp      []net.IPAddr
		err       error
	)
	if dnsLookupOptions.IpAddress != nil {
		/*
			Perform a RDNS lookup for the current IP address to get
			the corresponding domain name.
		*/
		shared.GDnsResults, err = resolver.LookupAddr(context.Background(), dnsLookupOptions.IpAddress.String())
		if err != nil {
			return
		}
	} else if dnsLookupOptions.Subdomain != "" {
		/*
			Perform a DNS lookup for the current subdomain to get the corresponding
			IP addresses and filter out old and inactive subdomains.
		*/
		temp, err = resolver.LookupIPAddr(context.Background(), dnsLookupOptions.Subdomain)
		if err != nil {
			return
		}
		// Convert []net.IPAddr to []string
		for idx := 0; idx < len(temp); idx++ {
			shared.GDnsResults = append(dnsLookup, temp[idx].String())
		}
	}
}
