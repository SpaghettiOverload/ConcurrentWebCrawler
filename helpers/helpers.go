package helpers

import (
	"ConcurrentWebCrawler/stopwords"
	"ConcurrentWebCrawler/structs"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkg/browser"
)

const MaxResults = 10

// PrintExecTime evaluates execution time and pre-set limitations, then prints informational results.
func PrintExecTime(execT, maxIndexingTime float64, MaxResults, maxRoutines int) {
	var Reset  = "\033[0m"
	var Red    = "\033[31m"
	var Green  = "\033[32m"

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{""})
	t.AppendRows([]table.Row{
		{"Indexing completed!"},
		{"Web-pages processed", MaxResults},
		{"Max routines used", maxRoutines},
		{"Execution time", fmt.Sprintf("%v seconds", execT)},
	})

	if execT > maxIndexingTime {
		percent := math.Abs((maxIndexingTime - execT) / execT * 100)
		t.AppendRows([]table.Row{
			{"	THIS IS:", fmt.Sprintf(Red + "~%v%% slower" + Reset +
				" than `maxIndexingTime`(%vs)", math.RoundToEven(percent), maxIndexingTime)},
		})
	} else {
		percent := math.Abs((execT - maxIndexingTime) / maxIndexingTime * 100)
		t.AppendRows([]table.Row{
			{"	THIS IS:", fmt.Sprintf(Green + "~%v%% faster" + Reset +
				" than `maxIndexingTime`(%vs)", math.RoundToEven(percent), maxIndexingTime)},
		})
	}
	t.SetStyle(table.StyleColoredBright)
	t.Render()
}

// WriteResultsToFile prettify the page structs and writes them in JSON format to a file.
func WriteResultsToFile(baseResultsFile string, results *[]structs.WebPage) {
	jsonString, _ := json.MarshalIndent(results, "", "\t")
	err := ioutil.WriteFile(baseResultsFile, jsonString, os.ModePerm)
	if err != nil {
		return
	}
}

// FilteredWords filters out words based on stop-words file criteria and general regex expression.
// It returns a map with only "good to go" words with their respective weight.
func FilteredWords(words *[]string, weight float64) map[string]float64{
	w := make(map[string]float64, 0)
	for _, word := range *words {
		word = strings.Trim(word, ",.()!:-")
		if matched, _ := regexp.MatchString(`^[a-zA-Z]{4,}$`, word); matched {
			word = strings.ToLower(word)
			if !stopwords.English[word] {
				w[word] += weight
			}
		}
	}
	return w
}

// SearchPrompt prints transition to PHASE II of the program, collects user input and returns it as a slice of keywords.
func SearchPrompt() []string {
	var words []string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nEnter search words (space separated): ")
		scanner.Scan()
		in := scanner.Text()
		if in != "" {
			wh := strings.Fields(in)
			for i:=0; i<len(wh); i++ {
				w := strings.ToLower(wh[i])
				words = append(words, w)
			}
			break
		}
	}
	return words
}

// SearchForKeywords searches for the given keywords in the DB and returns relevant results or error if none.
func SearchForKeywords(baseResultsFile string, searchWords []string) ([]structs.UserResultsPage, error) {
	jsonDB := readFromJsonDB(baseResultsFile)
	temp := make(map[string]structs.UserResultsPage, 0)
	n := float64(len(searchWords))

	for _, page := range jsonDB {
		url := page.Url
		for _, word := range searchWords {
			if weight, ok := page.Words[word]; ok {
				if _, inside := temp[url]; inside {
					v := temp[url]
					v.Relevance += 100 / n
					v.TotalWeight += weight
					temp[url] = v
				} else {
					temp[url] = structs.UserResultsPage{
						Url: url,
						Title:     page.Title,
						Meta:      page.Meta,
						Relevance: 100 / n,
						TotalWeight: weight,
					}
				}
			}
		}
	}

	if len(temp) == 0 {
		return []structs.UserResultsPage{}, errors.New("no matches found")
	}
	searchResults := make([]structs.UserResultsPage, 0)
	for _, value := range temp { searchResults = append(searchResults, value) }
	return searchResults, nil
}

// FormatResults prettifies the search results and generates an HTML report from them.
func FormatResults(searchResults []structs.UserResultsPage) {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"#", "Title", "Info", "Relevancy", "URL"})
	for i, result := range searchResults {
		if i == MaxResults {
			break
		}
		t.AppendRows([]table.Row{
			{i+1, result.Title, result.Meta, fmt.Sprintf("%v %%", math.Round(result.Relevance)), result.Url},
			{" ", " ", " ", " ", " "},
		})
	}
	t.Style().HTML = table.HTMLOptions{
		CSSClass:    "Results report",
		EmptyColumn: "&nbsp;",
		EscapeText:  true,
		Newline:     "<br/>",
	}
	html := t.RenderHTML()
	writeAndOpenSearchReport(html)
}

// writeAndOpenSearchReport Generates a html document with the search results and opens it with the users default browser.
func writeAndOpenSearchReport(html string) {
	fileName:= "report.html"
	f, _ := os.Create(fileName)
	defer f.Close()
	_, _ = f.WriteString(html)
	fmt.Printf("\nReport file created: %v\n\n", fileName)
	err := browser.OpenFile(fileName)
	if err != nil {
		return
	}
}

// readFromJsonDB opens up the previously generated json file and creates inmemory DB.
func readFromJsonDB(baseResultsFile string) []structs.WebPage {
	var resultsFromJson []structs.WebPage
	f, _ := ioutil.ReadFile(baseResultsFile)
	err := json.Unmarshal(f, &resultsFromJson)
	if err != nil {
		return nil
	}
	return resultsFromJson
}