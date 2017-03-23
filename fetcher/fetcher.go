/*
Package fetcher downloads and parses web pages into sets of links.
*/
package fetcher

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

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
	wf.url = url
	resp, err := http.Get(url)
	if err != nil {
		wf.err = err
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		wf.err = err
		return
	}
	links := collectlinks.All(resp.Body)
	for _, link := range links {
		fmt.Println(link)
	}
	wf.err = errors.New("mock")
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
