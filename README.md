<p align="center">
  <img src="https://github.com/fhAnso/Sentinel/blob/main/assets/logo.png" />
</p>

# Sentinel - X-Platform Subdomain Emumeration
### Description:
```txt
This program is designed to fetch, filter and validate subdomains from a specific host.
The Sentinel project replaces the platform-dependent script "uma.sh" and makes it possible
to enumerate subdomains of a target passively using external services or
directly via brute-force using a custom wordlist. The results will be saved among each 
other to provide a quick solution for further processing.
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
- Specify the target and request subdomains
```
<sentinel> -t example.com
```
#### Or simply `run` the <sentinel> `executable` without args to see the available `options`

```txt
By default, Sentinel will create 3 output files. The output files are divided into subdomains, 
IPv4 and IPv6 addresses. These are each filtered so that they can be used directly afterwards. 
```

### Example Output:
```txt
 ===[ Sentinel, v1.2.0 ]===

[*] Using passive enum method
[*] Formatting db entries..
[*] Sending GET request to endpoints..

 ===[ www.example.com (2606:2800:21f:cb07:6820:80da:af6b:8b2c, 93.184.215.14)

[*] 5 subdomains obtained, 1 displayed 
[*] Finished in 1.4153683s
```

# License
Sentinel is published under the ![MIT](https://github.com/fhAnso/Sentinel/blob/main/LICENSE) license