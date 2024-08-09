package lib

var Help = `
Sentinel
********
===[ Description
 This program is designed to fetch, filter and validate subdomains from a specific host.
 Subdomains can be enumerated by requesting external ressources (passive) or via brute 
 force (direct). 

 Sentinel is published under the MIT license: 
 	https://github.com/fhAnso/Sentinel/blob/main/LICENSE

===[ Usage
 -v  Verbose output
 -t  Set the target domain name. 
	 [passive] (Without -w) request subdomains from external endpoints
	 Example: targetdomain.xyz 
 -w  [direct] Use wordlist to bruteforce subdomains of the target
 -oS Specify the output (.txt) file path for subdomains
 -o4 Specify the output (.txt) file path for IPv4 addresses
 -o6 Specify the output (.txt) file path for IPv6 addresses
 -c  Send GET request to retrieve the HTTP status code of every passive enumeration result
 -e  Exclude HTTP response codes (comma seperated)
 -f  Filter for specific HTTP response codes (comma seperated)
 -s  Display only subdomains which can be resolved to IP addresses
`
