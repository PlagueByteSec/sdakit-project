package utils

import (
	"github.com/PlagueByteSec/sdakit-project/v2/internal/shared"
)

func IsPassiveEnumeration(args *shared.Args) bool {
	return args.Domain != "" && !args.DnsLookup && !args.EnableVHostEnum && args.IpAddress == ""
}

func IsActiveEnumeration(args *shared.Args) bool {
	return args.WordlistPath != "" && !args.DnsLookup && args.IpAddress == ""
}

func IsDnsEnumeration(args *shared.Args) bool {
	return args.DnsLookup && args.Domain != ""
}

func IsVHostEnumeration(args *shared.Args) bool {
	return args.EnableVHostEnum && args.Domain != "" && args.IpAddress != ""
}

func IsHttpHeaderAnalysis(args *shared.Args) bool {
	return args.AnalyseHeaderSingle && args.Domain == "" && args.WordlistPath == ""
}
