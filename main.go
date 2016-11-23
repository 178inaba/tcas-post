package main

import (
	"math/rand"
	"time"

	"github.com/178inaba/tcas-post/conf"
	"github.com/178inaba/twitcasting"
	log "github.com/Sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	hostNameArg  = kingpin.Arg("host", "Broadcast host name.").Required().String()
	confPathFlag = kingpin.Flag("config", "Config toml file path.").Default("etc/conf.toml").Short('c').String()
)

func init() {
	kingpin.Parse()
}

func main() {
	hostName := *hostNameArg
	confPath := *confPathFlag

	cf, err := conf.LoadConf(confPath)
	if err != nil {
		log.Fatalf("LoadConf error: %v.", err)
	}

	c, err := twitcasting.NewClient(cf.Username, cf.Password)
	if err != nil {
		log.Fatalf("NewClient error: %v.", err)
	}

	err = c.Auth()
	if err != nil {
		log.Fatalf("Auth error: %v.", err)
	}

	rand.Seed(time.Now().UnixNano())
	for {
		time.Sleep(time.Minute * 1)

		movieID, err := c.GetMovieID(hostName)
		if err != nil {
			log.Errorf("GetMovieID error: %v.", err)
			continue
		}

		comment := cf.Comments[rand.Intn(len(cf.Comments))]
		err = c.PostComment(comment, hostName, movieID)
		if err != nil {
			log.Errorf("PostComment error: %v.", err)
		} else {
			log.Infof("PostComment success!: %s", comment)
		}
	}
}
