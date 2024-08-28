package main

import (
	"database/sql"
	"fmt"

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
		var level int
		if err := rows.Scan(&pk, &title, &url, &ip, &alive, &level); err != nil {
			panic(err)
		}
		fmt.Printf("ID: %d, Title: %s, URL: %s, IP: %s, Alive: %t, Level: %d\n", pk, title, url, ip, alive, level)
	}

	return db, nil
}

// After the companies are found the findlinks function will call InsertTarget to add a db record
func InsertTarget(title string, url string, ip string, alive bool, level int) (*sql.DB, error) {
	config := LoadConfig()
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", config.user, config.dbname, config.password, config.sslmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Perform a sample query
	rows, err := db.Query(fmt.Sprintf("INSERT INTO targets (title, url, ip, alive, level) VALUES ('%s', '%s', '%s', %t, %d)", title, url, ip, alive, level))
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	SelectRecords()

	return db, nil
}
