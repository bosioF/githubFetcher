package main

import (
	"fmt"
	"os"

	"githubFetcher/internal/events"
	"githubFetcher/internal/help"
)

func main() {
	if len(os.Args) < 2 {
		help.PrintUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "-e", "-events", "--events":
		events.HandleEvents()
	case "-u", "-user", "--user":
		events.HandleUser()
	case "-h", "-help", "--help":
		help.PrintHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		help.PrintUsage()
	}
}
