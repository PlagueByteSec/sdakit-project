package utils

import (
	"github.com/PlagueByteSec/Sentinel/v2/internal/shared"
)

func IsPassiveEnumeration(args *shared.Args) bool {
	return args.WordlistPath == "" && args.RDnsLookupFilePath == "" && args.Domain != "" && !args.AnalyseHeaderSingle
}

func IsActiveEnumeration(args *shared.Args) bool {
	return args.WordlistPath != "" && !args.DnsLookup && args.RDnsLookupFilePath == ""
}

func IsDnsEnumeration(args *shared.Args) bool {
	return args.DnsLookup && args.WordlistPath != "" && args.RDnsLookupFilePath == "" && args.Domain != ""
}

func IsRDnsEnumeration(args *shared.Args) bool {
	return args.RDnsLookupFilePath != "" && args.WordlistPath == "" && args.Domain == ""
}

func IsPingFromFile(args *shared.Args) bool {
	return args.PingSubdomainsFile != "" && args.Domain == "" && args.WordlistPath == ""
}

func IsHttpHeaderAnalysis(args *shared.Args) bool {
	return args.AnalyseHeaderSingle && args.Subdomain != "" && args.Domain == "" && args.WordlistPath == "" && args.RDnsLookupFilePath == ""
}
