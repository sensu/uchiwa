package main

import (
	"flag"

	"github.com/palourde/logger"
	"github.com/sensu/uchiwa/auth"
	"github.com/sensu/uchiwa/uchiwa"
)

func main() {
	configFile := flag.String("c", "./config.json", "Full or relative path to the configuration file")
	publicPath := flag.String("p", "public", "Full or relative path to the public directory")
	flag.Parse()

	config, err := uchiwa.LoadConfig(*configFile)
	if err != nil {
		logger.Fatal(err)
	}

	uchiwa.New(config)
	go uchiwa.Fetch(config.Uchiwa.Refresh)

	authentication := auth.New()

	if config.Uchiwa.Auth == "simple" {
		authentication.Simple(config.Uchiwa.User, config.Uchiwa.Pass)
	} else {
		authentication.None()
	}

	uchiwa.WebServer(config, publicPath, authentication)
}
