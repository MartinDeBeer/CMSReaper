package main

import (
	"encoding/json"
	"fmt"
)

type SiteInfo struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	IP    string `json:"ip"`
	Alive string `json:"alive"`
	CDN   string `json:"cdn"`
}

type html struct {
	Title title `xml:"head>title"`
}

type title struct {
	Text string `xml:",innerxml"`
}

func GetSiteInfo(flag string) string {

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
			fmt.Println(site)
			var siteInfo SiteInfo
			if err := json.Unmarshal([]byte(site), &siteInfo); err != nil {
				fmt.Println("Error decoding JSON:", err)
				break
			}
			Recon(siteInfo.Url)
		}
		return ""
	default:
		fmt.Println("Usage: cdnreaper -db [options]")
		return "Usage: cdnreaper -db [options]"
	}

}
