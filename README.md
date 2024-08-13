<p align="center">
  <img src="https://github.com/fhAnso/Sentinel/blob/main/assets/logo.png" />
</p>

# Sentinel - X-Platform Subdomain Enumeration
### Description:
```txt
This program is designed to fetch, filter and validate subdomains 
from a specific host. The enumeration can be done by either using the passive
method using external ressources or directly via brute-force. Further 
information can be queried for the each subdomain. This information 
includes HTTP header analysis (which server etc.), HTTP status 
code to find out if and how the subdomain is reachable and the possibility 
to perform a port scan for each subdomain. The output can be customized as 
required and the results (divided into: IPv4, IPv6 and subdomains) are automatically 
saved for further processing. 
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
# All results will be saved in the "output" directory.
# Passive:
./bin/sentinel -d example.com -s -c -a -p 1-1000 -f 200,401,403
# Direct:
./bin/sentinel -d example.com -s -a -p 1-1000 -w /wordlists/subdomains.txt
```
#### Or simply `run` the <sentinel> `executable` without args to see the available `options`

```txt
By default, Sentinel will create 3 output files. The output files are 
divided into subdomains, IPv4 and IPv6 addresses. 
```

#### Options:
| Flags | Argument Type | Description |
| ----- | ----------- | ------------|
| -d | string | Specify the taget domain eg. example.com (default: passive) |
| -w | string | Use direct method by specifying the wordlist |
| -oS | string | Specify the output file path for subdomains |
| -o4 | string | Specify the output file path for IPv4 addresses |
| -o6 | string | Specify the output file path for IPv6 addresses |
| -c | - | Display the HTTP status code of each subdomain |
| -e | string | Exclude HTTP status codes from results |
| -f | string | Filter specific HTTP status codes from results |
| -s | - | Display only subdomains which can be resolved to IP addresses |
| -a | - | Analyze HTTP header of each subdomain (server etc.) |
| -p | string | Scan subdomains for open ports in range |
| -x | string | Extend endpoint DB with custom list (.txt) |
| -t | int | Specify the request timeout |
| -r | - | Route all requests through TOR |

# License
Sentinel is published under the ![MIT](https://github.com/fhAnso/Sentinel/blob/main/LICENSE) license