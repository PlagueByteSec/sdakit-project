package lib

import (
	"bufio"
	"os"
)

var (
	GPool    Pool
	GStdout  *bufio.Writer
	GVerbose bool
)

func init() {
	GPool.PoolInit()
	GStdout = bufio.NewWriter(os.Stdout)
}
