package main

import (
	"log"
	"time"
)

const CRAWL_SEED = "https://google.com"

func Crawler(urlChannel <-chan string, resultChannel chan<- WebPage) {
	for url := range urlChannel {
		var result, ok = ReadWebPage(url)
		if ok {
			resultChannel <- result
			log.Println("Crawled: " + result.Url)
		}
		time.Sleep(5 * time.Second)
	}
}

func RunOrchestrator() {
	var seenUrls = make(map[string]bool)

	var urlChannel = make(chan string)
	var resultChannel = make(chan WebPage, 25)

	for range 25 {
		go Crawler(urlChannel, resultChannel)
	}

	urlChannel <- CRAWL_SEED

	for result := range resultChannel {
		for url := range result.Links {
			if !seenUrls[url] {
				urlChannel <- url
				seenUrls[url] = true
			}
		}
	}
}
