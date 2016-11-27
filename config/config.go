package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	configDirName   = ".tcpost"
	accountFileName = "account.toml"
)

// Exist is ...
type Exist interface {
	Exist() bool
}

type notExistError struct {
	msg string
}

func wrapNotExistError(err error) *notExistError {
	return &notExistError{msg: err.Error()}
}

func (e *notExistError) Error() string { return e.msg }

func (e *notExistError) Exist() bool { return false }

// Account is ...
type Account struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// LoadAccount is ...
func LoadAccount() (*Account, error) {
	accountFilePath, err := getAccountFilePath()
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(accountFilePath)
	if err != nil {
		return nil, wrapNotExistError(err)
	}

	var account *Account
	_, err = toml.DecodeFile(accountFilePath, &account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func getAccountFilePath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, configDirName, accountFileName), nil
}
