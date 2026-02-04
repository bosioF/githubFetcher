package fetch

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const BaseEndpoint = "https://api.github.com"

func FetchEvents(user string, perPage int) (string, error) {
	url := ""
	if perPage > 0 {
		url = BaseEndpoint + "/users/" + user + "/events?per_page=" + strconv.Itoa(perPage)
	} else {
		url = BaseEndpoint + "/users/" + user + "/events"
	}
	return GetReq(url)
}

func FetchUser(user string) (string, error) {
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
