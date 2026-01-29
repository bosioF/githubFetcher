package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const BaseEndpoint = "https://api.github.com"

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	}
	Public    bool   `json:"public"`
	CreatedAt string `json:"created_at"`
}

type User struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	Company     string `json:"company"`
	Blog        string `json:"blog"`
	Location    string `json:"location"`
	Email       string `json:"email"`
	Bio         string `json:"bio"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "-e", "-events", "--events":
		handleEvents()
	case "-u", "-user", "--user":
		handleUser()
	case "-h", "-help", "--help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func handleEvents() {
	username, showPrivate, filterByName, err := parseEventsArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := fetchEvents(username)
	if err != nil {
		fmt.Println(err)
		return
	}

	var events []Event
	err = json.Unmarshal([]byte(resp), &events)
	if err != nil {
		fmt.Println("json error:", err)
		return
	}

	validEvents := 0
	for i, e := range events {
		e.CreatedAt = strings.Replace(e.CreatedAt, "T", " ", 1)
		e.CreatedAt = strings.Replace(e.CreatedAt, "Z", " ", 1)

		if !showPrivate && !e.Public {
			continue
		}
		if filterByName != "" && e.Repo.Name != filterByName {
			if i == len(events)-1 && validEvents == 0 {
				fmt.Printf("No valid events found\n")
			}
			continue
		}
		validEvents++
		fmt.Printf("Event %d {Type: %s, Repo: %s, Public: %t, Created_At: %s}\n",
			i+1, e.Type, e.Repo.Name, e.Public, e.CreatedAt)
	}
}

func handleUser() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: program -user <username>")
		return
	}

	username := os.Args[2]
	resp, err := fetchUser(username)
	if err != nil {
		fmt.Println(err)
		return
	}

	var user User
	err = json.Unmarshal([]byte(resp), &user)
	if err != nil {
		fmt.Println("json error:", err)
		return
	}

	fmt.Printf("\nGitHub User: %s\n", user.Login)
	fmt.Printf("Name: %s\n", user.Name)
	if user.Company != "" {
		fmt.Printf("Company: %s\n", user.Company)
	}
	if user.Location != "" {
		fmt.Printf("Location: %s\n", user.Location)
	}
	if user.Blog != "" {
		fmt.Printf("Blog: %s\n", user.Blog)
	}
	if user.Bio != "" {
		fmt.Printf("Bio: %s\n", user.Bio)
	}
	fmt.Printf("Public Repos: %d\n", user.PublicRepos)
	fmt.Printf("Followers: %d\n", user.Followers)
	fmt.Printf("Following: %d\n", user.Following)
}

func parseEventsArgs() (username string, showPrivate bool, filterByName string, err error) {
	args := os.Args[2:]

	if len(args) < 1 || len(args) > 3 {
		return "", false, "", fmt.Errorf("usage: program -events <username> [y/n for private] [repo filter]")
	}

	username = args[0]

	showPrivateChoice := "n"
	if len(args) >= 2 {
		showPrivateChoice = args[1]
	}
	if showPrivateChoice != "y" && showPrivateChoice != "n" {
		return "", false, "", fmt.Errorf("invalid argument: use 'y' or 'n' for showing private events")
	}
	showPrivate = showPrivateChoice == "y"

	if len(args) == 3 {
		filterByName = args[2]
	}

	return username, showPrivate, filterByName, nil
}

func fetchEvents(user string) (string, error) {
	url := BaseEndpoint + "/users/" + user + "/events"
	return GetReq(url)
}

func fetchUser(user string) (string, error) {
	url := BaseEndpoint + "/users/" + user
	return GetReq(url)
}

func GetReq(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github returned %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if len(bodyBytes) == 2 {
		return "", fmt.Errorf("user does not exist")
	}

	return string(bodyBytes), nil
}

func printUsage() {
	fmt.Println("Usage: program <command> [arguments]")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  -e, -events, --events    Show user events")
	fmt.Println("  -u, -user, --user        Show user information")
	fmt.Println("  -h, -help, --help        Show this help message")
	fmt.Println("\nUse 'program -help' for more details")
}

func printHelp() {
	printUsage()
	fmt.Println("\nExamples:")
	fmt.Println("  program -events torvalds")
	fmt.Println("  program -events torvalds y")
	fmt.Println("  program -events torvalds n linux")
	fmt.Println("  program -user torvalds")
}
