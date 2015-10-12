// Copyright 2014 Benoit Caron and contributors
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package sensu implements simple methods to interact with Sensu API.

Usage:
	sensu := New("Sensu test API", "", "http://your.SENSU_SERVER_URL.tld"), 15, "username", "secret")
	events = sensu.GetEvents()

The sensu object expose methods that corresponds to the Sensu API:
 http://sensuapp.org/docs/0.13/api_overview

The package follow the form GetApiName() or GetApiName(args, args) where the api
 call wants arguments.

Eg:
	// match http://sensu/clients/
	clients, err := sensu.GetClients()
	// match http://sensu/clients/my_hostname
	a_client,err  := sensu.GetClient("my_hostname")
	// match http://sensu/clients/my_hostname/history
	a_history,err := sensu.GetClientHistory("my_hostname")

Methods returns a data structure and an error object that is nil on successfull query.

API methods that returns an array of results will return a []interface{}.
Methods that return a JSON object will return a map[string]interface{}.
Methods that returns no results on success will return nil, or non-nil (an error object) on errors

No distinction is made yet between server errors (HTTP Status 500) and missing
objects (HTTP status 404)

*/
package sensu

import (
	"crypto/tls"
	"net/http"
	"time"
	//"bytes"
)

// NoLimit is used as a limit parameter
const NoLimit int = -1

// DefaultLimit is used as the default limit parameter for endpoint that supports pagination
const DefaultLimit int = 250

// Sensu struct contains the API details
type Sensu struct {
	Name    string
	Path    string
	URL     string
	Timeout int
	User    string
	Pass    string
	Client  http.Client
}

// New Initialize a new Sensu API
func New(name string, path string, url string, timeout int, username string, password string, insecure bool) *Sensu {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}

	client := http.Client{Timeout: time.Duration(timeout) * time.Second, Transport: tr}

	return &Sensu{name, path, url, timeout, username, password, client}
}
