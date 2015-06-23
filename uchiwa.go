package main

import (
	"flag"

	"github.com/palourde/logger"
	"github.com/sensu/uchiwa/uchiwa"
	"github.com/sensu/uchiwa/uchiwa/auth"
	"github.com/sensu/uchiwa/uchiwa/config"
)

func main() {
	configFile := flag.String("c", "./config.json", "Full or relative path to the configuration file")
	publicPath := flag.String("p", "public", "Full or relative path to the public directory")
	flag.Parse()

	config, err := config.Load(*configFile)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Debug("Debug mode enabled")

	u := uchiwa.Init(config)

	authentication := auth.New()
	if config.Uchiwa.Auth == "simple" {
		authentication.Simple(config.Uchiwa.Users)
	} else {
		authentication.None()
	}

	u.WebServer(publicPath, authentication)
}
