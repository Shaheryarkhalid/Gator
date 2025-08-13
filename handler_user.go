package main

import (
	"context"
	"fmt"
	"strings"
	"time"
	"github.com/Shaheryarkhalid/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogout(s *state, _ command) error {
	err := s.config.RemoveUser()
	if err != nil {
		return fmt.Errorf("Error logging out user: %w", err)
	}
	fmt.Println("User logged out successfully.")
	return nil
}
func handlerUsers(s *state, _ command, _ database.User) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error happened while trying to get users from database.%w", err)
	}
	if len(users) == 0 {
		return fmt.Errorf("No user has been registered yet.")

	}
	for _, user := range users {
		txt := fmt.Sprintf("* %v", user)
		if user == s.config.CurrentUserName {
			txt += " (current)"
		}
		fmt.Println(txt)
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 || strings.Trim(cmd.args[0], " ") == "" {
		return fmt.Errorf("Login commmad expects a single argument user name.")
	}
	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("User by the name '%v' does not exist.", cmd.args[0])
	}
	err = s.config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Error setting user name: %w", err)
	}
	fmt.Println("User has been set.")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 || strings.Trim(cmd.args[0], " ") == "" {
		return fmt.Errorf("Error: Name must be passed as an argument. \nUsage: gator register <name> \n")
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err == nil {
		return fmt.Errorf("User by the name '%v' already exists,", cmd.args[0])
	}
	newUser := database.User{
		ID:        uuid.New(),
		Name:      cmd.args[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	newUser, err = s.db.CreateUser(context.Background(), database.CreateUserParams(newUser))
	if err != nil {
		return fmt.Errorf("Error happened while trying to create User: %w", err)
	}
	fmt.Println("User created successfully.")
	err = s.config.SetUser(newUser.Name)
	if err != nil {
		return err
	}
	fmt.Printf("%v loged in.\n", newUser.Name)
	return nil

}
