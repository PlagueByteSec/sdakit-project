package lib

import (
	"Sentinel/lib/utils"
	"bufio"
	"os"
)

func init() {
	/*
		Project initialization: make every pool at startup and open
		a stream writer to stdout.
	*/
	utils.GPool.PoolInit()
	utils.GStdout = bufio.NewWriter(os.Stdout)
}
