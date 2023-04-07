package resolver

import (
	o "github.com/0xRyuk/ipx/internal/output"
	"github.com/0xRyuk/ipx/internal/util"
	"net"
	"regexp"
	"sync"
	"time"
)

func Client(work <-chan string, wg *sync.WaitGroup, resolvers []string, timeout time.Duration, countSuccess *int, ipv4 *regexp.Regexp, verbose bool) map[string]map[string]bool {

	ipAddrs := make(map[string]map[string]bool)
	defer wg.Done()
	// create resolver object to perform DNS queries.
	resolver := NewResolver(
		resolvers,
		timeout,
	)
	var localCount int
	var errors []error

	for hostname := range work {
		ipAddrs[hostname] = make(map[string]bool)
		ips, err := net.LookupHost(hostname)
		if err != nil {
			errors = append(errors, err)
		}

		for _, ipAddr := range ips {
			if !ipv4.MatchString(ipAddr) {
				continue
			}
			ipAddrs[hostname][ipAddr] = true
		}
		domain := util.HasSuffix(hostname, ".")
		rr, err := resolver.Resolve(domain)
		*countSuccess += 1
		if err != nil {
			errors = append(errors, err)
		}

		aRecords := util.ParseARecord(rr)
		for _, ip := range aRecords {
			ipAddrs[hostname][ip] = true
		}
		o.PrintStdout(ipAddrs, errors, verbose)
	}

	*countSuccess += localCount
	return ipAddrs
}
