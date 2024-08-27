package utils

import (
	"Sentinel/lib/shared"
	"net"
	"os"
	"os/signal"
	"regexp"
)

func InArgList(httpCode string, list []string) bool {
	/*
		An HTTP status code will be compared against those
		specified by the -e or -f flag. This function is used
		to filter the output for customization.
	*/
	for _, code := range list {
		if httpCode == code {
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

func InterruptListenerInit() {
	/*
		Create a channel to receive interrupt signals from the OS.
		The goroutine continuously listens for an interrupt signal
		(Ctrl+C) and handles the interruption.
	*/
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for range sigChan {
			SentinelExit(shared.SentinelExitParams{
				ExitCode:    0,
				ExitMessage: "\n\nG0oDBy3!",
				ExitError:   nil,
			})
		}
	}()
}

func ResetSlice(slice *[]string) {
	if len(*slice) >= 1 && (*slice)[0] == "" {
		*slice = []string{}
	}
}
