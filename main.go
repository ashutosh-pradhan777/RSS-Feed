package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ashutosh-pradhan777/RSS-Feed/internal/config"
	"github.com/ashutosh-pradhan777/RSS-Feed/internal/database"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

type state struct {
	db   *database.Queries
	conf *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commandmap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {

	usercmd := cmd.name
	cmdfunc, ok := c.commandmap[usercmd]
	if !ok {
		return errors.New("No such command")
	}

	err := cmdfunc(s, cmd)
	if err != nil {
		return err
	}

	return nil

}

func (c *commands) register(name string, f func(*state, command) error) {

	c.commandmap[name] = f
}

func handlerLogin(s *state, cmd command) error {

	if cmd.args == nil {
		return errors.New("login handler expects a single argument, <username>")
	}

	if len(cmd.args) != 1 {
		return errors.New("Username is required.")
	}

	_, errx := s.db.GetUser(context.Background(), cmd.args[0])
	if errx == sql.ErrNoRows {
		fmt.Printf("Error: %v", errx)
		os.Exit(1)
	}

	err := s.conf.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been set.\n", cmd.args[0])
	return nil

}

func handlerRegister(s *state, cmd command) error {

	if cmd.args == nil {
		return errors.New("register handler expects a single argument, <username>")
	}

	if len(cmd.args) != 1 {
		return errors.New("Username is required.")
	}

	ctx := context.Background()
	userparams := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Name: cmd.args[0],
	}

	_, err := s.db.GetUser(ctx, cmd.args[0])

	// if err != nil {

	// 	if err == sql.ErrNoRows {
	// 		user, err := s.db.CreateUser(ctx, userparams)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		err2 := s.conf.SetUser(cmd.args[0])
	// 		if err2 != nil {
	// 			return err2
	// 		}

	// 		fmt.Println("User Created.")
	// 		fmt.Printf("%v", user)

	// 	} else {

	// 		os.Exit(1)
	// 	}
	// }

	if err == nil {
		fmt.Printf("username already taken")
		os.Exit(1)
	}

	if err != sql.ErrNoRows {
		return err
	}

	user, err := s.db.CreateUser(ctx, userparams)
	if err != nil {
		return err
	}

	err2 := s.conf.SetUser(cmd.args[0])
	if err2 != nil {
		return err2
	}

	fmt.Println("User Created.")
	fmt.Printf("%v", user)

	return nil

}

func main() {

	conf, err := config.Read()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	var currstate state
	currstate.conf = &conf

	var allcommands commands
	allcommands.commandmap = make(map[string]func(*state, command) error)

	allcommands.register("login", handlerLogin)
	allcommands.register("register", handlerRegister)

	args := os.Args
	if len(args) < 2 {
		fmt.Printf("%v\n", errors.New("Too few arguments"))
		os.Exit(1)
	}
	cmd := command{
		name: args[1],
		args: args[2:],
	}

	db, err := sql.Open("postgres", currstate.conf.DBURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	currstate.db = dbQueries

	err2 := allcommands.run(&currstate, cmd)
	if err2 != nil {
		fmt.Printf("Error: %v\n", err2)
		os.Exit(1)
	}

}
