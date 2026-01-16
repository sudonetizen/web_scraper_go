package main

import (
	"fmt"
	"os"
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
	
	baseURL := args[0]
	fmt.Println("starting crawl of:", baseURL)

	pages := map[string]int{}
	crawlPage(baseURL, baseURL, pages)

	fmt.Println(len(pages))

	for k, v := range pages {
		fmt.Printf("%v - %v\n", k, v)
	}
}
