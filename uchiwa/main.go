package uchiwa

import (
	"sync"
	"time"

	"github.com/sensu/uchiwa/uchiwa/config"
	"github.com/sensu/uchiwa/uchiwa/daemon"
	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/sensu"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// Uchiwa structure is used to manage Uchiwa
type Uchiwa struct {
	Config       *config.Config
	Daemon       *daemon.Daemon
	Data         *structs.Data
	Datacenters  *[]sensu.Sensu
	Mu           *sync.RWMutex
	PublicConfig *config.Config
}

// Init method initializes the Sensu structure with the provided configuration and start the Uchiwa daemon
func Init(c *config.Config) *Uchiwa {

	// Get the datacenters
	datacenters := initDatacenters(c)

	d := &daemon.Daemon{
		Data:        &structs.Data{},
		Datacenters: datacenters,
		Enterprise:  c.Uchiwa.Enterprise,
	}

	u := &Uchiwa{
		Config:       c,
		Daemon:       d,
		Data:         &structs.Data{},
		Datacenters:  datacenters,
		Mu:           &sync.RWMutex{},
		PublicConfig: c.GetPublic(),
	}

	// start Uchiwa daemon and listen for results over data channel
	interval := c.Uchiwa.Refresh
	data := make(chan *structs.Data, 1)
	go d.Start(interval, data)
	go u.listener(interval, data)

	return u
}

// initDatacenters initializes the Datacenters struct by initalizing each
// datacenter based on the provided configuration and by associating multiple
// APIs for the same datacenter for failover/load balancing purposes.
func initDatacenters(c *config.Config) *[]sensu.Sensu {
	var datacenters []sensu.Sensu

OUTER:
	for _, api := range c.Sensu {
		// Initialize the API
		dc := sensu.API{
			CloseRequest:      api.Advanced.CloseRequest,
			DisableKeepAlives: api.Advanced.DisableKeepAlives,
			Insecure:          api.Insecure,
			Pass:              api.Pass,
			Path:              api.Path,
			Timeout:           api.Timeout,
			Tracing:           api.Advanced.Tracing,
			URL:               api.URL,
			User:              api.User,
		}
		dc.Init()

		// Do we already have a datacenter with the same name as this API?
		for i, datacenter := range datacenters {
			if datacenter.Name == api.Name {
				// Add this API to the corresponding datacenter
				datacenter.APIs = append(datacenter.APIs, dc)
				datacenters[i] = datacenter

				continue OUTER
			}
		}
		// At this point we didn't find any datacenter with the same name
		// so we will create a new one and add it to the datacenters slice
		datacenter := sensu.Sensu{Name: api.Name}
		datacenter.APIs = append(datacenter.APIs, dc)
		datacenters = append(datacenters, datacenter)
	}

	return &datacenters
}

// listener listens on the data channel for messages from the daemon
// and updates the Data struct with latest results from the Sensu datacenters
func (u *Uchiwa) listener(interval int, data chan *structs.Data) {
	for {
		select {
		case result := <-data:
			logger.Trace("Received results on the 'data' channel")
			u.Mu.Lock()
			u.Data = result
			u.Mu.Unlock()

			// sleep during the interval
			timer := time.NewTimer(time.Second * time.Duration(interval))
			<-timer.C
		default:
			// sleep during 1 second
			timer := time.NewTimer(time.Second * 1)
			<-timer.C
		}
	}
}
