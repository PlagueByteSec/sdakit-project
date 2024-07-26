package lib

var Help = `
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
 -o Specify the output (.txt) file where every regex match from the main pool will be saved
=== Not implemented yet
 -c Send GET request to retrieve the HTTP status code of every entry
 -p Send ICMP packages to ping all entries
`
