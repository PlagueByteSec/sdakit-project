## Usage Examples

- All listed flags are also compatible with Windows, only the paths must be adjusted

### Basic Enumeration

- **Use basic `passive` enumeration method**
```bash
./bin/sentinel -d example.com
```

- **Use basic `active` enumeration method by brute-forcing possible subdomains**
```bash
./bin/sentinel -d example.com -w /wordlists/SubdomainListA.txt
```

- **Use basic `active` enumeration method by trying to resolve possible subdomains to IP addresses**
```bash
./bin/sentinel -d example.com -w /wordlists/SubdomainListA.txt -dns
```

### Response Filter

- **Display only subdomains which responded with HTTP codes 200,401,403**
```bash
./bin/sentinel -d example.com -c -f 200,401,403
```

- **Exclude subdomains which responded with HTTP status codes n/a,501,404**
```bash
./bin/sentinel -d example.com -c -e n/a,501,404
```
### Subdomain Analysis

- **Analyze the HTTP response header, and run a TCP scan for the ports 22,80,443,8080**
```bash
./bin/sentinel -d example.com -c -a -p 22,80,443,8080
```

- **Test reachability of subdomains which responded with HTTP code n/a** 
```bash
./bin/sentinel -d example.com -c -f n/a -pS
```

- **Scan each subdomain with status code 200,301,302,303 for common misconfigurations**
```bash
./bin/sentinel -d example.com -c -f 200,301,302,303 -mT -aR
```

- **Determine the purpose of each subdomain, analyze the response headers, and run a port scan from 1-1000**
```bash
./bin/sentinel -d example.com -c -dP -a -p 1-1000
```

### Output
- **Disable the automated generation of all output files**
```bash
./bin/sentinel -d example.com -w /wordlists/SubdomainListA.txt -dns -dO
```

- **Specify the path of the output directory**
```bash
./bin/sentinel -d example.com -nP /home/$USER/Documents/
```

<div align="center">
<a href="#">Home</a>
</div>
