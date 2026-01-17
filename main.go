package main

import (
	"os"
	"fmt"
	"time"
	"sync"
	"net/url"
)


func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	} 
	
	pagesData := map[string]PageData{}
	//pagesMap := map[string]int{}
	baseURL := args[0]
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	/*fmt.Println("starting crawl of:", baseURL)
	startS := time.Now()

	crawlPage(baseURL, baseURL, pagesMap)

	elapsedS := time.Since(startS)
	fmt.Printf("time: %v\n", elapsedS)
	fmt.Println(len(pagesMap))

	*/


	maxConcurrency := 10
	cfg := &config{
		pages: pagesData,
		baseURL: parsedBaseURL,
		mu: &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg: &sync.WaitGroup{},
	}
	
	fmt.Println("starting crawl of:", baseURL)
	start := time.Now()

	cfg.wg.Add(1)
	go cfg.crawlPage(baseURL)
	cfg.wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("time: %v\n", elapsed)
	fmt.Println(len(pagesData))

	/*
	for k, v := range pagesData {
		fmt.Printf("%v - %v\n", k, v)
	}
	*/
}
