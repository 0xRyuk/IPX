package resolver

import (
	"fmt"
	"net"
	"regexp"
	"sync"
	"time"

	log "github.com/0xRyuk/ipx/internal/logger"
	u "github.com/0xRyuk/ipx/internal/util"
)

func Client(work <-chan string, wg *sync.WaitGroup, resolvers []string, timeout time.Duration, countSuccess *int, IPv4 *regexp.Regexp, verbose bool) map[string]bool {
	errors := make(chan error)
	ipAddrs := make(map[string]bool)
	defer wg.Done()
	// create resolver object to perform DNS queries.
	resolver := NewResolver(
		resolvers,
		timeout,
	)
	var localCount int

	for hostname := range work {
		// resolve hostename using local resolver.
		ips, err := net.LookupHost(hostname)
		go func() {
			if err != nil {
				// send error through the error channel.
				errors <- err
			}
		}()
		for _, ipAddr := range ips {
			if !IPv4.MatchString(ipAddr) {
				continue
			}
			ipAddrs[ipAddr] = true
		}
		domain := u.HasSuffix(hostname, ".")
		// Perform a DNS query for the specified domain.
		rr, err := resolver.Resolve(domain)
		*countSuccess += 1
		go func() {
			if err != nil {
				*countSuccess -= 1
				errors <- err
			}
		}()
		aRecords := u.ParseARecord(rr)
		for _, ip := range aRecords {
			ipAddrs[ip] = true
		}
		for ip := range ipAddrs {
			if verbose {
				log.Info(u.FormatStr(ip))
				go func() {
					for err := range errors {
						log.Error(err)
					}
				}()
			} else {
				fmt.Println(ip)
			}
		}
	}

	*countSuccess += localCount
	return ipAddrs
}
