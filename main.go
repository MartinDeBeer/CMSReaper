package main

import (
	"flag"
	"fmt"
)

func main() {
	// findLinks()
	fmt.Println("\033[2J")
	fmt.Println(" $$$$$$\\  $$\\      $$\\  $$$$$$\\  $$$$$$$\\                                                    ")
	fmt.Println("$$  __$$\\ $$$\\    $$$ |$$  __$$\\ $$  __$$\\                                                   ")
	fmt.Println("$$ /  \\__|$$$$\\  $$$$ |$$ /  \\__|$$ |  $$ | $$$$$$\\   $$$$$$\\   $$$$$$\\   $$$$$$\\   $$$$$$\\  ")
	fmt.Println("$$ |      $$\\$$\\$$ $$ |\\$$$$$$\\  $$$$$$$  |$$  __$$\\  \\____$$\\ $$  __$$\\ $$  __$$\\ $$  __$$\\ ")
	fmt.Println("$$ |      $$ \\$$$  $$ | \\____$$\\ $$  __$$< $$$$$$$$ | $$$$$$$ |$$ /  $$ |$$$$$$$$ |$$ |  \\__|")
	fmt.Println("$$ |  $$\\ $$ |\\$  /$$ |$$\\   $$ |$$ |  $$ |$$   ____|$$  __$$ |$$ |  $$ |$$   ____|$$ |      ")
	fmt.Println("\\$$$$$$  |$$ | \\_/ $$ |\\$$$$$$  |$$ |  $$ |\\$$$$$$$\\ \\$$$$$$$ |$$$$$$$  |\\$$$$$$$\\ $$ |      ")
	fmt.Println(" \\______/ \\__|     \\__| \\______/ \\__|  \\__| \\_______| \\_______|$$  ____/  \\_______|\\__|      ")
	fmt.Println("                                                               $$ |                          ")
	fmt.Println("                                                               $$ |                          ")
	fmt.Println("                                                               \\__|                          ")

	// Check if flag for Google is ste
	dbFlag := flag.String("db", "", "-db Use either a database on the local machine or build a db from Google search yes|no")
	// var customDBFlag = flag.String("-dbname", "cdnreaper", "-dbname Name of custom database. Default <cdnreaper>")

	flag.Parse()
	fmt.Println(GetSiteInfo(*dbFlag))

}
