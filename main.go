package main

import (
	"os"

	"github.com/178inaba/tcpost/poster"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	p, err := poster.NewPoster()

	app := cli.NewApp()
	app.Name = "tcpost"
	app.HelpName = app.Name
	app.Usage = "Post comment to TwitCasting."
	app.Version = "1.0.0"
	app.Action = p.Action

	err = app.Run(os.Args)
	if err != nil {
		log.Fatalf("Run error: %v", err)
	}
}
