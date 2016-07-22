package uchiwa

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sensu/uchiwa/uchiwa/authentication"
	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// FilterAggregates is a function that filters aggregates
var FilterAggregates func(aggregates *[]interface{}, token *jwt.Token) []interface{}

// FilterChecks is a function that filters checks
var FilterChecks func(checks *[]interface{}, token *jwt.Token) []interface{}

// FilterClients is a function that filters clients
var FilterClients func(clients *[]interface{}, token *jwt.Token) []interface{}

// FilterAggregates is a function that filters datacenters
var FilterDatacenters func(datacenters []*structs.Datacenter, token *jwt.Token) []*structs.Datacenter

// FilterEvents is a function that filters events
var FilterEvents func(events *[]interface{}, token *jwt.Token) []interface{}

// FilterAggregates is a function that filters aggregates
var FilterStashes func(stashes *[]interface{}, token *jwt.Token) []interface{}

// FilterAggregates is a function that filters aggregates
var FilterSubscriptions func(subscriptions *[]string, token *jwt.Token) []string

// FilterGetRequest is a function that filters GET requests.
var FilterGetRequest func(string, *jwt.Token) bool

// FilterPostRequest is a function that filters POST requests.
var FilterPostRequest func(*jwt.Token, *interface{}) bool

// FilterSensuDataData is a function that filters Sensu Data.
var FilterSensuData func(*jwt.Token, *structs.Data) *structs.Data

