package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
)

type SiteInfo struct {
	Title string `json:"title"`
	IP    string `json:"ip"`
	Alive string `json:"alive"`
}

type html struct {
	Title title `xml:"head>title"`
}

type title struct {
	Text string `xml:",innerxml"`
}

func GetSiteInfo(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	domainPattern := regexp.MustCompile(`https?://([^/]+)`)
	domain := domainPattern.FindStringSubmatch(url)[1]
	ip, err := net.LookupIP(domain)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	h := html{}
	b := []byte(string(body))

	decoder := xml.NewDecoder(bytes.NewBuffer(b))
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose
	decoder.Entity = xml.HTMLEntity
	erro := decoder.Decode(&h)
	if erro != nil {
		fmt.Println("No Title")
		// return
	}

	fmt.Println(strings.TrimSpace(h.Title.Text))
	// status := resp.StatusCode
	// fmt.Println(string(body))
	fmt.Println(resp.Header)
	fmt.Println(ip)
	// fmt.Println(status)
}
