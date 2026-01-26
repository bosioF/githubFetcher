package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const BASE_ENDPOINT = "https://api.github.com/users/"

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	}
	Public    bool   `json:"public"`
	CreatedAt string `json:"created_at"`
}

func main() {
	var events []Event
	var showPrivate = false
	var filterByName = ""
	var validEvents = 0

	username, showPrivate, filterByName, err := parseArgs()
	if err != nil {
		fmt.Println(err)
	}

	resp, err := fetchGithub(username)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := parseJSON(resp, &events)
	if err != nil {
		fmt.Println("json error: ", err)
		return
	}

	for i, e := range res {
		e.CreatedAt = strings.Replace(e.CreatedAt, "T", " ", 1)
		e.CreatedAt = strings.Replace(e.CreatedAt, "Z", " ", 1)
		if !showPrivate && !e.Public {
			continue
		}
		if filterByName != "" && e.Repo.Name != filterByName {
			if i == len(res)-1 && validEvents == 0 {
				fmt.Printf("No valid events found")
			}
			continue
		}
		validEvents++
		fmt.Printf("Event %d {Type: %s, Repo: %s, Public: %t, Created_At: %s}\n", i+1, e.Type, e.Repo.Name, e.Public, e.CreatedAt)
	}
}

func fetchGithub(user string) (string, error) {
	userApiEndpoint := BASE_ENDPOINT + user + "/events"

	resp, err := http.Get(userApiEndpoint)
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

func parseJSON(jsonData string, events *[]Event) ([]Event, error) {

	err := json.Unmarshal([]byte(jsonData), events)
	if err != nil {
		fmt.Println("err decoding json")
		return *events, err
	}

	return *events, nil
}

func parseArgs() (username string, showPrivate bool, filterByName string, err error) {
	args := os.Args

	if len(args) < 2 || len(args) > 4 {
		return "", false, "", fmt.Errorf("wrong args number. usage: ./main.exe <username> <show private? (y/n) default: false> [filter by name]")
	}

	username = os.Args[1]

	showPrivateChoice := "n"
	if len(args) >= 3 {
		showPrivateChoice = os.Args[2]
	}
	if showPrivateChoice != "y" && showPrivateChoice != "n" {
		return "", false, "", fmt.Errorf("invalid argument [2]")
	}
	showPrivate = showPrivateChoice == "y"

	if len(args) == 4 {
		filterByName = os.Args[3]
	}

	return username, showPrivate, filterByName, nil
}
