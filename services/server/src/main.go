package main

import (
	"fmt"
	"log"
)

const CRAWL_SEED = "https://markmcgranaghan.com"

func main() {
	var result, ok = Crawl(CRAWL_SEED)
	if !ok {
		log.Fatalln(":(")
	}

	fmt.Printf("%#v\n", result)
}
