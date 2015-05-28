package uchiwa

import (
	"sync"
	"time"

	"github.com/palourde/logger"
	"github.com/sensu/uchiwa/uchiwa/config"
	"github.com/sensu/uchiwa/uchiwa/daemon"
	"github.com/sensu/uchiwa/uchiwa/sensu"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// Uchiwa structure is used to manage Uchiwa
type Uchiwa struct {
	Config       *config.Config
	Daemon       *daemon.Daemon
	Data         *structs.Data
	Datacenters  *[]sensu.Sensu
	Mu           sync.Mutex
	PublicConfig *config.Config
}

// Init method initializes the Sensu structure with the provided configuration and start the Uchiwa daemon
func Init(c *config.Config) *Uchiwa {

	// build datacenters list from configuration
	datacenters := make([]sensu.Sensu, len(c.Sensu))
	for i, datacenter := range c.Sensu {
		datacenter := sensu.New(datacenter.Name, datacenter.Path, datacenter.URL, datacenter.Timeout, datacenter.User, datacenter.Pass, datacenter.Insecure)
		datacenters[i] = *datacenter
	}

	d := &daemon.Daemon{
		Data:        &structs.Data{},
		Datacenters: &datacenters,
	}

	u := &Uchiwa{
		Config:       c,
		Daemon:       d,
		Data:         &structs.Data{},
		Datacenters:  &datacenters,
		PublicConfig: c.GetPublic(),
	}

	// start Uchiwa daemon and listen for results over data channel
	interval := c.Uchiwa.Refresh
	data := make(chan *structs.Data, 1)
	go d.Start(interval, data)
	go u.listener(interval, data)

	return u
}

// listener method listens on the data channel for messages from the daemon
// and updates the Data structure with results from the Sensu APIs
func (u *Uchiwa) listener(interval int, data chan *structs.Data) {
	for {
		select {
		case result := <-data:
			logger.Debug("Received results on the 'data' channel")

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
