package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jasonwashburn/gator/internal/config"
)

type state struct {
	Config *config.ConfigFile
}

type command struct {
	command string
	args    []string
}

type commands struct {
	allCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.allCommands[cmd.command]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.command)
	}

	return handler(s, cmd)
}

func (c *commands) register(name string, handler func(*state, command) error) {
	c.allCommands[name] = handler
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login requires a username")
	}

	user := cmd.args[0]
	s.Config.SetUser(user)
	fmt.Printf("Logged in as %s\n", user)
	return nil
}

func main() {
	configFile, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	s := &state{
		Config: &configFile,
	}

	commands := &commands{
		allCommands: make(map[string]func(*state, command) error),
	}
	commands.register("login", handlerLogin)
	userArgs := os.Args
	if len(userArgs) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	command := command{
		command: userArgs[1],
		args:    userArgs[2:],
	}

	err = commands.run(s, command)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
