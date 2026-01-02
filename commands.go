package main

import "fmt"

type commands struct{
	availableCommands map[string] func(*state, command) error
}


func (c *commands)run(s *state, cmd command)error{
	commandHandler, ok := c.availableCommands[cmd.name]
	if !ok {
		return fmt.Errorf("'%v' command does not exist.", cmd.name)
	}
	err := commandHandler(s, cmd)
	if err != nil {
		return fmt.Errorf("Error: Trying to run command '%v'\n%v\n", cmd.name, err)
	}
	return nil 
}


func (c *commands)register(name string, f func(*state, command) error){
	c.availableCommands[name] = f
}

