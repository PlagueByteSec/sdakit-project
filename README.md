<p align="center">
  <img src="https://github.com/PlagueByteSec/Sentinel/blob/main/assets/logoTransparent.png" alt="logoTransparent" width="500" height="250" />
  <br>
  <img src="https://img.shields.io/github/stars/PlagueByteSec/Sentinel?style=social" alt="stars" />
  <img src="https://img.shields.io/github/v/release/PlagueByteSec/Sentinel" alt="version" />
  <img src="https://img.shields.io/github/license/PlagueByteSec/Sentinel" alt="license" />
</p>

## :book: Description
```txt
This project is designed to enumerate, filter, and validate subdomains for a specified
target domain. For each identified subdomain, Sentinel collects comprehensive information, including 
HTTP header analysis (such as server software, HSTS, and used technology), the 
option to conduct port scans, checks for common misconfigurations, and attempts to determine 
the purpose of the server (e.g., Mail, API, etc.). This project also determines the subdomains
availability by requesting HTTP status codes and ping probes.

The CLI is designed for clarity and ease of use, providing a structured overview of the
results. Output is customizable and is organized into categories: IPv4, IPv6, subdomains,
and summaries. All findings are automatically saved for further processing.
```

## :paperclip: Overview

- ***Discovery***: Easily identify and analyze subdomains of a target domain:

| Method | Description | Required Flags | Example |
|--------|-------------|----------------|---------|  
| Passive | Use public services like cert transparency logs etc. | -d | ./sentinel -d example.com |
| Active (Direct) | Brute-force subdomains by sending HTTP requests to the target domain, and analyze the status codes | -d, -w | ./sentinel -d -w /path/to/wordlist |
| Active (DNS) | Built possible subdomains, and try to resolve them to IP addresses (A, AAAA) | -d, -w, -dns | ./sentinel -d example.com -w /path/to/wordlist -dns |

- ***Result Filtering***: Filter results to avoid unneccessary result overwhelming.
- ***Port Scanning***: Scanning the target subdomain for open ports.
- ***Response Analysis***: Analyze responses to identify the tech stack, common misconfigurations, and security weaknesses.
- ***Categorized Output***: Generate organized output for easy review and reporting.

## :wrench: Build and Usage

- If Sentinel needs to be used on `Windows`, make sure to add the `.exe` file extension:
```cmd
go build -o .\bin\sentinel.exe 
```

- This command can be used to compile Sentinel on `Linux`:
```bash
go build -o ./bin/sentinel 
```

- `Display` the available `options` by simply specifying the `-h` flag
```txt
===[ The Sentinel Project, By PlagueByteSec
        
	Contact: plaguebyte.sec@keemail.me
	License: MIT (https://github.com/PlagueByteSec/Sentinel/blob/main/LICENSE)
																		 
 ===[ Description

	This program is designed to fetch, filter and validate subdomains from a target domain.

 ===[ Options ]===

 -v	Verbose output
 -d	Set the target domain name
		[passive] (Without -w) request subdomains from external endpoints
		Example: targetdomain.xyz 
 -s	Set the target subdomain 
		Example: sub.targetdomain.xyz 
 -w	[active] Use wordlist to bruteforce subdomains of the target
 -dns	Use wordlist (-w) and resolve subdomains by querying a DNS

 ===[ ANALYSIS

 -c	Send HTTP request to retrieve the status code 
 -a	Analyze HTTP headers of each subdomain (server, csp, software, ...)
 -aH	Display all HTTP headers of response from -a and -aS
 -aS	Analyse the HTTP response from a subdomain (specified by -s)
 -p	Scan subdomains for open ports (comma seperated or from-to)
 -pS	Ping subdomains (privileged execution required)
 -pC	Ping subdomains from file (privileged execution required)
 -rF	Read IP addresses from file and perform RDNS lookup
 -dP	Analyse subdomain to determine its purpose (mail, API, ...)
 -mT	Test subdomain for common weaknesses (CORS, header injections, ...)

 ===[ FILTERS

 -e	Exclude HTTP response codes (comma seperated)
 -f	Filter for specific HTTP response codes (comma seperated)
 
 ===[ SETTINGS

 -x	Extend endpoint DB with custom list (.txt)
 -r	Route all requests through TOR: 127.0.0.1:9050, SOCKS5
 -m	Set HTTP request method (default: GET)
 -rD	Set HTTP request delay in ms
 -t	Specify the HTTP request timeout
 -dnsT	Set the timeout for DNS queries
 -dnsC	Specify a custom DNS server address (ip:port)
 -pC	Specify Ping count (default=2)
- aR	Follow redirects: 301, 302, 303, ...

 ===[ OUTPUT

 -dO	Disable all output file streams
 -nP	Specify the directory path for all output files
```

## :question: Which External Libraries are Used?

- Sentinel is using the the `go-version` library from [hashicorp](https://github.com/hashicorp/go-version) to compare local and remote versions.
- To be able to determine which ports are open, the `nmap` library from [Ullaakut](https://github.com/Ullaakut/nmap) is utilized.
- If subdomain reachability needs to be verified by a ping probe, `pro-bing` from [prometheus-community](https://github.com/prometheus-community/pro-bing) comes into play.

## :warning: IMPORTANT NOTICE

```txt
THIS TOOL IS INTENDED FOR ETHICAL USE ONLY. BY USING THE SENTINEL PROJECT, YOU
AGREE TO USE IT RESPONSIBLY AND LEGALLY. I AM NOT RESPONSIBLE FOR YOUR ACTIONS.
USE AT YOUR OWN RISK.
```

## :sparkler: Some Examples

- **Analyze the HTTP response header, and run a TCP scan for the ports 22,80,443,8080**
```bash
./bin/sentinel -d example.com -c -a -p 22,80,443,8080
```
- **Display only subdomains which responded with HTTP codes 200,401,403**
```bash
./bin/sentinel -d example.com -c -f 200,401,403
```
- **Avoid output of responses with HTTP status codes n/a,501,404**
```bash
./bin/sentinel -d example.com -c -e n/a,501,404
```
- **Send a ping to subdomains responded with HTTP code n/a** 
```bash
./bin/sentinel -d example.com -c -f n/a -pS
```
- **Check subdomains with status code 200,301,302,303 for common misconfigurations (follow redirections: true)**
```bash
./bin/sentinel -d example.com -c -f 200,301,302,303 -mT -aR
```
- **Determine the subdomain purpose, analyze the response header, and run a port scan from 1-1000**
```bash
./bin/sentinel -d example.com -c -dP -a -p 1-1000
```

## :memo: To-Do List

- [ ] Output summary to PDF
- [ ] Add other subdomains found in HTTP response to pool
- [ ] Implement functions:
  - [ ] HeaderInjection
  - [ ] RequestSmuggling
- [ ] Extend purpose detection

## LICENSE
```txt
MIT License

Copyright (c) 2024 fhAnso

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```