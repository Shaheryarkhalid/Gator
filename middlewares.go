package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/Shaheryarkhalid/Gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			if strings.Contains(fmt.Sprintf("%v", err), "no rows") {
				return fmt.Errorf("Please login to continue.")
			}
			return fmt.Errorf("Unable to get current user: %w", err)
		}
		return handler(s, cmd, currentUser)
	}
}
