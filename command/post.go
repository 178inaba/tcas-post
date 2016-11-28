package command

import (
	"fmt"

	"github.com/178inaba/tcpost/config"
	"github.com/178inaba/twitcasting"
	log "github.com/Sirupsen/logrus"
	"github.com/howeyc/gopass"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// newTCClient create TwitCasting client.
func newTCClient() (*twitcasting.Client, error) {
	ac, err := config.LoadAccount()
	if err != nil {
		if _, ok := err.(config.Exist); !ok {
			return nil, errors.Wrap(err, "load account error")
		}

		ac, err = inputAccount()
		ac.Save()
	}

	client, err := twitcasting.NewClient(ac.Username, ac.Password)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func inputAccount() (*config.Account, error) {
	var username, password string

	// Login username from stdin.
	fmt.Print("Username: ")
	fmt.Scanln(&username)

	// Password from stdin.
	fmt.Print("Password: ")
	pBytes, err := gopass.GetPasswd()
	if err != nil {
		return nil, errors.Wrap(err, "input password error")
	}

	password = string(pBytes)

	return &config.Account{Username: username, Password: password}, nil
}

// Post is ...
func Post(c *cli.Context) error {
	target := c.Args().Get(0)
	if target == "" {
		return cli.ShowAppHelp(c)
	}

	comment := c.Args().Get(1)
	if comment == "" {
		return cli.NewExitError("post comment not found.", failExitStatusCode)
	}

	client, err := newTCClient()
	if err != nil {
		return cli.NewExitError(err.Error(), failExitStatusCode)
	}

	if err := client.Auth(); err != nil {
		config.RemoveAccountFile()
		return cli.NewExitError(err.Error(), failExitStatusCode)
	}

	movieID, err := client.GetMovieID(target)
	if err != nil {
		return cli.NewExitError(err.Error(), failExitStatusCode)
	}

	if err = client.PostComment(comment, target, movieID); err != nil {
		return cli.NewExitError(err.Error(), failExitStatusCode)
	}

	log.Infof("Post comment: %s", comment)

	return nil
}
