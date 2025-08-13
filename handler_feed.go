package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shaheryarkhalid/Gator/internal/database"
	"github.com/google/uuid"
	"strings"
	"time"
)

func handlerFollowing(s *state, _ command, user database.User) error {
	currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Unable to get current user: %w", err)
	}
	feedfollows, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("Error getting the feeds followed by currentUser: %w", err)
	}
	if len(feedfollows) == 0 {
		fmt.Printf("No feeds followed by current user '%v'\n", user.Name)
		return nil
	}
	fmt.Println("Feeds followed by the current user: ")
	for _, feedfollow := range feedfollows {
		fmt.Printf("    -- %v\n", feedfollow.FeedName)
	}
	return nil
}

func handlerUnFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 || strings.Trim(cmd.args[0], " ") == "" {
		return fmt.Errorf("Unfollow command expects one argument 'url'\n usage: Gator unfollow <url> ")
	}
	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil || feed.Url == "" {
		return fmt.Errorf("Given Url doesnot match any feed in database. %v", err)
	}
	if feed.UserID == user.ID {
		err := s.db.DeleteFeedById(context.Background(), feed.ID)
		if err != nil {
			return fmt.Errorf("Error while trying to delete given feed.%v", err)
		}
		fmt.Printf("Given feed '%v' removed successfully.\n", url)
		return nil
	}
	feedfollow, err := s.db.GetFeedFollowByFeedIdAndUserId(context.Background(), database.GetFeedFollowByFeedIdAndUserIdParams{FeedID: feed.ID, UserID: user.ID})
	if err != nil {
		fmt.Printf("Error happened while trying to get feed follow for given url: %v\n", err)
	}
	if feedfollow.UserID == user.ID {
		err := s.db.DeleteFeedFollowById(context.Background(), feedfollow.ID)
		if err != nil {
			return fmt.Errorf("Error deleting the feed follow record for given url: %v", err)
		}
		fmt.Printf("Given feed '%v' removed successfully.\n", url)
		return nil
	}
	return fmt.Errorf("Given url '%v' does not belong to current user '%v'", url, user.Name)
}
func handlerFollow(s *state, cmd command, _ database.User) error {
	if len(cmd.args) < 1 || strings.Trim(cmd.args[0], " ") == "" {
		return fmt.Errorf("follow command expects one argument 'url'\n usage: Gator follow <url> ")
	}
	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Error getting the feed by Url: %v", err)
	}
	if feed.Name == "" {
		return fmt.Errorf("No Feed by provided url: '%v'", url)
	}
	currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}
	feedFollow := database.FeedFollow{
		ID:        uuid.New(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	insertedFeedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams(feedFollow))
	if err != nil {
		return fmt.Errorf("Already following the feed: %w", err)
	}
	fmt.Printf("User: %v \n", insertedFeedFollow.UserName)
	fmt.Printf("Feed: %v \n", insertedFeedFollow.FeedName)
	return nil
}
func handlerFeeds(s *state, _ command, _ database.User) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting the feeds: %w", err)
	}
	if len(feeds) == 0 {
		fmt.Println("No feeds added yet. Run below command to add feed. \n - Gator addfeed <feed name> <url>")
		return nil
	}
	formatedFeeds, err := json.MarshalIndent(feeds, "", "  ")
	fmt.Println(string(formatedFeeds))
	return nil
}

func handlerAddFeed(s *state, cmd command, _ database.User) error {
	if len(cmd.args) < 2 || strings.Trim(cmd.args[0], " ") == "" || strings.Trim(cmd.args[1], " ") == "" {
		return fmt.Errorf("addfeed command expects two arguments 'name', 'url'\n usage: Gator addfeed <name> <url>")
	}
	if s.config.CurrentUserName == "" {
		return fmt.Errorf("Please login or register first to create feed.")
	}
	currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}
	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		Name:      cmd.args[0],
		Url:       cmd.args[1],
	}
	createdFeed, err := s.db.CreateFeed(context.Background(), newFeed)
	if err != nil {
		return fmt.Errorf("Error creating the feed: %w", err)
	}
	feedFollow := database.FeedFollow{
		ID:        uuid.New(),
		UserID:    currentUser.ID,
		FeedID:    createdFeed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams(feedFollow))
	if err != nil {
		return fmt.Errorf("Error creating feed follow record: %w", err)
	}

	formatedFeed, err := json.MarshalIndent(createdFeed, "", "  ")
	if err != nil {
		return fmt.Errorf("Error formating the newFeed: %w", err)
	}

	fmt.Println(string(formatedFeed))
	return nil
}
