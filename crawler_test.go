package main

import (
	"testing"

	"github.com/nerophon/crawler/fetcher"
)

var url = "http://www.wiprodigital.com"

func TestSanity(t *testing.T) {
	if 1 != 1 {
		t.Errorf("TestSanity(), cosmic rays changing bits in yo RAM!")
	}
}

func BenchmarkFetch(b *testing.B) {
	// run the Fetch b.N times
	for n := 0; n < b.N; n++ {
		f := fetcher.New()
		f.Fetch(url)
	}
}

func BenchmarkCrawl(b *testing.B) {
	// run the Crawl b.N times
	for n := 0; n < b.N; n++ {
		crawlCommand(url)
	}
}
