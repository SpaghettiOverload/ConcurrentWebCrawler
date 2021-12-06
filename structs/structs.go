package structs

import (
	"sync"
)

type WebPage struct {
	Url string
	Title string
	Meta string
	Words map[string]float64
	ParsedChars int
}

type UserResultsPage struct {
	Url string
	Title string
	Meta string
	Relevance float64
	TotalWeight float64
}

type UserResultsPageList []UserResultsPage

func (e UserResultsPageList) Len() int {
	return len(e)
}

func (e UserResultsPageList) Less(i, j int) bool {
	if e[i].Relevance > e[j].Relevance {
		return true
	}
	if e[i].Relevance < e[j].Relevance {
		return false
	}
	return e[i].TotalWeight > e[j].TotalWeight
}

func (e UserResultsPageList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

type Results struct {
	mu sync.Mutex
	SP []WebPage
	L  []string
}

func (c *Results) UpdatePages(s *WebPage) {
	c.mu.Lock()
	c.SP = append(c.SP, *s)
	c.mu.Unlock()
}

func (c *Results) UpdateLinks(s *[]string) {
	c.mu.Lock()
	c.L = append(c.L, *s...)
	c.mu.Unlock()
}

func (c *Results) CleanLinks() {
	c.mu.Lock()
	c.L = make([]string, 0)
	c.mu.Unlock()
}
