package shared

import (
	"net"
	"net/http"
	"os"
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
	Verbose            bool
	Domain             string // target domain
	OutFileSubdoms     string // custom subdomains output path
	OutFileIPv4        string // custom IPv4 output path
	OutFileIPv6        string // custom IPv6 output path
	OutFileJSON        string // custom JSON output path
	NewOutputDirPath   string // custom output dir path
	HttpCode           bool
	WordlistPath       string
	ExcHttpCodes       string // results to hide (specified by HTTP status code)
	FilHttpCodes       string // results to display (specified by HTTP status code)
	SubOnlyIp          bool
	AnalyzeHeader      bool
	PortScan           string // port range
	DbExtendPath       string // File path containing endpoints
	Timeout            int    // in seconds
	TorRoute           bool
	DnsLookup          bool
	DnsLookupCustom    string // Custom DNS server (args)
	DnsLookupTimeout   int
	HttpRequestDelay   int    // in milliseconds
	RDnsLookupFilePath string // IP address file path
	DisableAllOutput   bool
	PingSubdomain      bool
	PingCount          int
	PingSubdomainsFile string // read subdomains from file and ping
}

type PoolBase struct {
	PoolIPv4Addresses []string
	PoolIPv6Addresses []string
	PoolSubdomains    []string
}

type EnumerationMethod struct {
	MethodKey string
	Action    func(*Args, *http.Client, *FilePaths)
}

type ExternsMethod struct {
	MethodKey string
	Action    func(*Args)
}

// ENUM
type DnsLookupOptions struct {
	Subdomain string
	IpAddress net.IP
}

type HttpHeaders struct {
	Server string
	Hsts   string
	PowBy  string
	Csp    string
}

// OUTPUT
type FilePaths struct {
	FilePathSubdomain string
	FilePathIPv4      string
	FilePathIPv6      string
	FilePathJSON      string
}

type FileStreams struct {
	Ipv4AddrStream  *os.File
	Ipv6AddrStream  *os.File
	SubdomainStream *os.File
}

type ParamsSetupFilesBase struct {
	FileParams *Params
	CliArgs    *Args
	FilePaths  *FilePaths
	Subdomain  string
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

// EXIT
type SentinelExitParams struct {
	ExitCode    int
	ExitMessage string
	ExitError   error
}