package lib

var Help = `
Sentinel
********
===[ Description
 This program is designed to fetch, filter and validate subdomains from a specific host.
 Subdomains can be enumerated by requesting external ressources (passive) or via brute 
 force (direct). 

 Sentinel is published under the MIT license: https://github.com/fhAnso/Sentinel

===[ Usage
 -v Verbose output
 -t Set the target domain name. 
	[passive] (Without -w) request subdomains from external endpoints
	Example: targetdomain.xyz 
 -o Specify the output (.txt) file 
 -c Send GET request to retrieve the HTTP status code of every passive enumeration result
 -w [direct] Use wordlist to bruteforce subdomains of the target
 -e Exclude response codes from bruteforce results (comma seperated)
 -s Only show subdomains that can be resolved to ip addresses 

===[ Not implemented yet
 -f Show only specific HTTP response codes from bruteforce results (comma seperated)
`
