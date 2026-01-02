package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Shaheryarkhalid/Gator/internal/databse"
)

func middlewareLoggedIn(st *state, handler func(s *state, cmd command, user databse.User) error) func(*state, command) error{
	return func(s *state, c command) error {
		if st.cnfg.CurrentUsername == "" {
			fmt.Println("User not logged in.")
			os.Exit(1)
		}
		user, err := st.db.GetUser(context.Background(), st.cnfg.CurrentUsername)
		if err != nil {
			fmt.Println("Error: Trying to get current user from database.")
			fmt.Println(err)
			os.Exit(1)
		}

		return handler(s , c , user)
	}

}
