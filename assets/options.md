# Available Options

#### The options are divided into

- `Core`: Specify the target domain, and enumeration method
- `Analysis`: Analyze the enumerated subdomains
- `Filters`: Ensure only the needed results are displayed
- `Settings`: Customizing the scanner behaviour
- `Output`: Decide what to do with the output files

### Argument Overview

#### Core
| Flag | Description |
|------|-------------|
| -v | Enable verbose output |
| -d [DOMAIN] | Specify the target domain |
| -i | Set the target IP address |
| -w [WORDLIST] | Use a wordlist for direct and DNS enumeration |
| -dns | Enable DNS enumeration (-w required) |
| -vhost | Enable VHOST enumeration | 

#### Analysis
| Flag | Description |
|------|-------------|
| -c | Display the HTTP status code of each subdomain |
| -a | Analyze the HTTP response from each subdomain |
| -aH | Expose all HTTP response headers of each subdomain |
| -aS | Inspect the HTTP response from a subdomain specified by -s |
| -p | Scan each subdomain for open ports |
| -dP | Determine the subdomain purpose |
| -mT | Scan each subdomain for common flaws |

#### Filters
| Flag | Description |
|------|-------------|
| -e | Exclude subdomains with specified HTTP response codes |
| -f | Display only subdomains returning the specified HTTP response codes |
| -fS |	Exclude response size(s) from VHost results |

#### Settings
| Flag | Description |
|------|-------------|
| -x | Extend the default endpoints for passive enumeration |
| -r | Route all traffic through TOR |
| -m | Specify the HTTP method used for enumeration |
| -rD | Specify the request delay of each HTTP request |
| -t | Specify the timeout for each HTTP request |
| -dnsT | Specify the timeout for DNS queries |
| -dnsC | Use a custom DNS server address |
| -aR | Follow HTTP redirects |

#### Output
| Flag | Description |
|------|-------------|
| -dO | Disable auto saving of any results |
| -nP | Use a custom directory path for all output files |