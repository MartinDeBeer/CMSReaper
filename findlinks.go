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
	apiKey := os.Getenv("GOOGLE_API_KEY") // Google API key set as environment variable for security
	cx := os.Getenv("CSE")                // Custom search engine key set as environment variable for security

	var sites []string
	// The search query
	query := "Companies in bloemfontein location:bloemfontein -list -directory -top -best -companies -site:*.gov.* -site:maps.google.com -site:facebook.* -site:tiktok.* -site:twitter.* -site:pinterest.*"

	// Create the request URL
	searchURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s", apiKey, cx, url.QueryEscape(query))

	// Make the HTTP request to goole
	resp, err := http.Get(searchURL)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var result SearchResult // Decode the JSON response
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

	for start := range calculatePages {
		/*
			We can only get 10 results per page, so we need to figure out how many results there are so we can work out how many pages there are.
			This makes sure that we start at on and then at the last result on the page + 1
		*/
		start = start*10 + 1
		fmt.Println(start)
		searchURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&start=%d", apiKey, cx, url.QueryEscape(query), start)
		resp, err := http.Get(searchURL)
		pattern := regexp.MustCompile(`http[s]?://[^/]+`) // This ensures that we end up on the main page.
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
			sites = append(sites, GetSiteInfo(url))
		}
		InsertTarget(sites)
		sites = nil
	}

	return &result, nil
}
