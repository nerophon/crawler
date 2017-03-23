/*
Package fetcher downloads and parses web pages into sets of links.
*/
package fetcher

import (
	"crypto/tls"
	"fmt"
	"net/http"
	pkgurl "net/url"

	"github.com/jackdanger/collectlinks"
)

// WebFetch encapsulates a Fetch operation and its result
type WebFetch struct {
	url           string
	urlStruct     *pkgurl.URL
	client        http.Client
	contentType   string
	allLinks      map[string]bool
	internalLinks map[string]bool
	externalLinks map[string]bool
	resourceLinks map[string]bool
	err           error
}

// URL returns the url of the fetch operation
func (wf *WebFetch) URL() string {
	return wf.url
}

// ContentType returns the content-type, if headers have been fetched
func (wf *WebFetch) ContentType() string {
	return wf.contentType
}

// AllLinks gets the complete uncategorised link list
// from a successful fetch operation, nil otherwise.
func (wf *WebFetch) AllLinks() map[string]bool {
	return wf.allLinks
}

// InternalLinks gets the domain-internal link list
// from a successful fetch operation, nil otherwise.
func (wf *WebFetch) InternalLinks() map[string]bool {
	return wf.internalLinks
}

// ExternalLinks gets the domain-external link list
// from a successful fetch operation, nil otherwise.
func (wf *WebFetch) ExternalLinks() map[string]bool {
	return wf.externalLinks
}

// ResourceLinks gets the resource link list
// from a successful fetch operation, nil otherwise.
func (wf *WebFetch) ResourceLinks() map[string]bool {
	return wf.resourceLinks
}

// Err returns the error from a fetch operation, if there was one.
func (wf *WebFetch) Err() error {
	return wf.err
}

// New creates and initializes a new Fetch struct
func New(url string) (*WebFetch, error) {
	wf := new(WebFetch)

	// validate url
	urlStruct, err := pkgurl.Parse(url)
	if err != nil {
		return nil, err
	}
	wf.url = url
	wf.urlStruct = urlStruct

	// disable security since this is just a test app
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	wf.client = http.Client{Transport: transport}

	return wf, nil
}

// FetchHeader downloads the page header and stores the content-type
func (wf *WebFetch) FetchHeader() {
	resp, err := wf.client.Head(wf.url)
	if err != nil {
		fmt.Printf("\nWarning: failed to make a HEAD request to a url:\n%v\n", err)
		return
	}
	defer resp.Body.Close()
	wf.contentType = resp.Header.Get("Content-Type")
}

// Fetch downloads the page body and parses it into a structure of links
func (wf *WebFetch) Fetch() {

	// download the html
	resp, err := wf.client.Get(wf.url)
	if err != nil {
		wf.err = err
		return
	}
	defer resp.Body.Close()

	// use an open-source library to parse it for links
	allLinks := collectlinks.All(resp.Body)

	// remove dupes, blanks, and malforms
	for _, newLink := range allLinks {
		if newLink == "" {
			continue
		}
		_, err := pkgurl.Parse(newLink)
		if err != nil {
			fmt.Printf("\nWarning: failed to validate a link as a proper url:\n%v\n", err)
			continue
		}
		wf.allLinks[newLink] = true
	}
}

// CategoriseLinkAsInternal puts the specified link into the internal category
func (wf *WebFetch) CategoriseLinkAsInternal(link string) {
	wf.internalLinks[link] = true
}

// CategoriseLinkAsExternal puts the specified link into the external category
func (wf *WebFetch) CategoriseLinkAsExternal(link string) {
	wf.externalLinks[link] = true
}

// CategoriseLinkAsResource puts the specified link into the resource category
func (wf *WebFetch) CategoriseLinkAsResource(link string) {
	wf.externalLinks[link] = true
}
