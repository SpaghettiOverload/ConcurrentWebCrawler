package main

import (
	"ConcurrentWebCrawler/helpers"
	"ConcurrentWebCrawler/structs"
	"ConcurrentWebCrawler/web"
	"github.com/briandowns/spinner"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

var (
	loading = spinner.New(spinner.CharSets[35], 100*time.Millisecond)
)

func main() {

	// todo INFO: Tested with command line:
	// "go run main.go https://gobyexample.com/range-over-channels 40 5 results.json 20"

	commandLineInput := os.Args
	initUrl := commandLineInput[1]                                     				// 1
	maxRoutines, _ := strconv.Atoi(commandLineInput[2])               				// 2
	maxIndexingTime, _ := strconv.ParseFloat(commandLineInput[3], 64) 		// 3 seconds
	baseResultsFile := commandLineInput[4]                           				// 4
	MaxResults, _ := strconv.Atoi(commandLineInput[5])               				// 5

	// PHASE I
	start := time.Now()
	loading.Start()
	results :=	web.CrawlUrls(initUrl, maxRoutines, MaxResults)
	elapsed := time.Since(start).Seconds()
	loading.Stop()

	helpers.PrintExecTime(elapsed, maxIndexingTime, MaxResults, maxRoutines)
	helpers.WriteResultsToFile(baseResultsFile, &results)

	// PHASE II
	searchWords := helpers.SearchPrompt()
	searchResults, err := helpers.SearchForKeywords(baseResultsFile, searchWords)
	if err != nil {
		log.Fatal("No results found")
	}

	// Sorting structures using by-then (by relevancy & then by weight)
	sort.Sort(structs.UserResultsPageList(searchResults))
	helpers.FormatResults(searchResults)
}
