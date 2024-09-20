package utils

import (
	"github.com/PlagueByteSec/sdakit-project/v2/internal/shared"
)

func IsPassiveEnumeration(args *shared.Args) bool {
	return args.WordlistPath == "" && args.Domain != "" && !args.AnalyseHeaderSingle && args.IpAddress == ""
}

func IsActiveEnumeration(args *shared.Args) bool {
	return args.WordlistPath != "" && !args.DnsLookup && args.IpAddress == ""
}

func IsDnsEnumeration(args *shared.Args) bool {
	return args.DnsLookup && args.WordlistPath != "" && args.Domain != ""
}

func IsVHostEnumeration(args *shared.Args) bool {
	return args.EnableVHostEnum && args.IpAddress != "" && args.Domain != "" && args.Subdomain == "" && args.WordlistPath != ""
}

func IsHttpHeaderAnalysis(args *shared.Args) bool {
	return args.AnalyseHeaderSingle && args.Subdomain != "" && args.Domain == "" && args.WordlistPath == ""
}
