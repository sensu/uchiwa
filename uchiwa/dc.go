package uchiwa

import (
	"fmt"
	"sync"
	"time"

	"github.com/bencaron/gosensu"
	"github.com/palourde/logger"
)

type results struct {
	Checks        []interface{}
	Clients       []interface{}
	Dc            []map[string]string
	Events        []interface{}
	Stashes       []interface{}
	Subscriptions []string
	mu            sync.Mutex
}

// PublicConfig contains the public configuration of Uchiwa (hidden user & pass)
var PublicConfig *Config

// Results is a results struct that contains all Sensu APIs data
var Results = new(results)

var tmpResults = new(results)
var datacenters []sensu.Sensu
var mutex sync.Mutex

// Get returns results struct
func (r *results) Get() *results {
	mutex.Lock()
	defer mutex.Unlock()
	return r
}

// Build retrieves all endpoints for every API
func Build(dcSlice *[]sensu.Sensu) {

	for _, api := range *dcSlice {
		Health.Sensu[api.Name] = map[string]string{"output": "connection refused"}

		// fetch sensu data from the API
		stashes, err := api.GetStashes()
		if err != nil {
			logger.Warning(err)
			continue
		}
		checks, err := api.GetChecks()
		if err != nil {
			logger.Warning(err)
			continue
		}
		clients, err := api.GetClients()
		if err != nil {
			logger.Warning(err)
			continue
		}
		events, err := api.GetEvents()
		if err != nil {
			logger.Warning(err)
			continue
		}
		info, err := api.Info()
		if err != nil {
			logger.Warning(err)
			continue
		}

		// add fetched data to results interface
		for _, v := range stashes {
			setDc(v, api.Name)
			tmpResults.Stashes = append(tmpResults.Stashes, v)
		}
		for _, v := range checks {
			setDc(v, api.Name)
			tmpResults.Checks = append(tmpResults.Checks, v)
		}
		for _, v := range clients {
			setDc(v, api.Name)
			tmpResults.Clients = append(tmpResults.Clients, v)
		}
		for _, v := range events {
			setDc(v, api.Name)
			tmpResults.Events = append(tmpResults.Events, v)
		}

		// build dc status
		d := Status(info, api.Name)
		d["checks"] = fmt.Sprintf("%d", len(checks))
		d["clients"] = fmt.Sprintf("%d", len(clients))
		d["events"] = fmt.Sprintf("%d", len(events))
		d["stashes"] = fmt.Sprintf("%d", len(stashes))
		tmpResults.Dc = append(tmpResults.Dc, d)
		Health.Sensu[api.Name] = map[string]string{"output": "ok"}
	}

	BuildEvents()
	BuildClients()
	BuildSubscriptions()

	mutex.Lock()
	defer mutex.Unlock()
	*Results = *tmpResults
}

// Fetch retrieves data from each API every t seconds
func Fetch(t int) {
	Build(&datacenters)
	duration := time.Duration(t) * time.Second
	for _ = range time.Tick(duration) {
		reset()
		Build(&datacenters)
	}
}

// New initialize all Sensu APIs
func New(c *Config) *[]sensu.Sensu {
	buildPublicConfig(c)
	datacenters = make([]sensu.Sensu, len(c.Sensu))
	for i, apiConf := range c.Sensu {
		api := sensu.New(apiConf.Name, apiConf.Path, apiConf.URL, apiConf.Timeout, apiConf.User, apiConf.Pass)
		datacenters[i] = *api
		name := apiConf.Name
		Health.Sensu[name] = map[string]string{"output": "ok"}
	}
	return &datacenters
}

func reset() {
	tmpResults = new(results)
}
