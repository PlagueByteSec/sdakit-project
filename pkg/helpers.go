package pkg

import (
	"net"
	"regexp"
	"sort"
)

func GetIpVersion(ipAddress string) int {
	parser := net.ParseIP(ipAddress)
	if parser == nil {
		return 0
	}
	if parser.To4() != nil {
		return 4
	}
	return 6
}

func IsInSlice(value string, slice []string) bool {
	/*
		An HTTP status code will be compared against those
		specified by the -e or -f flag. This function is used
		to filter the output for customization.
	*/
	sort.Strings(slice)
	idx := sort.SearchStrings(slice, value)
	return idx < len(slice) && slice[idx] == value
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

func Tern[T any](condition bool, value T, alt T) T {
	if condition {
		return value
	}
	return alt
}
