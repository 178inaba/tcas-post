package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func TestLoadConf(t *testing.T) {
	var existConfigDir, existAccountFile bool
	accountFilePath, err := getAccountFilePath()
	if err == nil {
		_, err := os.Stat(accountFilePath)
		if err == nil {
			existConfigDir, existAccountFile = true, true
			os.Rename(accountFilePath, accountFilePath+".tmp")
		} else {
			_, err := os.Stat(filepath.Dir(accountFilePath))
			if err == nil {
				existConfigDir = true
			}
		}
	}

	_, err = LoadAccount()
	assert.Error(t, err)

	ex, ok := err.(Exist)
	assert.True(t, ok)
	assert.False(t, ex.Exist())

	if !existConfigDir {
		os.Mkdir(filepath.Dir(accountFilePath), os.ModePerm)
	}

	accountFile, err := os.Create(accountFilePath)
	assert.NoError(t, err)

	uname, pass := "test", "pass"
	err = toml.NewEncoder(accountFile).Encode(Account{Username: uname, Password: pass})
	accountFile.Close()
	assert.NoError(t, err)

	account, err := LoadAccount()
	assert.NoError(t, err)
	assert.Equal(t, uname, account.Username)
	assert.Equal(t, pass, account.Password)

	os.Remove(accountFilePath)

	if existAccountFile {
		os.Rename(accountFilePath+".tmp", accountFilePath)
	} else if !existConfigDir {
		os.Remove(filepath.Dir(accountFilePath))
	}
}
