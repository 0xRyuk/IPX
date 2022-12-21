package global

import (
	"fmt"
	"regexp"
	"time"
)

const (
	// DefaultDNSServer is the default DNS server that the resolver will use if none is specified.
	DefaultDNSServer = "8.8.8.8:53"

	// DefaultTimeout is the default timeout for DNS queries.
	DefaultTimeout = 500 * time.Millisecond
)

// Variables to hold counting data on runtime.
var (
	Total int
	Count int
)

// IPX Options
var (
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
	Resolved      []string
)

// Color variables
var (
	BG     = "\033[47m"   //BG
	B      = "\033[0;90m" //Black
	G      = "\033[1;32m" //Green
	R      = "\033[1;31m" //Red
	Y      = "\033[1;33m" //Yellow
	W      = "\033[1;97m" //White
	IC     = "\033[0;96m" //High Intensty Cyan
	Rst    = "\033[1;m"   //Reset
	Banner = []string{fmt.Sprintf("\t%s  ____ __  %s_,  ,\n\t%s ( /( /  )%s( |,' \n\t%s  /  /--' %s  +   \n\t%s_/_ /    %s_,'|__ \n\t\t%sv0.1.0 %sbeta\n%s", IC, R, IC, R, IC, R, IC, R, W, Y, Rst),
		fmt.Sprintf("%s\t_ ___  %s_  _ \n\t%s| |__]  %s\\/\n\t%s| |    %s_/\\_ \n\t%s\tv0.1.0%s beta%s", IC, R, IC, R, IC, R, W, Y, Rst)} // Fancy banner(s)
)
