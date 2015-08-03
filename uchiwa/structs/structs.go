package structs

import "time"

// AuditLog is a structure for holding a log of the audit
type AuditLog struct {
	Date       time.Time `json:"date"`
	Action     string    `json:"action"`
	Level      string    `json:"level"`
	Output     string    `json:"output,omitempty"`
	RemoteAddr string    `json:"remoteaddr"`
	URL        string    `json:"url,omitempty"`
	User       string    `json:"user"`
}

// Data is a structure for holding public data fetched from the Sensu APIs and exposed by the endpoints
type Data struct {
	Aggregates    []interface{}
	Checks        []interface{}
	Clients       []interface{}
	Dc            []*Datacenter
	Events        []interface{}
	Health        Health
	Results       []interface{} `json:"-"`
	Stashes       []interface{}
	Subscriptions []string
}

// Datacenter is a structure for holding the information about a datacenter
type Datacenter struct {
	Name  string         `json:"name"`
	Info  Info           `json:"info"`
	Stats map[string]int `json:"stats"`
}

// Dashboard is a struct for holding Dashboard content
type Dashboard struct {
	Checks      DashboardStats `json:"checks"`
	Clients     DashboardStats `json:"clients"`
	Datacenters DashboardStats `json:"datacenters"`
}

// DashboardStats is a struct for holding statistics about a dashboard element
type DashboardStats struct {
	Critical    int `json:"critical"`
	Warning     int `json:"warning"`
	Unknown     int `json:"unknown"`
	Silenced    int `json:"silenced"`
	Healthy     int `json:"healthy"`
	Unsupported int `json:"unsupported"`
	Total       int `json:"total"`
}

// Generic is a structure for holding a generic element
type Generic struct {
	Dc string `json:"dc"`
}

// GenericCheck is a structure for holding a generic check
type GenericCheck struct {
	Dc          string   `json:"dc"`
	Output      string   `json:"output"`
	Status      int      `json:"status"`
	Subscribers []string `json:"subscribers"`
}

// GenericClient is a structure for holding a generic client
type GenericClient struct {
	Dc            string   `json:"dc"`
	Name          string   `json:"name"`
	Subscriptions []string `json:"subscriptions"`
}

// GenericEvent is a structure for holding a generic event
type GenericEvent struct {
	Check  GenericCheck  `json:"check"`
	Client GenericClient `json:"client"`
	Dc     string        `json:"dc"`
}

// Health is a structure for holding health informaton about Sensu & Uchiwa
type Health struct {
	Sensu  map[string]SensuHealth `json:"sensu"`
	Uchiwa string                 `json:"uchiwa"`
}

// SensuHealth is a structure for holding health information about a specific sensu datacenter
type SensuHealth struct {
	Output string `json:"output"`
}

// Info is a structure for holding the /info API information
type Info struct {
	Redis     Redis     `json:"redis"`
	Sensu     Sensu     `json:"sensu"`
	Transport transport `json:"transport"`
}

// Redis is a structure for holding the redis status
type Redis struct {
	Connected bool `json:"connected"`
}

// Sensu is a structure for holding the sensu version
type Sensu struct {
	Version string `json:"version"`
}

type transport struct {
	Connected  bool            `json:"connected"`
	Keepalives transportStatus `json:"keepalives"`
	Results    transportStatus `json:"results"`
}

type transportStatus struct {
	Messages  int `json:"messages"`
	Consumers int `json:"consumers"`
}
