package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Shaheryarkhalid/Gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	var err error
	if len(cmd.args) < 1 || strings.Trim(cmd.args[0], " ") == "" {
		limit = 2
	} else {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			limit = 2
		}
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limit)})
	if err != nil {
		return fmt.Errorf("Error getting user posts: %w", err)
	}
	if len(posts) == 0 {
		fmt.Println("User has no posts. Please try following new feeds, if already added try running aggregate command to fetch.")
		return nil
	}
	fmt.Println("Posts from user followed feeds:  \n")
	for _, post := range posts {
		fmt.Println("{")
		fmt.Printf("   Title: %v\n", post.Title)
		fmt.Printf("   Description: %v\n", post.Description)
		fmt.Printf("   Published At: %v\n", post.PublishedAt)
		fmt.Println("}")
	}
	return nil
}
