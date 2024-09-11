package utils

import (
	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
)

func IsPassiveEnumeration(args *shared.Args) bool {
	return args.WordlistPath == "" && args.RDnsLookupFilePath == "" && args.Domain != "" &&
		!args.AnalyseHeaderSingle && args.IpAddress == ""
}

func IsActiveEnumeration(args *shared.Args) bool {
	return args.WordlistPath != "" && !args.DnsLookup && args.RDnsLookupFilePath == "" && args.IpAddress == ""
}

func IsDnsEnumeration(args *shared.Args) bool {
	return args.DnsLookup && args.WordlistPath != "" && args.RDnsLookupFilePath == "" && args.Domain != ""
}

func IsVHostEnumeration(args *shared.Args) bool {
	return args.EnableVHostEnum && args.IpAddress != "" && args.Domain != "" && args.Subdomain == "" &&
		args.WordlistPath != "" && args.RDnsLookupFilePath == ""
}

func IsRDnsEnumeration(args *shared.Args) bool {
	return args.RDnsLookupFilePath != "" && args.WordlistPath == "" && args.Domain == ""
}

func IsPingFromFile(args *shared.Args) bool {
	return args.PingSubdomainsFile != "" && args.Domain == "" && args.WordlistPath == ""
}

func IsHttpHeaderAnalysis(args *shared.Args) bool {
	return args.AnalyseHeaderSingle && args.Subdomain != "" && args.Domain == "" &&
		args.WordlistPath == "" && args.RDnsLookupFilePath == ""
}
