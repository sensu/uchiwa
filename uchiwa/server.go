package uchiwa

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sensu/uchiwa/uchiwa/authentication"
	"github.com/sensu/uchiwa/uchiwa/authorization"
	"github.com/sensu/uchiwa/uchiwa/filters"
	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// Authorization contains the available authorization methods
var Authorization authorization.Authorization

// Filters contains the available filters for the Sensu data
var Filters filters.Filters

// aggregateHandler serves the /aggregates/:name[...] endpoint
func (u *Uchiwa) aggregateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" && r.Method != "DELETE" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resources := strings.Split(r.URL.Path, "/")
	if len(resources) < 3 || resources[2] == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	name := resources[2]
	token := authentication.GetJWTFromContext(r)

	// Get the datacenter name, passed as a query string
	dc := r.URL.Query().Get("dc")

	if dc == "" {
		aggregates, err := u.findAggregate(name)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusNotFound)
			return
		}

		u.Mu.Lock()
		visibleAggregates := Filters.Aggregates(&aggregates, token)
		u.Mu.Unlock()

		if len(visibleAggregates) > 1 {
			// Create header
			w.Header().Add("Accept-Charset", "utf-8")
			w.Header().Add("Content-Type", "application/json")

			// If GZIP compression is not supported by the client
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.WriteHeader(http.StatusMultipleChoices)

				encoder := json.NewEncoder(w)
				if err = encoder.Encode(visibleAggregates); err != nil {
					http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
					return
				}

				return
			}

			w.Header().Add("Content-Encoding", "gzip")
			w.WriteHeader(http.StatusMultipleChoices)

			gz := gzip.NewWriter(w)
			defer gz.Close()
			if err = json.NewEncoder(gz).Encode(visibleAggregates); err != nil {
				http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
				return
			}

			return
		}

		c, ok := aggregates[0].(map[string]interface{})
		if !ok {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		dc, ok = c["dc"].(string)
		if !ok {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}

	unauthorized := Filters.GetRequest(dc, token)
	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	// Are we responding to a /aggregates/:name request?
	if len(resources) == 3 {
		if r.Method == "DELETE" {
			err := u.DeleteAggregate(name, dc)
			if err != nil {
				http.Error(w, fmt.Sprint(err), 500)
				return
			}
			return
		}

		aggregate, err := u.GetAggregate(name, dc)
		if err != nil {
			http.Error(w, fmt.Sprint(err), 500)
			return
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(aggregate); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
		return
	}

	var data *[]interface{}
	var err error

	if len(resources) == 4 {
		// We are responding to a /aggregates/:name/[checks|clients] request

		if resources[3] == "checks" {
			data, err = u.GetAggregateChecks(name, dc)
			if err != nil {
				http.Error(w, fmt.Sprint(err), 500)
				return
			}
		} else if resources[3] == "clients" {
			data, err = u.GetAggregateClients(name, dc)
			if err != nil {
				http.Error(w, fmt.Sprint(err), 500)
				return
			}
		} else {
			http.Error(w, fmt.Sprint(err), http.StatusNotFound)
			return
		}

	} else if len(resources) == 5 {
		// We are responding to a /aggregates/:name/results/:severity request
		severity := resources[4]
		data, err = u.GetAggregateResults(name, severity, dc)
		if err != nil {
			http.Error(w, fmt.Sprint(err), 500)
			return
		}
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}

	return
}

// aggregatesHandler serves the /aggregates endpoint
func (u *Uchiwa) aggregatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := authentication.GetJWTFromContext(r)

	u.Mu.Lock()
	aggregates := Filters.Aggregates(&u.Data.Aggregates, token)
	u.Mu.Unlock()

	if len(aggregates) == 0 {
		aggregates = make([]interface{}, 0)
	}

	// Create header
	w.Header().Add("Accept-Charset", "utf-8")
	w.Header().Add("Content-Type", "application/json")

	// If GZIP compression is not supported by the client
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(aggregates); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Encoding", "gzip")

	gz := gzip.NewWriter(w)
	defer gz.Close()
	if err := json.NewEncoder(gz).Encode(aggregates); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}

	return
}

// checksHandler serves the /checks endpoint
func (u *Uchiwa) checksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := authentication.GetJWTFromContext(r)

	u.Mu.Lock()
	checks := Filters.Checks(&u.Data.Checks, token)
	u.Mu.Unlock()

	if len(checks) == 0 {
		checks = make([]interface{}, 0)
	}

	// Create header
	w.Header().Add("Accept-Charset", "utf-8")
	w.Header().Add("Content-Type", "application/json")

	// If GZIP compression is not supported by the client
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(checks); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Encoding", "gzip")

	gz := gzip.NewWriter(w)
	defer gz.Close()
	if err := json.NewEncoder(gz).Encode(checks); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
	return
}

