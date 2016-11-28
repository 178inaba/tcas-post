package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAccount(t *testing.T) {
	var existConfigDir bool
	accountFilePath, err := getAccountFilePath()
	if err == nil {
		_, err := os.Stat(accountFilePath)
		if err == nil {
			existConfigDir = true
			err := os.Rename(accountFilePath, accountFilePath+".tmp")
			defer os.Rename(accountFilePath+".tmp", accountFilePath)
			assert.NoError(t, err)
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
		err := os.Mkdir(filepath.Dir(accountFilePath), os.ModePerm)
		defer os.Remove(filepath.Dir(accountFilePath))
		assert.NoError(t, err)
	}

	uname, pass := "test", "pass"
	saveAc := &Account{Username: uname, Password: pass}
	err = saveAc.Save()
	defer os.Remove(accountFilePath)
	assert.NoError(t, err)
	assert.Equal(t, uname, saveAc.Username)
	assert.Equal(t, pass, saveAc.Password)

	account, err := LoadAccount()
	assert.NoError(t, err)
	assert.Equal(t, uname, account.Username)
	assert.Equal(t, pass, account.Password)
}
