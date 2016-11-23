package main

import (
	"github.com/178inaba/tcas-post/client"
	"github.com/178inaba/tcas-post/conf"
	log "github.com/Sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	hostName = kingpin.Arg("host", "Broadcast host name.").Required().String()
	confPath = kingpin.Flag("config", "Config toml file path.").Default("etc/conf.toml").Short('c').String()
)

func init() {
	kingpin.Parse()
}

func main() {
	cf, err := conf.LoadConf(*confPath)
	if err != nil {
		log.Fatal(err)
	}

	c, err := client.NewClient(cf.Username, cf.Password)
	if err != nil {
		log.Fatal("Fail create client.")
	}

	err = c.Auth()
	if err != nil {
		log.Fatal(err)
	}

	movieID, err := c.GetMovieID(*hostName)
	if err != nil {
		log.Fatal(err)
	}

	lastCommentID, err := c.GetLastCommentID(*hostName, movieID)
	if err != nil {
		log.Fatal(err)
	}

	log.Info(lastCommentID)
}
