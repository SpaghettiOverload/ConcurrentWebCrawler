package web

import (
	"ConcurrentWebCrawler/helpers"
	"ConcurrentWebCrawler/structs"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// pushToAGoroutine takes URL and pass it to Crawl function for further processing.
//Results (a page struct and a slice of new URLs) are pushed to the results structure.
func pushToAGoroutine(goUrl string, results *structs.Results, wg *sync.WaitGroup, workDone <- chan struct{}) {
	p, l, err := Crawl(goUrl)
	if err == nil {
		results.UpdatePages(&p)
		results.UpdateLinks(&l)
	}
	<- workDone
	wg.Done()
}

// CrawlUrls reads the currently available URLs from the workingList, process each of them by pushing it into a
//goroutine where the Crawl function extracts data, then top-up the workingList with any new links.
// This process continues until maxResults is reached. It returns a slice of structures from every parsed URL.
func CrawlUrls(initUrl string, maxRoutines int, maxResults int) []structs.WebPage {
	var wg sync.WaitGroup
	maxGoroutines := make(chan struct{}, maxRoutines)
	defer close(maxGoroutines)
	seen := make(map[string]bool)
	workingList := []string{initUrl}
	results := structs.Results{SP:make([]structs.WebPage,0), L: make([]string, 0)}

	for len(workingList) > 0 && maxResults > 0 {
		for _, link := range workingList {

			if !seen[link] {
				seen[link] = true
				maxResults--

				wg.Add(1)
				maxGoroutines <- struct{}{}
				go pushToAGoroutine(link, &results, &wg, maxGoroutines)

				if maxResults == 0 {
					break
				}
			}
		}
		wg.Wait()
		workingList = results.L
		results.CleanLinks()
	}
	return results.SP
}

// isUrl validates whether an extracted href address is a valid URL.
func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// evaluateWeight check tags and returns proper weight for the content.
func evaluateWeight(tag string) float64 {
	weight := 1.0
	switch tag {
	case "h1":
		weight = 4
	case "h2":
		weight = 3.5
	case "h3":
		weight = 3
	case "h4":
		weight = 2.5
	case "h5":
		weight = 2
	case "h6":
		weight = 1.5
	}
	return weight
}

// Get takes a URL, makes a GET call and returns a goquery document for parsing.
func Get(link string) (*goquery.Document, error){
	resp, err := http.Get(link)
		if err != nil {
			return nil, err
		}
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	return doc, nil
}

// findPageTitle parses a goquery document to obtain and return the title (if such).
func findPageTitle(doc *goquery.Document) string {
	var title string
	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		title = s.Find("title").Text()
	})
	return title
}

// findPageMetaInfo parses a goquery document to obtain and return the meta info (if such).
func findPageMetaInfo(doc *goquery.Document) string {
	var meta string
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); strings.EqualFold(name, "description") {
			meta, _ = s.Attr("content")
		}
	})
	return meta
}

// Crawl uses a URL link to obtain goquery document, then raises a page structure instance and parses the document.
// Returns populated page structure and slice of new links found at the given URL web-page or error if parsing was not possible.
func Crawl(link string) (structs.WebPage, []string, error) {
	doc, err := Get(link)
		if err != nil {
			return structs.WebPage{}, nil, err
		}
	temp := make(map[string]struct{})
	links := make([]string, 0)
	page := structs.WebPage{Words: make(map[string]float64, 0)}
	page.Url = link
	page.Title = findPageTitle(doc)
	page.Meta = findPageMetaInfo(doc)

	doc.Find("h1, h2, h3, h4, h5, h6, p, ol, li, th, td, a").Each(func(_ int, link *goquery.Selection) {
		content := link.Text()
		tag := link.Nodes[0].Data

		if tag == "a" {
			href, _ := link.Attr("href")
			if isUrl(href) {
				temp[href] = struct{}{}
			}
		} else {
			processContent(tag, content, &page)
		}
	})
	for l := range temp {
		links = append(links, l)
	}
	if page.ParsedChars == 0 {
		return page, links, errors.New("page not parsed")
	}
	return page, links, nil
}

// processContent processes non-link content data records the final words into a page struct.
func processContent(tag string, content string, page *structs.WebPage) {
	contentLength := len(content)
	weight := evaluateWeight(tag)
	if page.ParsedChars + contentLength > 5000 {
		cut := (page.ParsedChars + contentLength) - 5000
		content = content[:cut]
		contentLength -= cut
	}
	page.ParsedChars += contentLength
	words := strings.Fields(content)
	filteredWords := helpers.FilteredWords(&words, weight)
	for word, w := range filteredWords {
		page.Words[word] += w
	}
}
