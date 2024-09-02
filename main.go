package main

import (
	"flag"
	"fmt"
)

func main() {
	// findLinks()

	// Check if flag for Google is ste
	dbFlag := flag.String("db", "", "-db Use either a database on the local machine or build a db from Google search yes|no")
	// var customDBFlag = flag.String("-dbname", "cdnreaper", "-dbname Name of custom database. Default <cdnreaper>")

	flag.Parse()
	fmt.Println(GetSiteInfo(*dbFlag))
}
