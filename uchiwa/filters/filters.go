package filters

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// Filters contains the different filtering methods based on the edition
type Filters interface {
	Aggregates(*[]interface{}, *jwt.Token) []interface{}
	Checks(*[]interface{}, *jwt.Token) []interface{}
	Clients(*[]interface{}, *jwt.Token) []interface{}
	Datacenters([]*structs.Datacenter, *jwt.Token) []*structs.Datacenter
	Events(*[]interface{}, *jwt.Token) []interface{}
	// NEED WORK
	GetRequest(string, *jwt.Token) bool
	Stashes(*[]interface{}, *jwt.Token) []interface{}
	Subscriptions(*[]string, *jwt.Token) []string
}

// Uchiwa represents an instance of the Filters interface for the community filters
type Uchiwa struct{}

// Aggregates filters based on role's datacenters
func (u *Uchiwa) Aggregates(data *[]interface{}, token *jwt.Token) []interface{} {
	return *data
}

// Checks filters based on role's datacenters and subscriptions
func (u *Uchiwa) Checks(data *[]interface{}, token *jwt.Token) []interface{} {
	return *data
}

// Clients filters based on role's datacenters and subscriptions
func (u *Uchiwa) Clients(data *[]interface{}, token *jwt.Token) []interface{} {
	return *data
}

// Datacenters filters based on role's datacenters
func (u *Uchiwa) Datacenters(data []*structs.Datacenter, token *jwt.Token) []*structs.Datacenter {
	return data
}

// Events filters based on role's datacenters and subscriptions
func (u *Uchiwa) Events(data *[]interface{}, token *jwt.Token) []interface{} {
	return *data
}

// Stashes filters based on role's datacenters
func (u *Uchiwa) Stashes(data *[]interface{}, token *jwt.Token) []interface{} {
	return *data
}

// Subscriptions filters based on role's subscriptions
func (u *Uchiwa) Subscriptions(data *[]string, token *jwt.Token) []string {
	return *data
}

// GetRequest is a function that filters GET requests.
func (u *Uchiwa) GetRequest(dc string, token *jwt.Token) bool {
	return false
}
