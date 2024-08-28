package shared

// MARKS
const (
	Placeholder  = "HOST"
	NotAvailable = "n/a"
)

// STANDARD
const (
	Passive        = "PASSIVE"
	Active         = "ACTIVE"
	Dns            = "DNS"
	HeaderAnalysis = "HTTP-HEADER-ANALYSIS"
)

// EXTERNS
const (
	RDns = "RDNS"
	Ping = "PING"
)

// REQUEST
const (
	VersionUrl       = "https://raw.githubusercontent.com/fhAnso/Sentinel/main/version.txt"
	DefaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"
	TorProxyUrl      = "socks5://127.0.0.1:9050"
)

// OUTPUT
const (
	LogFileName     = "sentinel.log"
	VersionFile     = "version.txt"
	LoggerOutputDir = "log"
	OutputDir       = "output"
)

const DefaultPermission = 0755
