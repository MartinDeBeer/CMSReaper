package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
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
		Link  string `json:"link"`
		Title string `json:"title"`
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
	_, err = file.Write([]byte(numberString))
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
func findLinks() ([]string, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY") // Google API key set as environment variable for security
	cx := os.Getenv("CSE")                // Custom search engine key set as environment variable for security
	start := 0
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
		log.Panic(Red+"Error decoding JSON:"+Reset, err)
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

	content, err := readNumberFromFile("query_count.txt")
	if err != nil {
		fmt.Println(err)
	}
	start = content
	for start < 100 {
		/*
			We can only get 10 results per page, so we need to figure out how many results there are so we can work out how many pages there are.
			This makes sure that we start at on and then at the last result on the page + 1

		*/

		fmt.Println(start)

		if start%100 == 0 && start >= 100 {
			fmt.Println(Red + "[-] We went over the quota" + Reset)
			break
		}

		writeNumberToFile("query_count.txt", start)

		searchURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&start=%d", apiKey, cx, url.QueryEscape(query), start)
		resp, err := http.Get(searchURL)
		pattern := regexp.MustCompile(`http[s]?://[^/]+`) // This ensures that we end up on the main page.
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Println(Red+"[-] Error decoding JSON:"+Reset, err)
			break
		}
		if resp.Body != nil {
			for _, item := range result.Items {
				url := pattern.FindString(item.Link)
				sites = append(sites, getExtraInfo(url))
			}
		}
		fmt.Println(sites)
		InsertTarget(sites)
		start = start + 10
	}
	return sites, nil
}

func getExtraInfo(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	defer resp.Body.Close()

	domainPattern := regexp.MustCompile(`https?://([^/]+)`)
	cdnPattern := regexp.MustCompile(`<meta\s+name="generator"\scontent="([^"]\w+)\s(\S+)".*`)
	wpNoVersion := regexp.MustCompile(`wp-content`)

	domain := domainPattern.FindStringSubmatch(url)[1]
	ip, err := net.LookupIP(domain)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	h := HTML{}
	b := []byte(string(body))

	decoder := xml.NewDecoder(bytes.NewBuffer(b))
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose
	decoder.Entity = xml.HTMLEntity
	erro := decoder.Decode(&h)
	var siteInfo SiteInfo
	if erro != nil {
		siteInfo.Title = strings.TrimSpace(domain)
	}

	siteInfo.Title = strings.TrimSpace(h.Title.Text)
	siteInfo.IP = ip[0].String()
	if resp.StatusCode == 200 {
		siteInfo.Alive = "true"
	} else {
		siteInfo.Alive = "false"
	}

	// Check for CMS patterns in the HTML content
	match := cdnPattern.FindStringSubmatch(string(body))
	wpNoVersionMatch := wpNoVersion.FindStringSubmatch(string(body))

	if len(match) < 2 {
		siteInfo.HasCMS = "false"
	} else {
		siteInfo.HasCMS = "true"
		siteInfo.CMSVersion = match[2]
		siteInfo.CMS = match[1]
		fmt.Printf(Blue+"[+] The CMS is: %s and the version is: %s\n"+Reset, siteInfo.CMS, siteInfo.CMSVersion)
	}
	if len(wpNoVersionMatch) > 0 && len(match) < 2 {
		siteInfo.HasCMS = "true"
		siteInfo.CMSVersion = "No Version"
		siteInfo.CMS = "WordPress"
		fmt.Printf(Blue+"[+] The CMS is: %s and the version is: %s\n"+Reset, siteInfo.CMS, siteInfo.CMSVersion)
	}

	siteInfo.Url = url
	siteInfoJson, err := json.Marshal(siteInfo)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	return string(siteInfoJson)
}
