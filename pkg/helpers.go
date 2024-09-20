package pkg

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

func BuildBanner(text string) string {
	lines := strings.Split(text, "\n")
	// Determine the maximum line length
	maxLineLength := 0
	for idx := 0; idx < len(lines); idx++ {
		line := lines[idx]
		if len(line) > maxLineLength {
			maxLineLength = len(line)
		}
	}
	frameLength := maxLineLength + 10
	frameLine := strings.Repeat("* ", frameLength/2)
	var result []string
	for idx := 0; idx < len(lines); idx++ {
		line := lines[idx]
		framedLine := fmt.Sprintf("*   %s   *", line)
		// Pad the line with if it's shorter than the max length
		if len(line) < maxLineLength {
			framedLine = fmt.Sprintf("*   %s%s   *", line, strings.Repeat(" ", maxLineLength-len(line)))
		}
		result = append(result, framedLine)
	}
	return fmt.Sprintf("%s\n%s\n%s", frameLine, strings.Join(result, "\n"), frameLine)
}

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
