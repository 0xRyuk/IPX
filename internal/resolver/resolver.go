package resolver

import (
	"fmt"
	"time"

	miekdns "github.com/miekg/dns"
)

// Resolver is a DNS resolver client that can work with multiple DNS servers.
type Resolver struct {
	servers []string
	timeout time.Duration
}

// NewResolver creates a new Resolver with the given list of DNS servers and timeout.
func NewResolver(servers []string, timeout time.Duration) *Resolver {
	return &Resolver{
		servers: servers,
		timeout: timeout,
	}
}

// Resolve performs a DNS query using the specified DNS servers and returns the response.

func (r *Resolver) Resolve(d string) ([]miekdns.RR, error) {

	q := miekdns.Question{Name: d, Qtype: miekdns.TypeA, Qclass: miekdns.ClassINET}
	c := new(miekdns.Client)
	c.Timeout = r.timeout

	var results []miekdns.RR

	// Perform the DNS query on each server in the list.
	for _, server := range r.servers {
		in, _, err := c.Exchange(&miekdns.Msg{
			MsgHdr: miekdns.MsgHdr{
				RecursionDesired: true,
			},
			Question: []miekdns.Question{q},
		}, server)

		if err != nil {
			// If there was an error, try the next server in the list.
			continue
		}

		// If we received a successful response, add the results to the list.
		if in != nil && in.Rcode == miekdns.RcodeSuccess {
			results = append(results, in.Answer...)
		}
	}

	// If we didn't get any results from the servers, return an error.
	if len(results) == 0 {
		return nil, fmt.Errorf("no successful response from DNS servers")
	}

	// Create a map to track which results we have seen.
	seen := make(map[string]bool)

	var finalResults []miekdns.RR
	for _, result := range results {
		// Create a string representation of the result.
		resultStr := result.String()

		// If we haven't seen this result before, add it to the final results slice.
		if !seen[resultStr] {
			finalResults = append(finalResults, result)

			// Mark the result as seen.
			seen[resultStr] = true
		}
	}

	return finalResults, nil
}
