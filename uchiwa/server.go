package uchiwa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/goji/httpauth"
	"github.com/palourde/logger"
)

func deleteClientHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(r.URL.String())
	i := u.Query().Get("id")
	d := u.Query().Get("dc")
	if i == "" || d == "" {
		http.Error(w, fmt.Sprint("Parameters 'id' and 'dc' are required"), 500)
	}

	err := DeleteClient(i, d)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
	}
}

func deleteStashHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data interface{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprint("Could not decode body"), 500)
	}

	err = DeleteStash(data)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
	}
}

func getClientHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(r.URL.String())
	i := u.Query().Get("id")
	d := u.Query().Get("dc")
	if i == "" || d == "" {
		http.Error(w, fmt.Sprint("Parameters 'id' and 'dc' are required"), 500)
	}

	c, err := GetClient(i, d)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
	} else {
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(c); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), 500)
		}
	}
}

func getConfigHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(PublicConfig); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), 500)
	}
}

func getSensuHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(Results.Get()); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), 500)
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
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), 500)
	}
}

func postEventHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data interface{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprint("Could not decode body"), 500)
	}

	err = ResolveEvent(data)

	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
	}
}

func postStashHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data interface{}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprint("Could not decode body"), 500)
	}

	err = CreateStash(data)

	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
	}
}

// WebServer starts the web server and serves GET & POST requests
func WebServer(config *Config, publicPath *string) {
	if config.Uchiwa.User != "" && config.Uchiwa.Pass != "" {
		http.Handle("/delete_client", httpauth.SimpleBasicAuth(config.Uchiwa.User, config.Uchiwa.Pass)(http.HandlerFunc(deleteClientHandler)))
		http.Handle("/delete_stash", httpauth.SimpleBasicAuth(config.Uchiwa.User, config.Uchiwa.Pass)(http.HandlerFunc(deleteStashHandler)))
		http.Handle("/get_client", httpauth.SimpleBasicAuth(config.Uchiwa.User, config.Uchiwa.Pass)(http.HandlerFunc(getClientHandler)))
		http.Handle("/get_config", httpauth.SimpleBasicAuth(config.Uchiwa.User, config.Uchiwa.Pass)(http.HandlerFunc(getConfigHandler)))
		http.Handle("/get_sensu", httpauth.SimpleBasicAuth(config.Uchiwa.User, config.Uchiwa.Pass)(http.HandlerFunc(getSensuHandler)))
		http.Handle("/post_event", httpauth.SimpleBasicAuth(config.Uchiwa.User, config.Uchiwa.Pass)(http.HandlerFunc(postEventHandler)))
		http.Handle("/post_stash", httpauth.SimpleBasicAuth(config.Uchiwa.User, config.Uchiwa.Pass)(http.HandlerFunc(postStashHandler)))
		http.Handle("/", httpauth.SimpleBasicAuth(config.Uchiwa.User, config.Uchiwa.Pass)(http.FileServer(http.Dir(*publicPath))))
	} else {
		http.Handle("/delete_client", http.HandlerFunc(deleteClientHandler))
		http.Handle("/delete_stash", http.HandlerFunc(deleteStashHandler))
		http.Handle("/get_client", http.HandlerFunc(getClientHandler))
		http.Handle("/get_config", http.HandlerFunc(getConfigHandler))
		http.Handle("/get_sensu", http.HandlerFunc(getSensuHandler))
		http.Handle("/post_event", http.HandlerFunc(postEventHandler))
		http.Handle("/post_stash", http.HandlerFunc(postStashHandler))
		http.Handle("/", http.FileServer(http.Dir(*publicPath)))
	}

	http.Handle("/health", http.HandlerFunc(healthHandler))
	http.Handle("/health/", http.HandlerFunc(healthHandler))

	listen := fmt.Sprintf("%s:%d", config.Uchiwa.Host, config.Uchiwa.Port)
	logger.Infof("Uchiwa is now listening on %s", listen)
	http.ListenAndServe(listen, nil)
}
