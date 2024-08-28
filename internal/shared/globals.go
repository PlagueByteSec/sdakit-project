package shared

import (
	"bufio"
	"log"
	"net"
	"time"
)

var (
	Glogger *log.Logger
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
