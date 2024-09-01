package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
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

// Function to write a number to a file
func writeNumberToFile(filePath string, number int) error {
	// Create or truncate the file
	file, err := os.Create(filePath)
	numberString := strconv.Itoa(number)
	if err != nil {
		return err
	}
	defer file.Close() // Ensure the file is closed after writing

	// Convert the number to a string and write it to the file
	_, err = file.WriteString(numberString)
	return err
}

// Function to read a number from a file
func readNumberFromFile(filePath string) (int, error) {
	// Read the entire file content
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return 0, err
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	// Convert the content to a string and trim whitespace/newlines
	contentStr := strings.TrimSpace(string(content))

	// Convert the string to an integer
	number, err := strconv.Atoi(contentStr)
	if err != nil {
		return 0, err
	}

	return number, nil
}

// This function connects to Google to find URLs and Titles for potential companies that we will want to pentest
func findLinks() (*SearchResult, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")          // Google API key set as environment variable for security
	cx := os.Getenv("CSE")                         // Custom search engine key set as environment variable for security
	queryFile, err := os.Create("query_count.txt") // Keeps count of the query count because we have limited queries with the Google API
	if err != nil {
		fmt.Println(err)
	}
	fileInfo, err := os.Stat("query_count.txt")
	flag := false
	if err != nil {
		fmt.Println(err)
	}
	defer queryFile.Close()
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
		log.Panic("Error decoding JSON:", err)
		return nil, err
	}
	// Print the search results
	var totalPages string = result.SearchInformation.TotalResults
	calculatePages, err := strconv.Atoi(totalPages)
	calculatePages = calculatePages / 10
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	for start := range calculatePages {
		/*
			We can only get 10 results per page, so we need to figure out how many results there are so we can work out how many pages there are.
			This makes sure that we start at on and then at the last result on the page + 1
		*/
		fmt.Println(start)
		if start%100 == 0 && start >= 100 {
			break
		}
		if fileInfo.Size() > 2 && !flag {

			content, err := readNumberFromFile("query_count.txt")
			if err != nil {
				fmt.Println(err)
			}

			start, err = strconv.Atoi(string(content))
			if err != nil {
				fmt.Println(err)
			}
			flag = true
		}
		writeNumberToFile("query_count.txt", start)
		start = start*10 + 1
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
			break
		}
		if resp.Body != nil {
			for _, item := range result.Items {
				url := pattern.FindString(item.Link)
				sites = append(sites, GetSiteInfo(url))
			}
		}
		InsertTarget(sites)
		sites = nil
	}

	return &result, nil
}
