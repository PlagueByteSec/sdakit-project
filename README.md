<p align="center">
  <img src="https://github.com/fhAnso/Sentinel/blob/main/assets/logo.png" />
</p>

# Sentinel - X-Platform Subdomain Emumeration
### Description:
```txt
This program is designed to fetch, filter and validate subdomains from a specific host.
The Sentinel project replaces the platform-dependent script "uma.sh" and makes it possible
to enumerate subdomains of a target passively using crt.sh and RapidDNS services or
directly via brute-force using a custom wordlist. The results will be saved among each 
other to provide a quick solution for further processing.
```

### Compile:
- `Build` the Sentinel `project`

`Windows`
```cmd
go build -o .\bin\sentinel.exe 
```
`Linux`
```bash
go build -o bin/sentinel 
```

### Usage:
- Request subdomains
```
<sentinel> -t example.com
```
#### Or simply `run` the <sentinel> `executable` without args to see the available `options`

### Example Output:
```txt
 ===[ Sentinel, v1.1.0 ]===

[*] Using passive enum method
[*] Formatting db entries..

1. Entry: https://crt.sh/?q=%25.HOST
 ===[ https://crt.sh/?q=%25.example.com

2. Entry: https://rapiddns.io/subdomain/HOST?full=1
 ===[ https://rapiddns.io/subdomain/example.com?full=1

[*] Using 2 endpoints
[*] Sending GET request to endpoints..

 ===[ support.example.com
 ===[ 20mail2.example.com
 ===[ www.example.com
 ===[ dev.example.com
 ===[ products.example.com

[*] 5 subdomains obtained. Finished in 1.458662222s
```

# License
Sentinel is published under the ![MIT](https://github.com/fhAnso/Sentinel/blob/main/LICENSE) license