// clientHandler serves the /clients/:client(/history) endpoint
func (u *Uchiwa) clientHandler(w http.ResponseWriter, r *http.Request) {
	// We only support DELETE & GET requests
	if r.Method != "DELETE" && r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := authentication.GetJWTFromContext(r)

	// Get the client name
	resources := strings.Split(r.URL.Path, "/")
	if len(resources) < 3 || resources[2] == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	name := resources[2]

	// Get the datacenter name, passed as a query string
	dc := r.URL.Query().Get("dc")

	if dc == "" {
		clients, err := u.findClient(name)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusNotFound)
			return
		}

		u.Mu.Lock()
		visibleClients := Filters.Clients(&clients, token)
		u.Mu.Unlock()

		if len(visibleClients) > 1 {
			// Create header
			w.Header().Add("Accept-Charset", "utf-8")
			w.Header().Add("Content-Type", "application/json")

			// If GZIP compression is not supported by the client
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.WriteHeader(http.StatusMultipleChoices)

				encoder := json.NewEncoder(w)
				if err = encoder.Encode(visibleClients); err != nil {
					http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
					return
				}

				return
			}

			w.Header().Add("Content-Encoding", "gzip")
			w.WriteHeader(http.StatusMultipleChoices)

			gz := gzip.NewWriter(w)
			defer gz.Close()
			if err = json.NewEncoder(gz).Encode(visibleClients); err != nil {
				http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
				return
			}

			return
		}

		c, ok := clients[0].(map[string]interface{})
		if !ok {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		dc, ok = c["dc"].(string)
		if !ok {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}

	// Verify that an authenticated user is authorized to access this resource
	unauthorized := Filters.GetRequest(dc, token)
	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	// DELETE on /clients/:client
	if r.Method == "DELETE" {
		err := u.DeleteClient(dc, name)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		return
	}

	// GET on /clients/:client/history
	if len(resources) == 4 {
		data, err := u.GetClientHistory(dc, name)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusNotFound)
			return
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(data); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}

		return
	}

	// GET on /clients/:client
	data, err := u.GetClient(dc, name)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusNotFound)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}

	return
}

// clientsHandler serves the /clients endpoint
func (u *Uchiwa) clientsHandler(w http.ResponseWriter, r *http.Request) {
	// We only support GET requests
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := authentication.GetJWTFromContext(r)

	u.Mu.Lock()
	clients := Filters.Clients(&u.Data.Clients, token)
	u.Mu.Unlock()

	if len(clients) == 0 {
		clients = make([]interface{}, 0)
	}

	// Create header
	w.Header().Add("Accept-Charset", "utf-8")
	w.Header().Add("Content-Type", "application/json")

	// If GZIP compression is not supported by the client
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(clients); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	defer gz.Close()
	if err := json.NewEncoder(gz).Encode(clients); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}

	return
}

// configHandler serves the /config endpoint
func (u *Uchiwa) configHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resources := strings.Split(r.URL.Path, "/")

	if len(resources) == 2 {
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(u.PublicConfig); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		if resources[2] == "auth" {
			fmt.Fprintf(w, "%s", u.PublicConfig.Uchiwa.Auth.Driver)
		} else if resources[2] == "users" {
			encoder := json.NewEncoder(w)
			if err := encoder.Encode(u.PublicConfig.Uchiwa.UsersOptions); err != nil {
				http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "", http.StatusNotFound)
			return
		}
	}
}

// datacentersHandler serves the /datacenters endpoint
func (u *Uchiwa) datacentersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := authentication.GetJWTFromContext(r)
	datacenters := Filters.Datacenters(u.Data.Dc, token)

	// Create header
	w.Header().Add("Accept-Charset", "utf-8")
	w.Header().Add("Content-Type", "application/json")

	// If GZIP compression is not supported by the client
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(datacenters); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Encoding", "gzip")

	gz := gzip.NewWriter(w)
	defer gz.Close()
	if err := json.NewEncoder(gz).Encode(datacenters); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
	return
}

