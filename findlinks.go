package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// Define a structure for the JSON response
type SearchResult struct {
	SearchInformation struct {
		TotalResults string `json:"totalResults"`
	} `json:"searchInformation"`
	Items []struct {
		Title string `json:"title"`
		Link  string `json:"link"`
	} `json:"items"`
}

func findLinks() (*SearchResult, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	cx := os.Getenv("CSE")

	// The search query
	query := "Companies in bloemfontein -list -directory -top -best -companies"

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
		fmt.Printf("Fuck")
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
			fmt.Printf("Title: %s\nLink: %s\n\n", item.Title, item.Link)
		}

	}

	// fmt.Println(result.Items)

	return &result, nil
}
