<p align="center">
  <img src="https://github.com/PlagueByteSec/sentinel-project/blob/main/assets/TheSentinelProject-LogoTransparent.png" alt="logoTransparent" width="500" height="250" />
  <br>
  <img src="https://img.shields.io/github/stars/PlagueByteSec/sentinel-project?style=social" alt="stars" />
  <img src="https://img.shields.io/github/v/release/PlagueByteSec/sentinel-project" alt="version" />
  <img src="https://img.shields.io/github/license/PlagueByteSec/sentinel-project" alt="license" />
</p>

## Description
```
The Sentinel Project is designed to assist security testers in the reconnaissance phase
by providing various methods for subdomain discovery and analysis. The main goal 
of this project is to make the process of subdomain enumeration as easy as possible 
by automatically performing general analysis, testing for typical flaws, 
determining the subdomain's purpose, and ensuring that all basic needs are met. All 
results will be sorted and saved for further processing.
```

## Overview

- ***Discovery***: Identify subdomains using the preferred option:

| Method | Description | Required Flags | Example |
|--------|-------------|----------------|---------|  
| Passive | Use public services like cert transparency logs, public DNS services etc. | -d [DOMAIN] | ./bin/sentinel -d example.com |
| Active (Direct) | Brute-force subdomains by sending HTTP requests to the target domain, and analyze the status codes | -d [DOMAIN], -w [wordlist] | ./bin/sentinel -d -w /path/to/wordlist |
| Active (DNS) | Built possible subdomains, and try to resolve them to IP addresses (A, AAAA) | -d [DOMAIN], -w [WORDLIST], -dns | ./bin/sentinel -d example.com -w /path/to/wordlist -dns |

- ***Result Filtering***: Filter results to avoid unneccessary output overwhelming.
- ***Subdomain Analysis***: Test each subdomain automatically for common security flaws and more.
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

### Getting Started

- [Setup](https://plaguebytesec.github.io/sentinel-project/pages/setup)
- [Usage](https://plaguebytesec.github.io/sentinel-project/pages/usage)
- [Examples](https://plaguebytesec.github.io/sentinel-project/pages/examples)

## License
The Sentinel Project is published under the [MIT](https://github.com/PlagueByteSec/sentinel-project/blob/main/LICENSE) license.
