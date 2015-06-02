package uchiwa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/palourde/logger"
	"github.com/sensu/uchiwa/uchiwa/auth"
)

type sensuFilterFn func()

func (u *Uchiwa) configAuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
	}
	fmt.Fprintf(w, "%s", u.PublicConfig.Uchiwa.Auth)
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
	unauthorized := filterGetRequest(dc, token)

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
	unauthorized := filterGetRequest(dc, token)

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
	unauthorized := filterGetRequest(dc, token)

	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	aggregate, err := u.GetAggregateByIssued(check, issued, dc)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
	} else {
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(aggregate); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		}
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
	unauthorized := filterGetRequest(dc, token)

	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	client, err := u.GetClient(id, dc)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusNotFound)
	} else {
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(client); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		}
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
	data := filterSensu(token, u.Data)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
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
	unauthorized := filterPostRequest(token, &data)

	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	err = u.ResolveEvent(data)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	}
}

func (u *Uchiwa) stashHandler(w http.ResponseWriter, r *http.Request) {
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
	unauthorized := filterGetRequest(data.Dc, token)

	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	err = u.PostStash(data)
	if err != nil {
		http.Error(w, "Could not create the stash", http.StatusNotFound)
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
	unauthorized := filterGetRequest(data.Dc, token)

	if unauthorized {
		http.Error(w, fmt.Sprint(""), http.StatusNotFound)
		return
	}

	err = u.DeleteStash(data)
	if err != nil {
		http.Error(w, "Could not create the stash", http.StatusNotFound)
	}
}

// WebServer starts the web server and serves GET & POST requests
func (u *Uchiwa) WebServer(publicPath *string, auth auth.Config) {
	// private endpoints
	http.Handle("/delete_client", auth.Authenticate(http.HandlerFunc(u.deleteClientHandler)))
	http.Handle("/get_aggregate", auth.Authenticate(http.HandlerFunc(u.getAggregateHandler)))
	http.Handle("/get_aggregate_by_issued", auth.Authenticate(http.HandlerFunc(u.getAggregateByIssuedHandler)))
	http.Handle("/get_client", auth.Authenticate(http.HandlerFunc(u.getClientHandler)))
	http.Handle("/get_config", auth.Authenticate(http.HandlerFunc(u.getConfigHandler)))
	http.Handle("/get_sensu", auth.Authenticate(http.HandlerFunc(u.getSensuHandler)))
	http.Handle("/post_event", auth.Authenticate(http.HandlerFunc(u.postEventHandler)))
	http.Handle("/stashes", auth.Authenticate(http.HandlerFunc(u.stashHandler)))
	http.Handle("/stashes/delete", auth.Authenticate(http.HandlerFunc(u.stashDeleteHandler)))

	// static files
	http.Handle("/", http.FileServer(http.Dir(*publicPath)))

	// public endpoints
	http.Handle("/config/auth", http.HandlerFunc(u.configAuthHandler))
	http.Handle("/health", http.HandlerFunc(u.healthHandler))
	http.Handle("/health/", http.HandlerFunc(u.healthHandler))
	http.Handle("/login", auth.GetIdentification())

	listen := fmt.Sprintf("%s:%d", u.Config.Uchiwa.Host, u.Config.Uchiwa.Port)
	logger.Infof("Uchiwa is now listening on %s", listen)
	http.ListenAndServe(listen, nil)
}
