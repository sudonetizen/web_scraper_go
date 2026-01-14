package main 

import (
	"fmt"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

func getH1FromHTML(html string) string {
	r := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		fmt.Errorf("error: %w", err)	
		return ""
	}

	h1 := doc.Find("h1").First().Text()
	return strings.TrimSpace(h1)

}

func getFirstParagraphFromHTML(html string) string {
	r := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		fmt.Errorf("error: %w", err)	
		return ""
	}

	docMain := doc.Find("main")
	if docMain.Length() != 0 {
		return docMain.Find("p").First().Text()
	}

	docP := doc.Find("p").First().Text()
	return strings.TrimSpace(docP)
}
