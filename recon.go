package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type VulnReport struct {
	URL     string   `json:"url"`
	Folders []string `json:"folders"`
}

func Recon(wordlist string, subdomainList string, url string, hasCMS bool, CMS string, CMSVersion string) (string, error) {
	folders := CrawlSite(url)
	fmt.Printf("URL: %s\nhasCMS: %t\nCMS: %s\nCMS Version: %s\n", url, hasCMS, CMS, CMSVersion)

	// Find all the folders using a wordlist
	fmt.Println("Brute Forcing")
	errorPattern := regexp.MustCompile(`(?i)Error|Oops|404|Not\sFound|Page\sIsn't\sAvailable`)
	if wordlist == "" {
		return "Usage: cdnreaper -db ['google | local'] -dw ['directory wordlist'] -sw ['subdomain wordlist']", nil
	} else {
		wordlistFile, err := os.ReadFile(wordlist)
		if err != nil {
			return "", err
		}

		content := strings.Split(string(wordlistFile), "\n")
		for _, line := range content {
			if strings.HasPrefix(line, "#") {
				continue
			}
			resp, err := http.Get(fmt.Sprintf("%s/%s", url, line))
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			defer resp.Body.Close()
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}

			body := string(bodyBytes)
			if errorPattern.MatchString(body) {
				fmt.Printf("%s/%s\n"+Red, url, line)
				continue
			} else {
				fmt.Printf("%s/%s\n"+Green, url, line)
				folders = append(folders, fmt.Sprintf("%s/%s", url, line))
			}
		}
	}

	// Crawl the website to find hrefs and follow them
	var vulnReport VulnReport

	vulnReport.URL = url
	vulnReport.Folders = folders

	foldersJSON, err := json.Marshal(vulnReport)
	if err != nil {
		return "", err
	}
	InsertFolders(url, foldersJSON)
	fmt.Printf("%s\n"+Reset, folders)

	// Download CMS and plugins to scan for vulns

	return "", nil
}

func CrawlSite(url string) []string {

	// Create an empty array that will hold all of the links which kinda equate to folders
	// Open the link
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer resp.Body.Close()

	// Decode the HTML so we can look for links

	doc, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}
	return ExtractLinks(doc)
}

func ExtractLinks(doc *html.Node) []string {
	var links []string

	if doc.Type == html.ElementNode && doc.Data == "a" {
		for _, attr := range doc.Attr {
			if attr.Key == "href" {
				links = append(links, attr.Val)
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, ExtractLinks(c)...)
	}
	return links
}

// Vulnerability Scanners

func WordPressVulnerabilityScanner(url string) ([]string, error) {

	// return [nil], nil
}
