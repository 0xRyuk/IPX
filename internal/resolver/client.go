package resolver

import (
	"fmt"
	"net"
	"regexp"
	"sync"
	"time"

	log "github.com/0xRyuk/ipx/internal/logger"
	"github.com/0xRyuk/ipx/internal/util"
)

func Client(work <-chan string, wg *sync.WaitGroup, resolvers []string, timeout time.Duration, countSuccess *int, ipv4 *regexp.Regexp, verbose bool) map[string]bool {
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

		ips, err := net.LookupHost(hostname)
		go func() {
			if err != nil {
				errors <- err
			}
		}()
		for _, ipAddr := range ips {
			if !ipv4.MatchString(ipAddr) {
				continue
			}
			ipAddrs[ipAddr] = true
		}
		domain := util.HasSuffix(hostname, ".")
		rr, err := resolver.Resolve(domain)
		*countSuccess += 1
		go func() {
			if err != nil {
				*countSuccess -= 1
				errors <- err
			}
		}()
		aRecords := util.ParseARecord(rr)
		for _, ip := range aRecords {
			ipAddrs[ip] = true
		}
		for ip := range ipAddrs {
			if verbose {
				log.Info(util.FormatStr(ip))
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
