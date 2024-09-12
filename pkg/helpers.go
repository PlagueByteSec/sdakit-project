package pkg

import (
	"fmt"
	"net"
	"regexp"
	"sort"
	"time"
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

func IsInSlice(value interface{}, slice interface{}) bool {
	switch sliceType := slice.(type) {
	case []string:
		if checkValue, ok := value.(string); ok {
			sort.Strings(sliceType)
			idx := sort.SearchStrings(sliceType, checkValue)
			return idx < len(sliceType) && sliceType[idx] == checkValue
		}
	case []int:
		if checkValue, ok := value.(int); ok {
			sort.Ints(sliceType)
			idx := sort.SearchInts(sliceType, checkValue)
			return idx < len(sliceType) && sliceType[idx] == checkValue
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

func PrintDots(subdomain string, dotChan <-chan struct{}) {
	for {
		select {
		case <-dotChan:
			return
		default:
			fmt.Print(".")
			time.Sleep(800 * time.Millisecond)
		}
	}
}
