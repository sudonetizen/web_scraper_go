package main 

import (
	"io"
	"fmt"
	"strings"
	"net/url"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	req.Header.Set("User-Agent", "MiniCrawler/1.0")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("error: %v\n", res.Status)
	}

	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("error: %v\n", contentType)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}
		
	return string(body), nil
}

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

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	r := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		fmt.Errorf("error: %w", err)	
		return []string{}, err
	}

	foundURLs := []string{}
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")

		if exists {
			readyURL, err := baseURL.Parse(href)
			if err != nil {
				fmt.Errorf("error: %w", err)
			} else {
				foundURLs = append(foundURLs, readyURL.String())
			}
		}

	})

	return foundURLs, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	r := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		fmt.Errorf("error: %w", err)	
		return []string{}, err
	}

	foundURLs := []string{}
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")	

		if exists {
			readyURL, err := baseURL.Parse(src)
			if err != nil {
				fmt.Errorf("error: %w", err)
			} else {
				foundURLs = append(foundURLs, readyURL.String())
			}
		}

	})


	return foundURLs, nil
}

func extractPageData(html, pageURL string) PageData {
	h1 := getH1FromHTML(html)
	fp := getFirstParagraphFromHTML(html)

	baseURL, err := url.Parse(pageURL)
	if err != nil {
		return PageData{
			URL: pageURL,
			H1: h1,
			FirstParagraph: fp,
			OutgoingLinks: nil,
			ImageURLs: nil,
		}
	}

	outLinks, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		fmt.Errorf("error: %v", err)
		return PageData{}
	}

	imgURLs, err := getImagesFromHTML(html, baseURL)
	if err != nil {
		fmt.Errorf("error: %v", err)
		return PageData{}
	}

	return PageData{
		URL: pageURL,
		H1: h1,
		FirstParagraph: fp,
		OutgoingLinks: outLinks,
		ImageURLs: imgURLs,
	}
}
