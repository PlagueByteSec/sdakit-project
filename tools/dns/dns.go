package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var fqdnPool []string

func parseFile(filePath string) error {
	stream, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Could not open file: %s\n%s", filePath, err)
	}
	defer stream.Close()
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		entry := scanner.Text()
		if len(entry) == 0 || strings.HasPrefix(entry, "#") || strings.HasPrefix(entry, "//") {
			continue
		}
		fqdnPool = append(fqdnPool, entry)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Scanner failed: %s", err)
	}
	return nil
}

func fqdnToTld(fqdn string) string {
	split := strings.Split(fqdn, ".")
	if len(split) < 2 {
		return fqdn
	}
	return strings.Join(split[len(split)-2:], ".")
}

func dnsLookup(fqdn string) {
	fmt.Println("[*] IP:")
	ip, err := net.LookupIP(fqdn)
	if err != nil {
		fmt.Println(" |  lookup for IP addresses failed:", err.Error())
	}
	var ipResult strings.Builder
	for idx := 0; idx < len(ip); idx++ {
		ip := ip[idx].String()
		parser := net.ParseIP(ip)
		if parser.To4() != nil {
			ipResult.WriteString(" |  A")
		} else {
			ipResult.WriteString("\r |  AAAA")
		}
		fmt.Printf("%s Entry: %s\n", ipResult.String(), ip)
	}
	fmt.Println("[*] MX:")
	mx, err := net.LookupMX(fqdn)
	if err != nil {
		fmt.Println(" |  lookup for MX entries failed:", err.Error())
	}
	for idx := 0; idx < len(mx); idx++ {
		fmt.Printf(" |  Host: %s, Pref: %d\n", mx[idx].Host, mx[idx].Pref)
	}
	fmt.Println("[*] CNAME:")
	cname, err := net.LookupCNAME(fqdn)
	if err != nil {
		fmt.Println(" |  lookup for CNAME failed:", err.Error())
	}
	fmt.Println(" | ", cname)
	fmt.Println("[*] TXT:")
	txt, err := net.LookupTXT(fqdn)
	if err != nil {
		fmt.Println(" |  lookup for TXT entries failed:", err.Error())
	}
	for idx := 0; idx < len(txt); idx++ {
		fmt.Println(" |  ", txt[idx])
	}
	fmt.Println("[*] NS:")
	ns, err := net.LookupNS(fqdn)
	if err != nil {
		fmt.Println(" |  lookup for NS entries failed:", err.Error())
	}
	for idx := 0; idx < len(ns); idx++ {
		fmt.Println(" |  ", ns[idx].Host)
	}
}

func main() {
	fmt.Println("\n\t< SDAkit - DNS lookup >")
	var path string
	flag.StringVar(&path, "f", "", "Specify path of list containing subdomains")
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Printf("Required argument MISSING -%s: %s\n", f.Name, f.Usage)
		})
	}
	parseFile(path)
	tld := fqdnToTld(fqdnPool[0])
	fmt.Println("\n[+]", tld)
	dnsLookup(tld)
	for idx := 0; idx < len(fqdnPool); idx++ {
		fqdn := fqdnPool[idx]
		fmt.Println("\n[+]", fqdn)
		dnsLookup(fqdn)
	}
	fmt.Println("\n[*] Finished")
}
