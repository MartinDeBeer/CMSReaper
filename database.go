package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	user     string
	dbname   string
	password string
	sslmode  string
}

func LoadConfig() DBConfig {
	return DBConfig{
		user:     "martin",
		dbname:   "cdnreaper",
		password: "Martin323",
		sslmode:  "disable",
	}
}

// This function returns all records in the database for a given table
func SelectRecords() ([]string, error) {
	results := []string{}
	var siteInfo SiteInfo
	config := LoadConfig()
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", config.user, config.dbname, config.password, config.sslmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Perform a sample query
	rows, err := db.Query("SELECT * FROM targets")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Iterate through the results
	for rows.Next() {
		var pk int
		var title string
		var url string
		var ip string
		var alive bool
		var cms string
		var has_cms bool
		var cms_version string
		if err := rows.Scan(&pk, &title, &url, &ip, &alive, &cms, &has_cms, &cms_version); err != nil {
			panic(err)
		}
		// fmt.Printf("ID: %d, Title: %s, URL: %s, IP: %s, Alive: %t, CDN: %t\n", pk, title, url, ip, alive, cdn)
		siteInfo.Title = title
		siteInfo.Url = url
		siteInfo.IP = ip
		siteInfo.Alive = strconv.FormatBool(alive)
		siteInfo.CMS = cms
		siteInfo.HasCMS = strconv.FormatBool(has_cms)
		siteInfo.CMSVersion = cms_version
		siteInfoJson, err := json.Marshal(siteInfo)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		results = append(results, string(siteInfoJson))
	}

	return results, nil
}

// After the companies are found the findlinks function will call InsertTarget to add a db record
func InsertTarget(sites []string) (*sql.DB, error) {
	config := LoadConfig()
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", config.user, config.dbname, config.password, config.sslmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	for _, site := range sites {
		var websiteProperty SiteInfo
		if err := json.Unmarshal([]byte(site), &websiteProperty); err != nil {
			fmt.Println(err)
		}
		alive, _ := strconv.ParseBool(websiteProperty.Alive)
		if !alive {
			continue
		}
		// Perform an insert query

		// First, create the temporary table
		createTmpTableQuery := `
			CREATE TEMPORARY TABLE IF NOT EXISTS tmp_data (title TEXT, url TEXT, ip TEXT, alive BOOLEAN, cms TEXT, has_cms BOOLEAN, cms_version TEXT);
		`

		_, err := db.Exec(createTmpTableQuery)
		if err != nil {
			log.Fatalf("Failed to create temporary table: %v", err)
		}
		fmt.Println("Created temporary table")
		// Next, insert the data into the temporary table
		insertTmpDataQuery := `
			INSERT INTO tmp_data (title, url, ip, alive, cms, cms_version, has_cms) VALUES ($1, $2, $3, $4, $5, $6, $7);
		`

		_, err = db.Exec(insertTmpDataQuery, websiteProperty.Title, websiteProperty.Url, websiteProperty.IP, alive, websiteProperty.CMS, websiteProperty.CMSVersion, websiteProperty.HasCMS)
		if err != nil {
			log.Fatalf("Failed to insert data into temporary table: %v", err)
		}
		// Finally, perform the conditional insert into the targets table
		insertIntoTargetsQuery := `
			INSERT INTO targets (title, url, ip, alive, cms, has_cms, cms_version)
			SELECT DISTINCT title, url, ip, alive, cms, has_cms, cms_version
			FROM tmp_data
			WHERE NOT EXISTS (
				SELECT 1 FROM targets
				WHERE targets.url = tmp_data.url
			);
		`

		_, err = db.Exec(insertIntoTargetsQuery)
		if err != nil {
			log.Fatalf("Failed to insert into targets: %v", err)
		}
	}

	return db, nil
}

func InsertFolders(url string, sites []byte) (*sql.DB, error) {
	config := LoadConfig()
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", config.user, config.dbname, config.password, config.sslmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var vulnReport VulnReport
	if err := json.Unmarshal(sites, &vulnReport); err != nil {
		fmt.Println(err)
	}
	// Perform an insert query

	for _, folder := range vulnReport.Folders {
		fmt.Println(folder)

		// First, create the temporary table
		createTmpTableQuery := `
			CREATE TEMPORARY TABLE IF NOT EXISTS tmp_data (pk SERIAL PRIMARY KEY, target_key INT, url TEXT, folder_name TEXT);
		`

		if _, err := db.Exec(createTmpTableQuery); err != nil {
			log.Fatalf("Failed to create temporary table: %v", err)
		}
		fmt.Println("Created temporary table")

		// Next, insert the data into the temporary table
		insertTmpDataQuery := `
			INSERT INTO tmp_data (target_key, url, folder_name) VALUES ((SELECT pk FROM targets WHERE url = $1 LIMIT 1), $1, $2);
		`

		_, err = db.Exec(insertTmpDataQuery, vulnReport.URL, folder)
		if err != nil {
			log.Fatalf("Failed to insert data into temporary table: %v", err)
		}
		// Finally, perform the conditional insert into the targets table
		insertIntoTargetsQuery := `
			INSERT INTO folders (target_key, url, folder_name)
			SELECT DISTINCT target_key, url, folder_name
			FROM tmp_data
			WHERE NOT EXISTS (
				SELECT 1 FROM folders
				WHERE folders.url = tmp_data.url
				AND folders.folder_name = tmp_data.folder_name
			);
		`

		_, err = db.Exec(insertIntoTargetsQuery)
		if err != nil {
			log.Fatalf("Failed to insert into targets: %v", err)
		}

		// SelectRecords()
	}

	return db, nil
}
