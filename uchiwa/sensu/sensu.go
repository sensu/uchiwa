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
	Path    string
	URL     string
	Timeout int
	User    string
	Pass    string
	Client  http.Client
}

// NewAPI initializes a new Sensu API struct
func NewAPI(path string, url string, timeout int, username string, password string, insecure bool) API {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}

	client := http.Client{Timeout: time.Duration(timeout) * time.Second, Transport: tr}

	return API{path, url, timeout, username, password, client}
}

// GetName returns the Name attribute
func (s *Sensu) GetName() string {
	return s.Name
}
