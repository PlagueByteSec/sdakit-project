package lib

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	Host        string
	OutFile     string
	HttpCode    bool
	PingResults bool
}

func CliParser() Args {
	host := flag.String("t", "", "Target host")
	outFile := flag.String("o", "default", "Output file")
	httpCode := flag.Bool("c", false, "Get HTTP status code of each entry")
	pingResults := flag.Bool("p", false, "Send ICMP packet to each entry")
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println(Help)
		os.Exit(-1)
	}
	args := Args{
		Host:        *host,
		OutFile:     *outFile,
		HttpCode:    *httpCode,
		PingResults: *pingResults,
	}
	return args
}
