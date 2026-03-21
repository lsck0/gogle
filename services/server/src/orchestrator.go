package main

import (
	"log"
	"maps"
	"slices"
	"time"

	"github.com/lsck0/gogle/src/collection"
)

const CRAWL_SEED = "https://reddit.com"

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
	const MAX_PARALLEL_CRAWLERS = 2

	urlQueue := collection.NewURLQueue()
	resultChannel := make(chan WebPage, MAX_PARALLEL_CRAWLERS)

	for range MAX_PARALLEL_CRAWLERS {
		go Crawler(urlQueue.GetStream(), resultChannel)
	}

	urlQueue.PushUrls(CRAWL_SEED)

	for result := range resultChannel {
		urlQueue.PushUrls(slices.Collect(maps.Keys(result.Links))...)
	}
}
