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
	pkgurl "net/url"
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

func crawlCommand(url string) {
	urlStruct, err := pkgurl.Parse(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "\nSorry, there was an error:\n", err)
		fmt.Println("")
		return
	}
	var domain = urlStruct.Host
	var pages = make(map[string]*Page)
	var newHTMLUrls = make(map[string]bool)
	newHTMLUrls[url] = true
	var knownResourceUrls = make(map[string]bool)

	// loop until return
	for {
		result := crawl(domain, pages, newHTMLUrls, knownResourceUrls)
		if result.err != nil {
			fmt.Fprintln(os.Stderr, "\nSorry, there was an error:\n", result.err)
			fmt.Println("")
			return
		}
		newHTMLUrls = result.newHTMLUrls
		knownResourceUrls = result.knownResourceUrls
		if len(newHTMLUrls) <= 0 {
			//no more links to crawl, so print results
			fmt.Println("Success!")
			fmt.Println("")
			printPageMap(pages)
			fmt.Println("")
			return
		}
	}
}

func printPageMap(pages map[string]*Page) {
	for _, page := range pages {
		printPage(page)
	}
}

func printPage(page *Page) {
	fmt.Printf("URL: %v\n", page.url)
	fmt.Printf("\tInternal Links: \n")
	printList(page.internalLinks)
	fmt.Printf("\tExternal Links: \n")
	printList(page.externalLinks)
	fmt.Printf("\tResource Links: \n")
	printList(page.resourceLinks)
}

func printList(list map[string]bool) {
	for key := range list {
		fmt.Printf("\t\t%v\n", key)
	}
}

// Page stores site map metadata about a single web page
type Page struct {
	url           string
	internalLinks map[string]bool
	externalLinks map[string]bool
	resourceLinks map[string]bool
}

// CrawlResult encapsulates the result of a crawl
type CrawlResult struct {
	newHTMLUrls       map[string]bool
	knownResourceUrls map[string]bool
	err               error
}

// crawl over a set of urls, adding page metadata into specified pages map;
// assumes all urls in the list have not yet been crawled;
// returns a new set of urls to crawl, filtering out those already in the map
func crawl(domain string, pages map[string]*Page, urls map[string]bool, knownResourceUrls map[string]bool) (res CrawlResult) {

	// download pages and parse into links
	c := make(chan Fetcher)
	for url := range urls {
		go fetch(url, c) // woof! woof!
	}

	// collate links
	foundLinks := make(map[string][]Fetcher)
	foundFetchers := make(map[string]Fetcher)
	for i := 0; i < len(urls); i++ {
		fetcher := <-c
		if fetcher.Err() != nil {
			res.err = fetcher.Err()
			return res
		}
		for link := range fetcher.AllLinks() {
			foundLinks[link] = append(foundLinks[link], fetcher)
		}
		foundFetchers[fetcher.URL()] = fetcher
	}

	// assess links
	remainingLinks := make(map[string][]Fetcher)
	for link, fetchers := range foundLinks {

		// we don't have to reasses links that have already been crawled
		// we know they are internal links and page structs exist for them
		_, exists := pages[link]
		if exists {
			for _, fetcher := range fetchers {
				fetcher.CategoriseLinkAsInternal(link)
			}
			continue
		}

		// put the link url into a structure for analysis
		linkStruct, _ := pkgurl.Parse(link)

		// compare hosts to see if the link is external
		if !strings.Contains(linkStruct.Host, domain) && !strings.Contains(domain, linkStruct.Host) {
			for _, fetcher := range fetchers {
				fetcher.CategoriseLinkAsExternal(link)
			}
			continue
		}

		// filter by scheme to match some (but not all) resources
		if linkStruct.Scheme != "http" && linkStruct.Scheme != "https" {
			for _, fetcher := range fetchers {
				fetcher.CategoriseLinkAsResource(link)
			}
			knownResourceUrls[link] = true
			continue
		}

		// unfortunately the only way to discriminate further
		// is to fetch the headers and check the content-type
		remainingLinks[link] = fetchers
		go fetchHeader(link, c) // woof! woof!
	}

	// reach back into the channel to retrieve fetchers with headers only
	newHTMLUrls := make(map[string]bool)
	for i := 0; i < len(remainingLinks); i++ {
		headFetcher := <-c
		if headFetcher.Err() != nil {
			// ignore link, nothing we can do
			continue
		}
		// finally, categorise by content-type
		if strings.Contains(headFetcher.ContentType(), "text/html") {
			for _, fetcher := range remainingLinks[headFetcher.URL()] {
				fetcher.CategoriseLinkAsInternal(headFetcher.URL())
				newHTMLUrls[headFetcher.URL()] = true // for next iteration
			}
			continue
		} else {
			for _, fetcher := range remainingLinks[headFetcher.URL()] {
				fetcher.CategoriseLinkAsResource(headFetcher.URL())
				knownResourceUrls[headFetcher.URL()] = true
			}
			continue
		}
	}

	// categorisation is complete, now we need to make a page for each found fetcher
	// and add it to the page map
	for url, fetcher := range foundFetchers {
		newPage := &Page{url, fetcher.InternalLinks(),
			fetcher.ExternalLinks(), fetcher.ResourceLinks()}
		pages[url] = newPage
	}

	res.newHTMLUrls = newHTMLUrls             // this tells us what to crawl next
	res.knownResourceUrls = knownResourceUrls // this saves us refetching headers
	return res
}

// Fetcher takes a URL, downloads the page body
// and parses it into a structure of links,
// filtering out duplicate & self links
type Fetcher interface {
	URL() string
	ContentType() string
	AllLinks() map[string]bool
	InternalLinks() map[string]bool
	ExternalLinks() map[string]bool
	ResourceLinks() map[string]bool
	Err() error

	FetchHeader()
	Fetch()
	CategoriseLinkAsInternal(string)
	CategoriseLinkAsExternal(string)
	CategoriseLinkAsResource(string)
}

func fetch(url string, c chan Fetcher) {
	f, _ := fetcher.New(url) // single point of contact
	// TODO handle error
	f.Fetch()
	c <- f
}

func fetchHeader(url string, c chan Fetcher) {
	f, _ := fetcher.New(url) // single point of contact
	// TODO handle error
	f.FetchHeader()
	c <- f
}
