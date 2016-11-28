package main

import (
	"os"

	"github.com/178inaba/tcpost/command"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "tcpost"
	app.HelpName = app.Name
	app.Usage = "Post comment to TwitCasting."
	app.Version = "1.0.0"
	app.Action = command.Post
	app.Commands = []cli.Command{
		{
			Name:   "logout",
			Usage:  "Logout TwitCasting.",
			Action: command.Logout,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Run error: %v", err)
	}
}
