<p align="center">
  <img src="https://github.com/fhAnso/Sentinel/blob/main/assets/logo.png" />
</p>

# Sentinel - X-Platform Subdomain Emumeration
### Description:
```txt
This program is designed to fetch, filter and validate subdomains 
from a specific host. This works by querying public services like DNS, 
certificate transparency logs, etc. The output can be adjusted as 
needed and the results are automatically saved for further processing. 
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
<sentinel> -t example.com
```
- extend the default enumeration
```bash
# Discover subdomains, display HTTP status codes, analyze 
# the header and display only resolvable subdomains.
# All results will be saved in the "output" directory.
./bin/sentinel -t example.com -s -c -a
```
#### Or simply `run` the <sentinel> `executable` without args to see the available `options`

```txt
By default, Sentinel will create 3 output files. The output files are 
divided into subdomains, IPv4 and IPv6 addresses. 
```

#### Options:
| Flags | Description |
| ----- | ----------- |
| -t | Specify the taget domain eg. example.com (default: passive) |
| -w | Use direct method by specifying the wordlist |
| -oS | Specify the output file path for subdomains |
| -o4 | Specify the output file path for IPv4 addresses |
| -o6 | Specify the output file path for IPv6 addresses |
| -c | Display the HTTP status code of each subdomain |
| -e | Exclude HTTP status codes from results |
| -f | Filter specific HTTP status codes from results |
| -s | Display only subdomains which can be resolved to IP addresses |
| -a | Analyze HTTP header of each subdomain (server etc.) |

# License
Sentinel is published under the ![MIT](https://github.com/fhAnso/Sentinel/blob/main/LICENSE) license