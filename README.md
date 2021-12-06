##Concurrent web-crawling and data aggregation using goroutines


Implement a Golang application that receives from command line 1) a URL of a web page, 2) a MaxRoutines number - the maximal number of goroutines, 3) MaxIndexingTime - the maximal allowed time for indexing, 4) base result files name (default "results.json”), 5) MaxResults - the maximum number of resources processed (crawled and indexed)

PHASE I: Indexing
1. Traverse the connected pages starting from the URL provided and following the hyperlinks in these web pages. Separate goroutines should be used to speed up the traversal process. You could use breadth-first-search, or custom heuristics to prioritize the search pages.

2. The number of goroutines working in parallel should not exceed MaxRoutines number in any moment during the program execution.

3. Extracts information about the keywords mentioned in each page using a custom extraction and weighting criteria (keywords and names found in headers and sub-headers - h1 and h2 - should be more relevant than the others, as well as the keywords found in the first 5000 characters of the web page body). You should suggest your own criteria.

4. Aggregate the extracted information about the found web pages. Output the page urls and necessary metadata to search in them to a JSON file (up to MaxResults resources indexed) named using the provided base result files name (default "results.json”).

5. The maximal search time should not exceed the specified MaxIndexingTime.

PHASE II: Search
6. Input a list of search keywords from console and find the most relevant web pages to keywords. Output the sorted list of most relevant top 10 resources with relevance percentage and metadata to console in human readable form. You can repeat this step with multiple keywords combinations.

* When indexing pages you can skip the so-called stop words, such as (please extend the list as appropriate):

a, is, the, an, and, are, as, at, be, but, by, for,

if, in, into, it, no, not, of, on, or, such, that, their,

then, there, these, they, this, to, was, will, with
