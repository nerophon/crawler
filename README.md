# Crawler

An engineering exercise implemented in Go.

A simple web crawler that visits all pages within a given domain, but does not follow  external links. It outputs a simple structured site map, showing for each page:

1. domain-internal page links
2. external page links
3. links to static content such as images

This entire project can be cloned directly from github via:
https://github.com/nerophon/crawler

<br>
##Prerequisites

1. The [__Go Programming Langugage__][0] must be installed to build, test, and install this software.

##Installation

1. Clone this project.
2. `cd` to the project directory
3. run `go install`

The software will be installed to the `$GOPATH/bin` directory by default.

##Testing & Benchmarking

This software includes unit tests. They can be run as per standard for Go tests:

1. `cd` to the source folder with test files you wish to run
2. run `go test`

Benchmarks exist for key steps in the process. These can be run from the root project directory, via the `crawler_test.go` file. I suggest running each benchmark separately, using the following  commands:

```
go test -bench=BenchmarkFetch -benchtime=7s
go test -bench=BenchmarkCrawl -benchtime=15s
```

Please be aware that this kind of benchmark could, if run without care, be interpreted as a DOS attack. The `benchtime` flag may need to be adjusted depending upon which website is being used in the test. I strongly advise NOT using commonly DOS'd websites such as those belonging to major corporations.

##Launching

1. `cd` to the install directory, usually `$GOPATH/bin`
2. run `./crawler`

##Operation

At the application command prompt, the following commands are available:

```
crawl [URL]		begin crawling the specified domain
help			show available commands
quit			exit the application
```

Press `ctrl-c` during a crawl to halt and force quit back to the OS command line.


[0]: https://golang.org/dl/