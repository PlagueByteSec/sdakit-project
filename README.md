<p align="center">
  <img src="https://github.com/fhAnso/Sentinel/blob/main/assets/logo.png" />
</p>

# Sentinel - Passive Subdomain Emumeration
### Description:
```txt
This program is designed to fetch, filter and validate subdomains from a specific host.
The Sentinal project replaces the platform-dependent script "uma.sh" and makes it possible
to passively enumerate subdomains of a target using crt.sh and RapidDNS services independently of
the OS. The results will be saved among each other to provide a quick solution for
further processing.
```

### Compile:
- `Build` the Sentinel `project`

`Windows`
```cmd
go build -o .\bin\sentinel.exe .\main.go
```
`Linux`
```bash
go build -o bin/sentinel main.go
```

### Usage:
- Request subdomains
```
<sentinel> -t example.com
```
#### Or simply `run` the <sentinels> `executable` without args to see the available `options`

### Example Output:
```txt
[+] support
    ╚► support.example.com
[+] 20mail2
    ╚► 20mail2.example.com
[+] www
    ╚► www.example.com
[+] m
    ╚► m.example.com
[+] dev
    ╚► dev.example.com
[+] products
    ╚► products.example.com
```

# License
Sentinels is published under the ![MIT](https://github.com/fhAnso/Sentinel/blob/main/LICENSE) license