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

```txt
The endpoints are read and processed from a simple database. After the results have 
arrived, they are checked and inserted into a main pool. From there, the program will 
perform various operations (IP Lookup etc.) for each entry. Further options are 
available as CLI parameters.
```

### Usage:
- Specify the target and request subdomains
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

3. Entry: https://jldc.me/anubis/subdomains/HOST
 ===[ https://jldc.me/anubis/subdomains/example.com

[*] Using 3 endpoints
[*] Sending GET request to endpoints..

 ===[ dev.example.com
 ===[ products.example.com
 ===[ support.example.com
 ===[ www.example.com (2606:2800:21f:cb07:6820:80da:af6b:8b2c, 93.184.215.14)
 ===[ m.example.com

[*] 5 subdomains obtained. Finished in 1.4153683s
```

# License
Sentinel is published under the ![MIT](https://github.com/fhAnso/Sentinel/blob/main/LICENSE) license