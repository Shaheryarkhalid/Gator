package main

import (
	"context"
	"fmt"
	"github.com/Shaheryarkhalid/Gator/internal/database"
	"strings"
	"time"
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
}

func handlerClear(s *state, cmd command, _ database.User) error {
	err := s.db.DeleteAllPosts(context.Background())
	if err != nil {
		return fmt.Errorf("Error removing the posts: %w", err)
	}
	fmt.Println("All posts removed successfully.")
	return nil
}
