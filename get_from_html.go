package main 

import (
	"io"
	"fmt"
	"sync"
	"strings"
	"net/url"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	_, ok := cfg.pages[normalizedURL]
	if ok {
		return false
	} 

	return true
}

func (cfg *config) addPage(normalizedURL, rawCurrentURL, htmlString string) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.pages[normalizedURL] = extractPageData(htmlString, rawCurrentURL)
}

func (cfg *config) crawlPage(rawCurrentURL string) error {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()


	parsedCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}

	if cfg.baseURL.Hostname() != parsedCurrentURL.Hostname() {
		return fmt.Errorf("base url: %v, current url: %v\n", cfg.baseURL.Hostname(), parsedCurrentURL)
	}


	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	isFirst := cfg.addPageVisit(normalizedURL)
	if !isFirst { return fmt.Errorf("visited\n")}

	htmlString, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	cfg.addPage(normalizedURL, rawCurrentURL, htmlString)

	allURLs, err := getURLsFromHTML(htmlString, parsedCurrentURL)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	allURLs2 := []string{}

	for _, url := range allURLs {
		if !strings.HasSuffix(url, ".xml") {
			allURLs2 = append(allURLs2, url)	
		}
	}


		
	for _, url := range allURLs2 {
		cfg.wg.Add(1)
		go func(url string) {
			err := cfg.crawlPage(url)
			if err != nil {return}
			fmt.Printf("crawled: %v\n", url)
		}(url)
	}

	return nil
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) error {
	parsedBaseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}

	parsedCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err	
	}

	if parsedBaseURL.Hostname() != parsedCurrentURL.Hostname() {
		return fmt.Errorf("base url: %v, current url: %v", parsedBaseURL, parsedCurrentURL)
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	_, ok := pages[normalizedURL]
	if ok {
		pages[normalizedURL] += 1
		return nil
	} 
	
	pages[normalizedURL] = 1

	htmlString, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	allURLs, err := getURLsFromHTML(htmlString, parsedBaseURL)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	for _, url := range allURLs {
		err := crawlPage(rawBaseURL, url, pages)
		if err != nil {
			continue
		}

		fmt.Printf("crawled: %v\n", url)
	}

	return nil

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
