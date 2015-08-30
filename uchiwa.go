package main

import (
	"flag"

	"github.com/sensu/uchiwa/uchiwa"
	"github.com/sensu/uchiwa/uchiwa/audit"
	"github.com/sensu/uchiwa/uchiwa/auth"
	"github.com/sensu/uchiwa/uchiwa/config"
	"github.com/sensu/uchiwa/uchiwa/filters"
	"github.com/sensu/uchiwa/uchiwa/logger"
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

	// Audit
	audit.Log = audit.LogMock

	// filters
	uchiwa.FilterAggregates = filters.FilterAggregates
	uchiwa.FilterChecks = filters.FilterChecks
	uchiwa.FilterClients = filters.FilterClients
	uchiwa.FilterDatacenters = filters.FilterDatacenters
	uchiwa.FilterEvents = filters.FilterEvents
	uchiwa.FilterStashes = filters.FilterStashes
	uchiwa.FilterSubscriptions = filters.FilterSubscriptions

	uchiwa.FilterGetRequest = filters.GetRequest
	uchiwa.FilterPostRequest = filters.PostRequest
	uchiwa.FilterSensuData = filters.SensuData

	u.WebServer(publicPath, authentication)
}
