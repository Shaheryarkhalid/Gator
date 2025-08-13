package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := handlerLogout(s, cmd)
	if err != nil {
		return fmt.Errorf("Error logging out user: %v", err)
	}
	err = s.db.DeleteAllUsers(context.Background())
	err = s.db.DeleteAllFeeds(context.Background())
	err = s.db.DeleteAllPosts(context.Background())
	if err != nil {
		return fmt.Errorf("Error reseting the database: %v", err)
	}
	fmt.Println("Database reset successfull.")
	return nil
}
