package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/PlagueByteSec/Sentinel/v1/internal/shared"
)

func EditDbEntries(args *shared.Args) ([]string, error) {
	/*
		All endpoints will be read from db.go and formatted by replacing
		the placeholder (HOST) with the target domain. If a text file containing
		custom endpoints is specified by the -x flag, those will be
		added to the existing entries.
	*/
	entries := make([]string, 0, len(shared.Db))
	for idx, entry := range shared.Db {
		endpoint := strings.Replace(entry, shared.Placeholder, args.Domain, 1)
		if args.Verbose {
			fmt.Fprintf(shared.GStdout, "\n%d. Entry: %s\n ===[ %s\n", idx+1, entry, endpoint)
		}
		entries = append(entries, endpoint)
	}
	if args.DbExtendPath != "" {
		PrintVerbose("\n[*] Extending endpoints..")
		stream, err := os.Open(args.DbExtendPath)
		if err != nil {
			shared.Glogger.Println(err)
			return nil, errors.New("failed to open file stream for: " + args.DbExtendPath)
		}
		defer stream.Close()
		scanner := bufio.NewScanner(stream)
		idx := 0
		for scanner.Scan() {
			entry := scanner.Text()
			if !strings.Contains(entry, shared.Placeholder) {
				fmt.Fprintln(shared.GStdout, "[-] Invalid pattern (HOST missing): "+entry)
				continue
			}
			endpoint := strings.Replace(entry, shared.Placeholder, args.Domain, 1)
			PrintVerbose("\n%d. X Entry: %s\n ===[ %s\n", idx+1, entry, endpoint)
			entries = append(entries, endpoint)
			idx++
		}
	}
	PrintVerbose("\n[*] Using %d endpoints\n", len(entries))
	return entries, nil
}
