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

	log "github.com/0xRyuk/ipx/internal/logger"
	"github.com/0xRyuk/ipx/internal/resolver"
	"github.com/0xRyuk/ipx/internal/util"
)

const (
	DefaultDNSServer = "8.8.8.8:53"
	DefaultTimeout = 500 * time.Millisecond
)


type Options struct {
	IPv4          *regexp.Regexp
	Host          string
	Filename      string
	SaveOutput    string
	ResolversFile string
	Resolvers     []string
	Timeout       time.Duration
	Threads       int
	Verbose       bool
	Plain         bool
	Count         int
}

var opts Options

func init() {
	opts.IPv4 = regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)
	flag.StringVar(&opts.Host, "d", "", "Set hostname to resolve (i.e. example.com)")
	flag.StringVar(&opts.Filename, "f", "", "Read a file containing hostnames to resolve (i.e. hosts.txt)")
	flag.BoolVar(&opts.Plain, "i", false, "Print only IP address (default false)")
	flag.BoolVar(&opts.Verbose, "v", false, "Turn on verbose mode (default off)")
	flag.StringVar(&opts.ResolversFile, "r", "", "Resolvers list (i.e. resolvers.txt)")
	flag.DurationVar(&opts.Timeout, "timeout", DefaultTimeout, "Set timeout")
	flag.IntVar(&opts.Threads, "t", 20, "Number of threads to utilize")
	flag.StringVar(&opts.SaveOutput, "o", "", "Save output to a text file")
}

func main() {
	flag.Parse()

	util.HandleExit()

	var (
		count int
		total int
		in       io.Reader
		resolved []string
	)

	if len(opts.Resolvers) < 1 {
		opts.Resolvers = append(opts.Resolvers, DefaultDNSServer)
		if opts.ResolversFile != "" {
			opts.Resolvers = append(util.SetResolver(opts.ResolversFile, DefaultDNSServer), opts.Resolvers...)
		}
	}

	start := time.Now()
	if opts.Plain {
		opts.Verbose = false
	} else {
		fmt.Println(util.Banner())
	}
	// Determine the input source for the domain names to be resolved
	stat, _ := os.Stdin.Stat()
	if (stat.Mode()&os.ModeCharDevice) != 0 && (opts.Host == "" && opts.Filename == "") {
		flag.Usage()
		os.Exit(1)
	}

	if opts.Host != "" {
		in = strings.NewReader(opts.Host)
	} else if opts.Filename != "" {
		file, err := os.Open(opts.Filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		in = file

	} else if (stat.Mode() & os.ModeCharDevice) == 0 {
		in = os.Stdin
	}

	work := make(chan string)

	buf := bufio.NewScanner(in)
	go func() {
		for buf.Scan() {
			work <- strings.ToLower(buf.Text())
			total += 1
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
	for i := 0; i < opts.Threads; i++ {
		wg.Add(1)
		go func() {
			// ipAddrMap := r.Client(work, &wg, g.Resolvers, g.DefaultTimeout, &g.Count, g.IPv4, g.Verbose)
			ipAddrMap := resolver.Client(work, &wg, opts.Resolvers, opts.Timeout, &count, opts.IPv4, opts.Verbose)
			mu.Lock()
			for ipAddr := range ipAddrMap {
				// The resolved slice is protected by a mutex lock to prevent data race conditions.
				resolved = append(resolved, ipAddr)
			}
			mu.Unlock()
		}()

	}

	wg.Wait()

	if opts.Verbose {
		time.Sleep(500 * time.Millisecond)
		log.Info("Finished in ", time.Since(start))
		log.Info("Total domains resolved ", count, "/", total)
		log.Info(len(resolved), " IP address(s) found")
	}

	if opts.SaveOutput != "" {
		util.SaveToFile(resolved, opts.SaveOutput)
	}
}
