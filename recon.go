package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

type VulnReport struct {
	URL string `json:"url"`
}

func Recon(url string, hasCMS bool, CMS string, CMSVersion string) (string, error) {
	fmt.Printf("URL: %s\nhasCMS: %t\nCMS: %s\nCMS Version: %s\n", url, hasCMS, CMS, CMSVersion)

	// Find all the folders using a wordlist
	CrawlSite(url)
	// Crawl the website to find hrefs and follow them

	// Brute force subdomains

	// Download CMS and plugins to scan for vulns

	return "", nil
}

func CrawlSite(url string) {

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
	fmt.Printf("%v\n\n", ExtractLinks(doc))
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
