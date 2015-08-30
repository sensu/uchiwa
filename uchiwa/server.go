package uchiwa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dgrijalva/jwt-go"
	"github.com/sensu/uchiwa/uchiwa/auth"
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

func (u *Uchiwa) aggregatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := auth.GetTokenFromContext(r)
	stashes := FilterAggregates(&u.Data.Aggregates, token)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(stashes); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

func (u *Uchiwa) checksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := auth.GetTokenFromContext(r)
	checks := FilterChecks(&u.Data.Checks, token)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(checks); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

func (u *Uchiwa) clientsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := auth.GetTokenFromContext(r)
	clients := FilterClients(&u.Data.Clients, token)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(clients); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

func (u *Uchiwa) configAuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
	}
	fmt.Fprintf(w, "%s", u.PublicConfig.Uchiwa.Auth)
}

func (u *Uchiwa) datacentersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := auth.GetTokenFromContext(r)
	datacenters := FilterDatacenters(u.Data.Dc, token)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(datacenters); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

func (u *Uchiwa) deleteClientHandler(w http.ResponseWriter, r *http.Request) {
	urlStruct, _ := url.Parse(r.URL.String())
	id := urlStruct.Query().Get("id")
	dc := urlStruct.Query().Get("dc")
	if id == "" || dc == "" {
		http.Error(w, fmt.Sprint("Parameters 'id' and 'dc' are required"), http.StatusNotFound)
		return
	}

	// verify that the authenticated user is authorized to access this resource
	token := auth.GetTokenFromContext(r)
	unauthorized := FilterGetRequest(dc, token)

	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	err := u.DeleteClient(id, dc)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
}

// events
func (u *Uchiwa) eventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := auth.GetTokenFromContext(r)
	events := FilterEvents(&u.Data.Events, token)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(events); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

func (u *Uchiwa) getAggregateHandler(w http.ResponseWriter, r *http.Request) {
	urlStruct, _ := url.Parse(r.URL.String())
	check := urlStruct.Query().Get("check")
	dc := urlStruct.Query().Get("dc")
	if check == "" || dc == "" {
		http.Error(w, fmt.Sprint("Parameters 'check' and 'dc' are required"), http.StatusNotFound)
		return
	}

	// verify that the authenticated user is authorized to access this resource
	token := auth.GetTokenFromContext(r)
	unauthorized := FilterGetRequest(dc, token)

	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	aggregate, err := u.GetAggregate(check, dc)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
	} else {
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(aggregate); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func (u *Uchiwa) getAggregateByIssuedHandler(w http.ResponseWriter, r *http.Request) {
	urlStruct, _ := url.Parse(r.URL.String())
	check := urlStruct.Query().Get("check")
	issued := urlStruct.Query().Get("issued")
	dc := urlStruct.Query().Get("dc")
	if check == "" || issued == "" || dc == "" {
		http.Error(w, fmt.Sprint("Parameters 'check', 'issued' and 'dc' are required"), http.StatusNotFound)
		return
	}

	// verify that the authenticated user is authorized to access this resource
	token := auth.GetTokenFromContext(r)
	unauthorized := FilterGetRequest(dc, token)

	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	aggregate, err := u.GetAggregateByIssued(check, issued, dc)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(aggregate); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

func (u *Uchiwa) getClientHandler(w http.ResponseWriter, r *http.Request) {
	urlStruct, _ := url.Parse(r.URL.String())
	id := urlStruct.Query().Get("id")
	dc := urlStruct.Query().Get("dc")
	if id == "" || dc == "" {
		http.Error(w, fmt.Sprint("Parameters 'id' and 'dc' are required"), http.StatusNotFound)
		return
	}

	// verify that the authenticated user is authorized to access this resource
	token := auth.GetTokenFromContext(r)
	unauthorized := FilterGetRequest(dc, token)

	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	client, err := u.GetClient(id, dc)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusNotFound)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(client); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

func (u *Uchiwa) getConfigHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(u.PublicConfig); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
	}
}

func (u *Uchiwa) getSensuHandler(w http.ResponseWriter, r *http.Request) {
	token := auth.GetTokenFromContext(r)
	data := FilterSensuData(token, u.Data)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

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
	}
}

