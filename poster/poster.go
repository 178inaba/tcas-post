package poster

import (
	"fmt"

	"github.com/178inaba/tcpost/config"
	"github.com/178inaba/twitcasting"
	log "github.com/Sirupsen/logrus"
	"github.com/howeyc/gopass"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const (
	failExitStatusCode = 1
)

// Poster is ...
type Poster struct {
	client *twitcasting.Client
}

// NewPoster is ...
func NewPoster() (*Poster, error) {
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
		return nil, errors.Wrap(err, "new twitcasting client error")
	}

	p := &Poster{client: client}

	return p, nil
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

// Action is ...
func (p *Poster) Action(c *cli.Context) error {
	target := c.Args().Get(0)
	comment := c.Args().Get(1)

	err := p.client.Auth()
	if err != nil {
		config.RemoveAccountFile()
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

	log.Infof("Post comment: %s", comment)

	return nil
}
