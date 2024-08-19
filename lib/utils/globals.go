package utils

import (
	"bufio"
	"log"
)

var (
	Glogger       *log.Logger
	GStdout       *bufio.Writer
	GPool         PoolBase
	GSubdomBase   SubdomainBase
	GSubdomAddrs  SubdomainIpAddresses
	GJsonResult   JsonResult
	GDisplayCount int
	GVerbose      bool
)
