package main

import (
	"log"
	"time"
)

const crawlSeedURL = "https://google.com"

func crawler(urlCh <-chan string, resCh chan<- WebPage) {
	for url := range urlCh {
		res, ok := ReadWebPage(url)
		if !ok {
			log.Printf("Failed crawling %q", url)
			continue
		}

		log.Printf("Succeeded crawling %q", url)
		go func() { resCh <- res }() // Send res from separate goroutine to prevent deadlock

		time.Sleep(5 * time.Second)
	}
}

func RunOrchestrator() {
	var urlCh = make(chan string)
	var resCh = make(chan WebPage)

	for range 25 {
		go crawler(urlCh, resCh)
	}
	urlCh <- crawlSeedURL

	seen := map[string]bool{crawlSeedURL: true}
	for res := range resCh {
		for url := range res.Links {
			if !seen[url] {
				seen[url] = true
				urlCh <- url
			}
		}
	}
}
