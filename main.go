package main

import (
	"Sentinel/lib"
	"flag"
	"fmt"
	"os"
	"strings"
)

var help = `
Sentinel Project 
****************
=== Description
 This program is designed to fetch, filter and validate subdomains from a specific host.
 The Sentinal project replaces the platform-dependent script "uma.sh" and makes it possible
 to passively enumerate subdomains of a target independently of the OS. The results will be 
 saved among each other to provide a quick solution for further processing.
=== Usage
 -t Set the target hostname. This host will be be send to "RapidDNS" and "crt.sh"
	Example: targethostname.xyz 
=== Not implemented yet
 -o Specify the output (.txt) file where every regex match from the main pool will be saved
 -c Send GET request to retrieve the HTTP status code of every entry
 -p Send ICMP packages to ping all entries
`

// Pool init and preparation
func Manager(host string) {
	pool := make(lib.Pool)
	for _, entry := range lib.Db {
		url := strings.Replace(entry, "HOST", host, 1)
		lib.Request(pool, host, url)
	}
	for result := range pool {
		fmt.Println(host, ": ", result)
	}
}

func main() {
	host := flag.String("t", "", "")
	//outFile := flag.String("o", "", "")
	//httpCode := flag.String("c", "", "")
	//pingResults := flag.String("p", "", "")
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println(help)
		os.Exit(-1)
	}

	if *host != "" {
		Manager(*host)
	}
}
