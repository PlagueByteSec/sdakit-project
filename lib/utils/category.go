package utils

import (
	"Sentinel/lib/shared"
)

func IsPassiveEnumeration(args *shared.Args) bool {
	return args.WordlistPath == "" && args.RDnsLookupFilePath == ""
}

func IsActiveEnumeration(args *shared.Args) bool {
	return args.WordlistPath != "" && !args.DnsLookup && args.RDnsLookupFilePath == ""
}

func IsDnsEnumeration(args *shared.Args) bool {
	return args.DnsLookup && args.WordlistPath != "" && args.RDnsLookupFilePath == ""
}

func IsRDnsEnumeration(args *shared.Args) bool {
	return args.RDnsLookupFilePath != "" && args.WordlistPath == "" && args.Domain == ""
}
