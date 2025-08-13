package main

import (
	"fmt"
	"strings"
)

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmnd command) error {
	cmnd.name = strings.Trim(cmnd.name, " ")
	if cmnd.name == "" {
		return fmt.Errorf("Command cannot be empty.")
	}
	method, ok := c.registeredCommands[cmnd.name]
	if !ok {
		return fmt.Errorf(`"%v" Command not found.`, cmnd.name)
	}
	err := method(s, cmnd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	cmndName := strings.Trim(name, " ")
	c.registeredCommands[cmndName] = f
}

type command struct {
	name string
	args []string
}
