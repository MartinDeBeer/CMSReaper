package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
)

func WordPressVulnerabilityScanner(url string) (string, error) {
	fmt.Println(Blue + "[+] Finding Plugins")
	themePattern := `\/wp-content\/themes\/([^\/]+)\/.*ver=([^']+)`
	pluginPattern := `\/wp-content\/plugins\/([^\/]+)\/.*ver=([^'|"]+)`

	httpClient := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(Red+"[-] Error creating request:"+Reset, err)
		return "Failed", nil
	}
	resp, err := httpClient.Do(req)

	if err != nil {
		fmt.Println(Red+"[-] Error fetching URL:"+Reset, err)
		return "Failed", nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 || resp.Body == http.NoBody {
		return "Failed", nil
	}
	if err != nil {
		fmt.Println(Red+"[-] Error reading response body:"+Reset, err)
		return "Failed", nil
	}
	// Compile regex patterns
	themeRegex := regexp.MustCompile(themePattern)

	// Find theme and version in URL
	themeMatches := themeRegex.FindStringSubmatch(string(body))
	pluginRegex := regexp.MustCompile(pluginPattern)

	// Find plugin names and versions in URL
	pluginMatches := pluginRegex.FindAllStringSubmatch(string(body), -1)

	if len(themeMatches[1]) < 2 || len(themeMatches[2]) < 2 {
		fmt.Println(Red + "[-] Theme or version not found in URL" + Reset)
		return "Failed", nil
	}
	var pluginName = ""
	var pluginZipFile = ""
	domainPattern := regexp.MustCompile(`https?://([^/]+)`)
	domain := domainPattern.FindStringSubmatch(url)[1]
	newpath := filepath.Join(".", fmt.Sprintf("%s/plugins", domain))
	if err := os.MkdirAll(newpath, os.ModePerm); err != nil {
		return "Failed", err
	}

	for _, match := range pluginMatches {
		pluginName = match[1]
		pluginZipFile = fmt.Sprintf("%s/%s.zip", newpath, pluginName)
		err := FindWPPlugins(match[1], match[2], fmt.Sprintf("%s/%s.zip", newpath, pluginName))
		if err != nil {
			fmt.Println("[-] Error downloading"+Reset, err)
			return "Failed", nil
		}
	}
	// Unzip the plugin ZIP file
	pluginDir, err := unzip(pluginZipFile, fmt.Sprintf("%s/%s", domain, pluginName))
	if err != nil {
		fmt.Println(Red+"[-] Error unzipping plugin:"+Reset, err)
		os.Remove(pluginZipFile)
		return "Failed", nil
	}
	err = analyzePlugin(pluginDir)
	if err != nil {
		fmt.Println(Red+"[-] Error analyzing plugin:"+Reset, err)
		os.Remove(pluginZipFile)
		return "Failed", nil
	}
	defer os.Remove(pluginZipFile) // Clean up the ZIP file after extraction

	// defer os.RemoveAll(pluginDir) // Clean up the extracted files after analysis

	// Perform static analysis on the plugin code

	return "Success", nil
	// return [nil], nil
}

func FindWPPlugins(pluginName string, pluginVersion string, path string) error {
	// Download the plugin to the current directory
	var downloaded []string
	// pluginVersionPattern := regexp.MustCompile(`^[^.]*\.[^.]*\.[^.]{1,2}$`)
	// versionMatches := pluginVersionPattern.FindStringSubmatch(pluginVersion)
	// if len(versionMatches) == 0 {
	// 	return nil
	// }
	if slices.Contains(downloaded, pluginName) {
		return nil
	}
	url := fmt.Sprintf("https://downloads.wordpress.org/plugins/%s.%s.zip", pluginName, pluginVersion)
	fmt.Printf(Green+"[+] Reaching out to %s\n", url+Reset)
	resp, err := http.Get(url)
	if resp.ContentLength <= 100 {
		fmt.Printf(Red+"[-] Empty plugin %s | Content Length %d\n"+Reset, pluginName, resp.ContentLength)
		os.Remove(pluginName)
		return err
	}
	fileNotFound := regexp.MustCompile("(?i)File not found")
	if resp.StatusCode != http.StatusOK || resp.Body == http.NoBody {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if fileNotFound.Match(body) {
		fmt.Printf(Red+"[-] Plugin %s Not Found\n"+Reset, pluginName)
		return err
	}
	defer resp.Body.Close()
	downloaded = append(downloaded, pluginName)

	pluginFile := path

	out, err := os.Create(pluginFile)
	fmt.Printf(Green+"[+] Creating Folder %s\n", pluginName+Reset)

	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	// fmt.Printf(Green+"[+] Copying %s Content to  %s\n", resp.Body, pluginName+Reset)
	if err != nil {
		fmt.Println("Failed")
		return err
	}
	return err
}

func unzip(src string, dest string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()
		if f.Name == dest+"/" {
			continue
		}
		// fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(f.Name, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(f.Name), os.ModePerm); err != nil {
				return "", err
			}
			outFile, err := os.Create(f.Name)
			if err != nil {
				return "", err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, rc)
			if err != nil {
				return "", err
			}
		}
	}

	return dest, nil
}

func analyzePlugin(pluginDir string) error {
	// Read all files in the plugin directory

	// Root dir = pluginDir
	files, err := os.ReadDir(pluginDir)
	if err != nil {
		return err
	}

	// Example: Search for potential XSS vulnerabilities in PHP files
	xssRegex := regexp.MustCompile(`<script[^>]*>`)
	sqlInjectionRegex := []string{
		// Pattern to find SQL keywords followed by superglobals
		`(SELECT|INSERT|UPDATE|DELETE|FROM|WHERE).*?(\$_(GET|POST|REQUEST|COOKIE|SERVER)\[.*?\])`,
		// Pattern to find concatenation of variables in query strings
		`(SELECT|INSERT|UPDATE|DELETE|FROM|WHERE).*?(\$[a-zA-Z_][a-zA-Z0-9_]*)\s*\.\s*(\$[a-zA-Z_][a-zA-Z0-9_]*|\$_(GET|POST|REQUEST|COOKIE|SERVER)\[.*?\])`,
		// Pattern to identify usage of deprecated mysql_query with superglobals
		`mysql_query\s*\(.*\$(GET|POST|REQUEST|COOKIE|SERVER)\[.*\].*\)`,
	}
	for _, file := range files {
		filePath := filepath.Join(pluginDir, file.Name())
		if file.IsDir() {
			// Recursively analyze the subdirectory
			err := analyzePlugin(filePath)
			if err != nil {
				return err
			}
		} else if filepath.Ext(file.Name()) == ".php" {
			content, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			if xssRegex.Match(content) {
				fmt.Printf(Blue+"[+] Potential XSS vulnerability found in %s\n"+Reset, filePath)
			}
			for _, match := range sqlInjectionRegex {
				re := regexp.MustCompile(match)
				matches := re.FindAllString(filePath, -1)
				fmt.Printf(Blue+"[+] Checking SQLi in %s\n", filePath)
				if len(matches) > 0 {
					fmt.Printf("Potential SQL injection vulnerabilities found for pattern: %s\n"+Reset, match)
					for _, match := range matches {
						fmt.Println(match)
					}
				}
			}
		}
	}
	return nil
}
