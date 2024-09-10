<p align="center">
  <img src="https://github.com/PlagueByteSec/Sentinel/blob/main/assets/logoTransparent.png" alt="logoTransparent" width="500" height="250" />
  <br>
  <img src="https://img.shields.io/github/stars/PlagueByteSec/Sentinel?style=social" alt="stars" />
  <img src="https://img.shields.io/github/v/release/PlagueByteSec/Sentinel" alt="version" />
  <img src="https://img.shields.io/github/license/PlagueByteSec/Sentinel" alt="license" />
</p>

## Description
```
This project is designed to enumerate, filter, and validate subdomains for a specified
target domain. For each identified subdomain, Sentinel collects comprehensive information, 
including HTTP header analysis, common flaws, open ports, the subdomain purpose and more.
The CLI is designed for clarity and ease of use, providing a structured overview of the
results. Output is customizable and is organized into categories: IPv4, IPv6, subdomains,
and summaries. All findings are automatically saved for further processing.
```

## Overview

- ***Discovery***: Identify subdomains using the preferred option:

| Method | Description | Required Flags | Example |
|--------|-------------|----------------|---------|  
| Passive | Use public services like cert transparency logs, public DNS services etc. | -d [DOMAIN] | ./bin/sentinel -d example.com |
| Active (Direct) | Brute-force subdomains by sending HTTP requests to the target domain, and analyze the status codes | -d [DOMAIN], -w [wordlist] | ./bin/sentinel -d -w /path/to/wordlist |
| Active (DNS) | Built possible subdomains, and try to resolve them to IP addresses (A, AAAA) | -d [DOMAIN], -w [WORDLIST], -dns | ./bin/sentinel -d example.com -w /path/to/wordlist -dns |

- ***Result Filtering***: Filter results to avoid unneccessary result overwhelming.
- ***Port Scanning***: Scanning the target subdomain for open ports.
- ***Subdomain Analysis***: Analyze responses to payloads and general requestes to identify the tech stack, common misconfigurations, and security weaknesses.
- ***Categorized Output***: Generate organized output for summary and further processing.

## Quick Start

- `Windows`:
```
.\build\Windows\build.bat
```

- `Linux`:
```bash
cmd="./build/Linux/build.sh";chmod +x $cmd && $cmd
```

- `Display` the available `options`:
```
./bin/sentinel -h
```

## NOTE

```
THE SENTINEL PROJECT IS INTENDED FOR ETHICAL AND EDUCATIONAL 
USE ONLY. THE DEVELOPERS ARE NOT RESPONSIBLE FOR YOUR ACTIONS. 

THINK BEFORE YOU TYPE.
```

## License
The Sentinel Project is published under the [MIT](https://github.com/PlagueByteSec/Sentinel/blob/main/LICENSE) license.
