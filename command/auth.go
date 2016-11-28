package command

import (
	"github.com/178inaba/tcpost/config"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

// Logout is ...
func Logout(c *cli.Context) error {
	account, err := config.LoadAccount()
	if err != nil {
		if _, ok := err.(config.Exist); ok {
			return cli.NewExitError("You are not logged in.", failExitStatusCode)
		}

		return cli.NewExitError(err.Error(), failExitStatusCode)
	}

	username := account.Username

	if err := config.RemoveAccountFile(); err != nil {
		return cli.NewExitError(err.Error(), failExitStatusCode)
	}

	log.Infof("Logout %s", username)

	return nil
}
