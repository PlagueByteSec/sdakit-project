package utils

import (
	"net"
)

// CORE
type Params struct {
	Domain               string
	Subdomain            string
	FilePathSubdomains   string
	FileContentSubdoms   string
	FilePathIPv4Addrs    string
	FilePathIPv6Addrs    string
	FileContentIPv4Addrs string
	FileContentIPv6Addrs string
}

type Args struct {
	Verbose          bool
	Domain           string
	OutFileSubdoms   string
	OutFileIPv4      string
	OutFileIPv6      string
	OutFileJSON      string
	NewOutputDirPath string
	HttpCode         bool
	WordlistPath     string
	ExcHttpCodes     string
	FilHttpCodes     string
	SubOnlyIp        bool
	AnalyzeHeader    bool
	PortScan         string
	DbExtendPath     string
	Timeout          int
	TorRoute         bool
	DnsLookup        bool
	DnsLookupCustom  string
	DnsLookupTimeout int
	HttpRequestDelay int
}

type PoolBase struct {
	PoolIPv4Addresses []string
	PoolIPv6Addresses []string
	PoolSubdomains    []string
}

// OUTPUT
type FilePaths struct {
	FilePathSubdomain string
	FilePathIPv4      string
	FilePathIPv6      string
	FilePathJSON      string
}

type SubdomainIpAddresses struct {
	IPv4 []net.IP `json:"ipv4Addresses"`
	IPv6 []net.IP `json:"ipv6Addresses"`
}

type SubdomainBase struct {
	Subdomain   []string `json:"subdomain"`
	OpenPorts   []int    `json:"openPorts"`
	IpAddresses SubdomainIpAddresses
}

type JsonResult struct {
	Subdomains []SubdomainBase `json:"subdomains"`
}

type FileExtension int

// EXIT
type SentinelExitParams struct {
	ExitCode    int
	ExitMessage string
	ExitError   error
}
