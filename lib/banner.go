package lib

var Help = `
Sentinel
********
===[ Description
 This program is designed to fetch, filter and validate subdomains from a specific host.
 The Sentinel project replaces the platform-dependent script "uma.sh" and makes it possible
 to passively enumerate subdomains of a target independently of the OS. The results will be 
 saved among each other to provide a quick solution for further processing.

===[ Usage
 -t Set the target hostname. 
	[passive] (Without -w) this hostname will be be send to "RapidDNS" and "crt.sh"
	Example: targethostname.xyz 
 -o Specify the output (.txt) file where every regex match from the main pool will be saved
 -c Send GET request to retrieve the HTTP status code of every passive enumeration result
 -w [direct] Use wordlist to bruteforce subdomains of the target
 -e Exclude response codes from bruteforce results (comma seperated)

===[ Not implemented yet
 -f Show only specific HTTP response codes from bruteforce results (comma seperated)
 -p Send ICMP packages to ping all entries
`
