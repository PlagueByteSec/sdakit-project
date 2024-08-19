<p align="center">
  <img src="https://github.com/fhAnso/Sentinel/blob/main/assets/logo.png" />
</p>

# Sentinel - X-Platform Subdomain Enumeration
### Description:
```txt
This project is designed to enumerate, filter, and validate subdomains for 
a specified host. Enumeration can be performed using either a passive method 
that leverages external resources or through direct brute-force technique. 
For each discovered subdomain, Sentinel retrieves detailed information, 
including HTTP header analysis (e.g. server), HTTP status codes to 
assess reachability, and the option to conduct a port scan. The output is 
customizable to meet user requirements, and results are categorized into 
IPv4, IPv6, and subdomains, with all findings automatically saved for 
subsequent processing.
```

### Build:
`Windows`
```cmd
go build -o .\bin\sentinel.exe 
```
`Linux`
```bash
go build -o bin/sentinel 
```

### Usage:
- specify the target and request subdomains
```
<sentinel> -d example.com
```
- extend the default enumeration
```bash
# Discover subdomains, display HTTP status codes, analyze 
# the header, display only resolvable subdomains and run a
# port scan against them.
# All results will be saved by default in the output directory.
# Passive:
./bin/sentinel -d example.com -s -c -a -p 1-1000 -f 200,401,403
# Active:
./bin/sentinel -d example.com -s -a -p 1-1000 -w /wordlists/subdomains.txt
```
#### Or simply `run` the <sentinel> `executable` without args to see the available `options`

```txt
By default, Sentinel will create 4 output files. The output files are 
divided into subdomains, IPv4/IPv6 addresses and a summary in JSON format. 
```

#### Options:
| Flags | Argument Type | Description |
| ----- | ----------- | ------------|
| -d | string | Specify the taget domain eg. example.com (default: passive) |
| -w | string | Use active method by specifying the wordlist |
| -oS | string | Specify the output file path for subdomains |
| -o4 | string | Specify the output file path for IPv4 addresses |
| -o6 | string | Specify the output file path for IPv6 addresses |
| -oJ | string | Specify the output file path for summary |
| -nP | string | Specify the output directory path for all output files |
| -c | bool | Display the HTTP status code of each subdomain |
| -e | string | Exclude HTTP status codes from results |
| -f | string | Filter specific HTTP status codes from results |
| -a | bool | Analyze HTTP header of each subdomain (server etc.) |
| -p | string | Scan subdomains for open ports in range |
| -x | string | Extend endpoint DB with custom list (.txt) |
| -t | int | Specify the request timeout |
| -r | bool | Route all requests through TOR |

# License
Sentinel is published under the ![MIT](https://github.com/fhAnso/Sentinel/blob/main/LICENSE) license