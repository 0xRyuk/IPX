package main

// This package defines the main entry point for the program.
import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	g "github.com/0xRyuk/ipx/internal/global"
	log "github.com/0xRyuk/ipx/internal/logger"
	r "github.com/0xRyuk/ipx/internal/resolver"
	u "github.com/0xRyuk/ipx/internal/util"
)

func init() {
	g.IPv4 = regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)
	flag.StringVar(&g.Host, "d", "", "Set hostname to resolve (i.e. example.com)")
	flag.StringVar(&g.Filename, "f", "", "Read a file containing hostnames to resolve (i.e. hosts.txt)")
	flag.BoolVar(&g.Plain, "i", false, "Print only IP address (default false)")
	flag.BoolVar(&g.Verbose, "v", false, "Turn on verbose mode (default off)")
	flag.StringVar(&g.ResolversFile, "r", "", "Resolvers list (i.e. resolvers.txt)")
	flag.DurationVar(&g.Timeout, "timeout", g.DefaultTimeout, "Set timeout")
	flag.IntVar(&g.Threads, "t", 20, "Number of threads to utilize")
	flag.StringVar(&g.SaveOutput, "o", "", "Save output to a text file")
}

func main() {
	// Parse command-line flags
	flag.Parse()
	// Set up a signal handler to clean up resources when the program is interrupted or terminated
	u.HandleExit()

	var (
		in       io.Reader
		resolved []string
	)

	if len(g.Resolvers) < 1 {
		// Use the default DNS server if no resolvers were provided
		g.Resolvers = append(g.Resolvers, g.DefaultDNSServer)
		if g.ResolversFile != "" {
			g.Resolvers = append(u.SetResolver(g.ResolversFile), g.Resolvers...)
		}
	}

	start := time.Now()
	if g.Plain {
		g.Verbose = false
	} else {
		fmt.Println(u.Banner())
	}
	// Determine the input source for the domain names to be resolved
	// If a domain name was provided as a command-line flag, use it as the input.
	// If a filename was provided, open the file and use it as the input.
	// If neither of these is provided and running in non-interactive mode, use standard input as the input source.

	stat, _ := os.Stdin.Stat()
	if (stat.Mode()&os.ModeCharDevice) != 0 && (g.Host == "" && g.Filename == "") {
		flag.Usage()
		os.Exit(1)
	}

	if g.Host != "" {
		in = strings.NewReader(g.Host)
	} else if g.Filename != "" {
		file, err := os.Open(g.Filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		in = file

	} else if (stat.Mode() & os.ModeCharDevice) == 0 {
		in = os.Stdin
	}

	// Create a work channel and a goroutine to read the domain names from the input source and send them to the work channel
	work := make(chan string)

	buf := bufio.NewScanner(in)
	go func() {
		for buf.Scan() {
			work <- strings.ToLower(buf.Text())
			g.Total += 1
		}
		if err := buf.Err(); err != nil {
			log.Fatal(err)
		}
		close(work)
	}()
	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	// Create goroutines to resolve domain names concurrently using multiple threads.
	for i := 0; i < g.Threads; i++ {
		wg.Add(1)
		go func() {
			ipAddrMap := r.Client(work, &wg, g.Resolvers, g.DefaultTimeout, &g.Count, g.IPv4, g.Verbose)
			mu.Lock()
			for ipAddr := range ipAddrMap {
				// The resolved slice is protected by a mutex lock to prevent data race conditions.
				resolved = append(resolved, ipAddr)
			}
			mu.Unlock()
		}()

	}

	wg.Wait()

	if g.Verbose {
		time.Sleep(500 * time.Millisecond)
		log.Info("Finished in ", time.Since(start))
		log.Info("Total domains resolved ", g.Count, "/", g.Total)
		log.Info(len(resolved), " IP address(s) found")
	}
	// If the SaveOutput flag is set, save the resolved slice to a file.
	if g.SaveOutput != "" {
		u.SaveToFile(resolved)
	}
}
