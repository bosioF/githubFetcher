package parsing

import (
	"flag"
	"fmt"
	"os"
)

func ParseEventsArgs() (username string, showPrivate bool, perPage int, filterByName string, err error) {
	eventsCmd := flag.NewFlagSet("events", flag.ExitOnError)
	showPrivateFlag := eventsCmd.Bool("private", false, "show private events")
	perPageFlag := eventsCmd.Int("count", 5, "number of events")
	filterFlag := eventsCmd.String("repo", "", "filter by repository name")

	err = eventsCmd.Parse(os.Args[2:])
	if err != nil {
		return "", false, 0, "", err
	}

	args := eventsCmd.Args()
	if len(args) < 1 {
		return "", false, 0, "", fmt.Errorf("usage: program -events <username> [-private] [-count N] [-repo name]")
	}

	username = args[0]
	showPrivate = *showPrivateFlag
	perPage = *perPageFlag
	filterByName = *filterFlag

	return username, showPrivate, perPage, filterByName, nil
}
