<h1 align="center">
  <br>
  <a href="https://github.com/0xRyuk/IPX"><img src="https://i.ibb.co/QNhS6k4/ipx-logo.png" alt="ipx-logo" width="250" border="0"></a>
</h1>

***IP eXtreme*** is an open-source DNS resolver tool that efficiently and quickly resolves domain names to their corresponding IP addresses. It leverages advanced resolving techniques and multiple concurrent threads to resolve DNS queries quickly and efficiently, making it well-suited for resolving large lists of domains.

The purpose of this project is to learn and gain practical experience with the Go programming language. The objective is to create a small project using Go that demonstrates an understanding of the language's syntax and core features.

The project is available on GitHub for easy collaboration and contribution.


## Features
- Fast and simple DNS resolution
- Multithreaded support for concurrent resolution of multiple hostnames
- Supports input from the command line, a file, or stdin
- Can save the output to a file
- Customizable timeout duration for DNS requests
- Custom resolver support
- Verbose mode with logging support

# Installation
`ipx` requires [go 1.19](https://go.dev/dl/) to install successfully. Run the following command to install the latest version:
```bash
go install -v github.com/0xRyuk/ipx/cmd/ipx@latest
```

# Building From Source
To build and run IPX locally, you'll need to have Go installed on your system. Once you have Go installed, full details of installation and setup can be found on the Go language [website](https://golang.org/doc/install).

You can clone the repository and run the following commands from the project directory:

```bash
git clone https://github.com/0xRyuk/ipx.git
```
**Note:** go v1.19 required to compile ipx binary.

## Install dependencies
ipx has external dependencies, so they need to be pulled in first.
```bash
go get -d -v ./...
```
## Compiling
```bash
go build -o bin/ ./cmd/ipx
```
This will build an ipx binary for you under `./bin` directory. If you want to install it in the $GOPATH/bin directory you can run

```bash
go install ./cmd/ipx
```
## Run IPX
```bash
./bin/ipx -f hostnames.txt
```
## Usage
```bash
ipx -h
```
This will display help usage. Here are all the options that IPX currently supports.

```console
Usage of ipx:
  -d string
        Set hostname to resolve (i.e. example.com)
  -f string
        Read a file containing hostnames to resolve (i.e. hosts.txt)
  -i Print only IP address (default false)
  -o string
        Save output to a text file
  -r string
        Resolvers list (i.e. resolvers.txt)
  -t int
        Number of threads to utilize (default 20)
  -timeout duration
        Set timeout (default 500ms)
  -v    Turn on verbose mode (default off)
```

To use IPX, you can either provide a hostname as a command-line argument or specify a file containing a list of hostnames to be resolved.

To resolve a single hostname, use the `-d` flag followed by the hostname:
```bash
ipx -d example.com
```
```console
        _ ___  _  _        
        | |__]  \/
        | |    _/\_        
                v0.1.0 beta

93.184.216.34
```

To resolve a list of hostnames from a file, use the `-f` flag followed by the path to the file:
```bash
ipx -f test_hostnames.txt
```
```console
        _ ___  _  _        
        | |__]  \/
        | |    _/\_        
                v0.1.0 beta

93.184.216.34
74.6.231.20
74.6.143.25
74.6.143.26
98.137.11.164
98.137.11.163
74.6.231.21
142.250.195.14
142.250.206.110
104.16.99.52
104.16.100.52
142.250.206.110
142.250.194.110
157.240.198.35
```
Users also can supply input using stdin

By default, IPX uses 20 threads for concurrent DNS resolution. You can change the number of threads used with the `-t` flag, which is useful for a large number of hostnames (tested up to 2000 Threads with good accuracy):
```bash
ipx -t 50 -f large_hostnames.txt
```
IPX also allows you to specify a custom timeout duration for DNS requests using the `-timeout` flag. The timeout is specified as a duration (e.g. 5s, 1m, etc.), and the default timeout is 5 seconds.
```bash
ipx -timeout 10s -f test_hostnames.txt
```
If you want to print IPX's verbose output, you can use the `-v` flag to enable verbose mode. In verbose mode, IPX will output the resolved IP addresses along with additional information in log format.
```bash
ipx -v -f test_hostnames.txt
```
The `-i` flag is a boolean flag that specifies whether to print only the IP address or not. When the -i flag is set to true, the program will only print the IP address of the hostname being resolved. This can be useful if you want to pipe the output of the program to another tool and only need the IP address.
```bash
ipx -i -f test_hostnames.txt | nmap -iL -
```
To save the output to a file in the desired format, you can use the `--format` flag with a value of _json_, _csv_, or _text_, and the `-o` flag to specify the output file. This is useful if you want to save the resolved IP addresses for later use, or if you want to redirect the output to another program for further processing.
```bash
ipx -f test_hostnames.txt -o resolved_ips --format=json
```
Finally, to use stdin or piped input with IPX, you can specify the `-` flag. This will tell the program to read input from the stdin or the pipe.

For example, to pipe a list of hostnames to the program, you can use the echo or cat command and the | (pipe) operator like this:
```bash
cat "test_hostnames.txt" | ipx -
```
You can also specify additional flags when using stdin or piped input. For example, to turn on verbose mode and print only the IP addresses, you can run the program like this:
```bash
cat test_hostnames.txt | ipx -v -t 50 -timeout 10s -o resolved_ips.txt
```
```console
        _ ___  _  _        
        | |__]  \/
        | |    _/\_        
                v0.1.0 beta

INFO[19:45:14] 93.184.216.34
INFO[19:45:14] 74.6.231.21
INFO[19:45:14] 142.250.193.46
INFO[19:45:14] 142.250.182.174
INFO[19:45:14] 98.137.11.163
INFO[19:45:14] 74.6.231.20
INFO[19:45:14] 74.6.143.25
INFO[19:45:14] 74.6.143.26
INFO[19:45:14] 98.137.11.164
INFO[19:45:14] 157.240.239.35
INFO[19:45:14] 104.16.99.52
INFO[19:45:14] 104.16.100.52
INFO[19:45:14] 157.240.198.35
INFO[19:45:14] 142.250.206.110
INFO[19:45:14] 216.58.221.46
INFO[19:45:14] Finished in 613.1446ms
INFO[19:45:14] Total domains resolved 6/6
INFO[19:45:14] 15 IP address(s) found
```

## Contributing
If you'd like to contribute to IPX, you can fork the repository and submit a pull request. contributions are welcome, whether they are bug fixes, improvements to the code, or new features.

### License
IPX is released under the MIT License.