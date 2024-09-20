package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
)

func analyseHeaders(uri string, headersOnly bool) error {
	parser, err := url.Parse(uri)
	if err != nil || parser.Scheme == "" {
		return fmt.Errorf("could not verify given address")
	}
	resp, err := http.Get(uri)
	if err != nil {
		return fmt.Errorf("failed to send GET: %s", err)
	}
	defer resp.Body.Close()
	for key, values := range resp.Header {
		for idx := 0; idx < len(values); idx++ {
			if headersOnly {
				fmt.Printf("%s: %s\n", key, values[idx])
				continue
			}
			value := values[idx]
			switch key {
			case "Server":
				fmt.Println("[+] Server:", value)
			case "X-XSS-Protection":
				fmt.Println("[*] XSS Protection:", value)
			case "X-Frame-Options":
				fmt.Println("[*] Clickjacking Protection:", value)
			case "Set-Cookie":
				fmt.Println("[*] Cookies:", value)
			case "X-Powered-By":
				fmt.Println("[+] Software:", value)
			case "Strict-Transport-Security":
				fmt.Println("[*] HSTS:", value)
			case "Content-Security-Policy":
				fmt.Println("[*] Content Security Policy:", value)
			case "X-Content-Type-Options":
				fmt.Println("[*] Content Type Options:", value)
			case "Referrer-Policy":
				fmt.Println("[*] Referrer Policy:", value)
			case "Feature-Policy", "Permissions-Policy":
				fmt.Println("[*] Feature Policy:", value)
			case "Access-Control-Allow-Origin":
				fmt.Println("[*] CORS:", value)
			case "X-Download-Options":
				fmt.Println("[*] Download Options:", value)
			case "X-Permitted-Cross-Domain-Policies":
				fmt.Println("[*] Cross-Domain Policies:", value)
			}
		}
	}
	return nil
}

func main() {
	fmt.Print("\n\t< SDAkit - HTTP Header Analyzer >\n\n")
	var (
		uri         string
		headersOnly bool
	)
	flag.StringVar(&uri, "u", "", "Specify the target URL")
	flag.BoolVar(&headersOnly, "o", false, "Show all headers response (no analysis)")
	flag.Parse()
	if flag.NFlag() < 1 {
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Printf("Required argument MISSING -%s: %s\n", f.Name, f.Usage)
		})
	}
	analyseHeaders(uri, headersOnly)
	fmt.Println("\n[*] Finished")
}
