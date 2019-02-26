// Copyright 2014 Benoit Caron and contributors
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sensu

import (
	"crypto/tls"
	"net/http"
	"time"
)

// NoLimit is used as a limit parameter
const NoLimit int = -1

// DefaultLimit is used as the default limit parameter for endpoint that supports pagination
const DefaultLimit int = 1000

// Sensu struct contains the name and all the APIs for a particular datacenter
type Sensu struct {
	Name string
	APIs []API
}

// API struct contains the details of a specific Sensu API
type API struct {
	CloseRequest      bool
	DisableKeepAlives bool
	Insecure          bool
	Pass              string
	Path              string
	Timeout           int
	Tracing           bool
	URL               string
	User              string
	Healthy           bool
	CheckingHealth    bool

	Client http.Client
}

// Init initializes a new Sensu API HTTP client
func (a *API) Init() {
	tr := &http.Transport{
		DisableKeepAlives: a.DisableKeepAlives,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: a.Insecure},
	}

	client := http.Client{Timeout: time.Duration(a.Timeout) * time.Second, Transport: tr}

	a.Client = client
}

// GetName returns the Name attribute
func (s *Sensu) GetName() string {
	return s.Name
}
