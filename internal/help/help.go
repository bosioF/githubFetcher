package help

import "fmt"

func PrintUsage() {
	fmt.Println("Usage: program -events <username> [y/n for private, defaults to n] [Number of events, defaults to 5] [repo filter, defaults to \"\"]")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  -e, -events, --events    Show user events")
	fmt.Println("  -u, -user, --user        Show user information")
	fmt.Println("  -h, -help, --help        Show this help message")
	fmt.Println("\nUse 'program -help' for more details")
}

func PrintHelp() {
	PrintUsage()
	fmt.Println("\nExamples:")
	fmt.Println("  program -events torvalds")
	fmt.Println("  program -events torvalds y")
	fmt.Println("  program -events torvalds n linux")
	fmt.Println("  program -events torvalds n 10 linux")
	fmt.Println("  program -user torvalds")
}
