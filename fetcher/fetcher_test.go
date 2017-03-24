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
		{"", nil, errors.New("url invalid: empty")},
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
			t.Errorf("NewWebFetch(%s), expected=%v, actual=%v", c.inURL, c.want, got)
		}
		if !reflect.DeepEqual(gotErr, c.wantError) {
			t.Errorf("NewWebFetch(%s), expectedError=%v, actualError=%v", c.inURL, c.wantError, gotErr)
		}
	}
}

func TestFixURL(t *testing.T) {
	//t.SkipNow()
	cases := []struct {
		inHref    string
		inBase    string
		want      string
		wantError error
	}{
		{"", "", "", errors.New("url invalid: empty")},
		{"http://www.google.com", "", "", errors.New("url invalid: empty")},
		{"", "http://www.google.com", "", errors.New("url invalid: empty")},
		{"http://slides.com", "http://www.google.com", "http://slides.com", nil},
		{"http://www.google.com", "http://www.google.com", "http://www.google.com", nil},
		{"/help", "http://www.google.com", "http://www.google.com/help", nil},
		{"/help", "http://google.com", "http://google.com/help", nil},
		{"/help?arg=some", "http://google.com", "http://google.com/help?arg=some", nil},
	}
	for _, c := range cases {
		got, gotErr := fixURL(c.inHref, c.inBase)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("FixURL(%s, %s), expected=%v, actual=%v", c.inHref, c.inBase, c.want, got)
		}
		if !reflect.DeepEqual(gotErr, c.wantError) {
			t.Errorf("FixURL(%s, %s), expectedError=%v, actualError=%v", c.inHref, c.inBase, c.wantError, gotErr)
		}
	}
}
