package main

import (
	"flag"
	"fmt"
)

const Reset = "\033[0m"
const Cyan = "\033[36m"
const Green = "\033[32m"
const Blue = "\033[34m"
const Yellow = "\033[33m"
const Red = "\033[31m"

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
	dirList := flag.String("dw", "", "Specify the wordlist to be used for searching for directories")
	subdomainList := flag.String("sw", "", "Specify the wordlist to be used for searching for directories")
	// var customDBFlag = flag.String("-dbname", "cdnreaper", "-dbname Name of custom database. Default <cdnreaper>")

	flag.Parse()
	fmt.Println(GetSiteInfo(*dbFlag, *dirList, *subdomainList))

}