// eventHandler serves the /events/:client/:check endpoint
func (u *Uchiwa) eventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resources := strings.Split(r.URL.Path, "/")
	if len(resources) != 4 {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	check := resources[3]
	client := resources[2]
	token := authentication.GetJWTFromContext(r)

	// Get the datacenter name, passed as a query string
	dc := r.URL.Query().Get("dc")

	if dc == "" {
		clients, err := u.findClient(client)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusNotFound)
			return
		}

		u.Mu.Lock()
		visibleClients := Filters.Clients(&clients, token)
		u.Mu.Unlock()

		if len(visibleClients) > 1 {
			// Create header
			w.Header().Add("Accept-Charset", "utf-8")
			w.Header().Add("Content-Type", "application/json")

			// If GZIP compression is not supported by the client
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.WriteHeader(http.StatusMultipleChoices)

				encoder := json.NewEncoder(w)
				if err = encoder.Encode(visibleClients); err != nil {
					http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
					return
				}
				return
			}

			w.Header().Add("Content-Encoding", "gzip")
			w.WriteHeader(http.StatusMultipleChoices)

			gz := gzip.NewWriter(w)
			defer gz.Close()
			if err = json.NewEncoder(gz).Encode(visibleClients); err != nil {
				http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
				return
			}

			return
		}

		c, ok := clients[0].(map[string]interface{})
		if !ok {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		dc, ok = c["dc"].(string)
		if !ok {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}

	unauthorized := Filters.GetRequest(dc, token)
	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	// DELETE on /events/:client/:check
	err := u.ResolveEvent(check, client, dc)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}

// eventsHandler serves the /events endpoint
func (u *Uchiwa) eventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := authentication.GetJWTFromContext(r)

	u.Mu.Lock()
	events := Filters.Events(&u.Data.Events, token)
	u.Mu.Unlock()

	if len(events) == 0 {
		events = make([]interface{}, 0)
	}

	// Create header
	w.Header().Add("Accept-Charset", "utf-8")
	w.Header().Add("Content-Type", "application/json")

	// If GZIP compression is not supported by the client
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(events); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Encoding", "gzip")

	gz := gzip.NewWriter(w)
	defer gz.Close()
	if err := json.NewEncoder(gz).Encode(events); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}

	return
}

// healthHandler serves the /health endpoint
func (u *Uchiwa) healthHandler(w http.ResponseWriter, r *http.Request) {
	var encoded []byte
	var err error
	returnCode := http.StatusOK

	if r.URL.Path[1:] == "health/sensu" {
		for _, sensu := range u.Data.Health.Sensu {
			if sensu.Output != "ok" {
				returnCode = http.StatusServiceUnavailable
			}
		}
		encoded, err = json.Marshal(u.Data.Health.Sensu)
	} else if r.URL.Path[1:] == "health/uchiwa" {
		if u.Data.Health.Uchiwa != "ok" {
			returnCode = http.StatusServiceUnavailable
		}
		encoded, err = json.Marshal(u.Data.Health.Uchiwa)
	} else {
		for _, sensu := range u.Data.Health.Sensu {
			if sensu.Output != "ok" {
				returnCode = http.StatusServiceUnavailable
			}
		}

		if u.Data.Health.Uchiwa != "ok" {
			returnCode = http.StatusServiceUnavailable
		}

		encoded, err = json.Marshal(u.Data.Health)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(returnCode)
	w.Write(encoded)
	return
}

// metricsHandler serves the /metrics endpoint
func (u *Uchiwa) metricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&u.Data.Metrics); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

// requestHandler serves the /request endpoint
func (u *Uchiwa) requestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var data structs.CheckExecution
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Could not decode body", http.StatusInternalServerError)
		return
	}

	// verify that the authenticated user is authorized to access this resource
	token := authentication.GetJWTFromContext(r)
	unauthorized := Filters.GetRequest(data.Dc, token)
	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	err = u.IssueCheckExecution(data)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	return
}

// resultsHandler serves the /results/:client/:check endpoint
func (u *Uchiwa) resultsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resources := strings.Split(r.URL.Path, "/")
	if len(resources) != 4 {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	check := resources[3]
	client := resources[2]
	token := authentication.GetJWTFromContext(r)

	// Get the datacenter name, passed as a query string
	dc := r.URL.Query().Get("dc")

	if dc == "" {
		clients, err := u.findClient(client)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusNotFound)
			return
		}

		u.Mu.Lock()
		visibleClients := Filters.Clients(&clients, token)
		u.Mu.Unlock()

		if len(visibleClients) > 1 {
			// Create header
			w.Header().Add("Accept-Charset", "utf-8")
			w.Header().Add("Content-Type", "application/json")

			// If GZIP compression is not supported by the client
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.WriteHeader(http.StatusMultipleChoices)

				encoder := json.NewEncoder(w)
				if err = encoder.Encode(visibleClients); err != nil {
					http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
					return
				}
				return
			}

			w.Header().Add("Content-Encoding", "gzip")
			w.WriteHeader(http.StatusMultipleChoices)

			gz := gzip.NewWriter(w)
			defer gz.Close()
			if err = json.NewEncoder(gz).Encode(visibleClients); err != nil {
				http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
				return
			}

			return
		}

		c, ok := clients[0].(map[string]interface{})
		if !ok {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		dc, ok = c["dc"].(string)
		if !ok {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}

	unauthorized := Filters.GetRequest(dc, token)
	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	err := u.DeleteCheckResult(check, client, dc)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}

