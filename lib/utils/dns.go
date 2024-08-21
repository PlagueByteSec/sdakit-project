package utils

import (
	"context"
	"net"
)

func DnsResolverInit(useCustomDnsServer bool) *net.Resolver {
	var resolver *net.Resolver
	switch useCustomDnsServer {
	case true:
		resolver = &net.Resolver{
			PreferGo: false,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial("udp", CustomDnsServer)
			},
		}
	case false:
		resolver = &net.Resolver{}
	}
	return resolver
}

func DnsLookups(resolver *net.Resolver, dnsLookupOptions DnsLookupOptions) {
	var (
		dnsLookup []string
		temp      []net.IPAddr
		err       error
	)
	if dnsLookupOptions.IpAddress != nil {
		/*
			Perform a DNS lookup for the current subdomain to get the corresponding
			IP addresses and filter out old and inactive subdomains.
		*/
		GDnsResults, err = resolver.LookupAddr(context.Background(), dnsLookupOptions.IpAddress.String())
		if err != nil {
			return
		}
	} else if dnsLookupOptions.Subdomain != "" {
		/*
			Perform a RDNS lookup for the current IP address to get
			the corresponding domain name.
		*/
		temp, err = resolver.LookupIPAddr(context.Background(), dnsLookupOptions.Subdomain)
		if err != nil {
			return
		}
		// Convert []net.IPAddr to []string
		for idx := 0; idx < len(temp); idx++ {
			GDnsResults = append(dnsLookup, temp[idx].String())
		}
	}
}