func (u *Uchiwa) postEventHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data interface{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprint("Could not decode body"), http.StatusInternalServerError)
		return
	}

	// verify that the authenticated user is authorized to access this resource
	token := auth.GetTokenFromContext(r)
	unauthorized := FilterPostRequest(token, &data)

	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	err = u.ResolveEvent(data)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	}
}

// results
func (u *Uchiwa) resultsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := auth.GetTokenFromContext(r)
	results := FilterChecks(&u.Data.Results, token)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(results); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

// stashes
func (u *Uchiwa) stashesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		token := auth.GetTokenFromContext(r)
		stashes := FilterStashes(&u.Data.Stashes, token)

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(stashes); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var data stash
		err := decoder.Decode(&data)
		if err != nil {
			http.Error(w, "Could not decode body", http.StatusInternalServerError)
			return
		}

		// verify that the authenticated user is authorized to access this resource
		token := auth.GetTokenFromContext(r)
		unauthorized := FilterGetRequest(data.Dc, token)

		if unauthorized {
			http.Error(w, fmt.Sprint(""), http.StatusNotFound)
			return
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

func (u *Uchiwa) stashDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var data stash
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Could not decode body", http.StatusInternalServerError)
		return
	}

	// verify that the authenticated user is authorized to access this resource
	token := auth.GetTokenFromContext(r)
	unauthorized := FilterGetRequest(data.Dc, token)

	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	err = u.DeleteStash(data)
	if err != nil {
		http.Error(w, "Could not create the stash", http.StatusNotFound)
	}
}

// events
func (u *Uchiwa) subscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := auth.GetTokenFromContext(r)
	subscriptions := FilterSubscriptions(&u.Data.Subscriptions, token)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(subscriptions); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

// WebServer starts the web server and serves GET & POST requests
func (u *Uchiwa) WebServer(publicPath *string, auth auth.Config) {

	// private endpoints
	http.Handle("/aggregates", auth.Authenticate(http.HandlerFunc(u.aggregatesHandler)))
	http.Handle("/checks", auth.Authenticate(http.HandlerFunc(u.checksHandler)))
	http.Handle("/clients", auth.Authenticate(http.HandlerFunc(u.clientsHandler)))
	http.Handle("/datacenters", auth.Authenticate(http.HandlerFunc(u.datacentersHandler)))
	http.Handle("/events", auth.Authenticate(http.HandlerFunc(u.eventsHandler)))
	http.Handle("/results", auth.Authenticate(http.HandlerFunc(u.resultsHandler)))
	http.Handle("/stashes", auth.Authenticate(http.HandlerFunc(u.stashesHandler)))
	http.Handle("/stashes/delete", auth.Authenticate(http.HandlerFunc(u.stashDeleteHandler)))
	http.Handle("/subscriptions", auth.Authenticate(http.HandlerFunc(u.subscriptionsHandler)))

	http.Handle("/delete_client", auth.Authenticate(http.HandlerFunc(u.deleteClientHandler)))
	http.Handle("/get_aggregate", auth.Authenticate(http.HandlerFunc(u.getAggregateHandler)))
	http.Handle("/get_aggregate_by_issued", auth.Authenticate(http.HandlerFunc(u.getAggregateByIssuedHandler)))
	http.Handle("/get_client", auth.Authenticate(http.HandlerFunc(u.getClientHandler)))
	http.Handle("/get_config", auth.Authenticate(http.HandlerFunc(u.getConfigHandler)))
	http.Handle("/get_sensu", auth.Authenticate(http.HandlerFunc(u.getSensuHandler)))
	http.Handle("/post_event", auth.Authenticate(http.HandlerFunc(u.postEventHandler)))

	// static files
	http.Handle("/", http.FileServer(http.Dir(*publicPath)))

	// public endpoints
	http.Handle("/config/auth", http.HandlerFunc(u.configAuthHandler))
	http.Handle("/health", http.HandlerFunc(u.healthHandler))
	http.Handle("/health/", http.HandlerFunc(u.healthHandler))
	http.Handle("/login", auth.GetIdentification())

	listen := fmt.Sprintf("%s:%d", u.Config.Uchiwa.Host, u.Config.Uchiwa.Port)
	logger.Infof("Uchiwa is now listening on %s", listen)
	logger.Fatal(http.ListenAndServe(listen, nil))
}
