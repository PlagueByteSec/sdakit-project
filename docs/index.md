<p align="center">
  <img src="https://github.com/PlagueByteSec/Sentinel/blob/main/assets/logoTransparent.png" alt="logoTransparent" width="500" height="250" />
  <br>
</p>

## Description
```txt
This project is designed to enumerate, filter, and validate subdomains for a specified
target domain. For each identified subdomain, Sentinel collects comprehensive information, 
including HTTP header analysis, common flaws, open ports, the subdomain purpose and more.

The CLI is designed for clarity and ease of use, providing a structured overview of the
results. Output is customizable and is organized into categories: IPv4, IPv6, subdomains,
and summaries. All findings are automatically saved for further processing.
```

## External Libraries

- Sentinel is using the the `go-version` library from [hashicorp](https://github.com/hashicorp/go-version) to compare local and remote versions.
- To be able to determine which ports are open, the `nmap` library from [Ullaakut](https://github.com/Ullaakut/nmap) is utilized.
- If subdomain reachability needs to be verified by a ping probe, `pro-bing` from [prometheus-community](https://github.com/prometheus-community/pro-bing) comes into play.

## Getting Started

- [Setup](https://github.com/PlagueByteSec/Sentinel/blob/main/docs/pages/setup.md)
- [Usage](https://github.com/PlagueByteSec/Sentinel/blob/main/docs/pages/usage.md)
- [Examples](https://github.com/PlagueByteSec/Sentinel/blob/main/docs/pages/examples.md)