package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type SiteInfo struct {
	Title      string `json:"title"`
	Url        string `json:"url"`
	IP         string `json:"ip"`
	Alive      string `json:"alive"`
	HasCMS     string `json:"cms_cdn"`
	CMSVersion string `json:"cms_version"`
	CMS        string `json:"cms"`
}

type HTML struct {
	Title title `xml:"head>title"`
}

type title struct {
	Text string `xml:",innerxml"`
}

func GetSiteInfo(flag string, wordlist string, subdomainList string) string {

	switch flag {
	case "google": // Get the data from Google
		sites, err := findLinks()
		if err != nil {
			return err.Error()
		}
		for _, site := range sites {
			fmt.Println(site)
		}
		return ""
	case "local": // Get the data from a local database
		sites, err := SelectRecords()
		if err != nil {
			return err.Error()
		}
		// var siteInfo SiteInfo
		for _, site := range sites {
			var siteInfo SiteInfo
			if err := json.Unmarshal([]byte(site), &siteInfo); err != nil {
				fmt.Println("Error decoding JSON:", err)
				break
			}
			hasCMS, err := strconv.ParseBool(siteInfo.HasCMS)
			if err != nil {
				log.Panic(err)
			}
			Recon(wordlist, subdomainList, siteInfo.Url, hasCMS, siteInfo.CMS, siteInfo.CMSVersion)
		}
		return ""
	default:
		fmt.Println("Usage: cmsreaper -db ['google | local'] -dw ['directory wordlist'] -sw ['subdomain wordlist']")
		return "Usage: cmsreaper -db ['google | local'] -dw ['directory wordlist'] -sw ['subdomain wordlist']"
	}

}
