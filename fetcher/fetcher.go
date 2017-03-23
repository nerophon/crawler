/*
Package fetcher downloads and parses web pages into sets of links.
*/
package fetcher

import (
	"crypto/tls"
	"fmt"
	"net/http"
	pkgurl "net/url"
	"strings"

	"github.com/jackdanger/collectlinks"
)

// WebFetch encapsulates a Fetch operation and its result
type WebFetch struct {
	url           string
	internalLinks []string
	externalLinks []string
	resourceLinks []string
	err           error
}

// New is a convenience function to return a new Fetch struct
func New() *WebFetch {
	return new(WebFetch)
}

// Fetch takes a URL, downloads the page body
// and parses it into a structure of links,
// filtering out duplicate & self links
func (wf *WebFetch) Fetch(url string) {
	// validate parent
	parent, err := pkgurl.Parse(url)
	if err != nil {
		wf.err = err
		return
	}
	wf.url = url

	// disable security since this is just a test app
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := http.Client{Transport: transport}

	// download the html
	resp, err := client.Get(url)
	if err != nil {
		wf.err = err
		return
	}
	defer resp.Body.Close()

	// use an open-source library to parse it for links
	links := collectlinks.All(resp.Body)

	// distribute the links into desired categories
	for _, link := range links {
		//filter out dupes and blanks
		for _, v := range wf.internalLinks {
			if link == v {
				continue
			}
		}
		for _, v := range wf.externalLinks {
			if link == v {
				continue
			}
		}
		for _, v := range wf.resourceLinks {
			if link == v {
				continue
			}
		}
		if link == "" {
			continue
		}

		// run rules
		child, err := pkgurl.Parse(link)
		if err != nil {
			fmt.Printf("\nWarning: failed to validate a link as a proper url:\n%v\n", err)
			continue
		}
		//fmt.Println(link)

		// compare hosts
		if !strings.Contains(child.Host, parent.Host) && !strings.Contains(parent.Host, child.Host) {
			wf.externalLinks = append(wf.externalLinks, link)
			//fmt.Println("external")
			continue
		}

		// filter by scheme
		if child.Scheme != "http" && child.Scheme != "https" {
			wf.resourceLinks = append(wf.resourceLinks, link)
			//fmt.Println("scheme not http(s), resource!")
			continue
		}

		// check content type by fetching the header
		// TODO this is the bottleneck of the crawler, and should be improved
		hdResp, err := client.Head(link)
		if err != nil {
			fmt.Printf("\nWarning: failed to make a HEAD request to a link:\n%v\n", err)
			continue
		}
		defer hdResp.Body.Close()
		contentType := hdResp.Header.Get("Content-Type")
		//fmt.Printf("Content-Type = %v", contentType)
		if strings.Contains(contentType, "text/html") {
			wf.internalLinks = append(wf.internalLinks, link)
			continue
		} else {
			wf.resourceLinks = append(wf.resourceLinks, link)
			continue
		}

	}
}

// URL returns the url of the fetch operation
func (wf *WebFetch) URL() string {
	return wf.url
}

// InternalLinks gets the domain-internal link list
// from a successful fetch operation, nil otherwise.
func (wf *WebFetch) InternalLinks() []string {
	return wf.internalLinks
}

// ExternalLinks gets the domain-external link list
// from a successful fetch operation, nil otherwise.
func (wf *WebFetch) ExternalLinks() []string {
	return wf.externalLinks
}

// ResourceLinks gets the resource link list
// from a successful fetch operation, nil otherwise.
func (wf *WebFetch) ResourceLinks() []string {
	return wf.resourceLinks
}

// Err returns the error from a fetch operation, if there was one.
func (wf *WebFetch) Err() error {
	return wf.err
}
