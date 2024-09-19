### Description
```txt
The SDAkit Project is designed to assist security testers in the reconnaissance phase
by providing various methods for subdomain discovery and analysis. The main goal 
of this project is to make the process of subdomain enumeration as easy as possible 
by automatically performing general analysis, testing for typical flaws, 
determining the subdomain's purpose, and ensuring that all basic needs are met. All 
results will be sorted and saved for further processing.
```

### External Libraries

- sdakit is using the the `go-version` library from [hashicorp](https://github.com/hashicorp/go-version) to compare local and remote versions.
- To be able to determine which ports are open, the `nmap` library from [Ullaakut](https://github.com/Ullaakut/nmap) is utilized.
- If subdomain reachability needs to be verified by a ping probe, `pro-bing` from [prometheus-community](https://github.com/prometheus-community/pro-bing) comes into play.

### Getting Started

- [Setup](https://plaguebytesec.github.io/sdakit-project/pages/setup)
- [Usage](https://plaguebytesec.github.io/sdakit-project/pages/usage)
- [Examples](https://plaguebytesec.github.io/sdakit-project/pages/examples)