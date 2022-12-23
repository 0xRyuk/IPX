package util

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"

	log "github.com/0xRyuk/ipx/internal/logger"
	d "github.com/miekg/dns"
)

// Color constants
const (
	BG     = "\033[47m"   //BG
	B      = "\033[0;90m" //Black
	G      = "\033[1;32m" //Green
	R      = "\033[1;31m" //Red
	Y      = "\033[1;33m" //Yellow
	W      = "\033[1;97m" //White
	IC     = "\033[0;96m" //High Intensty Cyan
	Rst    = "\033[1;m"   //Reset
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

func ParseARecord(rr []d.RR) []string {
	seen := make(map[string]bool)
	var ipAddrs []string
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

func SetResolver(filename, DefaultDNSServer string) []string {
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

	for s.Scan() {
		r := s.Text()
		if strings.Contains(r, DefaultDNSServer) {
			continue
		}
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
	return fmt.Sprintf("%s%s%s", G, str, Rst)
}

func SaveToFile(resolved []string, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	for _, line := range resolved {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			log.Panic(err)
		}
	}
}

func Banner() string {
	banner := []string{fmt.Sprintf("\t%s  ____ __  %s_,  ,\n\t%s ( /( /  )%s( |,' \n\t%s  /  /--' %s  +   \n\t%s_/_ /    %s_,'|__ \n\t\t%sv0.1.0 %sbeta\n%s", IC, R, IC, R, IC, R, IC, R, W, Y, Rst),
		fmt.Sprintf("%s\t_ ___  %s_  _ \n\t%s| |__]  %s\\/\n\t%s| |    %s_/\\_ \n\t%s\tv0.1.0%s beta%s", IC, R, IC, R, IC, R, W, Y, Rst)} // Fancy banner(s)
	rand.Seed(time.Now().Unix())
	return banner[rand.Intn(len(banner))]
}
