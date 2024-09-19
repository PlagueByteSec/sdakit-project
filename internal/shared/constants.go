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
	VHost          = "VHOST"
	HeaderAnalysis = "HTTP-HEADER-ANALYSIS"
)

// EXTERNS
const (
	RDns = "RDNS"
	Ping = "PING"
)

// REQUEST
const (
	VersionUrl       = "https://raw.githubusercontent.com/PlagueByteSec/sdakit-project/main/version.txt"
	DefaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"
	TorProxyUrl      = "socks5://127.0.0.1:9050"
)

// OUTPUT
const (
	LogFileName     = "sdakit.log"
	LoggerOutputDir = "log"
	OutputDir       = "output"
	VersionFile     = "version.txt"
)
const DefaultPermission = 0755
