package util

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"

	g "github.com/0xRyuk/ipx/internal/global"
	log "github.com/0xRyuk/ipx/internal/logger"
	d "github.com/miekg/dns"
)

func HandleExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Warn("Keyboard interrupt: Exiting!")
			os.Exit(0)
		}
	}()
}

// Function to parse A Record and return a string slice of IP addresses.
func ParseARecord(rr []d.RR) []string {
	seen := make(map[string]bool)
	var ipAddrs []string
	// Print the response.
	for _, r := range rr {
		switch t := r.(type) {
		case *d.A:
			tA := t.A.String()
			if !seen[tA] {

				ipAddrs = append(ipAddrs, t.A.String())
				seen[tA] = true
			}
		}
	}
	return ipAddrs
}

func SetResolver(filename string) []string {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	if err != nil {
		log.Fatal(err)
	}
	var resolvers []string
	// Split the file contents on newline characters and return the resulting slice

	for s.Scan() {
		r := s.Text()
		if strings.Contains(r, g.DefaultDNSServer) {
			continue
		}
		// Append each line to the slice.
		resolvers = append(resolvers, HasSuffix(r, ":53"))
	}
	return resolvers
}

func HasSuffix(str, suffix string) string {
	if !strings.HasSuffix(str, suffix) {
		str += suffix
	}
	return str
}

func FormatStr(str string) string {
	return fmt.Sprintf("%s%s%s", g.G, str, g.Rst)
}

func SaveToFile(resolved []string) {
	// Open the file for writing
	file, err := os.Create(g.SaveOutput)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	// Write each string in the slice to the file
	// fmt.Println(g.Resolved)
	for _, line := range resolved {
		// fmt.Println("line",line)
		_, err := file.WriteString(line + "\n")
		if err != nil {
			log.Panic(err)
		}
	}
}

// Print banner randomly
func Banner() string {
	rand.Seed(time.Now().Unix())
	return g.Banner[rand.Intn(len(g.Banner))]
}
