package main

import "fmt"

type VulnReport struct {
	URL string `json:"url"`
}

func Recon(url string) (string, error) {

	// Figure out what the CDN is if there is one


	// Find all the folders using a wordlist


	// Crawl the website to find hrefs and follow them


	// Brute force subdomains


	// Download CMS and plugins to scan for vulns


	fmt.Printf("Starting recon on %s\n", url)

	return "", nil
}


func GetCMS() {
	
}