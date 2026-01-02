package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Shaheryarkhalid/Gator/internal/databse"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerLogin(s *state , cmd command)error{
	if len(cmd.arguments) ==0 || cmd.arguments[0] == "" || cmd.arguments[0] == " "{
		fmt.Println("Login command must be given a username.")
		os.Exit(1)
	}
	username := cmd.arguments[0]
	_, err :=s.db.GetUser(context.Background(), username)
	if err != nil {
		os.Exit(1)
		fmt.Println(err)
	}
	err = s.cnfg.SetUser(username)
	if err != nil {
		return err 
	}
	fmt.Printf("%v logged in.\n",username )
	return  nil
}

func handlerRegister(s *state , cmd command)error{
	if len(cmd.arguments) < 1{
		return fmt.Errorf("Unknown name given to register command.")
	}
	username := cmd.arguments[0]
	if username == "" || username == " "{
		return fmt.Errorf("Unknown name given to register command.")
	}
	user, err := s.db.CreateUser(context.Background(), databse.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: username})
	if err != nil {

		if strings.Contains(err.Error(), "exist"){
			fmt.Println("User Already exists")
			os.Exit(1)
		}

		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505"{
				fmt.Println("User Already exists")
				os.Exit(1)
			}
		}
		return err
		
	}
	err = s.cnfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("'%v' created successfully.\n", user.Name)
	return nil 
}

func handlerReset(s *state, _ command)error{
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil 
}

func handlerUsers(s *state, _ command)error{
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, user := range users{
		fmt.Printf("%v", user.Name)
		if user.Name == s.cnfg.CurrentUsername{
			fmt.Printf(" (current)\n")
		}
	}

	return nil
}
func handlerAgg(s *state, c command)error{
	if len(c.arguments) < 1{
		fmt.Println("time_between_reqs must be given to agg command.")
		os.Exit(1)
	}
	t, err := time.ParseDuration(c.arguments[0])
	if err != nil {
		fmt.Println("time_between_reqs must be valid time  like 1s, 1m, 1h.")
		os.Exit(1)
	}
	fmt.Printf("Collecting feeds every %v\n", t.String())
	tiker := time.NewTicker(t)
	for  range tiker.C{
		scrapeFeeds(s)
	}
	return nil
}

func handlerAddFeed(s *state, c command, user databse.User)error{
	if len(c.arguments) < 2{
		fmt.Println("\"name\" and \"url\" must be given to addfeed command.")
		os.Exit(1)
	}
	name:= c.arguments[0]
	url := c.arguments[1]
	if name == "" || name == " " || url == "" || url == " "{
		fmt.Println("\"name\" and \"url\" must be given to addfeed command.")
		os.Exit(1)
	}
	user, err:=  s.db.GetUser(context.Background(), s.cnfg.CurrentUsername)
	if err != nil {
		fmt.Println("User Must be longed in to add feed.")
		os.Exit(1)
	}
	feedParams := databse.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		fmt.Println("Error: Trying to add feed to db")
		fmt.Println(err)
		os.Exit(1)
	}
	feedFollowParams := databse.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:  user.ID,
		FeedID:  feed.ID,
	}
	_ , err = s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		fmt.Println("Error: Trying to create follow record for given feed")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(feed)
	return  nil
}

func handlerFeeds(s *state, c command)error{
	feeds, err := s.db.GetFeedsByUsers(context.Background())
	if err != nil {
		fmt.Println("Error: Trying to get feeds from databse.")
		fmt.Println(err)
		os.Exit(1)
	}
	for _, feed := range feeds{
		fmt.Printf("Feed Name: %v\n", feed.Name)
		fmt.Printf("Feed Url: %v\n", feed.Url)
		if feed.UserName == s.cnfg.CurrentUsername{
			fmt.Printf("Created By: %v(current)\n", feed.UserName)
			continue
		}
		fmt.Printf("Created By: %v\n", feed.UserName)
	}

	return nil
}

func handlerFollow(s *state, c command, user databse.User)error{
	if len(c.arguments) < 1{
		fmt.Println("url must be specified as command argument.")
		os.Exit(1)
	}
	feed, err := s.db.GetFeedsByUrl(context.Background(), c.arguments[0])
	if err != nil {
		fmt.Println("Given url does not belong to any feed in databse.")
		fmt.Println(err)
		os.Exit(1)
	}
	feedFollowParams := databse.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID: feed.ID,
		UserID: user.ID,
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		fmt.Println("Error: Trying to creat feed follow record in databse.")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Feed: %v\n", feedFollow.FeedName)
	fmt.Printf("Created By: %v\n", feedFollow.UserName)
	return nil
}

func handlerFollowing(s *state, c command, user databse.User)error{
	user_feeds, err:= s.db.GetFeedFollowsForUser(context.Background(), s.cnfg.CurrentUsername)
	if err != nil {
		fmt.Println("Error: Trying to get Feed follows for current user.")
		fmt.Println(err)
		os.Exit(1)
	}
	if len(user_feeds) == 0{
		fmt.Println("No feeds found for current user.")
	}
	for _, feedFollow := range user_feeds{
		fmt.Printf("Feed: %v\n", feedFollow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, c command, user databse.User)error{
	if len(c.arguments) < 1{
		fmt.Println("Url must  be given to unfollow command.")
		os.Exit(1)
	}
	feedUrl := c.arguments[0]
	if feedUrl == "" || feedUrl == " "{
		fmt.Println("Url must  be given to unfollow command.")
		os.Exit(1)
	}
	feed, err := s.db.GetFeedsByUrl(context.Background(), feedUrl)
	if err != nil {
		fmt.Println("Url not found in databse")
		fmt.Println(err)
		os.Exit(1)
	}
	deleteFeedFollowParams := databse.DeleteFeedFollowParams{
		FeedID: feed.ID,
		UserID:  user.ID,
	}
	err = s.db.DeleteFeedFollow(context.Background(), deleteFeedFollowParams)
	if err != nil {
		fmt.Println("Error: Trying to delete follow record for given feed url")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("\"%v\" unfollowed successfully.\n", feedUrl)
	return nil
}
func handlerBrowse(s *state, c command, user databse.User)error{
	limit := 2
	if len(c.arguments) > 0{
		var err error	
		limit, err =strconv.Atoi( c.arguments[0] )
		if err!= nil {
			fmt.Println("Wrong limit was given, must be a number")
			os.Exit(1)
		}
	}
	postsForUserParams := databse.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), postsForUserParams)
	if err != nil {
		fmt.Println("Error: Trying to get posts from databse.")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(posts)
	return nil
}
