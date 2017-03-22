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

##Testing

This software includes extensive unit tests. They can be run as per standard for Go tests:

1. `cd` to the source folder with test files you wish to run
2. run `go test`

##Launching

1. `cd` to the install directory, usually `$GOPATH/bin`
2. run `./crawler`

##Operation

TODO


[0]: https://golang.org/dl/