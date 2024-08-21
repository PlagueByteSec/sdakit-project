package utils

import (
	"bufio"
	"log"
	"net"
	"time"
)

var (
	Glogger          *log.Logger
	GStdout          *bufio.Writer
	GPool            PoolBase
	GSubdomBase      SubdomainBase
	GSubdomAddrs     SubdomainIpAddresses
	GJsonResult      JsonResult
	GDisplayCount    int
	GVerbose         bool
	GObtainedCounter int = 0
	GAllCounter      int = 0
	GStreams         FileStreams
	GStartTime       time.Time = time.Now()
	GDnsResolver     *net.Resolver
	GDnsResults      []string
)

var (
	CustomDnsServer string // IP address
	RequestDelay    int
)
