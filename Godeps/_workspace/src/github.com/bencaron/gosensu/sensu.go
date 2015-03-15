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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	//"bytes"
)

// Sensu struct contains the API details
type Sensu struct {
	Name     string
	Path     string
	URL      string
	Timeout  int
	User     string
	Pass     string
	Client   http.Client
}

// NoLimit do not specify a limit parameter
const NoLimit int = -1

// NoOffset do not specify an offset parameter
const NoOffset int = -1

// New Initialize a new Sensu API
func New(name string, path string, url string, timeout int, username string, password string, insecure bool) *Sensu {	
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
	
	client := http.Client{Timeout: time.Duration(timeout) * time.Second, Transport: tr}
	
	return &Sensu{name, path, url, timeout, username, password, client}
}

// Health The health endpoint checks to see if the api can connect to redis and rabbitmq. It takes parameters for minimum consumers and maximum messages and checks rabbitmq.
func (s *Sensu) Health(consumers int, messages int) (map[string]interface{}, error) {
	return s.get(fmt.Sprintf("health/%d/%d", consumers, messages))
}

// Info Will return the Sensu version along with rabbitmq and redis information.
func (s *Sensu) Info() (map[string]interface{}, error) {
	return s.get("info")
}

// get ...
func (s *Sensu) get(endpoint string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s", s.URL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Parsing error: %q returned: %v", err, err)
	}
	res, err := s.doHTTP(req)
	if err != nil {
		return nil, fmt.Errorf("API call to %q returned: %v", url, err)
	}
	return s.doJSON(res)
}

// getList Construct an API call and return the list of results
//  LIMITATION: the limit/offset is currently ignored and all results are sent back
func (s *Sensu) getList(endpoint string, limit int, offset int) ([]interface{}, error) {

	/*
		args := ""
		//ERROR GET LIST TODO deal with limit %d and offset %d", limit, offset
		if limit != NOLIMIT {
			args = fmt.Sprintf("%slimit=%d", args, limit)
		}
		if offset != NOOFFSET {
			args = fmt.Sprintf("%soffset=%d", args, limit)
		}
	*/
	url := fmt.Sprintf("%s/%s", s.URL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("URL Parsing error: %q returned: %v", url, err)
	}
	res, err := s.doHTTP(req)
	if err != nil {
		return nil, fmt.Errorf("API call to %q returned: %v", url, err)
	}
	return s.doJSONArray(res)
}

func (s *Sensu) doHTTP(req *http.Request) ([]byte, error) {

	if s.User != "" && s.Pass != "" {
		req.SetBasicAuth(s.User, s.Pass)
	}

	res, err := s.Client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("%v", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("Parsing response body returned: %v", err)
	}
	return body, nil
}

// doJsonArray Unmarshall JSON expecting an array
func (s *Sensu) doJSONArray(body []byte) ([]interface{}, error) {
	var results []interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("Parsing JSON-encoded response body: %v", err)
	}
	return results, nil
}

// doJsonArray Unmarshall JSON expecting a map
func (s *Sensu) doJSON(body []byte) (map[string]interface{}, error) {
	var results map[string]interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("Parsing JSON-encoded response body: %v", err)
	}
	return results, nil
}

// Post to endpoint
func (s *Sensu) post(endpoint string) (map[string]interface{}, error) {
	// Call a List with magic value of limit 0 and offset 0

	//ERROR GET LIST TODO deal with limit %d and offset %d", limit, offset

	url := fmt.Sprintf("%s/%s", s.URL, endpoint)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Parsing error: %q returned: %v", url, err)
	}

	res, err := s.doHTTP(req)
	if err != nil {
		return nil, fmt.Errorf("API call to %q returned: %v", url, err)
	}
	return s.doJSON(res)
}

// postPayload to endpoint
func (s *Sensu) postPayload(endpoint string, payload string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s", s.URL, endpoint)

	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("%s\n\n", payload)))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))

	if err != nil {
		return nil, fmt.Errorf("Parsing error: %q returned: %v", url, err)
	}
	res, err := s.doHTTP(req)
	if err != nil {
		return nil, fmt.Errorf("API call to %q returned: %v", url, err)
	}
	return s.doJSON(res)
}

// Delete resource
func (s *Sensu) delete(endpoint string) error {
	url := fmt.Sprintf("%s/%s", s.URL, endpoint)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Parsing error: %q returned: %v", err, err)
	}

	if s.User != "" && s.Pass != "" {
		req.SetBasicAuth(s.User, s.Pass)
	}

	res, err := s.Client.Do(req)

	if err != nil {
		return fmt.Errorf("API call to %q returned: %v", url, err)
	}
	defer res.Body.Close()

	if err != nil {
		return fmt.Errorf("%v", err)
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("%v", res.Status)
	}
	return nil
}
