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

func Client(work <-chan string, wg *sync.WaitGroup, resolvers []string, timeout time.Duration, countSuccess *int, ipv4 *regexp.Regexp, verbose bool) map[string]map[string]bool {

	errors := make(chan error)
	ipAddrs := make(map[string]map[string]bool)
	defer wg.Done()
	// create resolver object to perform DNS queries.
	resolver := NewResolver(
		resolvers,
		timeout,
	)
	var localCount int

	for hostname := range work {
		ipAddrs[hostname] = make(map[string]bool)
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
			ipAddrs[hostname][ipAddr] = true
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
			ipAddrs[hostname][ip] = true
		}
		for hostname, ipAddrMap := range ipAddrs {
			if verbose {
				var ips []string //list of ips to print in verbose mode
				for ip := range ipAddrMap {
					ips = append(ips, ip)
				}
				log.Info(util.FormatStr(hostname+" "), ips)
				go func() {
					for err := range errors {
						log.Error(err)
					}
				}()
			} else {
				for ip := range ipAddrMap {
					fmt.Println(ip)
				}
			}
		}
	}

	*countSuccess += localCount
	return ipAddrs
}
