package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jasonwashburn/gator/internal/config"
	"github.com/jasonwashburn/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.ConfigFile
	db  *database.Queries
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

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("list does not take any arguments")
	}

	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	currentUser := s.cfg.CurrentUserName
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("reset does not take any arguments")
	}

	if err := s.db.ResetUsers(context.Background()); err != nil {
		return fmt.Errorf("failed to reset users: %w", err)
	}

	fmt.Println("Users reset")
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login requires a username")
	}

	userName := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.Name == "" {
		return fmt.Errorf("user not found")
	}

	s.cfg.SetUser(userName)
	fmt.Printf("Logged in as %s\n", userName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("register requires a username")
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	fmt.Printf("User created: %s\nUser Info: %+v\n", user.Name, user)
	return nil
}

func main() {
	configFile, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	s := &state{
		cfg: &configFile,
	}

	db, err := sql.Open("postgres", s.cfg.DbURL)
	if err != nil {
		log.Fatal(err)
	}
	s.db = database.New(db)

	commands := &commands{
		allCommands: make(map[string]func(*state, command) error),
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
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
