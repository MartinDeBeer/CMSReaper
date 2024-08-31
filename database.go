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
func SelectRecords() (*sql.DB, error) {
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
		if err := rows.Scan(&pk, &title, &url, &ip, &alive); err != nil {
			panic(err)
		}
		fmt.Printf("ID: %d, Title: %s, URL: %s, IP: %s, Alive: %t\n", pk, title, url, ip, alive)
	}

	return db, nil
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
		// Perform an insert query

		// First, create the temporary table
		createTmpTableQuery := `
			CREATE TEMPORARY TABLE IF NOT EXISTS tmp_data (title TEXT, url TEXT, ip TEXT, alive BOOLEAN);
		`

		_, err := db.Exec(createTmpTableQuery)
		if err != nil {
			log.Fatalf("Failed to create temporary table: %v", err)
		}
		fmt.Println("Created temporary table")
		// Next, insert the data into the temporary table
		insertTmpDataQuery := `
			INSERT INTO tmp_data (title, url, ip, alive) VALUES ($1, $2, $3, $4);
		`

		_, err = db.Exec(insertTmpDataQuery, websiteProperty.Title, websiteProperty.Url, websiteProperty.IP, alive)
		if err != nil {
			log.Fatalf("Failed to insert data into temporary table: %v", err)
		}
		fmt.Println("Added data to temporary table")
		// Finally, perform the conditional insert into the targets table
		insertIntoTargetsQuery := `
			INSERT INTO targets (title, url, ip, alive)
			SELECT DISTINCT title, url, ip, alive
			FROM tmp_data
			WHERE NOT EXISTS (
				SELECT 1 FROM targets
				WHERE targets.title = tmp_data.title
				AND targets.url = tmp_data.url
				AND targets.ip = tmp_data.ip
				AND targets.alive = tmp_data.alive
			);
		`

		_, err = db.Exec(insertIntoTargetsQuery)
		if err != nil {
			log.Fatalf("Failed to insert into targets: %v", err)
		}
		fmt.Println("Added data to targets from tmp_data")
	}

	SelectRecords()

	return db, nil
}
