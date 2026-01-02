package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Shaheryarkhalid/Gator/internal/config"
	"github.com/Shaheryarkhalid/Gator/internal/databse"
	_ "github.com/lib/pq"
)

type state 	struct{
	cnfg *config.Config
	db *databse.Queries
}

type command struct{
	name string
	arguments []string
}
func main(){
	s := state{}
	cnfg := config.Config{}
	cnfg, err := cnfg.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s.cnfg = &cnfg
	db, err := sql.Open("postgres", cnfg.DbUrl)
	s.db = databse.New(db)
	if err != nil {
		fmt.Printf("Error: Trying to establish connection to the database. Please check your connection string.\n%v\n", err)
		os.Exit(1)
	}
	cmds := commands{
		availableCommands: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(&s, handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(&s,  handlerFollow))
	cmds.register("following", middlewareLoggedIn(&s, handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(&s, handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(&s, handlerBrowse))

	args := os.Args
	if len(args) < 2{
		fmt.Println("No arguments passed to the command.")
		os.Exit(1)
	}
	userCmd := command{
		name: args[1], 
		arguments: args[2:],
	}
	err = cmds.run(&s, userCmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
