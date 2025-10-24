package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

func (c *commands) run(s *state, cmd command) error {
	if s == nil {
		return fmt.Errorf("run: state is nil")
	}

	handler, ok := c.all[cmd.name]
	if !ok {
		return fmt.Errorf("run: attempted to run unregistered command %s", cmd.name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.all[name] = f
}
