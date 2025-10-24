package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jfcisco/gator/internal/config"
	_ "github.com/lib/pq"
)

type commands struct {
	all map[string]func(*state, command) error
}

func configureHandlers() *commands {
	handlers := commands{all: make(map[string]func(*state, command) error)}

	handlers.register("login", handlerLogin)
	handlers.register("register", handlerRegister)
	handlers.register("reset", handlerReset)
	handlers.register("users", handlerUsers)
	handlers.register("agg", handlerAgg)
	handlers.register("feeds", handlerFeeds)

	// Current user specific commands
	handlers.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	handlers.register("follow", middlewareLoggedIn(handlerFollow))
	handlers.register("following", middlewareLoggedIn(handlerFollowing))
	handlers.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	handlers.register("browse", middlewareLoggedIn(handlerBrowse))

	return &handlers
}

func main() {
	cfg := config.Read()
	stat := &state{config: &cfg}
	db, err := stat.LoadDb(&cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("error opening database: %w", err))
	}
	defer db.Close()

	handlers := configureHandlers()

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Please enter a command")
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}
	if err := handlers.run(stat, cmd); err != nil {
		log.Fatal(fmt.Errorf("error: %w", err))
	}
}
