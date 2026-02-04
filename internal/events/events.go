package events

import (
	"encoding/json"
	"fmt"
	"githubFetcher/internal/fetch"
	"githubFetcher/internal/parsing"
	"os"
	"strings"
)

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

func HandleEvents() {
	username, showPrivate, perPage, filterByName, err := parsing.ParseEventsArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := fetch.FetchEvents(username, perPage)
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

func HandleUser() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: program -user <username>")
		return
	}

	username := os.Args[2]
	resp, err := fetch.FetchUser(username)
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
