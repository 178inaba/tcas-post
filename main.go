package main

import (
	"fmt"
	"os"

	"github.com/178inaba/twitcasting"
	log "github.com/Sirupsen/logrus"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
)

const (
	failExitStatusCode = 1
)

type poster struct {
	client *twitcasting.Client
}

func main() {
	var username, password string

	// Login username from stdin.
	fmt.Print("Username: ")
	fmt.Scanln(&username)

	// Password from stdin.
	fmt.Print("Password: ")
	pBytes, err := gopass.GetPasswd()
	if err != nil {
		log.Fatalf("Input Password error: %v.", err)
	}

	password = string(pBytes)

	client, err := twitcasting.NewClient(username, password)
	if err != nil {
		log.Fatalf("NewClient error: %v.", err)
	}

	p := &poster{client: client}

	app := cli.NewApp()
	app.Name = "tcpost"
	app.HelpName = app.Name
	app.Usage = "Post comment to TwitCasting."
	app.Version = "1.0.0"
	app.Action = p.action

	err = app.Run(os.Args)
	if err != nil {
		log.Fatalf("Run error: %v", err)
	}
}

func (p *poster) action(c *cli.Context) error {
	target := c.Args().Get(0)
	comment := c.Args().Get(1)

	err := p.client.Auth()
	if err != nil {
		return cli.NewExitError(err.Error(), failExitStatusCode)
	}

	movieID, err := p.client.GetMovieID(target)
	if err != nil {
		return cli.NewExitError(err.Error(), failExitStatusCode)
	}

	err = p.client.PostComment(comment, target, movieID)
	if err != nil {
		return cli.NewExitError(err.Error(), failExitStatusCode)
	}

	log.Infof("PostComment success!: %s", comment)

	return nil
}
