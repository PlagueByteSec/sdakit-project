package pkg

import (
	"net"
	"regexp"
)

func IsInSlice(value string, slice []string) bool {
	/*
		An HTTP status code will be compared against those
		specified by the -e or -f flag. This function is used
		to filter the output for customization.
	*/
	for _, entry := range slice {
		if value == entry {
			return true
		}
	}
	return false
}

func GetIpVersion(ipAddress string) int {
	var ipVersion int
	if ip := net.ParseIP(ipAddress); ip != nil {
		if ip.To4() != nil {
			ipVersion = 4
		} else {
			ipVersion = 6
		}
	}
	return ipVersion
}

func IsValidDomain(domain string) bool {
	/*
		Verify the target domain by checking it against a regex and
		performing a DNS lookup. The domain will be considered invalid
		if no IP addresses can be enumerated.
	*/
	regex := `^(?i)([a-z0-9](-?[a-z0-9])*)+(\.[a-z]{2,})+$`
	regexCompile := regexp.MustCompile(regex)
	if regexCompile.MatchString(domain) {
		ipAddrs, err := net.LookupIP(domain)
		if err != nil {
			return false
		}
		if len(ipAddrs) != 0 {
			return true
		}
	}
	return false
}

func ResetSlice(slice *[]string) {
	if len(*slice) >= 1 && (*slice)[0] == "" {
		*slice = []string{}
	}
}