// stashHandler serves the /stashes/:path endpoint
func (u *Uchiwa) stashHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resources := strings.Split(r.URL.Path, "/")
	if len(resources) < 2 || resources[2] == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	path := strings.Join(resources[2:], "/")
	token := authentication.GetJWTFromContext(r)

	// Get the datacenter name, passed as a query string
	dc := r.URL.Query().Get("dc")

	if dc == "" {
		stashes, err := u.findStash(path)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusNotFound)
			return
		}

		u.Mu.Lock()
		visibleStashes := Filters.Stashes(&stashes, token)
		u.Mu.Unlock()

		if len(visibleStashes) > 1 {
			// Create header
			w.Header().Add("Accept-Charset", "utf-8")
			w.Header().Add("Content-Type", "application/json")

			// If GZIP compression is not supported by the client
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.WriteHeader(http.StatusMultipleChoices)

				encoder := json.NewEncoder(w)
				if err = encoder.Encode(visibleStashes); err != nil {
					http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
					return
				}
				return
			}

			w.Header().Add("Content-Encoding", "gzip")
			w.WriteHeader(http.StatusMultipleChoices)

			gz := gzip.NewWriter(w)
			defer gz.Close()
			if err = json.NewEncoder(gz).Encode(visibleStashes); err != nil {
				http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
				return
			}

			return
		}

		c, ok := stashes[0].(map[string]interface{})
		if !ok {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		dc, ok = c["dc"].(string)
		if !ok {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}

	unauthorized := Filters.GetRequest(dc, token)
	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	err := u.DeleteStash(dc, path)
	if err != nil {
		logger.Warningf("Could not delete the stash '%s': %s", path, err)
		http.Error(w, "Could not create the stash", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}

// silencedHandler serves the /silenced endpoint
func (u *Uchiwa) silencedHandler(w http.ResponseWriter, r *http.Request) {
	token := authentication.GetJWTFromContext(r)

	if r.Method == "GET" || r.Method == "HEAD" {
		// GET on /silenced
		u.Mu.Lock()
		silenced := Filters.Silenced(&u.Data.Silenced, token)
		u.Mu.Unlock()

		if len(silenced) == 0 {
			silenced = make([]interface{}, 0)
		}

		// Create header
		w.Header().Add("Accept-Charset", "utf-8")
		w.Header().Add("Content-Type", "application/json")

		// If GZIP compression is not supported by the client
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			encoder := json.NewEncoder(w)
			if err := encoder.Encode(silenced); err != nil {
				http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
				return
			}
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(w)
		defer gz.Close()
		if err := json.NewEncoder(gz).Encode(silenced); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}

		return
	} else if r.Method == "POST" {
		// POST on /silenced
		decoder := json.NewDecoder(r.Body)
		var data silence
		err := decoder.Decode(&data)
		if err != nil {
			http.Error(w, "Could not decode body", http.StatusInternalServerError)
			return
		}

		// verify that the authenticated user is authorized to access this resource
		unauthorized := Filters.GetRequest(data.Dc, token)
		if unauthorized {
			http.Error(w, fmt.Sprint(""), http.StatusNotFound)
			return
		}

		if token != nil && token.Claims["Username"] != nil {
			data.Creator = token.Claims["Username"].(string)
		}

		resources := strings.Split(r.URL.Path, "/")
		if len(resources) > 2 && resources[2] == "clear" {
			err = u.ClearSilenced(data)
			if err != nil {
				http.Error(w, "Could not clear from entry in the silenced registry", http.StatusNotFound)
				return
			}
			return
		}

		if u.Config.Uchiwa.UsersOptions.DisableNoExpiration && data.Expire < 1 {
			http.Error(w, "Open-ended silence entries are disallowed", http.StatusNotFound)
			return
		}

		if u.Config.Uchiwa.UsersOptions.RequireSilencingReason && data.Reason == "" {
			http.Error(w, "A reason must be provided for every silence entry", http.StatusNotFound)
			return
		}

		err = u.PostSilence(data)
		if err != nil {
			http.Error(w, "Could not create the entry in the silenced registry", http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

// stashesHandler serves the /stashes endpoint
func (u *Uchiwa) stashesHandler(w http.ResponseWriter, r *http.Request) {
	token := authentication.GetJWTFromContext(r)

	if r.Method == "GET" || r.Method == "HEAD" {
		// GET on /stashes
		u.Mu.Lock()
		stashes := Filters.Stashes(&u.Data.Stashes, token)
		u.Mu.Unlock()

		if len(stashes) == 0 {
			stashes = make([]interface{}, 0)
		}

		// Create header
		w.Header().Add("Accept-Charset", "utf-8")
		w.Header().Add("Content-Type", "application/json")

		// If GZIP compression is not supported by the client
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			encoder := json.NewEncoder(w)
			if err := encoder.Encode(stashes); err != nil {
				http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
				return
			}
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(w)
		defer gz.Close()
		if err := json.NewEncoder(gz).Encode(stashes); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}

		return
	} else if r.Method == "POST" {
		// POST on /stashes
		decoder := json.NewDecoder(r.Body)
		var data stash
		err := decoder.Decode(&data)
		if err != nil {
			http.Error(w, "Could not decode body", http.StatusInternalServerError)
			return
		}

		// verify that the authenticated user is authorized to access this resource
		unauthorized := Filters.GetRequest(data.Dc, token)
		if unauthorized {
			http.Error(w, fmt.Sprint(""), http.StatusNotFound)
			return
		}

		if token != nil && token.Claims["Username"] != nil {
			data.Content["username"] = token.Claims["Username"]
		}

		err = u.PostStash(data)
		if err != nil {
			http.Error(w, "Could not create the stash", http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

// subscriptionsHandler serves the /subscriptions endpoint
func (u *Uchiwa) subscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := authentication.GetJWTFromContext(r)

	u.Mu.Lock()
	subscriptions := Filters.Subscriptions(&u.Data.Subscriptions, token)
	u.Mu.Unlock()

	if len(subscriptions) == 0 {
		subscriptions = make([]structs.Subscription, 0)
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(subscriptions); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

// noCacheHandler sets the proper headers to prevent any sort of caching for the
// index.html file, served as /
func noCacheHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("cache-control", "no-cache, no-store, must-revalidate")
			w.Header().Set("pragma", "no-cache")
			w.Header().Set("expires", "0")
		}
		next.ServeHTTP(w, r)
	})
}

// WebServer starts the web server and serves GET & POST requests
func (u *Uchiwa) WebServer(publicPath *string, auth authentication.Config) {
	// Private endpoints
	http.Handle("/aggregates", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.aggregatesHandler))))
	http.Handle("/aggregates/", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.aggregateHandler))))
	http.Handle("/checks", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.checksHandler))))
	http.Handle("/clients", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.clientsHandler))))
	http.Handle("/clients/", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.clientHandler))))
	http.Handle("/config", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.configHandler))))
	http.Handle("/datacenters", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.datacentersHandler))))
	http.Handle("/events", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.eventsHandler))))
	http.Handle("/events/", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.eventHandler))))
	http.Handle("/request", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.requestHandler))))
	http.Handle("/results/", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.resultsHandler))))
	http.Handle("/silenced", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.silencedHandler))))
	http.Handle("/silenced/clear", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.silencedHandler))))
	http.Handle("/stashes", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.stashesHandler))))
	http.Handle("/stashes/", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.stashHandler))))
	http.Handle("/subscriptions", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.subscriptionsHandler))))
	if u.Config.Uchiwa.Enterprise == false {
		http.Handle("/metrics", auth.Authenticate(Authorization.Handler(http.HandlerFunc(u.metricsHandler))))
	}

	// Static files
	http.Handle("/", noCacheHandler(http.FileServer(http.Dir(*publicPath))))

	// Public endpoints
	http.Handle("/config/", http.HandlerFunc(u.configHandler))
	http.Handle("/health", http.HandlerFunc(u.healthHandler))
	http.Handle("/health/", http.HandlerFunc(u.healthHandler))
	http.Handle("/login", auth.Login())

	listen := fmt.Sprintf("%s:%d", u.Config.Uchiwa.Host, u.Config.Uchiwa.Port)
	logger.Warningf("Uchiwa is now listening on %s", listen)

	if u.Config.Uchiwa.SSL.CertFile != "" && u.Config.Uchiwa.SSL.KeyFile != "" {
		logger.Fatal(http.ListenAndServeTLS(listen, u.Config.Uchiwa.SSL.CertFile, u.Config.Uchiwa.SSL.KeyFile, nil))
	}

	logger.Fatal(http.ListenAndServe(listen, nil))
}
