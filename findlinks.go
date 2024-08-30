package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

// Define a structure for the JSON response
type SearchResult struct {
	SearchInformation struct {
		TotalResults string `json:"totalResults"`
	} `json:"searchInformation"`
	Items []struct {
		Link string `json:"link"`
	} `json:"items"`
}

// This function connects to Google to find URLs and Titles for potential companies that we will want to pentest
func findLinks() (*SearchResult, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	cx := os.Getenv("CSE")

	// The search query
	// query := "location:bloemfontein \"Software Company\" site:*.co.za -site:downdetector.* -site:crown.co.za -site:ethekwini.co.za -site:spurmtbleague.* -site:yep.co.za -filetype:pdf -filetype:doc -filetype:docx -filetype:ppt -filetype:pptx -filetype:xls -filetype:xlsx -site:facebook.* -site:linkedin.* -site:tiktok.* -site:wikipedia.org -site:twitter.* -site:pinterest.* -list -directory -best -top -jobs|job -careers|career -position|positions -news"

	// query := "location:bloemfontein \"Software Company\" -jobs|job -careers|career -position|positions -news"

	query := "Companies in bloemfontein location:bloemfontein -list -directory -top -best -companies -site:*.gov.* -site:maps.google.com -site:facebook.* -site:tiktok.* -site:twitter.* -site:pinterest.*"

	// Create the request URL
	searchURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s", apiKey, cx, url.QueryEscape(query))

	// Make the HTTP request
	resp, err := http.Get(searchURL)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}
	// Print the search results
	var totalPages string = result.SearchInformation.TotalResults
	calculatePages, err := strconv.Atoi(totalPages)
	calculatePages = calculatePages / 10
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(calculatePages)
	for page := range calculatePages {
		if page == 2 {
			fmt.Println("Reached limit")
			break
		}
		fmt.Printf("Page %d\n\n", page+1)
		searchURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&start=%d", apiKey, cx, url.QueryEscape(query), page)
		resp, err := http.Get(searchURL)
		pattern := regexp.MustCompile(`http[s]?://[^/]+`)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return nil, err
		}
		for _, item := range result.Items {
			url := pattern.FindString(item.Link)
			GetSiteInfo(url)
		}

	}

	return &result, nil
}
