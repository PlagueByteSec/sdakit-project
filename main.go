package main

import (
	"Sentinel/lib"
	"bufio"
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

// Read database file (URL.txt) and process the entries
func Manager(host string) {
	filePath := lib.DefaultPath()
	if lib.FileExist(filePath) {
		pool := make(lib.Pool)
		stream, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("ERROR: failed to open file: %s\n%s\n", filePath, err)
			os.Exit(-1)
		}
		defer stream.Close()
		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "#") {
				url := strings.Replace(line, "HOST", host, 1)
				lib.Request(pool, host, url)
			}
		}
	} else {
		fmt.Printf("ERROR: file not found in default path: %s\n", filePath)
		os.Exit(-1)
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
