package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	probing "github.com/prometheus-community/pro-bing"
)

var addressPool []string

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
		addressPool = append(addressPool, entry)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Scanner failed: %s", err)
	}
	return nil
}

func pingAddress(address string, pingCount int) error {
	fmt.Println(" |  ")
	pinger, err := probing.NewPinger(address)
	if err != nil {
		return fmt.Errorf("[-] failed to initialize pinger: %s", err)
	}
	pinger.Count = pingCount
	pinger.SetPrivileged(true)
	pinger.OnRecv = func(pkt *probing.Packet) {
		fmt.Printf(" |  %d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}
	pinger.OnDuplicateRecv = func(pkt *probing.Packet) {
		fmt.Printf(" |  %d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.TTL)
	}
	pinger.OnFinish = func(stats *probing.Statistics) {
		fmt.Printf(" |  \n |  --- %s ping statistics ---\n", stats.Addr)
		fmt.Printf(" |  %d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf(" |  round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}
	if err = pinger.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Println("\n\t< SDAkit - List Pinger >")
	var (
		path  string
		count int
	)
	flag.StringVar(&path, "f", "", "Specify path of list containing subdomains")
	flag.IntVar(&count, "c", 2, "Set ping count")
	flag.Parse()
	if flag.NFlag() < 1 {
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Printf("Required argument MISSING -%s: %s\n", f.Name, f.Usage)
		})
	}
	parseFile(path)
	for idx := 0; idx < len(addressPool); idx++ {
		address := addressPool[idx]
		fmt.Println("\n[*]", address)
		if err := pingAddress(address, count); err != nil {
			fmt.Println("[-]", err.Error())
			os.Exit(-1)
		}
	}
	fmt.Println("\n[*] Finished")
}
