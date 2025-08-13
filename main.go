package main

import (
	"database/sql"
	"fmt"
	"github.com/Shaheryarkhalid/Gator/internal/config"
	"github.com/Shaheryarkhalid/Gator/internal/database"
	_ "github.com/lib/pq"
	"os"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	dbQueries := database.New(db)
	if err != nil {
		fmt.Println("Error happened while trying to connect to the database.")
		os.Exit(1)
	}
	programState := state{
		dbQueries,
		&cfg,
	}
	programCommands := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	programCommands.register("register", handlerRegister)
	programCommands.register("login", handlerLogin)
	programCommands.register("logout", handlerLogout)
	programCommands.register("users", middlewareLoggedIn(handlerUsers))

	programCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	programCommands.register("feeds", middlewareLoggedIn(handlerFeeds))
	programCommands.register("follow", middlewareLoggedIn(handlerFollow))
	programCommands.register("unfollow", middlewareLoggedIn(handlerUnFollow))
	programCommands.register("following", middlewareLoggedIn(handlerFollowing))

	programCommands.register("agg", middlewareLoggedIn(handlerAgg))
	programCommands.register("browse", middlewareLoggedIn(handlerBrowse))
	programCommands.register("clear", middlewareLoggedIn(handlerClear))
	programCommands.register("reset", handlerReset)
	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Println("Command name cannot be empty.")
		os.Exit(1)
	}
	cmd := command{name: arguments[1], args: arguments[2:]}
	err = programCommands.run(&programState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
