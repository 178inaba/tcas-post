package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	configDirName   = ".tcpost"
	accountFileName = "account.toml"
	encryptKey      = "Go TwitCasting post comment app."
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

	if _, err = os.Stat(accountFilePath); err != nil {
		return nil, wrapNotExistError(err)
	}

	var account *Account
	if _, err = toml.DecodeFile(accountFilePath, &account); err != nil {
		return nil, err
	}

	if err := account.decrypt(); err != nil {
		return nil, err
	}

	return account, nil
}

// RemoveAccountFile ...
func RemoveAccountFile() error {
	accountFilePath, err := getAccountFilePath()
	if err != nil {
		return err
	}

	if err := os.Remove(accountFilePath); err != nil {
		return err
	}

	return nil
}

func getAccountFilePath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, configDirName, accountFileName), nil
}

// Save is ...
func (a *Account) Save() error {
	accountFilePath, err := getAccountFilePath()
	if err != nil {
		return err
	}

	if _, err = os.Stat(filepath.Dir(accountFilePath)); err != nil {
		if err := os.Mkdir(filepath.Dir(accountFilePath), os.ModePerm); err != nil {
			return err
		}
	}

	accountFile, err := os.Create(accountFilePath)
	if err != nil {
		return err
	}

	defer accountFile.Close()

	if err := a.encrypt(); err != nil {
		return err
	}

	if err = toml.NewEncoder(accountFile).Encode(a); err != nil {
		return err
	}

	if err := a.decrypt(); err != nil {
		return err
	}

	return nil
}

func (a *Account) encrypt() error {
	block, err := aes.NewCipher([]byte(encryptKey))
	if err != nil {
		return err
	}

	a.Username, err = encrypt(block, a.Username)
	if err != nil {
		return err
	}

	a.Password, err = encrypt(block, a.Password)
	if err != nil {
		return err
	}

	return nil
}

func encrypt(block cipher.Block, plaintext string) (string, error) {
	cipherBytes := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherBytes[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherBytes[aes.BlockSize:], []byte(plaintext))

	return hex.EncodeToString(cipherBytes), nil
}

func (a *Account) decrypt() error {
	block, err := aes.NewCipher([]byte(encryptKey))
	if err != nil {
		return err
	}

	a.Username, err = decrypt(block, a.Username)
	if err != nil {
		return err
	}

	a.Password, err = decrypt(block, a.Password)
	if err != nil {
		return err
	}

	return nil
}

func decrypt(block cipher.Block, ciphertext string) (string, error) {
	cipherBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plainBytes := make([]byte, len(cipherBytes[aes.BlockSize:]))
	stream := cipher.NewCTR(block, cipherBytes[:aes.BlockSize])
	stream.XORKeyStream(plainBytes, cipherBytes[aes.BlockSize:])

	return string(plainBytes), nil
}
