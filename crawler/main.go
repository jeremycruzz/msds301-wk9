// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	URL   string   `json:"url"`
	Links []string `json:"links"`
}

var visited = make(map[string]bool)
var mutex = &sync.Mutex{}

func main() {
	startURL := "https://en.wikipedia.org/wiki/Special:Random"
	results := make([]Page, 0)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go crawl(startURL, 2, wg, &results)

	wg.Wait()

	file, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		fmt.Printf("JSON marshalling failed: %s", err)
		return
	}

	err = ioutil.WriteFile("results.json", file, 0644)
	if err != nil {
		fmt.Printf("File writing failed: %s", err)
	}
}

func crawl(url string, depth int, wg *sync.WaitGroup, results *[]Page) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	mutex.Lock()
	if visited[url] {
		mutex.Unlock()
		return
	}
	visited[url] = true
	mutex.Unlock()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching URL %s: %s\n", url, err)
		return
	}
	defer resp.Body.Close()

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("Error parsing HTML for URL %s: %s\n", url, err)
		return
	}

	page := Page{URL: url}
	document.Find("a[href]").Each(func(index int, element *goquery.Selection) {
		link, exists := element.Attr("href")
		if exists {
			absoluteLink := resolveRelative(link, url)
			if isWikipediaArticle(absoluteLink) {
				page.Links = append(page.Links, absoluteLink)
				if depth > 1 {
					wg.Add(1)
					go crawl(absoluteLink, depth-1, wg, results)
				}
			}
		}
	})

	*results = append(*results, page)
}

func resolveRelative(link, base string) string {
	resolvedLink, err := resp.Request.URL.Parse(link)
	if err != nil {
		return ""
	}
	return resolvedLink.String()
}

func isWikipediaArticle(url string) bool {
	return true // Implement your logic to check if the URL is a Wikipedia article
}

// Note: This implementation lacks proper URL normalization and article check logic.
