package main

import "fmt"


func main() {
	fmt.Println("Hello, World!")

	_, err := normalizeURL("https://blog.boot.dev/path/")
	if err != nil {
		fmt.Println("not worked")
	}

	input_body := "<html><body><h1>Test Title</h1></body></html>"
	header1 := getH1FromHTML(input_body)
	fmt.Println(header1)
}
