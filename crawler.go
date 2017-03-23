/*
A simple web crawler that visits all pages within
a given domain, but does not follow external links.
It outputs a simple structured site map, showing for each page:
1. domain-internal page links
2. external page links
3. links to static content such as images
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nerophon/crawler/fetcher"
)

func main() {
	fmt.Print("\nWelcome to Crawler!\n")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter a command.\n" +
		"Typing \"help\" will show a command list.\n\n")
	commandLoop(reader)
}

func commandLoop(reader *bufio.Reader) {
	fmt.Print("> ")
	text, error := reader.ReadString('\n')
	if error != nil || len(text) == 0 {
		fmt.Printf("Sorry, there was an input error:\n%d\n", error)
		return
	}
	trimmed := strings.Trim(text, "\n ")
	fields := strings.Fields(trimmed)
	numFields := len(fields)
	if numFields <= 0 {
		commandLoop(reader)
		return
	}
	switch fields[0] {
	case "quit":
		fmt.Println("Quitting.")
		fmt.Println("")
		return
	case "help":
		fmt.Println("")
		fmt.Println("***Command List***")
		fmt.Println("")
		fmt.Println("\nhelp\t\tshows this command list")
		fmt.Println("crawl [URL]\tcrawls a single domain starting at [URL]")
		fmt.Println("quit\t\tquits the application")
		fmt.Println("")
	case "crawl":
		if numFields == 2 {
			crawlCommand(fields[1])
		} else {
			fmt.Println("Please specify only one URL.")
			fmt.Println("")
		}
	default:
		fmt.Println("Sorry, command not understood.")
		fmt.Println("")
	}
	commandLoop(reader)
}

// Fetcher takes a URL, downloads the page body
// and parses it into a structure of links,
// filtering out duplicate & self links
type Fetcher interface {
	Fetch(url string)
	URL() string
	InternalLinks() []string
	ExternalLinks() []string
	ResourceLinks() []string
	Err() (e error)
}

func fetch(url string, c chan Fetcher) {
	var f = fetcher.New() // single point of contact
	f.Fetch(url)
	c <- f
}

func crawlCommand(url string) {
	var pages = make(map[string]*Page)
	var newLinks = []string{url}
	// loops until return
	for {
		result := crawl(pages, newLinks)
		if result.err != nil {
			fmt.Fprintln(os.Stderr, "\nSorry, there was an error:\n", result.err)
			fmt.Println("")
			return
		}
		newLinks = result.newLinks
		if len(newLinks) <= 0 {
			//no more links to crawl, so print results
			fmt.Println("Success!")
			fmt.Printf("\n%v\n", pages)
			fmt.Println("")
			return
		}
	}
}

// Page stores site map metadata about a single web page
type Page struct {
	url           string
	internalLinks []string
	externalLinks []string
	resourceLinks []string
}

// CrawlResult encapsulates the result of a crawl
type CrawlResult struct {
	newLinks []string
	err      error
}

// crawl over a set of urls, adding page metadata into specified pages map;
// returns a new set of urls to crawl, filtering out those already in the map
func crawl(pages map[string]*Page, urls []string) (r CrawlResult) {
	c := make(chan Fetcher)
	for _, link := range urls {
		go fetch(link, c) // woof! woof!
	}
	foundLinks := make([]string, 0)
	for i := 0; i < len(urls); i++ {
		f := <-c
		if f.Err() != nil {
			r.err = f.Err()
			return r
		}
		newPage := &Page{f.URL(), f.InternalLinks(),
			f.ExternalLinks(), f.ResourceLinks()}
		pages[f.URL()] = newPage
		// combine fetchResult links
		foundLinks = append(foundLinks, newPage.internalLinks...)
	}
	//now filter out already crawled pages
	for _, link := range foundLinks {
		_, exists := pages[link]
		if !exists {
			r.newLinks = append(r.newLinks, link)
		}
	}
	return r
}
