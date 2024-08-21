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
IPv4, IPv6, subdomains and summary, with all findings automatically saved for 
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
`Examples`:
```bash
# Discover subdomains, display HTTP status codes, analyze 
# the header and run a port scan against them.
# All results will be saved by default in the output directory.
# Passive:
./bin/sentinel -d example.com -c -a -p 1-1000 -f 200,401,403
# Active:
./bin/sentinel -d example.com -a -p 1-1000 -w /wordlists/subdomains.txt
# DNS:
./bin/sentinel -d example.com -dns -w /wordlists/subdomains.txt -a -c -p 1-1000
```
#### Or simply `run` the <sentinel> `executable` without args to see the available `options`

```txt
By default, Sentinel will create 4 output files. The output files are 
divided into subdomains, IPv4/IPv6 addresses and a summary in JSON format. 
```

# License
Sentinel is published under the ![MIT](https://github.com/fhAnso/Sentinel/blob/main/LICENSE) license