package conf

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Conf is ...
type Conf struct {
	Username string   `toml:"username"`
	Password string   `toml:"password"`
	Comments []string `toml:"comments"`
}

// LoadConf is ...
func LoadConf(path string) (*Conf, error) {
	confBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf *Conf
	_, err = toml.Decode(string(confBytes), &conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