// aggregateHandler serves the /aggregates/:check/:issued endpoint
func (u *Uchiwa) aggregateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
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
		checks, err := u.findCheck(name)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusNotFound)
			return
		}

		visibleChecks := FilterChecks(&checks, token)

		if len(visibleChecks) > 1 {
			// Create header
			w.Header().Add("Accept-Charset", "utf-8")
			w.Header().Add("Content-Type", "application/json")

			// If GZIP compression is not supported by the client
			w.WriteHeader(http.StatusMultipleChoices)

			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				encoder := json.NewEncoder(w)
				if err := encoder.Encode(visibleChecks); err != nil {
					http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
					return
				}
				return
			}

			w.Header().Add("Content-Encoding", "gzip")
			w.WriteHeader(http.StatusMultipleChoices)

			gz := gzip.NewWriter(w)
			defer gz.Close()
			if err := json.NewEncoder(gz).Encode(visibleChecks); err != nil {
				http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
				return
			}

			return
		}

		c, ok := checks[0].(map[string]interface{})
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

	unauthorized := FilterGetRequest(dc, token)
	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	var aggregate *map[string]interface{}
	var err error

	if len(resources) == 3 {
		aggregate, err = u.GetAggregate(name, dc)
		if err != nil {
			http.Error(w, fmt.Sprint(err), 500)
			return
		}
	} else {
		issued := resources[3]
		aggregate, err = u.GetAggregateByIssued(name, issued, dc)
		if err != nil {
			http.Error(w, fmt.Sprint(err), 500)
			return
		}
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(aggregate); err != nil {
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
	aggregates := FilterAggregates(&u.Data.Aggregates, token)
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
	checks := FilterChecks(&u.Data.Checks, token)
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

		visibleClients := FilterClients(&clients, token)

		if len(visibleClients) > 1 {
			// Create header
			w.Header().Add("Accept-Charset", "utf-8")
			w.Header().Add("Content-Type", "application/json")

			// If GZIP compression is not supported by the client
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.WriteHeader(http.StatusMultipleChoices)

				encoder := json.NewEncoder(w)
				if err := encoder.Encode(visibleClients); err != nil {
					http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
					return
				}

				return
			}

			w.Header().Add("Content-Encoding", "gzip")
			w.WriteHeader(http.StatusMultipleChoices)

			gz := gzip.NewWriter(w)
			defer gz.Close()
			if err := json.NewEncoder(gz).Encode(visibleClients); err != nil {
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
	unauthorized := FilterGetRequest(dc, token)
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

	clients := FilterClients(&u.Data.Clients, token)
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
	datacenters := FilterDatacenters(u.Data.Dc, token)

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

		visibleClients := FilterClients(&clients, token)

		if len(visibleClients) > 1 {
			// Create header
			w.Header().Add("Accept-Charset", "utf-8")
			w.Header().Add("Content-Type", "application/json")

			// If GZIP compression is not supported by the client
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.WriteHeader(http.StatusMultipleChoices)

				encoder := json.NewEncoder(w)
				if err := encoder.Encode(visibleClients); err != nil {
					http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
					return
				}
				return
			}

			w.Header().Add("Content-Encoding", "gzip")
			w.WriteHeader(http.StatusMultipleChoices)

			gz := gzip.NewWriter(w)
			defer gz.Close()
			if err := json.NewEncoder(gz).Encode(visibleClients); err != nil {
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

	unauthorized := FilterGetRequest(dc, token)
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
	events := FilterEvents(&u.Data.Events, token)
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
	encoder := json.NewEncoder(w)
	var err error
	if r.URL.Path[1:] == "health/sensu" {
		err = encoder.Encode(u.Data.Health.Sensu)
	} else if r.URL.Path[1:] == "health/uchiwa" {
		err = encoder.Encode(u.Data.Health.Uchiwa)
	} else {
		err = encoder.Encode(u.Data.Health)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
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
	unauthorized := FilterGetRequest(data.Dc, token)
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

		visibleClients := FilterClients(&clients, token)

		if len(visibleClients) > 1 {
			// Create header
			w.Header().Add("Accept-Charset", "utf-8")
			w.Header().Add("Content-Type", "application/json")

			// If GZIP compression is not supported by the client
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.WriteHeader(http.StatusMultipleChoices)

				encoder := json.NewEncoder(w)
				if err := encoder.Encode(visibleClients); err != nil {
					http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
					return
				}
				return
			}

			w.Header().Add("Content-Encoding", "gzip")
			w.WriteHeader(http.StatusMultipleChoices)

			gz := gzip.NewWriter(w)
			defer gz.Close()
			if err := json.NewEncoder(gz).Encode(visibleClients); err != nil {
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

	unauthorized := FilterGetRequest(dc, token)
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

		visibleStashes := FilterStashes(&stashes, token)

		if len(visibleStashes) > 1 {
			// Create header
			w.Header().Add("Accept-Charset", "utf-8")
			w.Header().Add("Content-Type", "application/json")

			// If GZIP compression is not supported by the client
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.WriteHeader(http.StatusMultipleChoices)

				encoder := json.NewEncoder(w)
				if err := encoder.Encode(visibleStashes); err != nil {
					http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
					return
				}
				return
			}

			w.Header().Add("Content-Encoding", "gzip")
			w.WriteHeader(http.StatusMultipleChoices)

			gz := gzip.NewWriter(w)
			defer gz.Close()
			if err := json.NewEncoder(gz).Encode(visibleStashes); err != nil {
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

	unauthorized := FilterGetRequest(dc, token)
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

// stashesHandler serves the /stashes endpoint
func (u *Uchiwa) stashesHandler(w http.ResponseWriter, r *http.Request) {
	token := authentication.GetJWTFromContext(r)

	if r.Method == "GET" || r.Method == "HEAD" {
		// GET on /stashes
		stashes := FilterStashes(&u.Data.Stashes, token)
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
		unauthorized := FilterGetRequest(data.Dc, token)
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
	subscriptions := FilterSubscriptions(&u.Data.Subscriptions, token)
	if len(subscriptions) == 0 {
		subscriptions = make([]string, 0)
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(subscriptions); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

// WebServer starts the web server and serves GET & POST requests
func (u *Uchiwa) WebServer(publicPath *string, auth authentication.Config) {
	// Private endpoints
	http.Handle("/aggregates", auth.Authenticate(http.HandlerFunc(u.aggregatesHandler)))
	http.Handle("/aggregates/", auth.Authenticate(http.HandlerFunc(u.aggregateHandler)))
	http.Handle("/checks", auth.Authenticate(http.HandlerFunc(u.checksHandler)))
	http.Handle("/clients", auth.Authenticate(http.HandlerFunc(u.clientsHandler)))
	http.Handle("/clients/", auth.Authenticate(http.HandlerFunc(u.clientHandler)))
	http.Handle("/config", auth.Authenticate(http.HandlerFunc(u.configHandler)))
	http.Handle("/datacenters", auth.Authenticate(http.HandlerFunc(u.datacentersHandler)))
	http.Handle("/events", auth.Authenticate(http.HandlerFunc(u.eventsHandler)))
	http.Handle("/events/", auth.Authenticate(http.HandlerFunc(u.eventHandler)))
	http.Handle("/request", auth.Authenticate(http.HandlerFunc(u.requestHandler)))
	http.Handle("/results/", auth.Authenticate(http.HandlerFunc(u.resultsHandler)))
	http.Handle("/stashes", auth.Authenticate(http.HandlerFunc(u.stashesHandler)))
	http.Handle("/stashes/", auth.Authenticate(http.HandlerFunc(u.stashHandler)))
	http.Handle("/subscriptions", auth.Authenticate(http.HandlerFunc(u.subscriptionsHandler)))
	if u.Config.Uchiwa.Enterprise == false {
		http.Handle("/metrics", auth.Authenticate(http.HandlerFunc(u.metricsHandler)))
	}

	// Static files
	http.Handle("/", http.FileServer(http.Dir(*publicPath)))

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
