package main

import (
	"fmt"
	"net/http"
	"strings"

	"os"

	"golang.org/x/net/html"
)

// Pull href attribute from a token
func getHref(t html.Token) (href string, ok bool) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return // returns the defined variables
}

// Scrape all http links from a page
func scrape(url string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)
	defer func() {
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		return
	}

	defer b.Close()

	b := resp.Body

	z := html.NewTokenizer(b)

	for {
		nt := z.Next()

		switch {
		case nt == html.ErrorToken:
			return
		case nt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			url, ok := getHref(t)
			if !ok {
				continue
			}

			hasProto := strings.Index(url, "http") == 0
			if hasProto {
				ch <- url
			}
		}
	}
}

func main() {
	foundUrls := make(map[string]bool)
	seedUrls := os.Args[1:]

	chUrls := make(chan string)
	chFinished := make(chan bool)

	for _, url := range seedUrls {
		go scrape(url, chUrls, chFinished)
	}

	for c := 0; c < len(seedUrls); {
		select {
		case url := <-chUrls:
			foundUrls[url] = true
		case <-chFinished:
			c++
		}
	}

	fmt.Println("\nFound", len(foundUrls), "unique urls:")
	for url := range foundUrls {
		fmt.Println(" - " + url)
	}

	close(chUrls)

}
