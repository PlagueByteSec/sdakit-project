package utils

import (
	"github.com/fhAnso/Sentinel/v1/internal/shared"
)

func IsPassiveEnumeration(args *shared.Args) bool {
	return args.WordlistPath == "" && args.RDnsLookupFilePath == "" && args.Domain != ""
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
