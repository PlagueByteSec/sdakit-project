package cli

var HelpBanner = `
Sentinel
********
 ===[ Description
	This program is designed to fetch, filter and validate subdomains from a specific host.
	Subdomains can be enumerated by requesting external resources or via brute 
	force. 

	Sentinel is published under the MIT license: 
		https://github.com/fhAnso/Sentinel/blob/main/LICENSE

 ===[ Options ]===

 -v     Verbose output
 -d     Set the target domain name
	    [passive] (Without -w) request subdomains from external endpoints
	    Example: targetdomain.xyz 
 -w     [active] Use wordlist to bruteforce subdomains of the target
 -dns   Use wordlist (-w) and resolve subdomains by querying a DNS

 ===[ ANALYSIS

 -c     Send GET request to retrieve the HTTP status 
 -a     Analyze HTTP headers of each subdomain (server etc.)
 -p     Scan subdomains for open ports in range

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

 ===[ OUTPUT

 -dO    Disable all output file streams
 -oS    Specify the output (.txt) file path for subdomains
 -o4    Specify the output (.txt) file path for IPv4 addresses
 -o6    Specify the output (.txt) file path for IPv6 addresses
 -nP    Specify the output directory path for all output files
 -oJ    Specify the output (.json) file path for summary
`
