package lib

var Help = `
Sentinel
********
===[ Description
 This program is designed to fetch, filter and validate subdomains from a specific host.
 Subdomains can be enumerated by requesting external resources or via brute 
 force. 

 Sentinel is published under the MIT license: 
 	https://github.com/fhAnso/Sentinel/blob/main/LICENSE

===[ Usage
 -v  Verbose output
 -d  Set the target domain name
	 [passive] (Without -w) request subdomains from external endpoints
	 Example: targetdomain.xyz 
 -w  [active] Use wordlist to bruteforce subdomains of the target
 -oS Specify the output (.txt) file path for subdomains
 -o4 Specify the output (.txt) file path for IPv4 addresses
 -o6 Specify the output (.txt) file path for IPv6 addresses
 -nP Specify the output directory path for all output files
 -oJ Specify the output (.json) file path for summary
 -c  Send GET request to retrieve the HTTP status code of every passive enumeration result
 -e  Exclude HTTP response codes (comma seperated)
 -f  Filter for specific HTTP response codes (comma seperated)
 -a  Analyze HTTP header of each subdomain (server etc.)
 -p  Scan subdomains for open ports in range
 -x  Extend endpoint DB with custom list (.txt)
 -t  Specify the request timeout
 -r  Route all requests through TOR
`
