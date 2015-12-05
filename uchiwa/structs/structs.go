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

// Auth struct contains the generic configuration and details
// about the authentication
type Auth struct {
	Driver     string
	PrivateKey string
	PublicKey  string
}

// CheckExecution struct contains the payload for issuing a
// check execution request to a Sensu API
type CheckExecution struct {
	Check       string   `json:"check"`
	Dc          string   `json:"dc"`
	Subscribers []string `json:"subscribers"`
}

// Data is a structure for holding public data fetched from the Sensu APIs and exposed by the endpoints
type Data struct {
	Aggregates    []interface{}
	Checks        []interface{}
	Clients       []interface{}
	Dc            []*Datacenter
	Events        []interface{}
	Health        Health
	Metrics       Metrics
	SEMetrics     SEMetrics
	SERawMetrics  SERawMetrics `json:"-"`
	Stashes       []interface{}
	Subscriptions []string
}

// Datacenter is a structure for holding the information about a datacenter
type Datacenter struct {
	Name  string         `json:"name"`
	Info  Info           `json:"info"`
	Stats map[string]int `json:"stats"`
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
	Status int    `json:"status"`
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

// Metrics is a structure for holding the metrics of the Sensu objects
type Metrics struct {
	Aggregates  StatusMetrics `json:"aggregates"`
	Checks      StatusMetrics `json:"checks"`
	Clients     StatusMetrics `json:"clients"`
	Datacenters StatusMetrics `json:"datacenters"`
	Events      StatusMetrics `json:"events"`
	Stashes     StatusMetrics `json:"stashes"`
}

// StatusMetrics is a structure for holding the status count
type StatusMetrics struct {
	Critical int `json:"critical"`
	Healthy  int `json:"healthy"`
	Silenced int `json:"silenced"`
	Total    int `json:"total"`
	Unknown  int `json:"unknown"`
	Warning  int `json:"warning"`
}

// SEMetrics is a structure for holding the Sensu Enterprise metrics
type SEMetrics struct {
	Clients         *SEMetric   `json:"clients"`
	Events          []*SEMetric `json:"events"`
	KeepalivesAVG60 *SEMetric   `json:"keepalives_avg_60"`
	Requests        *SEMetric   `json:"requests"`
	Results         *SEMetric   `json:"results"`
}

// SEMetric is a structure for holding a Sensu Enterprise metric
type SEMetric struct {
	Data []XY   `json:"data"`
	Name string `json:"name"`
}

// SERawMetrics ...
type SERawMetrics struct {
	Clients         []*SERawMetric
	Events          []*SERawMetric
	KeepalivesAVG60 []*SERawMetric
	Requests        []*SERawMetric
	Results         []*SERawMetric
}

// SERawMetric ...
type SERawMetric struct {
	Name   string
	Points [][]interface{} `json:"points"`
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

// XPagination is a structure for holding the x-pagination HTTP header
type XPagination struct {
	Limit  int
	Offset int
	Total  int
}

// XY is a structure for holding the coordinates of Sensu Enterprise metrics points
type XY struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
