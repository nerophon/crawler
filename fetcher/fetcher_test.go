package fetcher

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNewWebFetch(t *testing.T) {
	//t.SkipNow()
	url1 := "http://www.google.com"
	urlStruct1, _ := url.Parse(url1)
	url2 := "www.google.com"
	//urlStruct2, _ := url.Parse(url2)
	// url3 := "google.com"
	// urlStruct3, _ := url.Parse(url3)
	// url4 := ".com"
	// urlStruct4, _ := url.Parse(url4)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{Transport: transport}
	cases := []struct {
		inURL     string
		want      *WebFetch
		wantError error
	}{
		{"", nil, errors.New("error")},
		{url1, &WebFetch{url1, urlStruct1, client,
			"", map[string]bool{}, map[string]bool{},
			map[string]bool{}, map[string]bool{}, nil}, nil},
		{url2, nil, errors.New("url invalid: no scheme")},
		// {url3, &WebFetch{url3, urlStruct3, client,
		// 	"", map[string]bool{}, map[string]bool{},
		// 	map[string]bool{}, map[string]bool{}, nil}, nil},
		// {url4, &WebFetch{url4, urlStruct4, client,
		// 	"", map[string]bool{}, map[string]bool{},
		// 	map[string]bool{}, map[string]bool{}, nil}, nil},
	}
	for _, c := range cases {
		got, gotErr := New(c.inURL)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("NewWebFetch(%s), expected=%v, actual=%v, expectedError=%v, actualError=%v", c.inURL, c.want, got, c.wantError, gotErr)
		}
	}

	// TODO test FetchHeader, Fetch, fixURL
}
