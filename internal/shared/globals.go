package shared

import (
	"bufio"
	"net"
	"time"
)

var (
	GStdout *bufio.Writer
)

var (
	GPoolBase    PoolBase
	GSubdomBase  SubdomainBase
	GSubdomAddrs SubdomainIpAddresses
	GJsonResult  JsonResult
)

var (
	GDnsResolver *net.Resolver
	GDnsResults  []string
)

var GStreams FileStreams

var (
	GDisplayCount     int
	GShowAllHeaders   bool
	GVerbose          bool
	GDisableAllOutput bool
	GObtainedCounter  int       = 0
	GAllCounter       int       = 0
	GStartTime        time.Time = time.Now()
)

var (
	CustomDnsServer string // IP address
	RequestDelay    int
)

var GIsExec int

// SUMMARY
var (
	GOpenPortsCount int
	GApiCount       int
	GMxCount        int
	GLoginCount     int
	GCorsCount      int
)

var (
	GTargetDomain string
	GScanMethod   string
)

var GReportPool = make(map[string]SetTestResults)

var (
	GCurrentIPv4Filename string
	GCurrentIPv6Filename string
)
