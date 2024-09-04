package cli

var HelpBanner = ` 													
===[ The Sentinel Project, By PlagueByteSec
        
    Contact: plaguebyte.sec@keemail.me
    License: MIT (https://github.com/PlagueByteSec/Sentinel/blob/main/LICENSE)
																		 
 ===[ Description

	This program is designed to fetch, filter and validate subdomains from a target domain.

 ===[ Options ]===

 -v     Verbose output
 -d		Set the target domain name
		[passive] (Without -w) request subdomains from external endpoints
		Example: targetdomain.xyz 
 -s		Set the target subdomain 
		Example: sub.targetdomain.xyz 
 -w     [active] Use wordlist to bruteforce subdomains of the target
 -dns   Use wordlist (-w) and resolve subdomains by querying a DNS

 ===[ ANALYSIS

 -c     Send HTTP request to retrieve the status code 
 -a     Analyze HTTP headers of each subdomain (server etc.)
 -aH	Display all headers of HTTP response from -a and -aS
 -aS	Analyse the HTTP response from a subdomain (specified by -s)
 -p     Scan subdomains for open ports in range
 -pS	Ping subdomains (privileged execution required)
 -pC	Ping subdomains from file (privileged execution required)
 -rF	Read IP addresses from file and perform RDNS lookup
 -dP	Analyse subdomain to determine its purpose (mail, API etc.)
 -mT	Test subdomain for common weaknesses (CORS, header injection etc.)

 ===[ FILTERS

 -e     Exclude HTTP response codes (comma seperated)
 -f     Filter for specific HTTP response codes (comma seperated)
 
 ===[ SETTINGS

 -x     Extend endpoint DB with custom list (.txt)
 -r     Route all requests through TOR
 -m		Set HTTP request method (default: GET)
 -rD    Set HTTP request delay in ms
 -t     Specify the request timeout
 -dnsT  Set the timeout for DNS queries
 -dnsC  Specify a custom DNS server address (ip:port)
 -pC	Specify Ping count (default=2)
- aR	Allow redirects

 ===[ OUTPUT

 -dO    Disable all output file streams
 -nP    Specify the output directory path for all output files
 -oS    Specify the output (.txt) file path for subdomains
 -o4    Specify the output (.txt) file path for IPv4 addresses
 -o6    Specify the output (.txt) file path for IPv6 addresses
 -oJ    Specify the output (.json) file path for summary
`
