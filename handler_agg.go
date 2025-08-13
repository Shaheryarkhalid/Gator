package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Shaheryarkhalid/Gator/internal/database"
)

func handlerAgg(s *state, cmd command, _ database.User) error {
	if len(cmd.args) < 1 || strings.Trim(cmd.args[0], " ") == "" {
		return fmt.Errorf("agg command expects one argument 'time_between_requests'\n usage: Gator agg <time_between_requests(1s,1m, 1h ... etc)> ")
	}
	duration, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Invalid duration given: %w", err)
	}
	fmt.Printf("Collecting feeds every. %v\n", duration)
	interval := time.Tick(duration)
	for ; ; <-interval {
		err = scrapeFeeds(s)
		if err != nil {
			if strings.Contains(err.Error(), "") {
				return err
			}
			fmt.Printf("Error happened: %v\n", err)

		}
	}
	// rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	// if err != nil {
	// 	return err
	// }
	// data, err := json.MarshalIndent(rssFeed, "", "  ")
	// fmt.Println(string(data))
	// return nil
}
