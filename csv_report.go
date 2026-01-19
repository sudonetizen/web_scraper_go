package main

import (
	"os"
	"fmt"
	"strings"
	"encoding/csv"
)

func writeCSVReport(pages map[string]PageData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("writeCSVReport got err: %v", err)	
	}
	defer file.Close()

	records := [][]string{
		{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"},
	}

	for _, page := range pages {
		records = append(records, []string{page.URL, page.H1, page.FirstParagraph, strings.Join(page.OutgoingLinks, ";"), strings.Join(page.ImageURLs, ";")})
	}

	w := csv.NewWriter(file)

	for _, record := range records {
		err := w.Write(record)
		if err != nil {
		return fmt.Errorf("w.Write got err: %v", err)	
		}
	}
	
	w.Flush()

	return nil

}
