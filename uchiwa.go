package main

import (
	"flag"

	"github.com/sensu/uchiwa/uchiwa"
	"github.com/sensu/uchiwa/uchiwa/audit"
	"github.com/sensu/uchiwa/uchiwa/auth"
	"github.com/sensu/uchiwa/uchiwa/config"
	"github.com/sensu/uchiwa/uchiwa/filters"
)

func main() {
	configFile := flag.String("c", "./config.json", "Full or relative path to the configuration file")
	configDir := flag.String("d", "", "Full or relative path to the configuration directory, or comma delimited directories")
	publicPath := flag.String("p", "public", "Full or relative path to the public directory")
	flag.Parse()

	config := config.Load(*configFile, *configDir)

	u := uchiwa.Init(config)

	authentication := auth.New(config.Uchiwa.Auth)
	if config.Uchiwa.Auth.Driver == "simple" {
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
