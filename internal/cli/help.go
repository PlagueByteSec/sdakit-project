package cli

var HelpBanner = ` 													
===[ The Sentinel Project	

	Author: fhAnso 														
	Contact: plaguebyte.sec@keemail.me									
	License: MIT (https://github.com/fhAnso/Sentinel/blob/main/LICENSE)	
																		 
 ===[ Description

	This program is designed to fetch, filter and validate subdomains from a target domain.

 ===[ Overview

	- Use external resources for passive enum 
	- Use wordlist for active enum (direct or DNS)
	- Analyse HTTP response headers
	- Filter results by HTTP response code
	- Perform RDNS lookup from IP list
	- Ping subdomains directly or from file
	- Scan port range on detected subdomain
	- Route all traffic through TOR
	- Automated output file generation

 ===[ Options ]===

 -v     Verbose output
 -d     Set the target domain name
	    [passive] (Without -w) request subdomains from external endpoints
	    Example: targetdomain.xyz 
 -s		Set the target subdomain 
		Example: sub.targetdomain.xyz 
 -w     [active] Use wordlist to bruteforce subdomains of the target
 -dns   Use wordlist (-w) and resolve subdomains by querying a DNS

 ===[ ANALYSIS

 -c     Send GET request to retrieve the HTTP status 
 -a     Analyze HTTP headers of each subdomain (server etc.)
 -p     Scan subdomains for open ports in range
 -pS	Ping subdomains (privileged execution required)
 -pC	Ping subdomains from file (privileged execution required)
 -rF	Read IP addresses from file and perform RDNS lookup
 -aS	Analyse the HTTP response from a subdomain (specified by -s)

 ===[ FILTERS

 -e     Exclude HTTP response codes (comma seperated)
 -f     Filter for specific HTTP response codes (comma seperated)
 
 ===[ SETTINGS

 -x     Extend endpoint DB with custom list (.txt)
 -t     Specify the request timeout
 -r     Route all requests through TOR
 -dnsC  Specify a custom DNS server address (ip:port)
 -dnsT  Set the timeout for DNS queries
 -rD    Set HTTP request delay in ms
 -pC	Specify Ping count (default=2)

 ===[ OUTPUT

 -dO    Disable all output file streams
 -oS    Specify the output (.txt) file path for subdomains
 -o4    Specify the output (.txt) file path for IPv4 addresses
 -o6    Specify the output (.txt) file path for IPv6 addresses
 -nP    Specify the output directory path for all output files
 -oJ    Specify the output (.json) file path for summary
`
