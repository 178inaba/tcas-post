package main

import (
	"github.com/178inaba/tcas-post/client"
	"github.com/178inaba/tcas-post/conf"
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

	c, err := client.NewClient(cf.Username, cf.Password)
	if err != nil {
		log.Fatalf("NewClient error: %v.", err)
	}

	err = c.Auth()
	if err != nil {
		log.Fatalf("Auth error: %v.", err)
	}

	movieID, err := c.GetMovieID(hostName)
	if err != nil {
		log.Fatalf("GetMovieID error: %v.", err)
	}

	err = c.PostComment("www", hostName, movieID)
	if err != nil {
		log.Fatalf("PostComment error: %v.", err)
	}
}
