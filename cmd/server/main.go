package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"redcellpartners.com/users-posts-api/commands/start"
)

func main() {
	var (
		err error
		app = cli.NewApp()
	)

	app.Description = "User and Posts REST API"
	app.Commands = []cli.Command{
		start.StartCommand(),
	}

	if err = app.Run(os.Args); err != nil {
		log.Println("error running cli app:", err.Error())
		os.Exit(1)
	}
}
