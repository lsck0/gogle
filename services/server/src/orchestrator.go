package main

import (
	"context"
	"log"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
)

const crawlSeedURL = "https://google.com"

func crawler(urlCh <-chan string, resCh chan<- WebPage) {
	for url := range urlCh {
		res, err := FetchWebPage(url)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("Succeeded crawling %q", url)
		go func() { resCh <- res }() // Send res from separate goroutine to prevent deadlock

		time.Sleep(5 * time.Second)
	}
}

func RunOrchestrator() {
	ctx := context.Background()

	var urlCh = make(chan string)
	var resCh = make(chan WebPage)

	for range 10 {
		go crawler(urlCh, resCh)
	}
	urlCh <- crawlSeedURL

	seen := map[string]bool{crawlSeedURL: true}
	for res := range resCh {
		for url := range res.Links {
			if !seen[url] {
				seen[url] = true
				go func() { urlCh <- url }()
			}
		}

		obj, _ := res.toModel()
		obj.InsertG(ctx, boil.Infer())
	}
}
