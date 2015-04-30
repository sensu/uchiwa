package uchiwa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/palourde/logger"
	"github.com/sensu/uchiwa/auth"
)

func configAuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
	}
	fmt.Fprintf(w, "%s", PublicConfig.Uchiwa.Auth)
}

func deleteClientHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(r.URL.String())
	i := u.Query().Get("id")
	d := u.Query().Get("dc")
	if i == "" || d == "" {
		http.Error(w, fmt.Sprint("Parameters 'id' and 'dc' are required"), http.StatusInternalServerError)
	}

	err := DeleteClient(i, d)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	}
}

func getAggregateHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(r.URL.String())
	c := u.Query().Get("check")
	d := u.Query().Get("dc")
	if c == "" || d == "" {
		http.Error(w, fmt.Sprint("Parameters 'check' and 'dc' are required"), 500)
	}

	a, err := GetAggregate(c, d)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
	} else {
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(a); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), 500)
		}
	}
}

func getAggregateByIssuedHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(r.URL.String())
	c := u.Query().Get("check")
	i := u.Query().Get("issued")
	d := u.Query().Get("dc")
	if c == "" || i == "" || d == "" {
		http.Error(w, fmt.Sprint("Parameters 'check', 'issued' and 'dc' are required"), 500)
	}

	a, err := GetAggregateByIssued(c, i, d)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
	} else {
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(a); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), 500)
		}
	}
}

func getClientHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(r.URL.String())
	i := u.Query().Get("id")
	d := u.Query().Get("dc")
	if i == "" || d == "" {
		http.Error(w, fmt.Sprint("Parameters 'id' and 'dc' are required"), http.StatusInternalServerError)
	}

	c, err := GetClient(i, d)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	} else {
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(c); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		}
	}
}

func getConfigHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(PublicConfig); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
	}
}

func getSensuHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(Results.Get()); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	var err error
	if r.URL.Path[1:] == "health/sensu" {
		err = encoder.Encode(Health.Sensu)
	} else if r.URL.Path[1:] == "health/uchiwa" {
		err = encoder.Encode(Health.Uchiwa)
	} else {
		err = encoder.Encode(Health)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
	}
}

func postEventHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data interface{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprint("Could not decode body"), http.StatusInternalServerError)
	}

	err = ResolveEvent(data)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	}
}

func stashHandler(w http.ResponseWriter, r *http.Request) {
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

	err = PostStash(data)
	if err != nil {
		http.Error(w, "Could not create the stash", http.StatusNotFound)
	}
}

func stashDeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	err = DeleteStash(data)
	if err != nil {
		http.Error(w, "Could not create the stash", http.StatusNotFound)
	}
}

// WebServer starts the web server and serves GET & POST requests
func WebServer(config *Config, publicPath *string, auth auth.Config) {
	// private endpoints
	http.Handle("/delete_client", auth.Authenticate(http.HandlerFunc(deleteClientHandler)))
	http.Handle("/get_aggregate", auth.Authenticate(http.HandlerFunc(getAggregateHandler)))
	http.Handle("/get_aggregate_by_issued", auth.Authenticate(http.HandlerFunc(getAggregateByIssuedHandler)))
	http.Handle("/get_client", auth.Authenticate(http.HandlerFunc(getClientHandler)))
	http.Handle("/get_config", auth.Authenticate(http.HandlerFunc(getConfigHandler)))
	http.Handle("/get_sensu", auth.Authenticate(http.HandlerFunc(getSensuHandler)))
	http.Handle("/post_event", auth.Authenticate(http.HandlerFunc(postEventHandler)))
	http.Handle("/stashes", auth.Authenticate(http.HandlerFunc(stashHandler)))
	http.Handle("/stashes/delete", auth.Authenticate(http.HandlerFunc(stashDeleteHandler)))

	// static files
	http.Handle("/", http.FileServer(http.Dir(*publicPath)))

	// public endpoints
	http.Handle("/config/auth", http.HandlerFunc(configAuthHandler))
	http.Handle("/health", http.HandlerFunc(healthHandler))
	http.Handle("/health/", http.HandlerFunc(healthHandler))
	http.Handle("/login", auth.GetIdentification())

	listen := fmt.Sprintf("%s:%d", config.Uchiwa.Host, config.Uchiwa.Port)
	logger.Infof("Uchiwa is now listening on %s", listen)
	http.ListenAndServe(listen, nil)
}
